# HslCommunication - HslCommunication.Profinet.Panasonic.Helper

> 分類頁數: 25



---
## HslCommunication.Profinet.Panasonic.Helper

[原文連結](http://api.hslcommunication.cn/html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Panasonic.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.Panasonic.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [MewtocolHelper](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm) | Mewtocol协议的辅助类信息 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MewtocolHelper 类

[原文連結](http://api.hslcommunication.cn/html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 构造函数](../html/daed9c09-67e1-2b15-5950-2805930cde0c.htm "MewtocolHelper 构造函数 ")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelper 类 |

Mewtocol协议的辅助类信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.Panasonic.HelperMewtocolHelper

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class MewtocolHelper
```

```
Public Class MewtocolHelper
```

```
public ref class MewtocolHelper
```

```
type MewtocolHelper =  class end
```

MewtocolHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [MewtocolHelper](daed9c09-67e1-2b15-5950-2805930cde0c.htm) | 初始化 MewtocolHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](e78daeb8-72d8-2403-794b-0fbede452393.htm) | 读取指定地址的原始数据，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Read the original data of the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [ReadAsync](3233fb37-7c23-f59a-b2dc-fc1fd738fbdf.htm) | 读取指定地址的原始数据，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Read the original data of the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Byte, String)](64f30e31-990b-905c-632b-e2d727b835ee.htm) | 读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Byte, String)](101cb77b-b52c-b91f-7124-67171c38166d.htm) | 批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Byte, String, UInt16)](c7b9b674-8ff2-bc4a-a893-e49a9668dafd.htm) | 批量读取松下PLC的位数据，按照字为单位，地址为 X0,X10,Y10，读取的长度为16的倍数  Read the bit data of Panasonic PLC in batches, the unit is word, the address is X0, X10, Y10, and the read length is a multiple of 16 |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Byte, String)](ae0a83b8-a1ec-9dc1-04f0-37a5cffa74e4.htm) | 读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Byte, String)](14b04612-0172-89fd-4e65-12569b11cd11.htm) | 批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Byte, String, UInt16)](118dfb75-779e-f41c-e457-e570db1735f9.htm) | 批量读取松下PLC的位数据，按照字为单位，地址为 X0,X10,Y10，读取的长度为16的倍数  Read the bit data of Panasonic PLC in batches, the unit is word, the address is X0, X10, Y10, and the read length is a multiple of 16 |
| 公共方法静态成员 | [ReadPlcType](de0359bb-2800-ccf9-ba8b-4d617d2afdd8.htm) | 读取PLC的型号信息  Read the model information of the PLC |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](3d2e6b51-292f-d104-0d59-24603db047b9.htm) | 往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](395acabb-39d1-5ff7-2f67-f488e62d3f70.htm) | 往指定的地址写入 Boolean 数组，地址举例 X0.0 Y0.0 R0.0 L0.0， 起始的位地址必须为16的倍数，写入的 Boolean 数组长度也为16的倍数。  Write the Boolean array to the specified address, address example: SR0.0 X0.0 Y0.0 R0.0 L0.0, the starting bit address must be a multiple of 16. Boolean The length of the array is also a multiple of 16. |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Byte)](2679d12f-4f12-b54c-6cd5-c82ecee975b1.htm) | 将数据写入到指定的地址里去，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Write data to the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](d4d2e6b4-3c6e-6318-d094-0b86689189de.htm) | 将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Write the Bool array value to the specified discrete address, one address corresponds to one bool value, the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](661ff8d1-c652-d7bf-782e-d79ec60fdc83.htm) | 往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](cc9a08bc-e17f-5db2-d99f-fadc2c115099.htm) | 往指定的地址写入 Boolean 数组，地址举例 X0.0 Y0.0 R0.0 L0.0， 起始的位地址必须为16的倍数，写入的 Boolean 数组长度也为16的倍数。  Write the Boolean array to the specified address, address example: SR0.0 X0.0 Y0.0 R0.0 L0.0, the starting bit address must be a multiple of 16. Boolean The length of the array is also a multiple of 16. |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Byte)](083baf45-f05f-fdbf-7cb8-17e0e6387352.htm) | 将数据写入到指定的地址里去，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Write data to the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](47b4bee5-97ac-23c5-ae26-ccace6519755.htm) | 将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Write the Bool array value to the specified discrete address, one address corresponds to one bool value, the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MewtocolHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/daed9c09-67e1-2b15-5950-2805930cde0c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 构造函数](../html/daed9c09-67e1-2b15-5950-2805930cde0c.htm "MewtocolHelper 构造函数 ")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelper 构造函数 |

初始化 [MewtocolHelper](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public MewtocolHelper()
```

```
Public Sub New
```

```
public:
MewtocolHelper()
```

```
new : unit -> MewtocolHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MewtocolHelper 方法

[原文連結](http://api.hslcommunication.cn/html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[Read 方法](../html/e78daeb8-72d8-2403-794b-0fbede452393.htm "Read 方法 ")

[ReadAsync 方法](../html/3233fb37-7c23-f59a-b2dc-fc1fd738fbdf.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4903b92e-f917-30f9-5f54-9ef5816b063e.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/48143ff1-208b-4c74-ac6f-d7fd22890db8.htm "ReadBoolAsync 方法 ")

[ReadPlcType 方法](../html/de0359bb-2800-ccf9-ba8b-4d617d2afdd8.htm "ReadPlcType 方法 ")

[Write 方法](../html/a96e986e-6858-131c-01a9-2edf0643db5a.htm "Write 方法 ")

[WriteAsync 方法](../html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelper 方法 |

[MewtocolHelper](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](e78daeb8-72d8-2403-794b-0fbede452393.htm) | 读取指定地址的原始数据，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Read the original data of the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [ReadAsync](3233fb37-7c23-f59a-b2dc-fc1fd738fbdf.htm) | 读取指定地址的原始数据，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Read the original data of the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Byte, String)](64f30e31-990b-905c-632b-e2d727b835ee.htm) | 读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Byte, String)](101cb77b-b52c-b91f-7124-67171c38166d.htm) | 批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Byte, String, UInt16)](c7b9b674-8ff2-bc4a-a893-e49a9668dafd.htm) | 批量读取松下PLC的位数据，按照字为单位，地址为 X0,X10,Y10，读取的长度为16的倍数  Read the bit data of Panasonic PLC in batches, the unit is word, the address is X0, X10, Y10, and the read length is a multiple of 16 |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Byte, String)](ae0a83b8-a1ec-9dc1-04f0-37a5cffa74e4.htm) | 读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Byte, String)](14b04612-0172-89fd-4e65-12569b11cd11.htm) | 批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Byte, String, UInt16)](118dfb75-779e-f41c-e457-e570db1735f9.htm) | 批量读取松下PLC的位数据，按照字为单位，地址为 X0,X10,Y10，读取的长度为16的倍数  Read the bit data of Panasonic PLC in batches, the unit is word, the address is X0, X10, Y10, and the read length is a multiple of 16 |
| 公共方法静态成员 | [ReadPlcType](de0359bb-2800-ccf9-ba8b-4d617d2afdd8.htm) | 读取PLC的型号信息  Read the model information of the PLC |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](3d2e6b51-292f-d104-0d59-24603db047b9.htm) | 往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](395acabb-39d1-5ff7-2f67-f488e62d3f70.htm) | 往指定的地址写入 Boolean 数组，地址举例 X0.0 Y0.0 R0.0 L0.0， 起始的位地址必须为16的倍数，写入的 Boolean 数组长度也为16的倍数。  Write the Boolean array to the specified address, address example: SR0.0 X0.0 Y0.0 R0.0 L0.0, the starting bit address must be a multiple of 16. Boolean The length of the array is also a multiple of 16. |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Byte)](2679d12f-4f12-b54c-6cd5-c82ecee975b1.htm) | 将数据写入到指定的地址里去，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Write data to the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](d4d2e6b4-3c6e-6318-d094-0b86689189de.htm) | 将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Write the Bool array value to the specified discrete address, one address corresponds to one bool value, the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](661ff8d1-c652-d7bf-782e-d79ec60fdc83.htm) | 往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](cc9a08bc-e17f-5db2-d99f-fadc2c115099.htm) | 往指定的地址写入 Boolean 数组，地址举例 X0.0 Y0.0 R0.0 L0.0， 起始的位地址必须为16的倍数，写入的 Boolean 数组长度也为16的倍数。  Write the Boolean array to the specified address, address example: SR0.0 X0.0 Y0.0 R0.0 L0.0, the starting bit address must be a multiple of 16. Boolean The length of the array is also a multiple of 16. |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Byte)](083baf45-f05f-fdbf-7cb8-17e0e6387352.htm) | 将数据写入到指定的地址里去，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Write data to the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](47b4bee5-97ac-23c5-ae26-ccace6519755.htm) | 将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Write the Bool array value to the specified discrete address, one address corresponds to one bool value, the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/e78daeb8-72d8-2403-794b-0fbede452393.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[Read 方法](../html/e78daeb8-72d8-2403-794b-0fbede452393.htm "Read 方法 ")

[ReadAsync 方法](../html/3233fb37-7c23-f59a-b2dc-fc1fd738fbdf.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4903b92e-f917-30f9-5f54-9ef5816b063e.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/48143ff1-208b-4c74-ac6f-d7fd22890db8.htm "ReadBoolAsync 方法 ")

[ReadPlcType 方法](../html/de0359bb-2800-ccf9-ba8b-4d617d2afdd8.htm "ReadPlcType 方法 ")

[Write 方法](../html/a96e986e-6858-131c-01a9-2edf0643db5a.htm "Write 方法 ")

[WriteAsync 方法](../html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperRead 方法 |

读取指定地址的原始数据，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  
Read the original data of the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function Read ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member Read : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
原始的字节数据的信息

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/3233fb37-7c23-f59a-b2dc-fc1fd738fbdf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[Read 方法](../html/e78daeb8-72d8-2403-794b-0fbede452393.htm "Read 方法 ")

[ReadAsync 方法](../html/3233fb37-7c23-f59a-b2dc-fc1fd738fbdf.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4903b92e-f917-30f9-5f54-9ef5816b063e.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/48143ff1-208b-4c74-ac6f-d7fd22890db8.htm "ReadBoolAsync 方法 ")

[ReadPlcType 方法](../html/de0359bb-2800-ccf9-ba8b-4d617d2afdd8.htm "ReadPlcType 方法 ")

[Write 方法](../html/a96e986e-6858-131c-01a9-2edf0643db5a.htm "Write 方法 ")

[WriteAsync 方法](../html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadAsync 方法 |

读取指定地址的原始数据，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  
Read the original data of the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function ReadAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
原始的字节数据的信息

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/4903b92e-f917-30f9-5f54-9ef5816b063e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[ReadBool 方法](../html/4903b92e-f917-30f9-5f54-9ef5816b063e.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, Byte, String)](../html/64f30e31-990b-905c-632b-e2d727b835ee.htm "ReadBool 方法 (IReadWriteDevice, Byte, String)")

[ReadBool 方法 (IReadWriteDevice, Byte, String[])](../html/101cb77b-b52c-b91f-7124-67171c38166d.htm "ReadBool 方法 (IReadWriteDevice, Byte, String[])")

[ReadBool 方法 (IReadWriteDevice, Byte, String, UInt16)](../html/c7b9b674-8ff2-bc4a-a893-e49a9668dafd.htm "ReadBool 方法 (IReadWriteDevice, Byte, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadBool 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Byte, String)](64f30e31-990b-905c-632b-e2d727b835ee.htm) | 读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Byte, String)](101cb77b-b52c-b91f-7124-67171c38166d.htm) | 批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Byte, String, UInt16)](c7b9b674-8ff2-bc4a-a893-e49a9668dafd.htm) | 批量读取松下PLC的位数据，按照字为单位，地址为 X0,X10,Y10，读取的长度为16的倍数  Read the bit data of Panasonic PLC in batches, the unit is word, the address is X0, X10, Y10, and the read length is a multiple of 16 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (IReadWriteDevice, Byte, String)

[原文連結](http://api.hslcommunication.cn/html/64f30e31-990b-905c-632b-e2d727b835ee.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[ReadBool 方法](../html/4903b92e-f917-30f9-5f54-9ef5816b063e.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, Byte, String)](../html/64f30e31-990b-905c-632b-e2d727b835ee.htm "ReadBool 方法 (IReadWriteDevice, Byte, String)")

[ReadBool 方法 (IReadWriteDevice, Byte, String[])](../html/101cb77b-b52c-b91f-7124-67171c38166d.htm "ReadBool 方法 (IReadWriteDevice, Byte, String[])")

[ReadBool 方法 (IReadWriteDevice, Byte, String, UInt16)](../html/c7b9b674-8ff2-bc4a-a893-e49a9668dafd.htm "ReadBool 方法 (IReadWriteDevice, Byte, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadBool 方法 (IReadWriteDevice, Byte, String) |

读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  
Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<bool> ReadBool(
	IReadWriteDevice plc,
	byte station,
	string address
)
```

```
Public Shared Function ReadBool ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String
) As OperateResult(Of Boolean)
```

```
public:
static OperateResult<bool>^ ReadBool(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address
)
```

```
static member ReadBool : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string -> OperateResult<bool> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[ReadBool 重载](4903b92e-f917-30f9-5f54-9ef5816b063e.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (IReadWriteDevice, Byte, String[])

[原文連結](http://api.hslcommunication.cn/html/101cb77b-b52c-b91f-7124-67171c38166d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[ReadBool 方法](../html/4903b92e-f917-30f9-5f54-9ef5816b063e.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, Byte, String)](../html/64f30e31-990b-905c-632b-e2d727b835ee.htm "ReadBool 方法 (IReadWriteDevice, Byte, String)")

[ReadBool 方法 (IReadWriteDevice, Byte, String[])](../html/101cb77b-b52c-b91f-7124-67171c38166d.htm "ReadBool 方法 (IReadWriteDevice, Byte, String[])")

[ReadBool 方法 (IReadWriteDevice, Byte, String, UInt16)](../html/c7b9b674-8ff2-bc4a-a893-e49a9668dafd.htm "ReadBool 方法 (IReadWriteDevice, Byte, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadBool 方法 (IReadWriteDevice, Byte, String) |

批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  
Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string[] address
)
```

```
Public Shared Function ReadBool ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String()
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ReadBool(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	array<String^>^ address
)
```

```
static member ReadBool : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string[] -> OperateResult<bool[]> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    等待读取的地址列表，数组长度不限制

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[ReadBool 重载](4903b92e-f917-30f9-5f54-9ef5816b063e.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (IReadWriteDevice, Byte, String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/c7b9b674-8ff2-bc4a-a893-e49a9668dafd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[ReadBool 方法](../html/4903b92e-f917-30f9-5f54-9ef5816b063e.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, Byte, String)](../html/64f30e31-990b-905c-632b-e2d727b835ee.htm "ReadBool 方法 (IReadWriteDevice, Byte, String)")

[ReadBool 方法 (IReadWriteDevice, Byte, String[])](../html/101cb77b-b52c-b91f-7124-67171c38166d.htm "ReadBool 方法 (IReadWriteDevice, Byte, String[])")

[ReadBool 方法 (IReadWriteDevice, Byte, String, UInt16)](../html/c7b9b674-8ff2-bc4a-a893-e49a9668dafd.htm "ReadBool 方法 (IReadWriteDevice, Byte, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadBool 方法 (IReadWriteDevice, Byte, String, UInt16) |

批量读取松下PLC的位数据，按照字为单位，地址为 X0,X10,Y10，读取的长度为16的倍数  
Read the bit data of Panasonic PLC in batches, the unit is word, the address is X0, X10, Y10, and the read length is a multiple of 16

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBool ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ReadBool(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBool : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<bool[]> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[ReadBool 重载](4903b92e-f917-30f9-5f54-9ef5816b063e.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/48143ff1-208b-4c74-ac6f-d7fd22890db8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[ReadBoolAsync 方法](../html/48143ff1-208b-4c74-ac6f-d7fd22890db8.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String)](../html/ae0a83b8-a1ec-9dc1-04f0-37a5cffa74e4.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String)")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String[])](../html/14b04612-0172-89fd-4e65-12569b11cd11.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String[])")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16)](../html/118dfb75-779e-f41c-e457-e570db1735f9.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadBoolAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Byte, String)](ae0a83b8-a1ec-9dc1-04f0-37a5cffa74e4.htm) | 读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Byte, String)](14b04612-0172-89fd-4e65-12569b11cd11.htm) | 批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Byte, String, UInt16)](118dfb75-779e-f41c-e457-e570db1735f9.htm) | 批量读取松下PLC的位数据，按照字为单位，地址为 X0,X10,Y10，读取的长度为16的倍数  Read the bit data of Panasonic PLC in batches, the unit is word, the address is X0, X10, Y10, and the read length is a multiple of 16 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 (IReadWriteDevice, Byte, String)

[原文連結](http://api.hslcommunication.cn/html/ae0a83b8-a1ec-9dc1-04f0-37a5cffa74e4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[ReadBoolAsync 方法](../html/48143ff1-208b-4c74-ac6f-d7fd22890db8.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String)](../html/ae0a83b8-a1ec-9dc1-04f0-37a5cffa74e4.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String)")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String[])](../html/14b04612-0172-89fd-4e65-12569b11cd11.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String[])")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16)](../html/118dfb75-779e-f41c-e457-e570db1735f9.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadBoolAsync 方法 (IReadWriteDevice, Byte, String) |

读取单个的地址信息的bool值，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  
Read the bool value of a single address, for example: SR0.0 X0.0 Y0.0 R0.0 L0.0

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<bool>> ReadBoolAsync(
	IReadWriteDevice plc,
	byte station,
	string address
)
```

```
Public Shared Function ReadBoolAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String
) As Task(Of OperateResult(Of Boolean))
```

```
public:
static Task<OperateResult<bool>^>^ ReadBoolAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address
)
```

```
static member ReadBoolAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string -> Task<OperateResult<bool>> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[ReadBoolAsync 重载](48143ff1-208b-4c74-ac6f-d7fd22890db8.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 (IReadWriteDevice, Byte, String[])

[原文連結](http://api.hslcommunication.cn/html/14b04612-0172-89fd-4e65-12569b11cd11.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[ReadBoolAsync 方法](../html/48143ff1-208b-4c74-ac6f-d7fd22890db8.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String)](../html/ae0a83b8-a1ec-9dc1-04f0-37a5cffa74e4.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String)")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String[])](../html/14b04612-0172-89fd-4e65-12569b11cd11.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String[])")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16)](../html/118dfb75-779e-f41c-e457-e570db1735f9.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadBoolAsync 方法 (IReadWriteDevice, Byte, String) |

批量读取松下PLC的位数据，传入一个读取的地址列表，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  
Batch read the bit data of Panasonic PLC, pass in a read address list, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string[] address
)
```

```
Public Shared Function ReadBoolAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String()
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
static Task<OperateResult<array<bool>^>^>^ ReadBoolAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	array<String^>^ address
)
```

```
static member ReadBoolAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string[] -> Task<OperateResult<bool[]>> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    等待读取的地址列表，数组长度不限制

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[ReadBoolAsync 重载](48143ff1-208b-4c74-ac6f-d7fd22890db8.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/118dfb75-779e-f41c-e457-e570db1735f9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[ReadBoolAsync 方法](../html/48143ff1-208b-4c74-ac6f-d7fd22890db8.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String)](../html/ae0a83b8-a1ec-9dc1-04f0-37a5cffa74e4.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String)")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String[])](../html/14b04612-0172-89fd-4e65-12569b11cd11.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String[])")

[ReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16)](../html/118dfb75-779e-f41c-e457-e570db1735f9.htm "ReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadBoolAsync 方法 (IReadWriteDevice, Byte, String, UInt16) |

批量读取松下PLC的位数据，按照字为单位，地址为 X0,X10,Y10，读取的长度为16的倍数  
Read the bit data of Panasonic PLC in batches, the unit is word, the address is X0, X10, Y10, and the read length is a multiple of 16

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBoolAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
static Task<OperateResult<array<bool>^>^>^ ReadBoolAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBoolAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[ReadBoolAsync 重载](48143ff1-208b-4c74-ac6f-d7fd22890db8.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadPlcType 方法 

[原文連結](http://api.hslcommunication.cn/html/de0359bb-2800-ccf9-ba8b-4d617d2afdd8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[Read 方法](../html/e78daeb8-72d8-2403-794b-0fbede452393.htm "Read 方法 ")

[ReadAsync 方法](../html/3233fb37-7c23-f59a-b2dc-fc1fd738fbdf.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4903b92e-f917-30f9-5f54-9ef5816b063e.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/48143ff1-208b-4c74-ac6f-d7fd22890db8.htm "ReadBoolAsync 方法 ")

[ReadPlcType 方法](../html/de0359bb-2800-ccf9-ba8b-4d617d2afdd8.htm "ReadPlcType 方法 ")

[Write 方法](../html/a96e986e-6858-131c-01a9-2edf0643db5a.htm "Write 方法 ")

[WriteAsync 方法](../html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperReadPlcType 方法 |

读取PLC的型号信息  
Read the model information of the PLC

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> ReadPlcType(
	IReadWriteDevice plc,
	byte station
)
```

```
Public Shared Function ReadPlcType ( 
	plc As IReadWriteDevice,
	station As Byte
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ ReadPlcType(
	IReadWriteDevice^ plc, 
	unsigned char station
)
```

```
static member ReadPlcType : 
        plc : IReadWriteDevice * 
        station : byte -> OperateResult<string> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    通信对象

station
:   类型：SystemByte  
    站号信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
PLC型号

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/a96e986e-6858-131c-01a9-2edf0643db5a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[Write 方法](../html/a96e986e-6858-131c-01a9-2edf0643db5a.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean)](../html/3d2e6b51-292f-d104-0d59-24603db047b9.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean)")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/395acabb-39d1-5ff7-2f67-f488e62d3f70.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/2679d12f-4f12-b54c-6cd5-c82ecee975b1.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

[Write 方法 (IReadWriteDevice, Byte, String[], Boolean[])](../html/d4d2e6b4-3c6e-6318-d094-0b86689189de.htm "Write 方法 (IReadWriteDevice, Byte, String[], Boolean[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](3d2e6b51-292f-d104-0d59-24603db047b9.htm) | 往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](395acabb-39d1-5ff7-2f67-f488e62d3f70.htm) | 往指定的地址写入 Boolean 数组，地址举例 X0.0 Y0.0 R0.0 L0.0， 起始的位地址必须为16的倍数，写入的 Boolean 数组长度也为16的倍数。  Write the Boolean array to the specified address, address example: SR0.0 X0.0 Y0.0 R0.0 L0.0, the starting bit address must be a multiple of 16. Boolean The length of the array is also a multiple of 16. |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Byte)](2679d12f-4f12-b54c-6cd5-c82ecee975b1.htm) | 将数据写入到指定的地址里去，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Write data to the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](d4d2e6b4-3c6e-6318-d094-0b86689189de.htm) | 将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Write the Bool array value to the specified discrete address, one address corresponds to one bool value, the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IReadWriteDevice, Byte, String, Boolean)

[原文連結](http://api.hslcommunication.cn/html/3d2e6b51-292f-d104-0d59-24603db047b9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[Write 方法](../html/a96e986e-6858-131c-01a9-2edf0643db5a.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean)](../html/3d2e6b51-292f-d104-0d59-24603db047b9.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean)")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/395acabb-39d1-5ff7-2f67-f488e62d3f70.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/2679d12f-4f12-b54c-6cd5-c82ecee975b1.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

[Write 方法 (IReadWriteDevice, Byte, String[], Boolean[])](../html/d4d2e6b4-3c6e-6318-d094-0b86689189de.htm "Write 方法 (IReadWriteDevice, Byte, String[], Boolean[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperWrite 方法 (IReadWriteDevice, Byte, String, Boolean) |

往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  
Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	bool value
)
```

```
Public Shared Function Write ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Boolean
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	bool value
)
```

```
static member Write : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : bool -> OperateResult 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemBoolean  
    数据值信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
返回是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[Write 重载](a96e986e-6858-131c-01a9-2edf0643db5a.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IReadWriteDevice, Byte, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/2679d12f-4f12-b54c-6cd5-c82ecee975b1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[Write 方法](../html/a96e986e-6858-131c-01a9-2edf0643db5a.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean)](../html/3d2e6b51-292f-d104-0d59-24603db047b9.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean)")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/395acabb-39d1-5ff7-2f67-f488e62d3f70.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/2679d12f-4f12-b54c-6cd5-c82ecee975b1.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

[Write 方法 (IReadWriteDevice, Byte, String[], Boolean[])](../html/d4d2e6b4-3c6e-6318-d094-0b86689189de.htm "Write 方法 (IReadWriteDevice, Byte, String[], Boolean[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperWrite 方法 (IReadWriteDevice, Byte, String, Byte) |

将数据写入到指定的地址里去，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  
Write data to the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	byte[] value
)
```

```
Public Shared Function Write ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member Write : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : byte[] -> OperateResult 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemByte  
    真实数据

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[Write 重载](a96e986e-6858-131c-01a9-2edf0643db5a.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IReadWriteDevice, Byte, String[], Boolean[])

[原文連結](http://api.hslcommunication.cn/html/d4d2e6b4-3c6e-6318-d094-0b86689189de.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[Write 方法](../html/a96e986e-6858-131c-01a9-2edf0643db5a.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean)](../html/3d2e6b51-292f-d104-0d59-24603db047b9.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean)")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/395acabb-39d1-5ff7-2f67-f488e62d3f70.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/2679d12f-4f12-b54c-6cd5-c82ecee975b1.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

[Write 方法 (IReadWriteDevice, Byte, String[], Boolean[])](../html/d4d2e6b4-3c6e-6318-d094-0b86689189de.htm "Write 方法 (IReadWriteDevice, Byte, String[], Boolean[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperWrite 方法 (IReadWriteDevice, Byte, String, Boolean) |

将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  
Write the Bool array value to the specified discrete address, one address corresponds to one bool value,
the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string[] address,
	bool[] value
)
```

```
Public Shared Function Write ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String(),
	value As Boolean()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	array<String^>^ address, 
	array<bool>^ value
)
```

```
static member Write : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string[] * 
        value : bool[] -> OperateResult 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    离散的地址列表

value
:   类型：SystemBoolean  
    bool数组值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[Write 重载](a96e986e-6858-131c-01a9-2edf0643db5a.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[WriteAsync 方法](../html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)](../html/661ff8d1-c652-d7bf-782e-d79ec60fdc83.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/cc9a08bc-e17f-5db2-d99f-fadc2c115099.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/083baf45-f05f-fdbf-7cb8-17e0e6387352.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])](../html/47b4bee5-97ac-23c5-ae26-ccace6519755.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperWriteAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](661ff8d1-c652-d7bf-782e-d79ec60fdc83.htm) | 往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](cc9a08bc-e17f-5db2-d99f-fadc2c115099.htm) | 往指定的地址写入 Boolean 数组，地址举例 X0.0 Y0.0 R0.0 L0.0， 起始的位地址必须为16的倍数，写入的 Boolean 数组长度也为16的倍数。  Write the Boolean array to the specified address, address example: SR0.0 X0.0 Y0.0 R0.0 L0.0, the starting bit address must be a multiple of 16. Boolean The length of the array is also a multiple of 16. |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Byte)](083baf45-f05f-fdbf-7cb8-17e0e6387352.htm) | 将数据写入到指定的地址里去，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  Write data to the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](47b4bee5-97ac-23c5-ae26-ccace6519755.htm) | 将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  Write the Bool array value to the specified discrete address, one address corresponds to one bool value, the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)

[原文連結](http://api.hslcommunication.cn/html/661ff8d1-c652-d7bf-782e-d79ec60fdc83.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[WriteAsync 方法](../html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)](../html/661ff8d1-c652-d7bf-782e-d79ec60fdc83.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/cc9a08bc-e17f-5db2-d99f-fadc2c115099.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/083baf45-f05f-fdbf-7cb8-17e0e6387352.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])](../html/47b4bee5-97ac-23c5-ae26-ccace6519755.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperWriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean) |

往指定的地址写入bool数据，地址举例：SR0.0 X0.0 Y0.0 R0.0 L0.0  
Write bool data to the specified address. Example address: SR0.0 X0.0 Y0.0 R0.0 L0.0

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	bool value
)
```

```
Public Shared Function WriteAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Boolean
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	bool value
)
```

```
static member WriteAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : bool -> Task<OperateResult> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemBoolean  
    数据值信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
返回是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[WriteAsync 重载](7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/cc9a08bc-e17f-5db2-d99f-fadc2c115099.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[WriteAsync 方法](../html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)](../html/661ff8d1-c652-d7bf-782e-d79ec60fdc83.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/cc9a08bc-e17f-5db2-d99f-fadc2c115099.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/083baf45-f05f-fdbf-7cb8-17e0e6387352.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])](../html/47b4bee5-97ac-23c5-ae26-ccace6519755.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperWriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean) |

往指定的地址写入 Boolean 数组，地址举例 X0.0 Y0.0 R0.0 L0.0，
起始的位地址必须为16的倍数，写入的 Boolean 数组长度也为16的倍数。  
Write the Boolean array to the specified address, address example: SR0.0 X0.0 Y0.0 R0.0 L0.0,
the starting bit address must be a multiple of 16. Boolean The length of the array is also a multiple of 16.

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	bool[] values
)
```

```
Public Shared Function WriteAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	values As Boolean()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	array<bool>^ values
)
```

```
static member WriteAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        values : bool[] -> Task<OperateResult> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

values
:   类型：SystemBoolean  
    数据值信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
返回是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[WriteAsync 重载](7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/083baf45-f05f-fdbf-7cb8-17e0e6387352.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[WriteAsync 方法](../html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)](../html/661ff8d1-c652-d7bf-782e-d79ec60fdc83.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/cc9a08bc-e17f-5db2-d99f-fadc2c115099.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/083baf45-f05f-fdbf-7cb8-17e0e6387352.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])](../html/47b4bee5-97ac-23c5-ae26-ccace6519755.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperWriteAsync 方法 (IReadWriteDevice, Byte, String, Byte) |

将数据写入到指定的地址里去，地址示例：D0 F0 K0 T0 C0, 地址支持携带站号的访问方式，例如：s=2;D100  
Write data to the specified address, address example: D0 F0 K0 T0 C0, the address supports carrying station number information, for example: s=2;D100

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	byte[] value
)
```

```
Public Shared Function WriteAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Byte()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member WriteAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : byte[] -> Task<OperateResult> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemByte  
    真实数据

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[WriteAsync 重载](7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])

[原文連結](http://api.hslcommunication.cn/html/47b4bee5-97ac-23c5-ae26-ccace6519755.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Panasonic.Helper](../html/d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm "HslCommunication.Profinet.Panasonic.Helper")

[MewtocolHelper 类](../html/6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm "MewtocolHelper 类")

[MewtocolHelper 方法](../html/4a07527e-1d2c-5155-f280-76f4bb31dce8.htm "MewtocolHelper 方法")

[WriteAsync 方法](../html/7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)](../html/661ff8d1-c652-d7bf-782e-d79ec60fdc83.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean)")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/cc9a08bc-e17f-5db2-d99f-fadc2c115099.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/083baf45-f05f-fdbf-7cb8-17e0e6387352.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])](../html/47b4bee5-97ac-23c5-ae26-ccace6519755.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String[], Boolean[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MewtocolHelperWriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean) |

将Bool数组值写入到指定的离散地址里，一个地址对应一个bool值，地址数组长度和值数组长度必须相等，地址支持X,Y,R,T,C,L, 举例：R1.0, X2.0, R3.A  
Write the Bool array value to the specified discrete address, one address corresponds to one bool value,
the length of the address array and the length of the value array must be equal, the address supports X, Y, R, T, C, L, for example: R1.0, X2.0, R3.A

**命名空间：**
 [HslCommunication.Profinet.Panasonic.Helper](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string[] address,
	bool[] value
)
```

```
Public Shared Function WriteAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String(),
	value As Boolean()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	array<String^>^ address, 
	array<bool>^ value
)
```

```
static member WriteAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string[] * 
        value : bool[] -> Task<OperateResult> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    离散的地址列表

value
:   类型：SystemBoolean  
    bool数组值

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MewtocolHelper 类](6d2f4fa6-2a5e-f654-b470-949a9fe0f88f.htm)

[WriteAsync 重载](7a0e5fc4-5ab1-e22d-ba38-5fc76e237370.htm)

[HslCommunication.Profinet.Panasonic.Helper 命名空间](d8fa61fb-ea95-000a-e1b6-172fd0ec4ff0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)