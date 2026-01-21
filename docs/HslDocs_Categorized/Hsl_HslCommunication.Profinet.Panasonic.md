# HslCommunication - HslCommunication.Profinet.Panasonic

> 分類頁數: 30



---
## HslCommunication.Profinet.Panasonic

[原文連結](http://api.hslcommunication.cn/html/65511854-a56d-47ca-f165-fe601044e6b1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMewtocol 类](../html/032b2c43-56fd-ac73-cf50-2cbc88888d43.htm "PanasonicMewtocol 类")

[PanasonicMewtocolOverTcp 类](../html/018a92fb-02be-5f0d-5c0a-643f0e4327d3.htm "PanasonicMewtocolOverTcp 类")

[PanasonicMewtocolServer 类](../html/77e19c2d-3eaa-b78b-6200-3181ea2fef09.htm "PanasonicMewtocolServer 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Panasonic 命名空间 |

[缺少 "N:HslCommunication.Profinet.Panasonic" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [PanasonicHelper](09335dfa-4192-3521-0b06-2dc4f9834070.htm) | 松下PLC的辅助类，提供了基本的辅助方法，用于解析地址，计算校验和，创建报文  The auxiliary class of Panasonic PLC provides basic auxiliary methods for parsing addresses, calculating checksums, and creating messages |
| 公共类 | [PanasonicMcNet](0653413c-b05c-09c8-dd20-ba8019a910c4.htm) | 松下PLC的数据读写类，基于MC协议的实现，具体的地址格式请参考备注说明  Data reading and writing of Panasonic PLC, based on the implementation of the MC protocol, please refer to the note for specific address format |
| 公共类 | [PanasonicMewtocol](032b2c43-56fd-ac73-cf50-2cbc88888d43.htm) | 松下PLC的数据交互协议，采用Mewtocol协议通讯，支持的地址列表参考api文档  The data exchange protocol of Panasonic PLC adopts Mewtocol protocol for communication. For the list of supported addresses, refer to the api document. |
| 公共类 | [PanasonicMewtocolOverTcp](018a92fb-02be-5f0d-5c0a-643f0e4327d3.htm) | 松下PLC的数据交互协议，采用Mewtocol协议通讯，基于Tcp透传实现的机制，支持的地址列表参考api文档  The data exchange protocol of Panasonic PLC adopts Mewtocol protocol for communication. It is based on the mechanism of Tcp transparent transmission. For the list of supported addresses, refer to the api document. |
| 公共类 | [PanasonicMewtocolServer](77e19c2d-3eaa-b78b-6200-3181ea2fef09.htm) | 松下Mewtocol协议的虚拟服务器，支持串口和网口的操作  Panasonic Mewtocol protocol virtual server, supports serial and network port operations |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicHelper 类

[原文連結](http://api.hslcommunication.cn/html/09335dfa-4192-3521-0b06-2dc4f9834070.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 构造函数](../html/8e6d41b6-8cfe-ee29-e4f1-aaf31af8b480.htm "PanasonicHelper 构造函数 ")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelper 类 |

松下PLC的辅助类，提供了基本的辅助方法，用于解析地址，计算校验和，创建报文  
The auxiliary class of Panasonic PLC provides basic auxiliary methods for parsing addresses, calculating checksums, and creating messages

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.PanasonicPanasonicHelper

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class PanasonicHelper
```

```
Public Class PanasonicHelper
```

```
public ref class PanasonicHelper
```

```
type PanasonicHelper =  class end
```

PanasonicHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [PanasonicHelper](8e6d41b6-8cfe-ee29-e4f1-aaf31af8b480.htm) | 初始化 PanasonicHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisAddress](f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm) | 解析数据地址，解析出地址类型，起始地址  Parse the data address, resolve the address type, start address |
| 公共方法静态成员 | [BuildReadCoils](795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm) | 创建读取多个bool值得报文命令 |
| 公共方法静态成员 | [BuildReadCommand](566c7123-63a5-72b8-299b-a88ccfce9fe3.htm) | 创建批量读取触点的报文指令  Create message instructions for batch reading contacts |
| 公共方法静态成员 | [BuildReadOneCoil](6b2b27ca-6776-7a69-2d9b-76430fbec193.htm) | 创建读取离散触点的报文指令  Create message instructions for reading discrete contacts |
| 公共方法静态成员 | [BuildReadPlcModel](ee306045-4817-0a20-f72a-d32e9011a07f.htm) | 构建获取PLC型号的报文命令 |
| 公共方法静态成员 | [BuildWriteCoils](02cf5c2a-be33-f158-b354-43c99581b852.htm) | 创建写入多个离散触点的报文指令 |
| 公共方法静态成员 | [BuildWriteCommand](65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm) | 创建批量读取触点的报文指令  Create message instructions for batch reading contacts |
| 公共方法静态成员 | [BuildWriteOneCoil](747f31f5-2829-8872-8e4e-cd84bd2c170f.htm) | 创建写入离散触点的报文指令  Create message instructions to write discrete contacts |
| 公共方法静态成员 | [CalculateComplexAddress](febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm) | 位地址转换方法，101等同于10.1等同于10\*16+1=161  Bit address conversion method, 101 is equivalent to 10.1 is equivalent to 10 \* 16 + 1 = 161 |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraActualBool](70877c84-e761-4c53-2fae-74e760f4ea3a.htm) | 检查从PLC反馈的数据，并返回正确的数据内容  Check the data feedback from the PLC and return the correct data content |
| 公共方法静态成员 | [ExtraActualData](ea826a6e-9724-df31-6cc2-6033aa4858cc.htm) | 检查从PLC反馈的数据，并返回正确的数据内容  Check the data feedback from the PLC and return the correct data content |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorDescription](f2f74729-74c2-4f34-5804-56c6288dfd31.htm) | 根据错误码获取到错误描述文本  Get the error description text according to the error code |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法静态成员 | [GetMcErrorDescription](d622432b-d1e5-2543-20f1-4400d2e88562.htm) | 根据MC的错误码去查找对象描述信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackPanasonicCommand](b94c7541-3fe0-03cd-6089-3f14b6730be0.htm) | 将松下的命令打包成带有%开头，CRC校验，CR结尾的完整的命令报文。如果参数 useExpandedHeader 设置为 Ture，则命令头使用 < 开头 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/8e6d41b6-8cfe-ee29-e4f1-aaf31af8b480.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 构造函数](../html/8e6d41b6-8cfe-ee29-e4f1-aaf31af8b480.htm "PanasonicHelper 构造函数 ")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelper 构造函数 |

初始化 [PanasonicHelper](09335dfa-4192-3521-0b06-2dc4f9834070.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public PanasonicHelper()
```

```
Public Sub New
```

```
public:
PanasonicHelper()
```

```
new : unit -> PanasonicHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicHelper 方法

[原文連結](http://api.hslcommunication.cn/html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelper 方法 |

[PanasonicHelper](09335dfa-4192-3521-0b06-2dc4f9834070.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisAddress](f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm) | 解析数据地址，解析出地址类型，起始地址  Parse the data address, resolve the address type, start address |
| 公共方法静态成员 | [BuildReadCoils](795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm) | 创建读取多个bool值得报文命令 |
| 公共方法静态成员 | [BuildReadCommand](566c7123-63a5-72b8-299b-a88ccfce9fe3.htm) | 创建批量读取触点的报文指令  Create message instructions for batch reading contacts |
| 公共方法静态成员 | [BuildReadOneCoil](6b2b27ca-6776-7a69-2d9b-76430fbec193.htm) | 创建读取离散触点的报文指令  Create message instructions for reading discrete contacts |
| 公共方法静态成员 | [BuildReadPlcModel](ee306045-4817-0a20-f72a-d32e9011a07f.htm) | 构建获取PLC型号的报文命令 |
| 公共方法静态成员 | [BuildWriteCoils](02cf5c2a-be33-f158-b354-43c99581b852.htm) | 创建写入多个离散触点的报文指令 |
| 公共方法静态成员 | [BuildWriteCommand](65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm) | 创建批量读取触点的报文指令  Create message instructions for batch reading contacts |
| 公共方法静态成员 | [BuildWriteOneCoil](747f31f5-2829-8872-8e4e-cd84bd2c170f.htm) | 创建写入离散触点的报文指令  Create message instructions to write discrete contacts |
| 公共方法静态成员 | [CalculateComplexAddress](febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm) | 位地址转换方法，101等同于10.1等同于10\*16+1=161  Bit address conversion method, 101 is equivalent to 10.1 is equivalent to 10 \* 16 + 1 = 161 |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraActualBool](70877c84-e761-4c53-2fae-74e760f4ea3a.htm) | 检查从PLC反馈的数据，并返回正确的数据内容  Check the data feedback from the PLC and return the correct data content |
| 公共方法静态成员 | [ExtraActualData](ea826a6e-9724-df31-6cc2-6033aa4858cc.htm) | 检查从PLC反馈的数据，并返回正确的数据内容  Check the data feedback from the PLC and return the correct data content |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorDescription](f2f74729-74c2-4f34-5804-56c6288dfd31.htm) | 根据错误码获取到错误描述文本  Get the error description text according to the error code |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法静态成员 | [GetMcErrorDescription](d622432b-d1e5-2543-20f1-4400d2e88562.htm) | 根据MC的错误码去查找对象描述信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackPanasonicCommand](b94c7541-3fe0-03cd-6089-3f14b6730be0.htm) | 将松下的命令打包成带有%开头，CRC校验，CR结尾的完整的命令报文。如果参数 useExpandedHeader 设置为 Ture，则命令头使用 < 开头 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AnalysisAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperAnalysisAddress 方法 |

解析数据地址，解析出地址类型，起始地址  
Parse the data address, resolve the address type, start address

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string, int> AnalysisAddress(
	string address
)
```

```
Public Shared Function AnalysisAddress ( 
	address As String
) As OperateResult(Of String, Integer)
```

```
public:
static OperateResult<String^, int>^ AnalysisAddress(
	String^ address
)
```

```
static member AnalysisAddress : 
        address : string -> OperateResult<string, int> 
```

#### 参数

address
:   类型：SystemString  
    数据地址

#### 返回值

类型：[OperateResult](f52f888f-5e8d-b0c4-2302-81e70c230a31.htm)String, Int32  
解析出地址类型，起始地址

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCoils 方法 

[原文連結](http://api.hslcommunication.cn/html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperBuildReadCoils 方法 |

创建读取多个bool值得报文命令

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<List<byte[]>> BuildReadCoils(
	byte station,
	string[] address
)
```

```
Public Shared Function BuildReadCoils ( 
	station As Byte,
	address As String()
) As OperateResult(Of List(Of Byte()))
```

```
public:
static OperateResult<List<array<unsigned char>^>^>^ BuildReadCoils(
	unsigned char station, 
	array<String^>^ address
)
```

```
static member BuildReadCoils : 
        station : byte * 
        address : string[] -> OperateResult<List<byte[]>> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    等待读取的地址数组

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)ListByte  
包含是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperBuildReadCommand 方法 |

创建批量读取触点的报文指令  
Create message instructions for batch reading contacts

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<List<byte[]>> BuildReadCommand(
	byte station,
	string address,
	ushort length,
	bool isBit
)
```

```
Public Shared Function BuildReadCommand ( 
	station As Byte,
	address As String,
	length As UShort,
	isBit As Boolean
) As OperateResult(Of List(Of Byte()))
```

```
public:
static OperateResult<List<array<unsigned char>^>^>^ BuildReadCommand(
	unsigned char station, 
	String^ address, 
	unsigned short length, 
	bool isBit
)
```

```
static member BuildReadCommand : 
        station : byte * 
        address : string * 
        length : uint16 * 
        isBit : bool -> OperateResult<List<byte[]>> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    地址信息

length
:   类型：SystemUInt16  
    数据长度

isBit
:   类型：SystemBoolean  
    是否进行位为单位

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)ListByte  
包含是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadOneCoil 方法 

[原文連結](http://api.hslcommunication.cn/html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperBuildReadOneCoil 方法 |

创建读取离散触点的报文指令  
Create message instructions for reading discrete contacts

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildReadOneCoil(
	byte station,
	string address
)
```

```
Public Shared Function BuildReadOneCoil ( 
	station As Byte,
	address As String
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadOneCoil(
	unsigned char station, 
	String^ address
)
```

```
static member BuildReadOneCoil : 
        station : byte * 
        address : string -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadPlcModel 方法 

[原文連結](http://api.hslcommunication.cn/html/ee306045-4817-0a20-f72a-d32e9011a07f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperBuildReadPlcModel 方法 |

构建获取PLC型号的报文命令

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildReadPlcModel(
	byte station
)
```

```
Public Shared Function BuildReadPlcModel ( 
	station As Byte
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadPlcModel(
	unsigned char station
)
```

```
static member BuildReadPlcModel : 
        station : byte -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取PLC型号的命令报文信息

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCoils 方法 

[原文連結](http://api.hslcommunication.cn/html/02cf5c2a-be33-f158-b354-43c99581b852.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperBuildWriteCoils 方法 |

创建写入多个离散触点的报文指令

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<List<byte[]>> BuildWriteCoils(
	byte station,
	string[] address,
	bool[] value
)
```

```
Public Shared Function BuildWriteCoils ( 
	station As Byte,
	address As String(),
	value As Boolean()
) As OperateResult(Of List(Of Byte()))
```

```
public:
static OperateResult<List<array<unsigned char>^>^>^ BuildWriteCoils(
	unsigned char station, 
	array<String^>^ address, 
	array<bool>^ value
)
```

```
static member BuildWriteCoils : 
        station : byte * 
        address : string[] * 
        value : bool[] -> OperateResult<List<byte[]>> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    等待写入的地址列表

value
:   类型：SystemBoolean  
    等待写入的值列表，长度应和地址长度一致

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)ListByte  
所有写入命令的报文列表

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperBuildWriteCommand 方法 |

创建批量读取触点的报文指令  
Create message instructions for batch reading contacts

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteCommand(
	byte station,
	string address,
	byte[] values
)
```

```
Public Shared Function BuildWriteCommand ( 
	station As Byte,
	address As String,
	values As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteCommand(
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ values
)
```

```
static member BuildWriteCommand : 
        station : byte * 
        address : string * 
        values : byte[] -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    设备站号

address
:   类型：SystemString  
    地址信息

values
:   类型：SystemByte  
    数据值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteOneCoil 方法 

[原文連結](http://api.hslcommunication.cn/html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperBuildWriteOneCoil 方法 |

创建写入离散触点的报文指令  
Create message instructions to write discrete contacts

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteOneCoil(
	byte station,
	string address,
	bool value
)
```

```
Public Shared Function BuildWriteOneCoil ( 
	station As Byte,
	address As String,
	value As Boolean
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteOneCoil(
	unsigned char station, 
	String^ address, 
	bool value
)
```

```
static member BuildWriteOneCoil : 
        station : byte * 
        address : string * 
        value : bool -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    地址信息

value
:   类型：SystemBoolean  
    bool值数组

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CalculateComplexAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperCalculateComplexAddress 方法 |

位地址转换方法，101等同于10.1等同于10\*16+1=161  
Bit address conversion method, 101 is equivalent to 10.1 is equivalent to 10 \* 16 + 1 = 161

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static int CalculateComplexAddress(
	string address,
	int fromBase = 16
)
```

```
Public Shared Function CalculateComplexAddress ( 
	address As String,
	Optional fromBase As Integer = 16
) As Integer
```

```
public:
static int CalculateComplexAddress(
	String^ address, 
	int fromBase = 16
)
```

```
static member CalculateComplexAddress : 
        address : string * 
        ?fromBase : int 
(* Defaults:
        let _fromBase = defaultArg fromBase 16
*)
-> int 
```

#### 参数

address
:   类型：SystemString  
    地址信息

fromBase (Optional)
:   类型：SystemInt32  
    倍率信息

#### 返回值

类型：Int32  
实际的位地址信息

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraActualBool 方法 

[原文連結](http://api.hslcommunication.cn/html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperExtraActualBool 方法 |

检查从PLC反馈的数据，并返回正确的数据内容  
Check the data feedback from the PLC and return the correct data content

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<bool[]> ExtraActualBool(
	byte[] response
)
```

```
Public Shared Function ExtraActualBool ( 
	response As Byte()
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ExtraActualBool(
	array<unsigned char>^ response
)
```

```
static member ExtraActualBool : 
        response : byte[] -> OperateResult<bool[]> 
```

#### 参数

response
:   类型：SystemByte  
    反馈信号

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
是否成功的结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraActualData 方法 

[原文連結](http://api.hslcommunication.cn/html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperExtraActualData 方法 |

检查从PLC反馈的数据，并返回正确的数据内容  
Check the data feedback from the PLC and return the correct data content

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> ExtraActualData(
	byte[] response,
	bool parseData = true
)
```

```
Public Shared Function ExtraActualData ( 
	response As Byte(),
	Optional parseData As Boolean = true
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ ExtraActualData(
	array<unsigned char>^ response, 
	bool parseData = true
)
```

```
static member ExtraActualData : 
        response : byte[] * 
        ?parseData : bool 
(* Defaults:
        let _parseData = defaultArg parseData true
*)
-> OperateResult<byte[]> 
```

#### 参数

response
:   类型：SystemByte  
    反馈信号

parseData (Optional)
:   类型：SystemBoolean  
    是否解析数据内容部分

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
是否成功的结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorDescription 方法 

[原文連結](http://api.hslcommunication.cn/html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperGetErrorDescription 方法 |

根据错误码获取到错误描述文本  
Get the error description text according to the error code

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetErrorDescription(
	int err
)
```

```
Public Shared Function GetErrorDescription ( 
	err As Integer
) As String
```

```
public:
static String^ GetErrorDescription(
	int err
)
```

```
static member GetErrorDescription : 
        err : int -> string 
```

#### 参数

err
:   类型：SystemInt32  
    错误代码

#### 返回值

类型：String  
字符信息

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetMcErrorDescription 方法 

[原文連結](http://api.hslcommunication.cn/html/d622432b-d1e5-2543-20f1-4400d2e88562.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperGetMcErrorDescription 方法 |

根据MC的错误码去查找对象描述信息

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetMcErrorDescription(
	int code
)
```

```
Public Shared Function GetMcErrorDescription ( 
	code As Integer
) As String
```

```
public:
static String^ GetMcErrorDescription(
	int code
)
```

```
static member GetMcErrorDescription : 
        code : int -> string 
```

#### 参数

code
:   类型：SystemInt32  
    错误码

#### 返回值

类型：String  
描述信息

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackPanasonicCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicHelper 类](../html/09335dfa-4192-3521-0b06-2dc4f9834070.htm "PanasonicHelper 类")

[PanasonicHelper 方法](../html/aa1c2b72-1051-c32d-216b-b0eb47f5f164.htm "PanasonicHelper 方法")

[AnalysisAddress 方法](../html/f7d8ab5e-8c6a-9a8c-ddea-2076dab00bbd.htm "AnalysisAddress 方法 ")

[BuildReadCoils 方法](../html/795ef8a2-9290-19e9-5f89-4e0b6ad63c0a.htm "BuildReadCoils 方法 ")

[BuildReadCommand 方法](../html/566c7123-63a5-72b8-299b-a88ccfce9fe3.htm "BuildReadCommand 方法 ")

[BuildReadOneCoil 方法](../html/6b2b27ca-6776-7a69-2d9b-76430fbec193.htm "BuildReadOneCoil 方法 ")

[BuildReadPlcModel 方法](../html/ee306045-4817-0a20-f72a-d32e9011a07f.htm "BuildReadPlcModel 方法 ")

[BuildWriteCoils 方法](../html/02cf5c2a-be33-f158-b354-43c99581b852.htm "BuildWriteCoils 方法 ")

[BuildWriteCommand 方法](../html/65b0cc98-20a1-5e81-2f06-e7cf8bdf7b20.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCoil 方法](../html/747f31f5-2829-8872-8e4e-cd84bd2c170f.htm "BuildWriteOneCoil 方法 ")

[CalculateComplexAddress 方法](../html/febd65af-0c75-9ba7-9b9c-8a88d060f3c6.htm "CalculateComplexAddress 方法 ")

[ExtraActualBool 方法](../html/70877c84-e761-4c53-2fae-74e760f4ea3a.htm "ExtraActualBool 方法 ")

[ExtraActualData 方法](../html/ea826a6e-9724-df31-6cc2-6033aa4858cc.htm "ExtraActualData 方法 ")

[GetErrorDescription 方法](../html/f2f74729-74c2-4f34-5804-56c6288dfd31.htm "GetErrorDescription 方法 ")

[GetMcErrorDescription 方法](../html/d622432b-d1e5-2543-20f1-4400d2e88562.htm "GetMcErrorDescription 方法 ")

[PackPanasonicCommand 方法](../html/b94c7541-3fe0-03cd-6089-3f14b6730be0.htm "PackPanasonicCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicHelperPackPanasonicCommand 方法 |

将松下的命令打包成带有%开头，CRC校验，CR结尾的完整的命令报文。如果参数 useExpandedHeader 设置为 Ture，则命令头使用 < 开头

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> PackPanasonicCommand(
	byte station,
	string cmd,
	bool useExpandedHeader
)
```

```
Public Shared Function PackPanasonicCommand ( 
	station As Byte,
	cmd As String,
	useExpandedHeader As Boolean
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ PackPanasonicCommand(
	unsigned char station, 
	String^ cmd, 
	bool useExpandedHeader
)
```

```
static member PackPanasonicCommand : 
        station : byte * 
        cmd : string * 
        useExpandedHeader : bool -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

cmd
:   类型：SystemString  
    松下的命令。例如 RCSR100F

useExpandedHeader
:   类型：SystemBoolean  
    设置是否使用扩展的命令头消息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
原始的字节数组的命令

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicHelper 类](09335dfa-4192-3521-0b06-2dc4f9834070.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicMcNet 类

[原文連結](http://api.hslcommunication.cn/html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 构造函数](../html/6841fcea-208c-6a3d-1dbd-7975ccf9a958.htm "PanasonicMcNet 构造函数 ")

[PanasonicMcNet 属性](../html/e3819fe2-e79d-c9a7-cdf1-2fcad2fb6457.htm "PanasonicMcNet 属性")

[PanasonicMcNet 方法](../html/1a9a79ba-af60-cae4-816c-d8ff275d5c34.htm "PanasonicMcNet 方法")

[PanasonicMcNet 字段](../html/6d8890da-c221-26cc-7bea-229ff00ad1e5.htm "PanasonicMcNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNet 类 |

松下PLC的数据读写类，基于MC协议的实现，具体的地址格式请参考备注说明  
Data reading and writing of Panasonic PLC, based on the implementation of the MC protocol, please refer to the note for specific address format

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
        [HslCommunication.Profinet.MelsecMelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)  
          HslCommunication.Profinet.PanasonicPanasonicMcNet

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class PanasonicMcNet : MelsecMcNet
```

```
Public Class PanasonicMcNet
	Inherits MelsecMcNet
```

```
public ref class PanasonicMcNet : public MelsecMcNet
```

```
type PanasonicMcNet =  
    class
        inherit MelsecMcNet
    end
```

PanasonicMcNet 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [PanasonicMcNet](8a89cbff-9483-171b-a929-b878572d3bb3.htm) | 实例化松下的的Qna兼容3E帧协议的通讯对象  Instantiate Panasonic's Qna compatible 3E frame protocol communication object |
| 公共方法 | [PanasonicMcNet(String, Int32)](98e3bd47-d9a1-9a8d-fbaa-4b87d4ddbe37.htm) | 指定ip地址及端口号来实例化一个松下的Qna兼容3E帧协议的通讯对象  Specify an IP address and port number to instantiate a Panasonic Qna compatible 3E frame protocol communication object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [EnableWriteBitToWordRegister](c933b3aa-a807-7364-0f0a-7ce1296085c5.htm) | 是否开启支持写入位到字寄存器的功能，该功能先读取字寄存器的字数据，然后修改其中的位，再写入回去，可能存在脏数据的风险  Whether to enable supporting the function of writing bit-to-word registers, which first reads the word data of the word register, then modifies the bits in it, and then writes back, there may be a risk of dirty data (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [McType](bf810248-dfa8-a546-5ec3-c92ec37722a4.htm) | 当前的MC协议的格式类型  The format type of the current MC protocol (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性 | [NetworkNumber](7eb9cb50-4c13-3151-627e-5e3cde56fe03.htm) | 网络号，通常为0  Network number, usually 0 (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性 | [NetworkStationNumber](d3f7232f-5184-fb8b-29a2-0a7007213561.htm) | 网络站号，通常为0  Network station number, usually 0 (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性 | [PLCNumber](f087a68e-3dc8-bc74-5a06-4294d9d2069a.htm) | PLC编号，如果是本站信息，则是 0xFF 值，其他站则根据实际的情况指定。  The PLC number, if it is the information of this station, is the 0xFF value, and other stations are specified according to the actual situation. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [TargetIOStation](ad8258fa-31ad-5240-417d-97db45ea311b.htm) | 请求目标模块的IO编号，默认是管理CPU，也就是 0x03FF 的值，如果需要访问其他的非管理CPU的时候，请参考手册进行配置相关的值  The IO number of the target module is the default management CPU, that is, the value of 0x03FF, if you need to access other non-management CPUs, please refer to the manual to configure the relevant values (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

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
| 公共方法 | [ErrorStateReset](dff5698a-7453-36db-0774-a7bf9249dfeb.htm) | LED 熄灭 出错代码初始化  LED off Error code initialization (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ErrorStateResetAsync](7ef2fa94-e871-2631-bed6-679a23591f92.htm) | LED 熄灭 出错代码初始化  LED off Error code initialization (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ExtractActualData](553f16db-7e99-d2da-265d-4d7f63a598e5.htm) | 从PLC反馈的数据中提取出实际的数据内容，需要传入反馈数据，是否位读取 (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](07258aec-2c66-b9d2-bbf3-b2f3eea3f1dd.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [McAnalysisAddress](23dea1af-6ed3-1b06-54e2-e65b7b32e173.htm) | 当前MC协议的分析地址的方法，对传入的字符串格式的地址进行数据解析。  The current MC protocol's address analysis method performs data parsing on the address of the incoming string format. (重写 [MelsecMcNetMcAnalysisAddress(String, UInt16, Boolean)](2e6fc94e-25bc-0ab8-5590-e334bc02148d.htm).) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](ed03f8d8-4e79-6ee1-a845-ea69a977dd03.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [Read(String, UInt16)](7c89c314-9fbf-a77a-c990-13ef2393c47d.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](c86b56d0-2661-6877-5cf9-ea980d039844.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](b5a3299a-2a1c-644b-aa1d-df14378eae96.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](43c50faf-993e-f3ac-2321-7f846ce458f3.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](a69ef7ef-c8a6-e815-cf29-586acdc173a5.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](4ced354d-18c9-6203-7ebb-ef056b07f6c2.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadExtend](f1f50a27-f4f5-ba75-0713-d64787cf5be0.htm) | **[商业授权]** 读取扩展的数据信息，需要在原有的地址，长度信息之外，输入扩展值信息 **[Authorization]** To read the extended data information, you need to enter the extended value information in addition to the original address and length information (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadExtendAsync](db4dd3cc-2cc6-8629-68b2-8e8233e11370.htm) | **[商业授权]** 读取扩展的数据信息，需要在原有的地址，长度信息之外，输入扩展值信息 **[Authorization]** To read the extended data information, you need to enter the extended value information in addition to the original address and length information (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
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
| 公共方法 | [ReadMemory](f58f02d7-a2ed-3958-3304-77df7fa05dd0.htm) | **[商业授权]** 读取缓冲寄存器的数据信息，地址直接为偏移地址 **[Authorization]** Read the data information of the buffer register, the address is directly the offset address (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadMemoryAsync](6f966786-0a21-46bd-c753-7d7159541458.htm) | **[商业授权]** 读取缓冲寄存器的数据信息，地址直接为偏移地址 **[Authorization]** Read the data information of the buffer register, the address is directly the offset address (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadPlcType](9f8aaf42-bf8c-764b-a3df-3a720da484c5.htm) | 读取PLC的型号信息，例如 Q02HCPU  Read PLC model information, such as Q02HCPU (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadPlcTypeAsync](716c2e02-7df6-18bd-c689-7eb57368ebd9.htm) | 读取PLC的型号信息，例如 Q02HCPU  Read PLC model information, such as Q02HCPU (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadRandom(String)](2b700dc2-727c-7541-cd80-8781ff171e71.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，需要自行解析数据  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the results, you need to parse the data yourself (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadRandom(String, UInt16)](1c644e2a-8674-57fd-cdc8-97d22edea925.htm) | 使用块读取PLC的数据信息，可以跨地址，跨类型组合，每个地址是任意的长度。收到结果后，需要自行解析数据，目前只支持字地址，比如D区，W区，R区，不支持X，Y，M，B，L等等  Read the data information of the PLC randomly. It can be combined across addresses and types. Each address is of any length. After receiving the results, you need to parse the data yourself. Currently, only word addresses are supported, such as D area, W area, R area. X, Y, M, B, L, etc (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadRandomAsync(String)](0ab9c6f2-7488-2f1d-feca-8bcb418739ab.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，需要自行解析数据  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the results, you need to parse the data yourself (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadRandomAsync(String, UInt16)](23278f65-11c2-0c87-5583-2540cdf16a2b.htm) | 使用块读取PLC的数据信息，可以跨地址，跨类型组合，每个地址是任意的长度。收到结果后，需要自行解析数据，目前只支持字地址，比如D区，W区，R区，不支持X，Y，M，B，L等等  Read the data information of the PLC randomly. It can be combined across addresses and types. Each address is of any length. After receiving the results, you need to parse the data yourself. Currently, only word addresses are supported, such as D area, W area, R area. X, Y, M, B, L, etc (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadRandomInt16](4338885c-8abf-8bfd-d553-9e57390f185c.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，自动转换为了short类型的数组  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the result, it is automatically converted to an array of type short. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadRandomInt16Async](0d680060-284c-3a6f-b245-c9855ba7a1b0.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，自动转换为了short类型的数组  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the result, it is automatically converted to an array of type short. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadRandomUInt16](c4d126ae-e24d-34b5-abf3-9e6e08492f45.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，自动转换为了ushort类型的数组  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the result, it is automatically converted to an array of type ushort. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadRandomUInt16Async](7faa94e9-a1f4-b321-8c1f-bc64c715a438.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，自动转换为了ushort类型的数组  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the result, it is automatically converted to an array of type ushort. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadSmartModule](00851b3c-e4b5-b306-d624-f14f9501d3df.htm) | **[商业授权]** 读取智能模块的数据信息，需要指定模块地址，偏移地址，读取的字节长度 **[Authorization]** To read the extended data information, you need to enter the extended value information in addition to the original address and length information (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadSmartModuleAsync](eecc3de1-8d08-3067-45b2-2f93e43b966f.htm) | **[商业授权]** 读取智能模块的数据信息，需要指定模块地址，偏移地址，读取的字节长度 **[Authorization]** To read the extended data information, you need to enter the extended value information in addition to the original address and length information (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadTags(String, UInt16)](d7e5e3e8-e3bf-2a02-7d90-561cf28138ea.htm) | **[商业授权]** 读取PLC的标签信息，需要传入标签的名称，读取的字长度，标签举例：A; label[1]; bbb[10,10,10] **[Authorization]** To read the label information of the PLC, you need to pass in the name of the label, the length of the word read, and an example of the label: A; label [1]; bbb [10,10,10] (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadTags(String, UInt16)](69e4bc2a-e235-1294-4653-b6a2dec1535c.htm) | **[商业授权]** 读取PLC的标签信息，需要传入标签的名称，读取的字长度，标签举例：A; label[1]; bbb[10,10,10] **[Authorization]** To read the label information of the PLC, you need to pass in the name of the label, the length of the word read, and an example of the label: A; label [1]; bbb [10,10,10] (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadTagsAsync(String, UInt16)](b2edcbab-1f74-1963-fc00-18aa8c8c02d7.htm) | **[商业授权]** 读取PLC的标签信息，需要传入标签的名称，读取的字长度，标签举例：A; label[1]; bbb[10,10,10] **[Authorization]** To read the label information of the PLC, you need to pass in the name of the label, the length of the word read, and an example of the label: A; label [1]; bbb [10,10,10] (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadTagsAsync(String, UInt16)](01ce89ac-b4c1-c2ed-25dc-662efae546e6.htm) | **[商业授权]** 读取PLC的标签信息，需要传入标签的名称，读取的字长度，标签举例：A; label[1]; bbb[10,10,10] **[Authorization]** To read the label information of the PLC, you need to pass in the name of the label, the length of the word read, and an example of the label: A; label [1]; bbb [10,10,10] (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
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
| 公共方法 | [RemoteReset](08535f06-ae9f-fca5-1df3-3299c505f76f.htm) | 远程Reset操作  Remote Reset Operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteResetAsync](cca35551-499a-b900-810f-e5bf5c50e3ec.htm) | 远程Reset操作  Remote Reset Operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteRun](2f2f4e7f-17e2-63fd-882d-59e20efa8ef7.htm) | 远程Run操作  Remote Run Operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteRunAsync](da1c67d7-6063-6ed7-d505-29e856404d54.htm) | 远程Run操作  Remote Run Operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteStop](e199a595-8ae0-3ee3-383d-cc7b5f81430b.htm) | 远程Stop操作  Remote Stop operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteStopAsync](06a7c9b0-0c6a-ce1f-77ea-047a4ea5e585.htm) | 远程Stop操作  Remote Stop operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](f2b1932b-0e3c-cf1f-51a5-ab87bc63da5e.htm) | (重写 [MelsecMcNetToString](77b9b8f7-46b7-1f45-fa1e-3e7b505b61e4.htm).) |
| 公共方法 | [UnpackResponseContent](49ea8ac8-f292-b606-2e97-9a7ea04d8400.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (重写 [MelsecMcNetUnpackResponseContent(Byte, Byte)](d9d2764f-86f6-440e-f361-d630b4533e74.htm).) |
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
| 公共方法 | [Write(String, Boolean)](db3d3c30-07f2-93ae-3d04-9d3a3d2c1ef5.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [Write(String, Byte)](c3d2bc96-974d-3a8f-5231-a57f904651bb.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
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
| 公共方法 | [WriteAsync(String, Boolean)](8ef6953f-494f-1b24-0f09-364133778b00.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](e477f621-d3dc-6d5f-8dae-2790108ae950.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

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

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicMcNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/6841fcea-208c-6a3d-1dbd-7975ccf9a958.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 构造函数](../html/6841fcea-208c-6a3d-1dbd-7975ccf9a958.htm "PanasonicMcNet 构造函数 ")

[PanasonicMcNet 构造函数](../html/8a89cbff-9483-171b-a929-b878572d3bb3.htm "PanasonicMcNet 构造函数 ")

[PanasonicMcNet 构造函数 (String, Int32)](../html/98e3bd47-d9a1-9a8d-fbaa-4b87d4ddbe37.htm "PanasonicMcNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNet 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [PanasonicMcNet](8a89cbff-9483-171b-a929-b878572d3bb3.htm) | 实例化松下的的Qna兼容3E帧协议的通讯对象  Instantiate Panasonic's Qna compatible 3E frame protocol communication object |
| 公共方法 | [PanasonicMcNet(String, Int32)](98e3bd47-d9a1-9a8d-fbaa-4b87d4ddbe37.htm) | 指定ip地址及端口号来实例化一个松下的Qna兼容3E帧协议的通讯对象  Specify an IP address and port number to instantiate a Panasonic Qna compatible 3E frame protocol communication object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicMcNet 类](0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicMcNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/8a89cbff-9483-171b-a929-b878572d3bb3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 构造函数](../html/6841fcea-208c-6a3d-1dbd-7975ccf9a958.htm "PanasonicMcNet 构造函数 ")

[PanasonicMcNet 构造函数](../html/8a89cbff-9483-171b-a929-b878572d3bb3.htm "PanasonicMcNet 构造函数 ")

[PanasonicMcNet 构造函数 (String, Int32)](../html/98e3bd47-d9a1-9a8d-fbaa-4b87d4ddbe37.htm "PanasonicMcNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNet 构造函数 |

实例化松下的的Qna兼容3E帧协议的通讯对象  
Instantiate Panasonic's Qna compatible 3E frame protocol communication object

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public PanasonicMcNet()
```

```
Public Sub New
```

```
public:
PanasonicMcNet()
```

```
new : unit -> PanasonicMcNet
```

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicMcNet 类](0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

[PanasonicMcNet 重载](6841fcea-208c-6a3d-1dbd-7975ccf9a958.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicMcNet 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/98e3bd47-d9a1-9a8d-fbaa-4b87d4ddbe37.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 构造函数](../html/6841fcea-208c-6a3d-1dbd-7975ccf9a958.htm "PanasonicMcNet 构造函数 ")

[PanasonicMcNet 构造函数](../html/8a89cbff-9483-171b-a929-b878572d3bb3.htm "PanasonicMcNet 构造函数 ")

[PanasonicMcNet 构造函数 (String, Int32)](../html/98e3bd47-d9a1-9a8d-fbaa-4b87d4ddbe37.htm "PanasonicMcNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNet 构造函数 (String, Int32) |

指定ip地址及端口号来实例化一个松下的Qna兼容3E帧协议的通讯对象  
Specify an IP address and port number to instantiate a Panasonic Qna compatible 3E frame protocol communication object

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public PanasonicMcNet(
	string ipAddress,
	int port
)
```

```
Public Sub New ( 
	ipAddress As String,
	port As Integer
)
```

```
public:
PanasonicMcNet(
	String^ ipAddress, 
	int port
)
```

```
new : 
        ipAddress : string * 
        port : int -> PanasonicMcNet
```

#### 参数

ipAddress
:   类型：SystemString  
    PLC的Ip地址

port
:   类型：SystemInt32  
    PLC的端口

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicMcNet 类](0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

[PanasonicMcNet 重载](6841fcea-208c-6a3d-1dbd-7975ccf9a958.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicMcNet 属性

[原文連結](http://api.hslcommunication.cn/html/e3819fe2-e79d-c9a7-cdf1-2fcad2fb6457.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 构造函数](../html/6841fcea-208c-6a3d-1dbd-7975ccf9a958.htm "PanasonicMcNet 构造函数 ")

[PanasonicMcNet 属性](../html/e3819fe2-e79d-c9a7-cdf1-2fcad2fb6457.htm "PanasonicMcNet 属性")

[PanasonicMcNet 方法](../html/1a9a79ba-af60-cae4-816c-d8ff275d5c34.htm "PanasonicMcNet 方法")

[PanasonicMcNet 字段](../html/6d8890da-c221-26cc-7bea-229ff00ad1e5.htm "PanasonicMcNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNet 属性 |

[PanasonicMcNet](0653413c-b05c-09c8-dd20-ba8019a910c4.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [EnableWriteBitToWordRegister](c933b3aa-a807-7364-0f0a-7ce1296085c5.htm) | 是否开启支持写入位到字寄存器的功能，该功能先读取字寄存器的字数据，然后修改其中的位，再写入回去，可能存在脏数据的风险  Whether to enable supporting the function of writing bit-to-word registers, which first reads the word data of the word register, then modifies the bits in it, and then writes back, there may be a risk of dirty data (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [McType](bf810248-dfa8-a546-5ec3-c92ec37722a4.htm) | 当前的MC协议的格式类型  The format type of the current MC protocol (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性 | [NetworkNumber](7eb9cb50-4c13-3151-627e-5e3cde56fe03.htm) | 网络号，通常为0  Network number, usually 0 (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性 | [NetworkStationNumber](d3f7232f-5184-fb8b-29a2-0a7007213561.htm) | 网络站号，通常为0  Network station number, usually 0 (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性 | [PLCNumber](f087a68e-3dc8-bc74-5a06-4294d9d2069a.htm) | PLC编号，如果是本站信息，则是 0xFF 值，其他站则根据实际的情况指定。  The PLC number, if it is the information of this station, is the 0xFF value, and other stations are specified according to the actual situation. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [TargetIOStation](ad8258fa-31ad-5240-417d-97db45ea311b.htm) | 请求目标模块的IO编号，默认是管理CPU，也就是 0x03FF 的值，如果需要访问其他的非管理CPU的时候，请参考手册进行配置相关的值  The IO number of the target module is the default management CPU, that is, the value of 0x03FF, if you need to access other non-management CPUs, please refer to the manual to configure the relevant values (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicMcNet 类](0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicMcNet 方法

[原文連結](http://api.hslcommunication.cn/html/1a9a79ba-af60-cae4-816c-d8ff275d5c34.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 方法](../html/1a9a79ba-af60-cae4-816c-d8ff275d5c34.htm "PanasonicMcNet 方法")

[McAnalysisAddress 方法](../html/23dea1af-6ed3-1b06-54e2-e65b7b32e173.htm "McAnalysisAddress 方法 ")

[ToString 方法](../html/f2b1932b-0e3c-cf1f-51a5-ab87bc63da5e.htm "ToString 方法 ")

[UnpackResponseContent 方法](../html/49ea8ac8-f292-b606-2e97-9a7ea04d8400.htm "UnpackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNet 方法 |

[PanasonicMcNet](0653413c-b05c-09c8-dd20-ba8019a910c4.htm) 类型公开以下成员。

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
| 公共方法 | [ErrorStateReset](dff5698a-7453-36db-0774-a7bf9249dfeb.htm) | LED 熄灭 出错代码初始化  LED off Error code initialization (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ErrorStateResetAsync](7ef2fa94-e871-2631-bed6-679a23591f92.htm) | LED 熄灭 出错代码初始化  LED off Error code initialization (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ExtractActualData](553f16db-7e99-d2da-265d-4d7f63a598e5.htm) | 从PLC反馈的数据中提取出实际的数据内容，需要传入反馈数据，是否位读取 (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](07258aec-2c66-b9d2-bbf3-b2f3eea3f1dd.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [McAnalysisAddress](23dea1af-6ed3-1b06-54e2-e65b7b32e173.htm) | 当前MC协议的分析地址的方法，对传入的字符串格式的地址进行数据解析。  The current MC protocol's address analysis method performs data parsing on the address of the incoming string format. (重写 [MelsecMcNetMcAnalysisAddress(String, UInt16, Boolean)](2e6fc94e-25bc-0ab8-5590-e334bc02148d.htm).) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](ed03f8d8-4e79-6ee1-a845-ea69a977dd03.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [Read(String, UInt16)](7c89c314-9fbf-a77a-c990-13ef2393c47d.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](c86b56d0-2661-6877-5cf9-ea980d039844.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](b5a3299a-2a1c-644b-aa1d-df14378eae96.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](43c50faf-993e-f3ac-2321-7f846ce458f3.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](a69ef7ef-c8a6-e815-cf29-586acdc173a5.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](4ced354d-18c9-6203-7ebb-ef056b07f6c2.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadExtend](f1f50a27-f4f5-ba75-0713-d64787cf5be0.htm) | **[商业授权]** 读取扩展的数据信息，需要在原有的地址，长度信息之外，输入扩展值信息 **[Authorization]** To read the extended data information, you need to enter the extended value information in addition to the original address and length information (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadExtendAsync](db4dd3cc-2cc6-8629-68b2-8e8233e11370.htm) | **[商业授权]** 读取扩展的数据信息，需要在原有的地址，长度信息之外，输入扩展值信息 **[Authorization]** To read the extended data information, you need to enter the extended value information in addition to the original address and length information (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
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
| 公共方法 | [ReadMemory](f58f02d7-a2ed-3958-3304-77df7fa05dd0.htm) | **[商业授权]** 读取缓冲寄存器的数据信息，地址直接为偏移地址 **[Authorization]** Read the data information of the buffer register, the address is directly the offset address (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadMemoryAsync](6f966786-0a21-46bd-c753-7d7159541458.htm) | **[商业授权]** 读取缓冲寄存器的数据信息，地址直接为偏移地址 **[Authorization]** Read the data information of the buffer register, the address is directly the offset address (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadPlcType](9f8aaf42-bf8c-764b-a3df-3a720da484c5.htm) | 读取PLC的型号信息，例如 Q02HCPU  Read PLC model information, such as Q02HCPU (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadPlcTypeAsync](716c2e02-7df6-18bd-c689-7eb57368ebd9.htm) | 读取PLC的型号信息，例如 Q02HCPU  Read PLC model information, such as Q02HCPU (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadRandom(String)](2b700dc2-727c-7541-cd80-8781ff171e71.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，需要自行解析数据  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the results, you need to parse the data yourself (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadRandom(String, UInt16)](1c644e2a-8674-57fd-cdc8-97d22edea925.htm) | 使用块读取PLC的数据信息，可以跨地址，跨类型组合，每个地址是任意的长度。收到结果后，需要自行解析数据，目前只支持字地址，比如D区，W区，R区，不支持X，Y，M，B，L等等  Read the data information of the PLC randomly. It can be combined across addresses and types. Each address is of any length. After receiving the results, you need to parse the data yourself. Currently, only word addresses are supported, such as D area, W area, R area. X, Y, M, B, L, etc (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadRandomAsync(String)](0ab9c6f2-7488-2f1d-feca-8bcb418739ab.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，需要自行解析数据  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the results, you need to parse the data yourself (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadRandomAsync(String, UInt16)](23278f65-11c2-0c87-5583-2540cdf16a2b.htm) | 使用块读取PLC的数据信息，可以跨地址，跨类型组合，每个地址是任意的长度。收到结果后，需要自行解析数据，目前只支持字地址，比如D区，W区，R区，不支持X，Y，M，B，L等等  Read the data information of the PLC randomly. It can be combined across addresses and types. Each address is of any length. After receiving the results, you need to parse the data yourself. Currently, only word addresses are supported, such as D area, W area, R area. X, Y, M, B, L, etc (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadRandomInt16](4338885c-8abf-8bfd-d553-9e57390f185c.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，自动转换为了short类型的数组  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the result, it is automatically converted to an array of type short. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadRandomInt16Async](0d680060-284c-3a6f-b245-c9855ba7a1b0.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，自动转换为了short类型的数组  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the result, it is automatically converted to an array of type short. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadRandomUInt16](c4d126ae-e24d-34b5-abf3-9e6e08492f45.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，自动转换为了ushort类型的数组  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the result, it is automatically converted to an array of type ushort. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadRandomUInt16Async](7faa94e9-a1f4-b321-8c1f-bc64c715a438.htm) | 随机读取PLC的数据信息，可以跨地址，跨类型组合，但是每个地址只能读取一个word，也就是2个字节的内容。收到结果后，自动转换为了ushort类型的数组  Randomly read PLC data information, which can be combined across addresses and types, but each address can only read one word, which is the content of 2 bytes. After receiving the result, it is automatically converted to an array of type ushort. (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadSmartModule](00851b3c-e4b5-b306-d624-f14f9501d3df.htm) | **[商业授权]** 读取智能模块的数据信息，需要指定模块地址，偏移地址，读取的字节长度 **[Authorization]** To read the extended data information, you need to enter the extended value information in addition to the original address and length information (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadSmartModuleAsync](eecc3de1-8d08-3067-45b2-2f93e43b966f.htm) | **[商业授权]** 读取智能模块的数据信息，需要指定模块地址，偏移地址，读取的字节长度 **[Authorization]** To read the extended data information, you need to enter the extended value information in addition to the original address and length information (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadTags(String, UInt16)](d7e5e3e8-e3bf-2a02-7d90-561cf28138ea.htm) | **[商业授权]** 读取PLC的标签信息，需要传入标签的名称，读取的字长度，标签举例：A; label[1]; bbb[10,10,10] **[Authorization]** To read the label information of the PLC, you need to pass in the name of the label, the length of the word read, and an example of the label: A; label [1]; bbb [10,10,10] (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadTags(String, UInt16)](69e4bc2a-e235-1294-4653-b6a2dec1535c.htm) | **[商业授权]** 读取PLC的标签信息，需要传入标签的名称，读取的字长度，标签举例：A; label[1]; bbb[10,10,10] **[Authorization]** To read the label information of the PLC, you need to pass in the name of the label, the length of the word read, and an example of the label: A; label [1]; bbb [10,10,10] (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadTagsAsync(String, UInt16)](b2edcbab-1f74-1963-fc00-18aa8c8c02d7.htm) | **[商业授权]** 读取PLC的标签信息，需要传入标签的名称，读取的字长度，标签举例：A; label[1]; bbb[10,10,10] **[Authorization]** To read the label information of the PLC, you need to pass in the name of the label, the length of the word read, and an example of the label: A; label [1]; bbb [10,10,10] (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [ReadTagsAsync(String, UInt16)](01ce89ac-b4c1-c2ed-25dc-662efae546e6.htm) | **[商业授权]** 读取PLC的标签信息，需要传入标签的名称，读取的字长度，标签举例：A; label[1]; bbb[10,10,10] **[Authorization]** To read the label information of the PLC, you need to pass in the name of the label, the length of the word read, and an example of the label: A; label [1]; bbb [10,10,10] (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
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
| 公共方法 | [RemoteReset](08535f06-ae9f-fca5-1df3-3299c505f76f.htm) | 远程Reset操作  Remote Reset Operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteResetAsync](cca35551-499a-b900-810f-e5bf5c50e3ec.htm) | 远程Reset操作  Remote Reset Operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteRun](2f2f4e7f-17e2-63fd-882d-59e20efa8ef7.htm) | 远程Run操作  Remote Run Operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteRunAsync](da1c67d7-6063-6ed7-d505-29e856404d54.htm) | 远程Run操作  Remote Run Operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteStop](e199a595-8ae0-3ee3-383d-cc7b5f81430b.htm) | 远程Stop操作  Remote Stop operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [RemoteStopAsync](06a7c9b0-0c6a-ce1f-77ea-047a4ea5e585.htm) | 远程Stop操作  Remote Stop operation (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](f2b1932b-0e3c-cf1f-51a5-ab87bc63da5e.htm) | (重写 [MelsecMcNetToString](77b9b8f7-46b7-1f45-fa1e-3e7b505b61e4.htm).) |
| 公共方法 | [UnpackResponseContent](49ea8ac8-f292-b606-2e97-9a7ea04d8400.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (重写 [MelsecMcNetUnpackResponseContent(Byte, Byte)](d9d2764f-86f6-440e-f361-d630b4533e74.htm).) |
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
| 公共方法 | [Write(String, Boolean)](db3d3c30-07f2-93ae-3d04-9d3a3d2c1ef5.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法 | [Write(String, Byte)](c3d2bc96-974d-3a8f-5231-a57f904651bb.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
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
| 公共方法 | [WriteAsync(String, Boolean)](8ef6953f-494f-1b24-0f09-364133778b00.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](e477f621-d3dc-6d5f-8dae-2790108ae950.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)。) |
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

[PanasonicMcNet 类](0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## McAnalysisAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/23dea1af-6ed3-1b06-54e2-e65b7b32e173.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 方法](../html/1a9a79ba-af60-cae4-816c-d8ff275d5c34.htm "PanasonicMcNet 方法")

[McAnalysisAddress 方法](../html/23dea1af-6ed3-1b06-54e2-e65b7b32e173.htm "McAnalysisAddress 方法 ")

[ToString 方法](../html/f2b1932b-0e3c-cf1f-51a5-ab87bc63da5e.htm "ToString 方法 ")

[UnpackResponseContent 方法](../html/49ea8ac8-f292-b606-2e97-9a7ea04d8400.htm "UnpackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNetMcAnalysisAddress 方法 |

当前MC协议的分析地址的方法，对传入的字符串格式的地址进行数据解析。  
The current MC protocol's address analysis method performs data parsing on the address of the incoming string format.

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<McAddressData> McAnalysisAddress(
	string address,
	ushort length,
	bool isBit
)
```

```
Public Overrides Function McAnalysisAddress ( 
	address As String,
	length As UShort,
	isBit As Boolean
) As OperateResult(Of McAddressData)
```

```
public:
virtual OperateResult<McAddressData^>^ McAnalysisAddress(
	String^ address, 
	unsigned short length, 
	bool isBit
) override
```

```
abstract McAnalysisAddress : 
        address : string * 
        length : uint16 * 
        isBit : bool -> OperateResult<McAddressData> 
override McAnalysisAddress : 
        address : string * 
        length : uint16 * 
        isBit : bool -> OperateResult<McAddressData>
```

#### 参数

address
:   类型：SystemString  
    地址信息

length
:   类型：SystemUInt16  
    数据长度

isBit
:   类型：SystemBoolean  
    当前是否读写bool操作

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)[McAddressData](9f68dd79-b0c1-e8b5-4524-821f83f59975.htm)  
解析后的数据信息

#### 实现

[IReadWriteMcMcAnalysisAddress(String, UInt16, Boolean)](83ac01c3-c0ca-d422-8656-482d6ef6821e.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicMcNet 类](0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/f2b1932b-0e3c-cf1f-51a5-ab87bc63da5e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 方法](../html/1a9a79ba-af60-cae4-816c-d8ff275d5c34.htm "PanasonicMcNet 方法")

[McAnalysisAddress 方法](../html/23dea1af-6ed3-1b06-54e2-e65b7b32e173.htm "McAnalysisAddress 方法 ")

[ToString 方法](../html/f2b1932b-0e3c-cf1f-51a5-ab87bc63da5e.htm "ToString 方法 ")

[UnpackResponseContent 方法](../html/49ea8ac8-f292-b606-2e97-9a7ea04d8400.htm "UnpackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNetToString 方法 |

[缺少 "M:HslCommunication.Profinet.Panasonic.PanasonicMcNet.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
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

[缺少 "M:HslCommunication.Profinet.Panasonic.PanasonicMcNet.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicMcNet 类](0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UnpackResponseContent 方法 

[原文連結](http://api.hslcommunication.cn/html/49ea8ac8-f292-b606-2e97-9a7ea04d8400.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 方法](../html/1a9a79ba-af60-cae4-816c-d8ff275d5c34.htm "PanasonicMcNet 方法")

[McAnalysisAddress 方法](../html/23dea1af-6ed3-1b06-54e2-e65b7b32e173.htm "McAnalysisAddress 方法 ")

[ToString 方法](../html/f2b1932b-0e3c-cf1f-51a5-ab87bc63da5e.htm "ToString 方法 ")

[UnpackResponseContent 方法](../html/49ea8ac8-f292-b606-2e97-9a7ea04d8400.htm "UnpackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNetUnpackResponseContent 方法 |

根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  
According to the message command returned by the other party, the command is basically unpacked, for example,
various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<byte[]> UnpackResponseContent(
	byte[] send,
	byte[] response
)
```

```
Public Overrides Function UnpackResponseContent ( 
	send As Byte(),
	response As Byte()
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ UnpackResponseContent(
	array<unsigned char>^ send, 
	array<unsigned char>^ response
) override
```

```
abstract UnpackResponseContent : 
        send : byte[] * 
        response : byte[] -> OperateResult<byte[]> 
override UnpackResponseContent : 
        send : byte[] * 
        response : byte[] -> OperateResult<byte[]>
```

#### 参数

send
:   类型：SystemByte  
    发送的原始报文数据

response
:   类型：SystemByte  
    设备方反馈的原始报文内容

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
返回拆包之后的报文信息，默认不进行任何的拆包操作

![](../icons/SectionExpanded.png)备注

在实际解包的操作过程中，通常对状态码，错误码等消息进行判断，如果校验不通过，将携带错误消息返回  
During the actual unpacking operation, the status code, error code and other messages are usually judged. If the verification fails, the error message will be returned.

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicMcNet 类](0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicMcNet 字段

[原文連結](http://api.hslcommunication.cn/html/6d8890da-c221-26cc-7bea-229ff00ad1e5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMcNet 类](../html/0653413c-b05c-09c8-dd20-ba8019a910c4.htm "PanasonicMcNet 类")

[PanasonicMcNet 构造函数](../html/6841fcea-208c-6a3d-1dbd-7975ccf9a958.htm "PanasonicMcNet 构造函数 ")

[PanasonicMcNet 属性](../html/e3819fe2-e79d-c9a7-cdf1-2fcad2fb6457.htm "PanasonicMcNet 属性")

[PanasonicMcNet 方法](../html/1a9a79ba-af60-cae4-816c-d8ff275d5c34.htm "PanasonicMcNet 方法")

[PanasonicMcNet 字段](../html/6d8890da-c221-26cc-7bea-229ff00ad1e5.htm "PanasonicMcNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMcNet 字段 |

[PanasonicMcNet](0653413c-b05c-09c8-dd20-ba8019a910c4.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicMcNet 类](0653413c-b05c-09c8-dd20-ba8019a910c4.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicMewtocol 类

[原文連結](http://api.hslcommunication.cn/html/032b2c43-56fd-ac73-cf50-2cbc88888d43.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMewtocol 类](../html/032b2c43-56fd-ac73-cf50-2cbc88888d43.htm "PanasonicMewtocol 类")

[PanasonicMewtocol 构造函数](../html/9b9eb9b9-abd0-f95d-3021-40514b5115d2.htm "PanasonicMewtocol 构造函数 ")

[PanasonicMewtocol 属性](../html/ba25ae80-075a-062c-2f5f-15c19cc3069c.htm "PanasonicMewtocol 属性")

[PanasonicMewtocol 方法](../html/2b675ff1-04d9-3154-cebf-d250fa7908ac.htm "PanasonicMewtocol 方法")

[PanasonicMewtocol 字段](../html/e163ec98-bd01-8aaa-c2a7-edc983c26345.htm "PanasonicMewtocol 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMewtocol 类 |

松下PLC的数据交互协议，采用Mewtocol协议通讯，支持的地址列表参考api文档  
The data exchange protocol of Panasonic PLC adopts Mewtocol protocol for communication. For the list of supported addresses, refer to the api document.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)  
        HslCommunication.Profinet.PanasonicPanasonicMewtocol

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class PanasonicMewtocol : DeviceSerialPort
```

```
Public Class PanasonicMewtocol
	Inherits DeviceSerialPort
```

```
public ref class PanasonicMewtocol : public DeviceSerialPort
```

```
type PanasonicMewtocol =  
    class
        inherit DeviceSerialPort
    end
```

PanasonicMewtocol 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [PanasonicMewtocol](9b9eb9b9-abd0-f95d-3021-40514b5115d2.htm) | 实例化一个默认的松下PLC通信对象，默认站号为0xEE  Instantiate a default Panasonic PLC communication object, the default station number is 0xEE |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](5554f563-5e29-699c-435b-c8e7ad5fa0e4.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](a4fcb477-64d2-8da7-c712-687f0261755b.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [IsClearCacheBeforeRead](81dc971f-1f86-8535-ae9a-586b6590a584.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](f2b8e4e3-7079-5240-4ff1-14c6c9b376e6.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](f800fd56-7d49-87c4-985b-103a38330078.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](dd4057ef-88fe-f7d0-d2c4-f6e15488fe78.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Station](d09ddd17-c8b3-1b30-353a-5f850a670861.htm) | PLC设备的目标站号，需要根据实际的设置来填写  The target station number of the PLC device needs to be filled in according to the actual settings |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

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
| 受保护的方法 | [GetNewNetMessage](7965ebf7-62c2-7e07-454b-34a847576e15.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
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
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](1bb89008-a1d3-8c0d-2467-57e5512e790e.htm) | 读取指定地址的原始数据，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Read the original data of the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](1990c6e5-00aa-df1c-cd4e-e7c7af8b9775.htm) | 读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0 (重写 [DeviceCommunicationReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm).) |
| 公共方法 | [ReadBool(String)](c592fb67-9f60-440d-2113-cbbd60656fdc.htm) | 批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法 | [ReadBool(String, UInt16)](ad733ba5-d529-b656-0836-d10bdde351fb.htm) | 批量读取松下PLC的位数据，按照字为单位，地址为 X0,X10,Y10，读取的长度为16的倍数  Read the bit data of Panasonic PLC in batches, the unit is word, the address is X0, X10, Y10, and the read length is a multiple of 16 (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](6073f7df-6433-f1f5-e4b0-8239f521efc6.htm) | 读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0 (重写 [DeviceCommunicationReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm).) |
| 公共方法 | [ReadBoolAsync(String)](8519e0c2-cee3-e74f-fbc7-364a9502903b.htm) | 批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
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
| 公共方法 | [ReadPlcType](49bcbe2b-8bd2-ee99-8102-67bb87ee6649.htm) | 读取PLC的型号信息  Read the model information of the PLC |
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
| 公共方法 | [SerialPortInni(ActionSerialPort)](3728a0f2-4553-0bc8-5c92-253c13e84f05.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String)](f82e0ba4-af2f-5165-dcde-e047c3351e74.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](433523a8-1606-f029-50db-68ef68034046.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](cefa7cd2-9c1f-b79d-9ea3-d4d27d7ddeae.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](46cfdab1-80b6-9db2-429d-cb1445f90f46.htm) | (重写 [DeviceSerialPortToString](774e7804-97ce-fb07-3764-39fb1ca2e1e7.htm).) |
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
| 公共方法 | [Write(String, Boolean)](78e98946-83ee-c9d3-cd75-48c23c04c48d.htm) | 往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0 (重写 [DeviceCommunicationWrite(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm).) |
| 公共方法 | [Write(String, Boolean)](9a5d5a2b-a59f-060b-8dec-cbefef733f36.htm) | 往指定的地址写入 Boolean 数组，地址举例 X0.0 Y0.0 R0.0 L0.0， 起始的位地址必须为16的倍数，写入的 Boolean 数组长度也为16的倍数。  Write the Boolean array to the specified address, address example: SR0.0 X0.0 Y0.0 R0.0 L0.0, the starting bit address must be a multiple of 16. Boolean The length of the array is also a multiple of 16. (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](f79999d2-31b6-4234-205a-3c89f7fcc30e.htm) | 将数据写入到指定的地址里去，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Write data to the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法 | [Write(String, Boolean)](34c96524-33df-db6f-b502-d3661f071e71.htm) | 将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Write the Bool array value to the specified discrete address, one address corresponds to one bool value, the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [WriteAsync(String, Boolean)](89ed09c3-bef3-0b44-0c2f-f0e7a43932c2.htm) | 往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0 (重写 [DeviceCommunicationWriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm).) |
| 公共方法 | [WriteAsync(String, Boolean)](0aa566c2-bb30-009b-8306-dbbb88d4dbc4.htm) | 将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Write the Bool array value to the specified discrete address, one address corresponds to one bool value, the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

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

地址支持携带站号的访问方式，例如：s=2;D100

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PanasonicMewtocol 构造函数 

[原文連結](http://api.hslcommunication.cn/html/9b9eb9b9-abd0-f95d-3021-40514b5115d2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic](../html/65511854-a56d-47ca-f165-fe601044e6b1.htm "HslCommunication.Profinet.Panasonic")

[PanasonicMewtocol 类](../html/032b2c43-56fd-ac73-cf50-2cbc88888d43.htm "PanasonicMewtocol 类")

[PanasonicMewtocol 构造函数](../html/9b9eb9b9-abd0-f95d-3021-40514b5115d2.htm "PanasonicMewtocol 构造函数 ")

[PanasonicMewtocol 属性](../html/ba25ae80-075a-062c-2f5f-15c19cc3069c.htm "PanasonicMewtocol 属性")

[PanasonicMewtocol 方法](../html/2b675ff1-04d9-3154-cebf-d250fa7908ac.htm "PanasonicMewtocol 方法")

[PanasonicMewtocol 字段](../html/e163ec98-bd01-8aaa-c2a7-edc983c26345.htm "PanasonicMewtocol 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PanasonicMewtocol 构造函数 |

实例化一个默认的松下PLC通信对象，默认站号为0xEE  
Instantiate a default Panasonic PLC communication object, the default station number is 0xEE

**命名空间：**
 [HslCommunication.Profinet.Panasonic](65511854-a56d-47ca-f165-fe601044e6b1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public PanasonicMewtocol(
	byte station = 238
)
```

```
Public Sub New ( 
	Optional station As Byte = 238
)
```

```
public:
PanasonicMewtocol(
	unsigned char station = 238
)
```

```
new : 
        ?station : byte 
(* Defaults:
        let _station = defaultArg station 238
*)
-> PanasonicMewtocol
```

#### 参数

station (Optional)
:   类型：SystemByte  
    站号信息，默认为0xEE

![](../icons/SectionExpanded.png)参见

#### 引用

[PanasonicMewtocol 类](032b2c43-56fd-ac73-cf50-2cbc88888d43.htm)

[HslCommunication.Profinet.Panasonic 命名空间](65511854-a56d-47ca-f165-fe601044e6b1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)