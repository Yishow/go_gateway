# HslCommunication - HslCommunication.Reflection

> 分類頁數: 30



---
## HslCommunication.Reflection

[原文連結](http://api.hslcommunication.cn/html/62abbc78-147f-7109-6253-bff7d7ce2461.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslMqttApiAttribute 类](../html/1b28f216-76a9-775b-c564-fcc28b2a4b7f.htm "HslMqttApiAttribute 类")

[HslMqttPermissionAttribute 类](../html/001cd451-c89d-ec5b-832f-7a05e556c192.htm "HslMqttPermissionAttribute 类")

[HslRedisHashFieldAttribute 类](../html/13118df4-a00f-12bb-887d-e78f8b46236b.htm "HslRedisHashFieldAttribute 类")

[HslRedisKeyAttribute 类](../html/0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm "HslRedisKeyAttribute 类")

[HslRedisListAttribute 类](../html/fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm "HslRedisListAttribute 类")

[HslRedisListItemAttribute 类](../html/2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm "HslRedisListItemAttribute 类")

[HslReflectionHelper 类](../html/f85879cd-39be-54b0-a86c-96dd7a667d09.htm "HslReflectionHelper 类")

[HslStructAttribute 类](../html/87e427a8-0c7b-43a8-c421-c27264eda90f.htm "HslStructAttribute 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Reflection 命名空间 |

[缺少 "N:HslCommunication.Reflection" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [HslAddressProperty](a405c5dc-044b-e26d-c270-b95b6e615c30.htm) | Hsl相关地址的属性信息 |
| 公共类 | [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm) | 应用于Hsl组件库读取的动态地址解析，具体用法为创建一个类，创建数据属性，如果这个属性需要绑定PLC的真实数据，就在属性的特性上应用本特性。  Applied to the dynamic address resolution read by the Hsl component library, the specific usage is to create a class and create data attributes. If this attribute needs to be bound to the real data of the PLC, this feature is applied to the characteristics of the attribute. |
| 公共类 | [HslMqttApiAttribute](1b28f216-76a9-775b-c564-fcc28b2a4b7f.htm) | 可以指定方法变成对外公开的API接口，如果方法不实现该特性，将不对外公开方法，无法获取相关的接口权限  You can specify the method to become an externally public API interface. If the method does not implement this feature, the method will not be publicly disclosed, and the related interface permissions cannot be obtained |
| 公共类 | [HslMqttPermissionAttribute](001cd451-c89d-ec5b-832f-7a05e556c192.htm) | 可以指定方法的权限内容，可以限定MQTT会话的ClientID信息或是UserName内容 |
| 公共类 | [HslRedisHashFieldAttribute](13118df4-a00f-12bb-887d-e78f8b46236b.htm) | 对应redis的一个哈希信息的内容 |
| 公共类 | [HslRedisKeyAttribute](0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm) | 对应redis的一个键值信息的内容 |
| 公共类 | [HslRedisListAttribute](fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm) | 对应redis的一个列表信息的内容 |
| 公共类 | [HslRedisListItemAttribute](2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm) | 对应redis的一个列表信息的内容 |
| 公共类 | [HslReflectionHelper](f85879cd-39be-54b0-a86c-96dd7a667d09.htm) | 反射的辅助类 |
| 公共类 | [HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm) | 结构体的字节偏移信息定义，用于方法 [ReadStructT(String, UInt16)](49853082-fbad-52e5-d756-0615cedb4b83.htm)读取字节数据并实现解析操作的  The byte offset information of the structure is defined, Method used for [ReadStructT(String, UInt16)](49853082-fbad-52e5-d756-0615cedb4b83.htm) read byte data and implement the parsing operation |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslAddressProperty 类

[原文連結](http://api.hslcommunication.cn/html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslAddressProperty 构造函数](../html/b286cd90-4ba6-f315-221c-e25c83e9d62d.htm "HslAddressProperty 构造函数 ")

[HslAddressProperty 属性](../html/17165e27-37b0-852d-32bf-4ac5745214cb.htm "HslAddressProperty 属性")

[HslAddressProperty 方法](../html/85cb8cad-8c72-405a-6b2e-eb6cd640a85a.htm "HslAddressProperty 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslAddressProperty 类 |

Hsl相关地址的属性信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.ReflectionHslAddressProperty

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class HslAddressProperty
```

```
Public Class HslAddressProperty
```

```
public ref class HslAddressProperty
```

```
type HslAddressProperty =  class end
```

HslAddressProperty 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HslAddressProperty](b286cd90-4ba6-f315-221c-e25c83e9d62d.htm) | 初始化 HslAddressProperty 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Buffer](4689beb8-35cb-5c3d-55c4-2dc00c3ed90c.htm) | 缓存的数据对象 |
| 公共属性 | [ByteLength](bcc9ff8e-3ddd-b0bf-7218-9644a46b8263.htm) | 读取的字节的长度信息 |
| 公共属性 | [ByteOffset](e0f69b73-b247-dded-4550-5d26e570b483.htm) | 起始的字节偏移信息 |
| 公共属性 | [DeviceAddressAttribute](15eb7c22-f438-8d92-8234-18a5f826076f.htm) | 该属性绑定的地址特性 |
| 公共属性 | [PropertyInfo](3e76daca-e07d-fa72-f8dc-1f832583d800.htm) | 地址绑定的属性信息 |

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

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslAddressProperty 构造函数 

[原文連結](http://api.hslcommunication.cn/html/b286cd90-4ba6-f315-221c-e25c83e9d62d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslAddressProperty 构造函数](../html/b286cd90-4ba6-f315-221c-e25c83e9d62d.htm "HslAddressProperty 构造函数 ")

[HslAddressProperty 属性](../html/17165e27-37b0-852d-32bf-4ac5745214cb.htm "HslAddressProperty 属性")

[HslAddressProperty 方法](../html/85cb8cad-8c72-405a-6b2e-eb6cd640a85a.htm "HslAddressProperty 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslAddressProperty 构造函数 |

初始化 [HslAddressProperty](a405c5dc-044b-e26d-c270-b95b6e615c30.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslAddressProperty()
```

```
Public Sub New
```

```
public:
HslAddressProperty()
```

```
new : unit -> HslAddressProperty
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslAddressProperty 类](a405c5dc-044b-e26d-c270-b95b6e615c30.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslAddressProperty 属性

[原文連結](http://api.hslcommunication.cn/html/17165e27-37b0-852d-32bf-4ac5745214cb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslAddressProperty 属性](../html/17165e27-37b0-852d-32bf-4ac5745214cb.htm "HslAddressProperty 属性")

[Buffer 属性](../html/4689beb8-35cb-5c3d-55c4-2dc00c3ed90c.htm "Buffer 属性 ")

[ByteLength 属性](../html/bcc9ff8e-3ddd-b0bf-7218-9644a46b8263.htm "ByteLength 属性 ")

[ByteOffset 属性](../html/e0f69b73-b247-dded-4550-5d26e570b483.htm "ByteOffset 属性 ")

[DeviceAddressAttribute 属性](../html/15eb7c22-f438-8d92-8234-18a5f826076f.htm "DeviceAddressAttribute 属性 ")

[PropertyInfo 属性](../html/3e76daca-e07d-fa72-f8dc-1f832583d800.htm "PropertyInfo 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslAddressProperty 属性 |

[HslAddressProperty](a405c5dc-044b-e26d-c270-b95b6e615c30.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Buffer](4689beb8-35cb-5c3d-55c4-2dc00c3ed90c.htm) | 缓存的数据对象 |
| 公共属性 | [ByteLength](bcc9ff8e-3ddd-b0bf-7218-9644a46b8263.htm) | 读取的字节的长度信息 |
| 公共属性 | [ByteOffset](e0f69b73-b247-dded-4550-5d26e570b483.htm) | 起始的字节偏移信息 |
| 公共属性 | [DeviceAddressAttribute](15eb7c22-f438-8d92-8234-18a5f826076f.htm) | 该属性绑定的地址特性 |
| 公共属性 | [PropertyInfo](3e76daca-e07d-fa72-f8dc-1f832583d800.htm) | 地址绑定的属性信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslAddressProperty 类](a405c5dc-044b-e26d-c270-b95b6e615c30.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Buffer 属性 

[原文連結](http://api.hslcommunication.cn/html/4689beb8-35cb-5c3d-55c4-2dc00c3ed90c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslAddressProperty 属性](../html/17165e27-37b0-852d-32bf-4ac5745214cb.htm "HslAddressProperty 属性")

[Buffer 属性](../html/4689beb8-35cb-5c3d-55c4-2dc00c3ed90c.htm "Buffer 属性 ")

[ByteLength 属性](../html/bcc9ff8e-3ddd-b0bf-7218-9644a46b8263.htm "ByteLength 属性 ")

[ByteOffset 属性](../html/e0f69b73-b247-dded-4550-5d26e570b483.htm "ByteOffset 属性 ")

[DeviceAddressAttribute 属性](../html/15eb7c22-f438-8d92-8234-18a5f826076f.htm "DeviceAddressAttribute 属性 ")

[PropertyInfo 属性](../html/3e76daca-e07d-fa72-f8dc-1f832583d800.htm "PropertyInfo 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslAddressPropertyBuffer 属性 |

缓存的数据对象

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Buffer { get; set; }
```

```
Public Property Buffer As Byte()
	Get
	Set
```

```
public:
property array<unsigned char>^ Buffer {
	array<unsigned char>^ get ();
	void set (array<unsigned char>^ value);
}
```

```
member Buffer : byte[] with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[HslAddressProperty 类](a405c5dc-044b-e26d-c270-b95b6e615c30.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ByteLength 属性 

[原文連結](http://api.hslcommunication.cn/html/bcc9ff8e-3ddd-b0bf-7218-9644a46b8263.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslAddressProperty 属性](../html/17165e27-37b0-852d-32bf-4ac5745214cb.htm "HslAddressProperty 属性")

[Buffer 属性](../html/4689beb8-35cb-5c3d-55c4-2dc00c3ed90c.htm "Buffer 属性 ")

[ByteLength 属性](../html/bcc9ff8e-3ddd-b0bf-7218-9644a46b8263.htm "ByteLength 属性 ")

[ByteOffset 属性](../html/e0f69b73-b247-dded-4550-5d26e570b483.htm "ByteOffset 属性 ")

[DeviceAddressAttribute 属性](../html/15eb7c22-f438-8d92-8234-18a5f826076f.htm "DeviceAddressAttribute 属性 ")

[PropertyInfo 属性](../html/3e76daca-e07d-fa72-f8dc-1f832583d800.htm "PropertyInfo 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslAddressPropertyByteLength 属性 |

读取的字节的长度信息

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ByteLength { get; set; }
```

```
Public Property ByteLength As Integer
	Get
	Set
```

```
public:
property int ByteLength {
	int get ();
	void set (int value);
}
```

```
member ByteLength : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[HslAddressProperty 类](a405c5dc-044b-e26d-c270-b95b6e615c30.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ByteOffset 属性 

[原文連結](http://api.hslcommunication.cn/html/e0f69b73-b247-dded-4550-5d26e570b483.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslAddressProperty 属性](../html/17165e27-37b0-852d-32bf-4ac5745214cb.htm "HslAddressProperty 属性")

[Buffer 属性](../html/4689beb8-35cb-5c3d-55c4-2dc00c3ed90c.htm "Buffer 属性 ")

[ByteLength 属性](../html/bcc9ff8e-3ddd-b0bf-7218-9644a46b8263.htm "ByteLength 属性 ")

[ByteOffset 属性](../html/e0f69b73-b247-dded-4550-5d26e570b483.htm "ByteOffset 属性 ")

[DeviceAddressAttribute 属性](../html/15eb7c22-f438-8d92-8234-18a5f826076f.htm "DeviceAddressAttribute 属性 ")

[PropertyInfo 属性](../html/3e76daca-e07d-fa72-f8dc-1f832583d800.htm "PropertyInfo 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslAddressPropertyByteOffset 属性 |

起始的字节偏移信息

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
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

[HslAddressProperty 类](a405c5dc-044b-e26d-c270-b95b6e615c30.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceAddressAttribute 属性 

[原文連結](http://api.hslcommunication.cn/html/15eb7c22-f438-8d92-8234-18a5f826076f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslAddressProperty 属性](../html/17165e27-37b0-852d-32bf-4ac5745214cb.htm "HslAddressProperty 属性")

[Buffer 属性](../html/4689beb8-35cb-5c3d-55c4-2dc00c3ed90c.htm "Buffer 属性 ")

[ByteLength 属性](../html/bcc9ff8e-3ddd-b0bf-7218-9644a46b8263.htm "ByteLength 属性 ")

[ByteOffset 属性](../html/e0f69b73-b247-dded-4550-5d26e570b483.htm "ByteOffset 属性 ")

[DeviceAddressAttribute 属性](../html/15eb7c22-f438-8d92-8234-18a5f826076f.htm "DeviceAddressAttribute 属性 ")

[PropertyInfo 属性](../html/3e76daca-e07d-fa72-f8dc-1f832583d800.htm "PropertyInfo 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslAddressPropertyDeviceAddressAttribute 属性 |

该属性绑定的地址特性

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslDeviceAddressAttribute DeviceAddressAttribute { get; set; }
```

```
Public Property DeviceAddressAttribute As HslDeviceAddressAttribute
	Get
	Set
```

```
public:
property HslDeviceAddressAttribute^ DeviceAddressAttribute {
	HslDeviceAddressAttribute^ get ();
	void set (HslDeviceAddressAttribute^ value);
}
```

```
member DeviceAddressAttribute : HslDeviceAddressAttribute with get, set
```

#### 属性值

类型：[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslAddressProperty 类](a405c5dc-044b-e26d-c270-b95b6e615c30.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PropertyInfo 属性 

[原文連結](http://api.hslcommunication.cn/html/3e76daca-e07d-fa72-f8dc-1f832583d800.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslAddressProperty 属性](../html/17165e27-37b0-852d-32bf-4ac5745214cb.htm "HslAddressProperty 属性")

[Buffer 属性](../html/4689beb8-35cb-5c3d-55c4-2dc00c3ed90c.htm "Buffer 属性 ")

[ByteLength 属性](../html/bcc9ff8e-3ddd-b0bf-7218-9644a46b8263.htm "ByteLength 属性 ")

[ByteOffset 属性](../html/e0f69b73-b247-dded-4550-5d26e570b483.htm "ByteOffset 属性 ")

[DeviceAddressAttribute 属性](../html/15eb7c22-f438-8d92-8234-18a5f826076f.htm "DeviceAddressAttribute 属性 ")

[PropertyInfo 属性](../html/3e76daca-e07d-fa72-f8dc-1f832583d800.htm "PropertyInfo 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslAddressPropertyPropertyInfo 属性 |

地址绑定的属性信息

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public PropertyInfo PropertyInfo { get; set; }
```

```
Public Property PropertyInfo As PropertyInfo
	Get
	Set
```

```
public:
property PropertyInfo^ PropertyInfo {
	PropertyInfo^ get ();
	void set (PropertyInfo^ value);
}
```

```
member PropertyInfo : PropertyInfo with get, set
```

#### 属性值

类型：PropertyInfo

![](../icons/SectionExpanded.png)参见

#### 引用

[HslAddressProperty 类](a405c5dc-044b-e26d-c270-b95b6e615c30.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslAddressProperty 方法

[原文連結](http://api.hslcommunication.cn/html/85cb8cad-8c72-405a-6b2e-eb6cd640a85a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslAddressProperty 类](../html/a405c5dc-044b-e26d-c270-b95b6e615c30.htm "HslAddressProperty 类")

[HslAddressProperty 构造函数](../html/b286cd90-4ba6-f315-221c-e25c83e9d62d.htm "HslAddressProperty 构造函数 ")

[HslAddressProperty 属性](../html/17165e27-37b0-852d-32bf-4ac5745214cb.htm "HslAddressProperty 属性")

[HslAddressProperty 方法](../html/85cb8cad-8c72-405a-6b2e-eb6cd640a85a.htm "HslAddressProperty 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslAddressProperty 方法 |

[HslAddressProperty](a405c5dc-044b-e26d-c270-b95b6e615c30.htm) 类型公开以下成员。

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

[HslAddressProperty 类](a405c5dc-044b-e26d-c270-b95b6e615c30.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 类

[原文連結](http://api.hslcommunication.cn/html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 构造函数](../html/1620635f-2d74-e9b9-534f-912318f69b52.htm "HslDeviceAddressAttribute 构造函数 ")

[HslDeviceAddressAttribute 属性](../html/8b5b4713-0385-460f-e503-27ff13b44feb.htm "HslDeviceAddressAttribute 属性")

[HslDeviceAddressAttribute 方法](../html/c8a0b24c-45bc-c5a0-29a5-7dab02512382.htm "HslDeviceAddressAttribute 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 类 |

应用于Hsl组件库读取的动态地址解析，具体用法为创建一个类，创建数据属性，如果这个属性需要绑定PLC的真实数据，就在属性的特性上应用本特性。  
Applied to the dynamic address resolution read by the Hsl component library, the specific usage is to create a class and create data attributes.
If this attribute needs to be bound to the real data of the PLC, this feature is applied to the characteristics of the attribute.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  SystemAttribute  
    HslCommunication.ReflectionHslDeviceAddressAttribute

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class HslDeviceAddressAttribute : Attribute
```

```
Public Class HslDeviceAddressAttribute
	Inherits Attribute
```

```
public ref class HslDeviceAddressAttribute : public Attribute
```

```
type HslDeviceAddressAttribute =  
    class
        inherit Attribute
    end
```

HslDeviceAddressAttribute 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HslDeviceAddressAttribute(String)](de789f37-5813-aa4a-6f45-644dc630e2f6.htm) | 实例化一个地址特性，指定地址信息，用于单变量的数据  Instantiate an address feature, specify the address information, for single variable data |
| 公共方法 | [HslDeviceAddressAttribute(String, Int32)](2531011f-d760-31f2-a1eb-9a2888f297b7.htm) | 实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取  Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays |
| 公共方法 | [HslDeviceAddressAttribute(String, Type)](6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm) | 实例化一个地址特性，指定地址信息，用于单变量的数据，并指定设备类型  Instantiate an address feature, specify address information, data for a single variable, and specify the device type |
| 公共方法 | [HslDeviceAddressAttribute(String, Int32, String)](fabe2455-6d3e-d9dc-e15e-089cee260dff.htm) | 实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取  Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays |
| 公共方法 | [HslDeviceAddressAttribute(String, Int32, Type)](4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm) | 实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取，并指定设备的类型，可用于不同种类的PLC  Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays, and specify the type of equipment, which can be used for different types of PLCs |
| 公共方法 | [HslDeviceAddressAttribute(String, Int32, Type, String)](4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm) | 实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取，并指定设备的类型，可用于不同种类的PLC  Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays, and specify the type of equipment, which can be used for different types of PLCs |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Address](dabc87a9-428d-0d19-196a-fb9d2e2c7880.htm) | 数据的地址信息，真实的设备的地址信息  Data address information, real device address information |
| 公共属性 | [DeviceType](6f9bdfca-b196-bb19-11eb-c53c9ccf6db7.htm) | 设备的类型，如果指定了特殊的PLC，那么该地址就可以支持多种不同PLC  The type of equipment, if a special PLC is specified, then the address can support a variety of different PLCs |
| 公共属性 | [Encoding](e88e527a-089b-a2da-3eb1-3eb4d46653a2.htm) | 如果关联了字符串类型的数据，则表示指定的字符编码，默认 ASCII 编码 |
| 公共属性 | [Length](a0ff7bae-d384-ba08-48a9-4bea9d3bafe3.htm) | 读取的数据长度  Length of data read |
| 公共属性 | TypeId | (继承自 Attribute。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Attribute。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetDataLength](69db5455-1f81-f149-9fb6-7bd791e1247a.htm) | 获取数据的数量信息，如果小于0，则返回1 Get the quantity information of the data, if it is less than 0, return 1 |
| 公共方法 | [GetEncoding](4fc9294e-b7b0-9140-a324-736da5b225da.htm) | 获取当前关联的编码信息，通常用于解析字符串的操作  Gets the encoding information of the current association, usually used for string parsing operations |
| 公共方法 | GetHashCode | (继承自 Attribute。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | IsDefaultAttribute | (继承自 Attribute。) |
| 公共方法 | Match | (继承自 Attribute。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](66e0577d-f1f8-0672-53ed-2e930524ede7.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 构造函数 

[原文連結](http://api.hslcommunication.cn/html/1620635f-2d74-e9b9-534f-912318f69b52.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 构造函数](../html/1620635f-2d74-e9b9-534f-912318f69b52.htm "HslDeviceAddressAttribute 构造函数 ")

[HslDeviceAddressAttribute 构造函数 (String)](../html/de789f37-5813-aa4a-6f45-644dc630e2f6.htm "HslDeviceAddressAttribute 构造函数 (String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32)](../html/2531011f-d760-31f2-a1eb-9a2888f297b7.htm "HslDeviceAddressAttribute 构造函数 (String, Int32)")

[HslDeviceAddressAttribute 构造函数 (String, Type)](../html/6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm "HslDeviceAddressAttribute 构造函数 (String, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, String)](../html/fabe2455-6d3e-d9dc-e15e-089cee260dff.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type)](../html/4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)](../html/4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HslDeviceAddressAttribute(String)](de789f37-5813-aa4a-6f45-644dc630e2f6.htm) | 实例化一个地址特性，指定地址信息，用于单变量的数据  Instantiate an address feature, specify the address information, for single variable data |
| 公共方法 | [HslDeviceAddressAttribute(String, Int32)](2531011f-d760-31f2-a1eb-9a2888f297b7.htm) | 实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取  Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays |
| 公共方法 | [HslDeviceAddressAttribute(String, Type)](6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm) | 实例化一个地址特性，指定地址信息，用于单变量的数据，并指定设备类型  Instantiate an address feature, specify address information, data for a single variable, and specify the device type |
| 公共方法 | [HslDeviceAddressAttribute(String, Int32, String)](fabe2455-6d3e-d9dc-e15e-089cee260dff.htm) | 实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取  Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays |
| 公共方法 | [HslDeviceAddressAttribute(String, Int32, Type)](4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm) | 实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取，并指定设备的类型，可用于不同种类的PLC  Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays, and specify the type of equipment, which can be used for different types of PLCs |
| 公共方法 | [HslDeviceAddressAttribute(String, Int32, Type, String)](4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm) | 实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取，并指定设备的类型，可用于不同种类的PLC  Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays, and specify the type of equipment, which can be used for different types of PLCs |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 构造函数 (String)

[原文連結](http://api.hslcommunication.cn/html/de789f37-5813-aa4a-6f45-644dc630e2f6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 构造函数](../html/1620635f-2d74-e9b9-534f-912318f69b52.htm "HslDeviceAddressAttribute 构造函数 ")

[HslDeviceAddressAttribute 构造函数 (String)](../html/de789f37-5813-aa4a-6f45-644dc630e2f6.htm "HslDeviceAddressAttribute 构造函数 (String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32)](../html/2531011f-d760-31f2-a1eb-9a2888f297b7.htm "HslDeviceAddressAttribute 构造函数 (String, Int32)")

[HslDeviceAddressAttribute 构造函数 (String, Type)](../html/6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm "HslDeviceAddressAttribute 构造函数 (String, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, String)](../html/fabe2455-6d3e-d9dc-e15e-089cee260dff.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type)](../html/4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)](../html/4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 构造函数 (String) |

实例化一个地址特性，指定地址信息，用于单变量的数据  
Instantiate an address feature, specify the address information, for single variable data

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslDeviceAddressAttribute(
	string address
)
```

```
Public Sub New ( 
	address As String
)
```

```
public:
HslDeviceAddressAttribute(
	String^ address
)
```

```
new : 
        address : string -> HslDeviceAddressAttribute
```

#### 参数

address
:   类型：SystemString  
    真实的地址信息

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslDeviceAddressAttribute 重载](1620635f-2d74-e9b9-534f-912318f69b52.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/2531011f-d760-31f2-a1eb-9a2888f297b7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 构造函数](../html/1620635f-2d74-e9b9-534f-912318f69b52.htm "HslDeviceAddressAttribute 构造函数 ")

[HslDeviceAddressAttribute 构造函数 (String)](../html/de789f37-5813-aa4a-6f45-644dc630e2f6.htm "HslDeviceAddressAttribute 构造函数 (String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32)](../html/2531011f-d760-31f2-a1eb-9a2888f297b7.htm "HslDeviceAddressAttribute 构造函数 (String, Int32)")

[HslDeviceAddressAttribute 构造函数 (String, Type)](../html/6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm "HslDeviceAddressAttribute 构造函数 (String, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, String)](../html/fabe2455-6d3e-d9dc-e15e-089cee260dff.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type)](../html/4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)](../html/4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 构造函数 (String, Int32) |

实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取  
Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslDeviceAddressAttribute(
	string address,
	int length
)
```

```
Public Sub New ( 
	address As String,
	length As Integer
)
```

```
public:
HslDeviceAddressAttribute(
	String^ address, 
	int length
)
```

```
new : 
        address : string * 
        length : int -> HslDeviceAddressAttribute
```

#### 参数

address
:   类型：SystemString  
    真实的地址信息

length
:   类型：SystemInt32  
    读取的数据长度

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslDeviceAddressAttribute 重载](1620635f-2d74-e9b9-534f-912318f69b52.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 构造函数 (String, Type)

[原文連結](http://api.hslcommunication.cn/html/6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 构造函数](../html/1620635f-2d74-e9b9-534f-912318f69b52.htm "HslDeviceAddressAttribute 构造函数 ")

[HslDeviceAddressAttribute 构造函数 (String)](../html/de789f37-5813-aa4a-6f45-644dc630e2f6.htm "HslDeviceAddressAttribute 构造函数 (String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32)](../html/2531011f-d760-31f2-a1eb-9a2888f297b7.htm "HslDeviceAddressAttribute 构造函数 (String, Int32)")

[HslDeviceAddressAttribute 构造函数 (String, Type)](../html/6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm "HslDeviceAddressAttribute 构造函数 (String, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, String)](../html/fabe2455-6d3e-d9dc-e15e-089cee260dff.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type)](../html/4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)](../html/4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 构造函数 (String, Type) |

实例化一个地址特性，指定地址信息，用于单变量的数据，并指定设备类型  
Instantiate an address feature, specify address information, data for a single variable, and specify the device type

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslDeviceAddressAttribute(
	string address,
	Type deviceType
)
```

```
Public Sub New ( 
	address As String,
	deviceType As Type
)
```

```
public:
HslDeviceAddressAttribute(
	String^ address, 
	Type^ deviceType
)
```

```
new : 
        address : string * 
        deviceType : Type -> HslDeviceAddressAttribute
```

#### 参数

address
:   类型：SystemString  
    真实的地址信息

deviceType
:   类型：SystemType  
    设备的地址信息

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslDeviceAddressAttribute 重载](1620635f-2d74-e9b9-534f-912318f69b52.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 构造函数 (String, Int32, String)

[原文連結](http://api.hslcommunication.cn/html/fabe2455-6d3e-d9dc-e15e-089cee260dff.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 构造函数](../html/1620635f-2d74-e9b9-534f-912318f69b52.htm "HslDeviceAddressAttribute 构造函数 ")

[HslDeviceAddressAttribute 构造函数 (String)](../html/de789f37-5813-aa4a-6f45-644dc630e2f6.htm "HslDeviceAddressAttribute 构造函数 (String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32)](../html/2531011f-d760-31f2-a1eb-9a2888f297b7.htm "HslDeviceAddressAttribute 构造函数 (String, Int32)")

[HslDeviceAddressAttribute 构造函数 (String, Type)](../html/6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm "HslDeviceAddressAttribute 构造函数 (String, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, String)](../html/fabe2455-6d3e-d9dc-e15e-089cee260dff.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type)](../html/4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)](../html/4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 构造函数 (String, Int32, String) |

实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取  
Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslDeviceAddressAttribute(
	string address,
	int length,
	string encoding
)
```

```
Public Sub New ( 
	address As String,
	length As Integer,
	encoding As String
)
```

```
public:
HslDeviceAddressAttribute(
	String^ address, 
	int length, 
	String^ encoding
)
```

```
new : 
        address : string * 
        length : int * 
        encoding : string -> HslDeviceAddressAttribute
```

#### 参数

address
:   类型：SystemString  
    真实的地址信息

length
:   类型：SystemInt32  
    读取的数据长度

encoding
:   类型：SystemString  
    如果是字符串类型的数据的话，就是字符串编码

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslDeviceAddressAttribute 重载](1620635f-2d74-e9b9-534f-912318f69b52.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 构造函数 (String, Int32, Type)

[原文連結](http://api.hslcommunication.cn/html/4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 构造函数](../html/1620635f-2d74-e9b9-534f-912318f69b52.htm "HslDeviceAddressAttribute 构造函数 ")

[HslDeviceAddressAttribute 构造函数 (String)](../html/de789f37-5813-aa4a-6f45-644dc630e2f6.htm "HslDeviceAddressAttribute 构造函数 (String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32)](../html/2531011f-d760-31f2-a1eb-9a2888f297b7.htm "HslDeviceAddressAttribute 构造函数 (String, Int32)")

[HslDeviceAddressAttribute 构造函数 (String, Type)](../html/6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm "HslDeviceAddressAttribute 构造函数 (String, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, String)](../html/fabe2455-6d3e-d9dc-e15e-089cee260dff.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type)](../html/4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)](../html/4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 构造函数 (String, Int32, Type) |

实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取，并指定设备的类型，可用于不同种类的PLC  
Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays,
and specify the type of equipment, which can be used for different types of PLCs

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslDeviceAddressAttribute(
	string address,
	int length,
	Type deviceType
)
```

```
Public Sub New ( 
	address As String,
	length As Integer,
	deviceType As Type
)
```

```
public:
HslDeviceAddressAttribute(
	String^ address, 
	int length, 
	Type^ deviceType
)
```

```
new : 
        address : string * 
        length : int * 
        deviceType : Type -> HslDeviceAddressAttribute
```

#### 参数

address
:   类型：SystemString  
    真实的地址信息

length
:   类型：SystemInt32  
    读取的数据长度

deviceType
:   类型：SystemType  
    设备类型

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslDeviceAddressAttribute 重载](1620635f-2d74-e9b9-534f-912318f69b52.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)

[原文連結](http://api.hslcommunication.cn/html/4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 构造函数](../html/1620635f-2d74-e9b9-534f-912318f69b52.htm "HslDeviceAddressAttribute 构造函数 ")

[HslDeviceAddressAttribute 构造函数 (String)](../html/de789f37-5813-aa4a-6f45-644dc630e2f6.htm "HslDeviceAddressAttribute 构造函数 (String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32)](../html/2531011f-d760-31f2-a1eb-9a2888f297b7.htm "HslDeviceAddressAttribute 构造函数 (String, Int32)")

[HslDeviceAddressAttribute 构造函数 (String, Type)](../html/6fd33c85-e4d5-b36f-8f31-a656d6cb3eac.htm "HslDeviceAddressAttribute 构造函数 (String, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, String)](../html/fabe2455-6d3e-d9dc-e15e-089cee260dff.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, String)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type)](../html/4b92efe8-2f82-2c47-6446-b4cdba20ca54.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type)")

[HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)](../html/4d261c85-f4ba-dd55-9a1a-fc4b923fbf4a.htm "HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 构造函数 (String, Int32, Type, String) |

实例化一个地址特性，指定地址信息和数据长度，通常应用于数组的批量读取，并指定设备的类型，可用于不同种类的PLC  
Instantiate an address feature, specify address information and data length, usually used in batch reading of arrays,
and specify the type of equipment, which can be used for different types of PLCs

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslDeviceAddressAttribute(
	string address,
	int length,
	Type deviceType,
	string encoding
)
```

```
Public Sub New ( 
	address As String,
	length As Integer,
	deviceType As Type,
	encoding As String
)
```

```
public:
HslDeviceAddressAttribute(
	String^ address, 
	int length, 
	Type^ deviceType, 
	String^ encoding
)
```

```
new : 
        address : string * 
        length : int * 
        deviceType : Type * 
        encoding : string -> HslDeviceAddressAttribute
```

#### 参数

address
:   类型：SystemString  
    真实的地址信息

length
:   类型：SystemInt32  
    读取的数据长度

deviceType
:   类型：SystemType  
    设备类型

encoding
:   类型：SystemString  
    如果是字符串类型的数据的话，就是字符串编码

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslDeviceAddressAttribute 重载](1620635f-2d74-e9b9-534f-912318f69b52.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 属性

[原文連結](http://api.hslcommunication.cn/html/8b5b4713-0385-460f-e503-27ff13b44feb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 属性](../html/8b5b4713-0385-460f-e503-27ff13b44feb.htm "HslDeviceAddressAttribute 属性")

[Address 属性](../html/dabc87a9-428d-0d19-196a-fb9d2e2c7880.htm "Address 属性 ")

[DeviceType 属性](../html/6f9bdfca-b196-bb19-11eb-c53c9ccf6db7.htm "DeviceType 属性 ")

[Encoding 属性](../html/e88e527a-089b-a2da-3eb1-3eb4d46653a2.htm "Encoding 属性 ")

[Length 属性](../html/a0ff7bae-d384-ba08-48a9-4bea9d3bafe3.htm "Length 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 属性 |

[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Address](dabc87a9-428d-0d19-196a-fb9d2e2c7880.htm) | 数据的地址信息，真实的设备的地址信息  Data address information, real device address information |
| 公共属性 | [DeviceType](6f9bdfca-b196-bb19-11eb-c53c9ccf6db7.htm) | 设备的类型，如果指定了特殊的PLC，那么该地址就可以支持多种不同PLC  The type of equipment, if a special PLC is specified, then the address can support a variety of different PLCs |
| 公共属性 | [Encoding](e88e527a-089b-a2da-3eb1-3eb4d46653a2.htm) | 如果关联了字符串类型的数据，则表示指定的字符编码，默认 ASCII 编码 |
| 公共属性 | [Length](a0ff7bae-d384-ba08-48a9-4bea9d3bafe3.htm) | 读取的数据长度  Length of data read |
| 公共属性 | TypeId | (继承自 Attribute。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Address 属性 

[原文連結](http://api.hslcommunication.cn/html/dabc87a9-428d-0d19-196a-fb9d2e2c7880.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 属性](../html/8b5b4713-0385-460f-e503-27ff13b44feb.htm "HslDeviceAddressAttribute 属性")

[Address 属性](../html/dabc87a9-428d-0d19-196a-fb9d2e2c7880.htm "Address 属性 ")

[DeviceType 属性](../html/6f9bdfca-b196-bb19-11eb-c53c9ccf6db7.htm "DeviceType 属性 ")

[Encoding 属性](../html/e88e527a-089b-a2da-3eb1-3eb4d46653a2.htm "Encoding 属性 ")

[Length 属性](../html/a0ff7bae-d384-ba08-48a9-4bea9d3bafe3.htm "Length 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttributeAddress 属性 |

数据的地址信息，真实的设备的地址信息  
Data address information, real device address information

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Address { get; }
```

```
Public ReadOnly Property Address As String
	Get
```

```
public:
property String^ Address {
	String^ get ();
}
```

```
member Address : string with get
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceType 属性 

[原文連結](http://api.hslcommunication.cn/html/6f9bdfca-b196-bb19-11eb-c53c9ccf6db7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 属性](../html/8b5b4713-0385-460f-e503-27ff13b44feb.htm "HslDeviceAddressAttribute 属性")

[Address 属性](../html/dabc87a9-428d-0d19-196a-fb9d2e2c7880.htm "Address 属性 ")

[DeviceType 属性](../html/6f9bdfca-b196-bb19-11eb-c53c9ccf6db7.htm "DeviceType 属性 ")

[Encoding 属性](../html/e88e527a-089b-a2da-3eb1-3eb4d46653a2.htm "Encoding 属性 ")

[Length 属性](../html/a0ff7bae-d384-ba08-48a9-4bea9d3bafe3.htm "Length 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttributeDeviceType 属性 |

设备的类型，如果指定了特殊的PLC，那么该地址就可以支持多种不同PLC  
The type of equipment, if a special PLC is specified, then the address can support a variety of different PLCs

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Type DeviceType { get; set; }
```

```
Public Property DeviceType As Type
	Get
	Set
```

```
public:
property Type^ DeviceType {
	Type^ get ();
	void set (Type^ value);
}
```

```
member DeviceType : Type with get, set
```

#### 属性值

类型：Type

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Encoding 属性 

[原文連結](http://api.hslcommunication.cn/html/e88e527a-089b-a2da-3eb1-3eb4d46653a2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 属性](../html/8b5b4713-0385-460f-e503-27ff13b44feb.htm "HslDeviceAddressAttribute 属性")

[Address 属性](../html/dabc87a9-428d-0d19-196a-fb9d2e2c7880.htm "Address 属性 ")

[DeviceType 属性](../html/6f9bdfca-b196-bb19-11eb-c53c9ccf6db7.htm "DeviceType 属性 ")

[Encoding 属性](../html/e88e527a-089b-a2da-3eb1-3eb4d46653a2.htm "Encoding 属性 ")

[Length 属性](../html/a0ff7bae-d384-ba08-48a9-4bea9d3bafe3.htm "Length 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttributeEncoding 属性 |

如果关联了字符串类型的数据，则表示指定的字符编码，默认 ASCII 编码

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Encoding { get; set; }
```

```
Public Property Encoding As String
	Get
	Set
```

```
public:
property String^ Encoding {
	String^ get ();
	void set (String^ value);
}
```

```
member Encoding : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Length 属性 

[原文連結](http://api.hslcommunication.cn/html/a0ff7bae-d384-ba08-48a9-4bea9d3bafe3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 属性](../html/8b5b4713-0385-460f-e503-27ff13b44feb.htm "HslDeviceAddressAttribute 属性")

[Address 属性](../html/dabc87a9-428d-0d19-196a-fb9d2e2c7880.htm "Address 属性 ")

[DeviceType 属性](../html/6f9bdfca-b196-bb19-11eb-c53c9ccf6db7.htm "DeviceType 属性 ")

[Encoding 属性](../html/e88e527a-089b-a2da-3eb1-3eb4d46653a2.htm "Encoding 属性 ")

[Length 属性](../html/a0ff7bae-d384-ba08-48a9-4bea9d3bafe3.htm "Length 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttributeLength 属性 |

读取的数据长度  
Length of data read

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int Length { get; }
```

```
Public ReadOnly Property Length As Integer
	Get
```

```
public:
property int Length {
	int get ();
}
```

```
member Length : int with get
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslDeviceAddressAttribute 方法

[原文連結](http://api.hslcommunication.cn/html/c8a0b24c-45bc-c5a0-29a5-7dab02512382.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 方法](../html/c8a0b24c-45bc-c5a0-29a5-7dab02512382.htm "HslDeviceAddressAttribute 方法")

[GetDataLength 方法](../html/69db5455-1f81-f149-9fb6-7bd791e1247a.htm "GetDataLength 方法 ")

[GetEncoding 方法](../html/4fc9294e-b7b0-9140-a324-736da5b225da.htm "GetEncoding 方法 ")

[ToString 方法](../html/66e0577d-f1f8-0672-53ed-2e930524ede7.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttribute 方法 |

[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Attribute。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetDataLength](69db5455-1f81-f149-9fb6-7bd791e1247a.htm) | 获取数据的数量信息，如果小于0，则返回1 Get the quantity information of the data, if it is less than 0, return 1 |
| 公共方法 | [GetEncoding](4fc9294e-b7b0-9140-a324-736da5b225da.htm) | 获取当前关联的编码信息，通常用于解析字符串的操作  Gets the encoding information of the current association, usually used for string parsing operations |
| 公共方法 | GetHashCode | (继承自 Attribute。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | IsDefaultAttribute | (继承自 Attribute。) |
| 公共方法 | Match | (继承自 Attribute。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](66e0577d-f1f8-0672-53ed-2e930524ede7.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetDataLength 方法 

[原文連結](http://api.hslcommunication.cn/html/69db5455-1f81-f149-9fb6-7bd791e1247a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 方法](../html/c8a0b24c-45bc-c5a0-29a5-7dab02512382.htm "HslDeviceAddressAttribute 方法")

[GetDataLength 方法](../html/69db5455-1f81-f149-9fb6-7bd791e1247a.htm "GetDataLength 方法 ")

[GetEncoding 方法](../html/4fc9294e-b7b0-9140-a324-736da5b225da.htm "GetEncoding 方法 ")

[ToString 方法](../html/66e0577d-f1f8-0672-53ed-2e930524ede7.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttributeGetDataLength 方法 |

获取数据的数量信息，如果小于0，则返回1
Get the quantity information of the data, if it is less than 0, return 1

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int GetDataLength()
```

```
Public Function GetDataLength As Integer
```

```
public:
int GetDataLength()
```

```
member GetDataLength : unit -> int 
```

#### 返回值

类型：Int32  
数据的个数

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetEncoding 方法 

[原文連結](http://api.hslcommunication.cn/html/4fc9294e-b7b0-9140-a324-736da5b225da.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 方法](../html/c8a0b24c-45bc-c5a0-29a5-7dab02512382.htm "HslDeviceAddressAttribute 方法")

[GetDataLength 方法](../html/69db5455-1f81-f149-9fb6-7bd791e1247a.htm "GetDataLength 方法 ")

[GetEncoding 方法](../html/4fc9294e-b7b0-9140-a324-736da5b225da.htm "GetEncoding 方法 ")

[ToString 方法](../html/66e0577d-f1f8-0672-53ed-2e930524ede7.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttributeGetEncoding 方法 |

获取当前关联的编码信息，通常用于解析字符串的操作  
Gets the encoding information of the current association, usually used for string parsing operations

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Encoding GetEncoding()
```

```
Public Function GetEncoding As Encoding
```

```
public:
Encoding^ GetEncoding()
```

```
member GetEncoding : unit -> Encoding 
```

#### 返回值

类型：Encoding  
字符编码信息

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/66e0577d-f1f8-0672-53ed-2e930524ede7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslDeviceAddressAttribute 类](../html/e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm "HslDeviceAddressAttribute 类")

[HslDeviceAddressAttribute 方法](../html/c8a0b24c-45bc-c5a0-29a5-7dab02512382.htm "HslDeviceAddressAttribute 方法")

[GetDataLength 方法](../html/69db5455-1f81-f149-9fb6-7bd791e1247a.htm "GetDataLength 方法 ")

[GetEncoding 方法](../html/4fc9294e-b7b0-9140-a324-736da5b225da.htm "GetEncoding 方法 ")

[ToString 方法](../html/66e0577d-f1f8-0672-53ed-2e930524ede7.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslDeviceAddressAttributeToString 方法 |

[缺少 "M:HslCommunication.Reflection.HslDeviceAddressAttribute.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
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

[缺少 "M:HslCommunication.Reflection.HslDeviceAddressAttribute.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[HslDeviceAddressAttribute 类](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslMqttApiAttribute 类

[原文連結](http://api.hslcommunication.cn/html/1b28f216-76a9-775b-c564-fcc28b2a4b7f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslMqttApiAttribute 类](../html/1b28f216-76a9-775b-c564-fcc28b2a4b7f.htm "HslMqttApiAttribute 类")

[HslMqttApiAttribute 构造函数](../html/61b386ec-7692-95df-d1ab-dc39fe311e33.htm "HslMqttApiAttribute 构造函数 ")

[HslMqttApiAttribute 属性](../html/814fd139-16f6-7e0e-51f4-533a37635e74.htm "HslMqttApiAttribute 属性")

[HslMqttApiAttribute 方法](../html/5fe84aa0-689a-9596-7140-e946cfaac14a.htm "HslMqttApiAttribute 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMqttApiAttribute 类 |

可以指定方法变成对外公开的API接口，如果方法不实现该特性，将不对外公开方法，无法获取相关的接口权限  
You can specify the method to become an externally public API interface. If the method does not implement this feature,
the method will not be publicly disclosed, and the related interface permissions cannot be obtained

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  SystemAttribute  
    HslCommunication.ReflectionHslMqttApiAttribute

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class HslMqttApiAttribute : Attribute
```

```
Public Class HslMqttApiAttribute
	Inherits Attribute
```

```
public ref class HslMqttApiAttribute : public Attribute
```

```
type HslMqttApiAttribute =  
    class
        inherit Attribute
    end
```

HslMqttApiAttribute 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HslMqttApiAttribute](8f265448-250e-0e52-7812-25f29f4bff3e.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [HslMqttApiAttribute(String)](e046712c-e81d-52a6-8fd6-15e181c1a3f6.htm) | 指定描述内容来实例化一个的对象  Specify the description content to instantiate an object |
| 公共方法 | [HslMqttApiAttribute(String, String)](28d51dc1-5fde-cc91-f817-ca1f074bae56.htm) | 指定接口的路由信息及描述内容来实例化一个的对象  Specify the routing information and description content of the interface to instantiate an object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ApiTopic](c9d9958e-45c2-3730-e7f4-9592f844f472.htm) | 当前指定的ApiTopic信息，如果当前的方法接口不指定别名，那么就使用当前的方法名称  The currently specified ApiTopic information, if the current method interface does not specify an alias, then the current method name is used |
| 公共属性 | [Description](65c71a0f-6bb9-003e-3480-5bcb7aa2d82a.htm) | 当前方法的注释内容  The comment content of the current method |
| 公共属性 | [HttpMethod](b2038d2c-3b40-2bbd-662a-33fda5793c6f.htm) | 如果当前的API接口是支持Http的请求方式，当前属性有效，例如GET,POST  If the current API interface is a request method that supports Http, the current attributes are valid, such as GET, POST |
| 公共属性 | [PropertyUnfold](4da4b638-110f-7115-ba98-a98951dc9ff7.htm) | 当前的属性是否需要展开API信息，默认不展开  Whether the current attribute needs to expand the API information, it is not expanded by default |
| 公共属性 | TypeId | (继承自 Attribute。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Attribute。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Attribute。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | IsDefaultAttribute | (继承自 Attribute。) |
| 公共方法 | Match | (继承自 Attribute。) |
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

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslMqttApiAttribute 构造函数 

[原文連結](http://api.hslcommunication.cn/html/61b386ec-7692-95df-d1ab-dc39fe311e33.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslMqttApiAttribute 类](../html/1b28f216-76a9-775b-c564-fcc28b2a4b7f.htm "HslMqttApiAttribute 类")

[HslMqttApiAttribute 构造函数](../html/61b386ec-7692-95df-d1ab-dc39fe311e33.htm "HslMqttApiAttribute 构造函数 ")

[HslMqttApiAttribute 构造函数](../html/8f265448-250e-0e52-7812-25f29f4bff3e.htm "HslMqttApiAttribute 构造函数 ")

[HslMqttApiAttribute 构造函数 (String)](../html/e046712c-e81d-52a6-8fd6-15e181c1a3f6.htm "HslMqttApiAttribute 构造函数 (String)")

[HslMqttApiAttribute 构造函数 (String, String)](../html/28d51dc1-5fde-cc91-f817-ca1f074bae56.htm "HslMqttApiAttribute 构造函数 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMqttApiAttribute 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HslMqttApiAttribute](8f265448-250e-0e52-7812-25f29f4bff3e.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [HslMqttApiAttribute(String)](e046712c-e81d-52a6-8fd6-15e181c1a3f6.htm) | 指定描述内容来实例化一个的对象  Specify the description content to instantiate an object |
| 公共方法 | [HslMqttApiAttribute(String, String)](28d51dc1-5fde-cc91-f817-ca1f074bae56.htm) | 指定接口的路由信息及描述内容来实例化一个的对象  Specify the routing information and description content of the interface to instantiate an object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMqttApiAttribute 类](1b28f216-76a9-775b-c564-fcc28b2a4b7f.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslMqttApiAttribute 构造函数 

[原文連結](http://api.hslcommunication.cn/html/8f265448-250e-0e52-7812-25f29f4bff3e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Reflection](../html/62abbc78-147f-7109-6253-bff7d7ce2461.htm "HslCommunication.Reflection")

[HslMqttApiAttribute 类](../html/1b28f216-76a9-775b-c564-fcc28b2a4b7f.htm "HslMqttApiAttribute 类")

[HslMqttApiAttribute 构造函数](../html/61b386ec-7692-95df-d1ab-dc39fe311e33.htm "HslMqttApiAttribute 构造函数 ")

[HslMqttApiAttribute 构造函数](../html/8f265448-250e-0e52-7812-25f29f4bff3e.htm "HslMqttApiAttribute 构造函数 ")

[HslMqttApiAttribute 构造函数 (String)](../html/e046712c-e81d-52a6-8fd6-15e181c1a3f6.htm "HslMqttApiAttribute 构造函数 (String)")

[HslMqttApiAttribute 构造函数 (String, String)](../html/28d51dc1-5fde-cc91-f817-ca1f074bae56.htm "HslMqttApiAttribute 构造函数 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMqttApiAttribute 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Reflection](62abbc78-147f-7109-6253-bff7d7ce2461.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslMqttApiAttribute()
```

```
Public Sub New
```

```
public:
HslMqttApiAttribute()
```

```
new : unit -> HslMqttApiAttribute
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMqttApiAttribute 类](1b28f216-76a9-775b-c564-fcc28b2a4b7f.htm)

[HslMqttApiAttribute 重载](61b386ec-7692-95df-d1ab-dc39fe311e33.htm)

[HslCommunication.Reflection 命名空间](62abbc78-147f-7109-6253-bff7d7ce2461.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)