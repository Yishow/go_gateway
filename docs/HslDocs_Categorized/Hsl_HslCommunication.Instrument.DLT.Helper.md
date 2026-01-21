# HslCommunication - HslCommunication.Instrument.DLT.Helper

> 分類頁數: 30



---
## HslCommunication.Instrument.DLT.Helper

[原文連結](http://api.hslcommunication.cn/html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Type 枚举](../html/3adc8897-273c-cb21-b98b-4af641b94e56.htm "DLT645Type 枚举")

[DLT698FcsHelper 类](../html/322dda37-945e-b527-d24a-fabfe833afaa.htm "DLT698FcsHelper 类")

[DLT698Helper 类](../html/9cda8cb9-a789-3ab8-8982-e7fedd62504c.htm "DLT698Helper 类")

[IDlt645 接口](../html/3aaea7ea-ccf9-0290-7ec5-641c25729089.htm "IDlt645 接口")

[IDlt698 接口](../html/b23ce2d3-e978-a355-994c-f99764772338.htm "IDlt698 接口")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Instrument.DLT.Helper 命名空间 |

[缺少 "N:HslCommunication.Instrument.DLT.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [DLT645Helper](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm) | DLT645相关的辅助类 |
| 公共类 | [DLT698FcsHelper](322dda37-945e-b527-d24a-fabfe833afaa.htm) | 校验部分的辅助方法 |
| 公共类 | [DLT698Helper](9cda8cb9-a789-3ab8-8982-e7fedd62504c.htm) | 698 协议的帮助类 |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [IDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm) | DLT645的接口实现 |
| 公共接口 | [IDlt698](b23ce2d3-e978-a355-994c-f99764772338.htm) | DLT698的接口实现 |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [DLT645Type](3adc8897-273c-cb21-b98b-4af641b94e56.htm) | DLT645的类型 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLT645Helper 类

[原文連結](http://api.hslcommunication.cn/html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 构造函数](../html/58a538a1-1e59-6fe9-d17f-641db6b51301.htm "DLT645Helper 构造函数 ")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645Helper 类 |

DLT645相关的辅助类

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Instrument.DLT.HelperDLT645Helper

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DLT645Helper
```

```
Public Class DLT645Helper
```

```
public ref class DLT645Helper
```

```
type DLT645Helper =  class end
```

DLT645Helper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DLT645Helper](58a538a1-1e59-6fe9-d17f-641db6b51301.htm) | 初始化 DLT645Helper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisBytesAddress](ccedd268-3156-f220-6f80-2cead5b8d003.htm) | 从用户输入的地址信息中解析出真实的地址及数据标识 |
| 公共方法静态成员 | [AnalysisIntegerAddress](cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm) | 从用户输入的地址信息中解析出真实的地址及数据标识 |
| 公共方法静态成员 | [BroadcastTime](dd48e0b5-7957-1d0a-442c-768473767de9.htm) | 广播指定的时间，强制从站与主站时间同步，传入DateTime时间对象，没有数据返回。  Broadcast the specified time, force the slave station to synchronize with the master station time, pass in the DateTime time object, and no data will be returned. |
| 公共方法静态成员 | [BroadcastTimeAsync](0b1242af-65e6-2c4d-9760-9a405d382ca7.htm) | 广播指定的时间，强制从站与主站时间同步，传入DateTime时间对象，没有数据返回。  Broadcast the specified time, force the slave station to synchronize with the master station time, pass in the DateTime time object, and no data will be returned. |
| 公共方法静态成员 | [BuildDlt645EntireCommand](2f3b71bb-0138-0674-64b5-81c1f8617812.htm) | 将指定的地址信息，控制码信息，数据域信息打包成完整的报文命令 |
| 公共方法静态成员 | [ChangeBaudRate](cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm) | 更改通信速率，波特率可选 600,1200,2400,4800,9600,19200，其他值无效，可以携带地址域信息，s=1;9600   Change the communication rate, the baud rate can be 600, 1200, 2400, 4800, 9600, 19200, other values are invalid, you can carry address domain information, s=1;9600 |
| 公共方法静态成员 | [ChangeBaudRateAsync](90f1d452-ac4c-23f4-fb70-ad1849840043.htm) | 更改通信速率，波特率可选 600,1200,2400,4800,9600,19200，其他值无效，可以携带地址域信息，s=1;9600   Change the communication rate, the baud rate can be 600, 1200, 2400, 4800, 9600, 19200, other values are invalid, you can carry address domain information, s=1;9600 |
| 公共方法静态成员 | [CheckReceiveDataComplete](992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm) | 判断DLT645的报文是否是完整的 |
| 公共方法静态成员 | [CheckResponse](2587a315-7949-2cae-44ac-f99ec0ceba53.htm) | 检查当前的DLT仪表设备反馈数据信息是否正确 |
| 公共方法静态成员 | [CheckResponseCS](d06c60d7-9201-e615-0780-6d4f00a7efbb.htm) | 检查设备返回的报文信息，是否校验码确认通过 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [FindHeadCode68H](1444533e-f9d8-ef99-ebab-67f62e99c929.htm) | 寻找0x68字节开头的位置信息 |
| 公共方法静态成员 | [FreezeCommand](0d285564-5729-5426-87ce-467f2cd7da85.htm) | 对设备发送冻结命令，默认点对点操作，地址域为 99999999999999 时为广播，数据域格式说明：MMDDhhmm(月日时分)， 99DDhhmm表示月为周期定时冻结，9999hhmm表示日为周期定时冻结，999999mm表示以小时为周期定时冻结，99999999表示瞬时冻结  Send a freeze command to the device, the default point-to-point operation, when the address field is 9999999999999, it is broadcast, and the data field format description: MMDDhhmm (month, day, hour and minute), 99DDhhmm means the month is the periodic fixed freeze, 9999hhmm means the day is the periodic periodic freeze, and 999999mm means the hour It is periodic timed freezing, 99999999 means instantaneous freezing |
| 公共方法静态成员 | [FreezeCommandAsync](4b7116ed-6aad-b74d-1284-6df9c44e8329.htm) | 对设备发送冻结命令，默认点对点操作，地址域为 99999999999999 时为广播，数据域格式说明：MMDDhhmm(月日时分)， 99DDhhmm表示月为周期定时冻结，9999hhmm表示日为周期定时冻结，999999mm表示以小时为周期定时冻结，99999999表示瞬时冻结  Send a freeze command to the device, the default point-to-point operation, when the address field is 9999999999999, it is broadcast, and the data field format description: MMDDhhmm (month, day, hour and minute), 99DDhhmm means the month is the periodic fixed freeze, 9999hhmm means the day is the periodic periodic freeze, and 999999mm means the hour It is periodic timed freezing, 99999999 means instantaneous freezing |
| 公共方法静态成员 | [Function1C](0847bf29-8b25-b11f-0ba2-815c53c07597.htm) | 功能码1C的操作，主要用来控制跳闸（控制类型1A），合闸允许（控制类型1B） |
| 公共方法静态成员 | [GetAddressByteFromString](2629bc4f-8a26-336c-0080-3440d4f37c0f.htm) | 将地址解析成BCD码的地址，并且扩充到12位，不够的补0操作 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. |
| 公共方法静态成员 | [ReadAddress](6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：149100007290  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法静态成员 | [ReadAddressAsync](65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：149100007290  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法静态成员 | [ReadAsync](4835b409-3524-b766-6d33-8df4aa5f6c60.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. |
| 公共方法静态成员 | [ReadDouble](ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm) | 读取指定地址的所有的double数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，然后全部转换为double数据信息  Read all the double data information of the specified address, in general, an address has only one data, but a small number of addresses exist multiple data, and then all converted to double data information |
| 公共方法静态成员 | [ReadDoubleAsync](0b933f21-d4c3-790a-8748-f4504694afc0.htm) | 读取指定地址的所有的double数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，然后全部转换为double数据信息  Read all the double data information of the specified address, in general, an address has only one data, but a small number of addresses exist multiple data, and then all converted to double data information |
| 公共方法静态成员 | [ReadStringArray](c4c395b8-e5c0-e974-c611-56ac6d76442f.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，例如 01-01-00-00 正向有功总需求及发生时间  Read all the string data information of the specified address, in general, there is only one data for one address, but there are multiple data for a small number of addresses, such as 01-01-00-00 Forward active total demand and occurrence time |
| 公共方法静态成员 | [ReadStringArrayAsync](cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，例如 01-01-00-00 正向有功总需求及发生时间  Read all the string data information of the specified address, in general, there is only one data for one address, but there are multiple data for a small number of addresses, such as 01-01-00-00 Forward active total demand and occurrence time |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IDlt645, String, String, String, Byte)](7ff3b579-6723-2ef5-b6ac-b646cb493846.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. |
| 公共方法静态成员 | [Write(IDlt645, String, String, String, String)](72105f34-878f-1335-1c37-cf2f1d888a77.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. |
| 公共方法静态成员 | [WriteAddress](9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：149100007290  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 149100007290 |
| 公共方法静态成员 | [WriteAddressAsync](a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：149100007290  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 149100007290 |
| 公共方法静态成员 | [WriteAsync(IDlt645, String, String, String, Byte)](5bb1110a-9eff-f921-c9e7-b0ae44a2c7c5.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. |
| 公共方法静态成员 | [WriteAsync(IDlt645, String, String, String, String)](9a08a056-3147-7db0-73ac-97cd81a038fb.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLT645Helper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/58a538a1-1e59-6fe9-d17f-641db6b51301.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 构造函数](../html/58a538a1-1e59-6fe9-d17f-641db6b51301.htm "DLT645Helper 构造函数 ")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645Helper 构造函数 |

初始化 [DLT645Helper](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DLT645Helper()
```

```
Public Sub New
```

```
public:
DLT645Helper()
```

```
new : unit -> DLT645Helper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLT645Helper 方法

[原文連結](http://api.hslcommunication.cn/html/ca652825-862f-2923-5313-d1b5b7340723.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645Helper 方法 |

[DLT645Helper](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisBytesAddress](ccedd268-3156-f220-6f80-2cead5b8d003.htm) | 从用户输入的地址信息中解析出真实的地址及数据标识 |
| 公共方法静态成员 | [AnalysisIntegerAddress](cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm) | 从用户输入的地址信息中解析出真实的地址及数据标识 |
| 公共方法静态成员 | [BroadcastTime](dd48e0b5-7957-1d0a-442c-768473767de9.htm) | 广播指定的时间，强制从站与主站时间同步，传入DateTime时间对象，没有数据返回。  Broadcast the specified time, force the slave station to synchronize with the master station time, pass in the DateTime time object, and no data will be returned. |
| 公共方法静态成员 | [BroadcastTimeAsync](0b1242af-65e6-2c4d-9760-9a405d382ca7.htm) | 广播指定的时间，强制从站与主站时间同步，传入DateTime时间对象，没有数据返回。  Broadcast the specified time, force the slave station to synchronize with the master station time, pass in the DateTime time object, and no data will be returned. |
| 公共方法静态成员 | [BuildDlt645EntireCommand](2f3b71bb-0138-0674-64b5-81c1f8617812.htm) | 将指定的地址信息，控制码信息，数据域信息打包成完整的报文命令 |
| 公共方法静态成员 | [ChangeBaudRate](cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm) | 更改通信速率，波特率可选 600,1200,2400,4800,9600,19200，其他值无效，可以携带地址域信息，s=1;9600   Change the communication rate, the baud rate can be 600, 1200, 2400, 4800, 9600, 19200, other values are invalid, you can carry address domain information, s=1;9600 |
| 公共方法静态成员 | [ChangeBaudRateAsync](90f1d452-ac4c-23f4-fb70-ad1849840043.htm) | 更改通信速率，波特率可选 600,1200,2400,4800,9600,19200，其他值无效，可以携带地址域信息，s=1;9600   Change the communication rate, the baud rate can be 600, 1200, 2400, 4800, 9600, 19200, other values are invalid, you can carry address domain information, s=1;9600 |
| 公共方法静态成员 | [CheckReceiveDataComplete](992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm) | 判断DLT645的报文是否是完整的 |
| 公共方法静态成员 | [CheckResponse](2587a315-7949-2cae-44ac-f99ec0ceba53.htm) | 检查当前的DLT仪表设备反馈数据信息是否正确 |
| 公共方法静态成员 | [CheckResponseCS](d06c60d7-9201-e615-0780-6d4f00a7efbb.htm) | 检查设备返回的报文信息，是否校验码确认通过 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [FindHeadCode68H](1444533e-f9d8-ef99-ebab-67f62e99c929.htm) | 寻找0x68字节开头的位置信息 |
| 公共方法静态成员 | [FreezeCommand](0d285564-5729-5426-87ce-467f2cd7da85.htm) | 对设备发送冻结命令，默认点对点操作，地址域为 99999999999999 时为广播，数据域格式说明：MMDDhhmm(月日时分)， 99DDhhmm表示月为周期定时冻结，9999hhmm表示日为周期定时冻结，999999mm表示以小时为周期定时冻结，99999999表示瞬时冻结  Send a freeze command to the device, the default point-to-point operation, when the address field is 9999999999999, it is broadcast, and the data field format description: MMDDhhmm (month, day, hour and minute), 99DDhhmm means the month is the periodic fixed freeze, 9999hhmm means the day is the periodic periodic freeze, and 999999mm means the hour It is periodic timed freezing, 99999999 means instantaneous freezing |
| 公共方法静态成员 | [FreezeCommandAsync](4b7116ed-6aad-b74d-1284-6df9c44e8329.htm) | 对设备发送冻结命令，默认点对点操作，地址域为 99999999999999 时为广播，数据域格式说明：MMDDhhmm(月日时分)， 99DDhhmm表示月为周期定时冻结，9999hhmm表示日为周期定时冻结，999999mm表示以小时为周期定时冻结，99999999表示瞬时冻结  Send a freeze command to the device, the default point-to-point operation, when the address field is 9999999999999, it is broadcast, and the data field format description: MMDDhhmm (month, day, hour and minute), 99DDhhmm means the month is the periodic fixed freeze, 9999hhmm means the day is the periodic periodic freeze, and 999999mm means the hour It is periodic timed freezing, 99999999 means instantaneous freezing |
| 公共方法静态成员 | [Function1C](0847bf29-8b25-b11f-0ba2-815c53c07597.htm) | 功能码1C的操作，主要用来控制跳闸（控制类型1A），合闸允许（控制类型1B） |
| 公共方法静态成员 | [GetAddressByteFromString](2629bc4f-8a26-336c-0080-3440d4f37c0f.htm) | 将地址解析成BCD码的地址，并且扩充到12位，不够的补0操作 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. |
| 公共方法静态成员 | [ReadAddress](6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：149100007290  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法静态成员 | [ReadAddressAsync](65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：149100007290  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法静态成员 | [ReadAsync](4835b409-3524-b766-6d33-8df4aa5f6c60.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. |
| 公共方法静态成员 | [ReadDouble](ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm) | 读取指定地址的所有的double数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，然后全部转换为double数据信息  Read all the double data information of the specified address, in general, an address has only one data, but a small number of addresses exist multiple data, and then all converted to double data information |
| 公共方法静态成员 | [ReadDoubleAsync](0b933f21-d4c3-790a-8748-f4504694afc0.htm) | 读取指定地址的所有的double数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，然后全部转换为double数据信息  Read all the double data information of the specified address, in general, an address has only one data, but a small number of addresses exist multiple data, and then all converted to double data information |
| 公共方法静态成员 | [ReadStringArray](c4c395b8-e5c0-e974-c611-56ac6d76442f.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，例如 01-01-00-00 正向有功总需求及发生时间  Read all the string data information of the specified address, in general, there is only one data for one address, but there are multiple data for a small number of addresses, such as 01-01-00-00 Forward active total demand and occurrence time |
| 公共方法静态成员 | [ReadStringArrayAsync](cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，例如 01-01-00-00 正向有功总需求及发生时间  Read all the string data information of the specified address, in general, there is only one data for one address, but there are multiple data for a small number of addresses, such as 01-01-00-00 Forward active total demand and occurrence time |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IDlt645, String, String, String, Byte)](7ff3b579-6723-2ef5-b6ac-b646cb493846.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. |
| 公共方法静态成员 | [Write(IDlt645, String, String, String, String)](72105f34-878f-1335-1c37-cf2f1d888a77.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. |
| 公共方法静态成员 | [WriteAddress](9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：149100007290  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 149100007290 |
| 公共方法静态成员 | [WriteAddressAsync](a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：149100007290  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 149100007290 |
| 公共方法静态成员 | [WriteAsync(IDlt645, String, String, String, Byte)](5bb1110a-9eff-f921-c9e7-b0ae44a2c7c5.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. |
| 公共方法静态成员 | [WriteAsync(IDlt645, String, String, String, String)](9a08a056-3147-7db0-73ac-97cd81a038fb.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AnalysisBytesAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/ccedd268-3156-f220-6f80-2cead5b8d003.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperAnalysisBytesAddress 方法 |

从用户输入的地址信息中解析出真实的地址及数据标识

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
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
	DLT645Type type,
	string address,
	string defaultStation,
	ushort length = 1
)
```

```
Public Shared Function AnalysisBytesAddress ( 
	type As DLT645Type,
	address As String,
	defaultStation As String,
	Optional length As UShort = 1
) As OperateResult(Of String, Byte())
```

```
public:
static OperateResult<String^, array<unsigned char>^>^ AnalysisBytesAddress(
	DLT645Type type, 
	String^ address, 
	String^ defaultStation, 
	unsigned short length = 1
)
```

```
static member AnalysisBytesAddress : 
        type : DLT645Type * 
        address : string * 
        defaultStation : string * 
        ?length : uint16 
(* Defaults:
        let _length = defaultArg length 1
*)
-> OperateResult<string, byte[]> 
```

#### 参数

type
:   类型：[HslCommunication.Instrument.DLT.HelperDLT645Type](3adc8897-273c-cb21-b98b-4af641b94e56.htm)  
    DLT的类型

address
:   类型：SystemString  
    用户输入的地址信息

defaultStation
:   类型：SystemString  
    默认的地址域

length (Optional)
:   类型：SystemUInt16  
    数据长度信息

#### 返回值

类型：[OperateResult](f52f888f-5e8d-b0c4-2302-81e70c230a31.htm)String, Byte  
解析结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AnalysisIntegerAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperAnalysisIntegerAddress 方法 |

从用户输入的地址信息中解析出真实的地址及数据标识

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string, int> AnalysisIntegerAddress(
	string address,
	string defaultStation
)
```

```
Public Shared Function AnalysisIntegerAddress ( 
	address As String,
	defaultStation As String
) As OperateResult(Of String, Integer)
```

```
public:
static OperateResult<String^, int>^ AnalysisIntegerAddress(
	String^ address, 
	String^ defaultStation
)
```

```
static member AnalysisIntegerAddress : 
        address : string * 
        defaultStation : string -> OperateResult<string, int> 
```

#### 参数

address
:   类型：SystemString  
    用户输入的地址信息

defaultStation
:   类型：SystemString  
    默认的地址域

#### 返回值

类型：[OperateResult](f52f888f-5e8d-b0c4-2302-81e70c230a31.htm)String, Int32  
解析结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BroadcastTime 方法 

[原文連結](http://api.hslcommunication.cn/html/dd48e0b5-7957-1d0a-442c-768473767de9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperBroadcastTime 方法 |

广播指定的时间，强制从站与主站时间同步，传入DateTime时间对象，没有数据返回。  
Broadcast the specified time, force the slave station to synchronize with the master station time,
pass in the DateTime time object, and no data will be returned.

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult BroadcastTime(
	IDlt645 dlt,
	DateTime dateTime
)
```

```
Public Shared Function BroadcastTime ( 
	dlt As IDlt645,
	dateTime As DateTime
) As OperateResult
```

```
public:
static OperateResult^ BroadcastTime(
	IDlt645^ dlt, 
	DateTime dateTime
)
```

```
static member BroadcastTime : 
        dlt : IDlt645 * 
        dateTime : DateTime -> OperateResult 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

dateTime
:   类型：SystemDateTime  
    时间对象

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BroadcastTimeAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperBroadcastTimeAsync 方法 |

广播指定的时间，强制从站与主站时间同步，传入DateTime时间对象，没有数据返回。  
Broadcast the specified time, force the slave station to synchronize with the master station time,
pass in the DateTime time object, and no data will be returned.

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> BroadcastTimeAsync(
	IDlt645 dlt,
	DateTime dateTime,
	Func<byte[], bool, bool, Task<OperateResult<byte[]>>> func
)
```

```
Public Shared Function BroadcastTimeAsync ( 
	dlt As IDlt645,
	dateTime As DateTime,
	func As Func(Of Byte(), Boolean, Boolean, Task(Of OperateResult(Of Byte())))
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ BroadcastTimeAsync(
	IDlt645^ dlt, 
	DateTime dateTime, 
	Func<array<unsigned char>^, bool, bool, Task<OperateResult<array<unsigned char>^>^>^>^ func
)
```

```
static member BroadcastTimeAsync : 
        dlt : IDlt645 * 
        dateTime : DateTime * 
        func : Func<byte[], bool, bool, Task<OperateResult<byte[]>>> -> Task<OperateResult> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

dateTime
:   类型：SystemDateTime  
    时间对象

func
:   类型：SystemFuncByte, Boolean, Boolean, Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  

    [缺少 "M:HslCommunication.Instrument.DLT.Helper.DLT645Helper.BroadcastTimeAsync(HslCommunication.Instrument.DLT.Helper.IDlt645,System.DateTime,System.Func{System.Byte[],System.Boolean,System.Boolean,System.Threading.Tasks.Task{HslCommunication.OperateResult{System.Byte[]}}})" 的 <param name="func"/> 文档]

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildDlt645EntireCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperBuildDlt645EntireCommand 方法 |

将指定的地址信息，控制码信息，数据域信息打包成完整的报文命令

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildDlt645EntireCommand(
	string address,
	byte control,
	byte[] dataArea
)
```

```
Public Shared Function BuildDlt645EntireCommand ( 
	address As String,
	control As Byte,
	dataArea As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildDlt645EntireCommand(
	String^ address, 
	unsigned char control, 
	array<unsigned char>^ dataArea
)
```

```
static member BuildDlt645EntireCommand : 
        address : string * 
        control : byte * 
        dataArea : byte[] -> OperateResult<byte[]> 
```

#### 参数

address
:   类型：SystemString  
    地址域信息，地址域由6个字节构成，每字节2位BCD码，地址长度可达12位十进制数。地址域支持锁位寻址，即从若干低位起，剩余高位补AAH作为通配符进行读表操作

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

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ChangeBaudRate 方法 

[原文連結](http://api.hslcommunication.cn/html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperChangeBaudRate 方法 |

更改通信速率，波特率可选 600,1200,2400,4800,9600,19200，其他值无效，可以携带地址域信息，s=1;9600   
Change the communication rate, the baud rate can be 600, 1200, 2400, 4800, 9600, 19200,
other values are invalid, you can carry address domain information, s=1;9600

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult ChangeBaudRate(
	IDlt645 dlt,
	string baudRate
)
```

```
Public Shared Function ChangeBaudRate ( 
	dlt As IDlt645,
	baudRate As String
) As OperateResult
```

```
public:
static OperateResult^ ChangeBaudRate(
	IDlt645^ dlt, 
	String^ baudRate
)
```

```
static member ChangeBaudRate : 
        dlt : IDlt645 * 
        baudRate : string -> OperateResult 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

baudRate
:   类型：SystemString  
    波特率的信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否更改成功

![](../icons/SectionExpanded.png)备注

对于DLT1997来说，只支持 300, 600, 2400, 4800, 9600

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ChangeBaudRateAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperChangeBaudRateAsync 方法 |

更改通信速率，波特率可选 600,1200,2400,4800,9600,19200，其他值无效，可以携带地址域信息，s=1;9600   
Change the communication rate, the baud rate can be 600, 1200, 2400, 4800, 9600, 19200,
other values are invalid, you can carry address domain information, s=1;9600

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> ChangeBaudRateAsync(
	IDlt645 dlt,
	string baudRate
)
```

```
Public Shared Function ChangeBaudRateAsync ( 
	dlt As IDlt645,
	baudRate As String
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ ChangeBaudRateAsync(
	IDlt645^ dlt, 
	String^ baudRate
)
```

```
static member ChangeBaudRateAsync : 
        dlt : IDlt645 * 
        baudRate : string -> Task<OperateResult> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

baudRate
:   类型：SystemString  
    波特率的信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否更改成功

![](../icons/SectionExpanded.png)备注

对于DLT1997来说，只支持 300, 600, 2400, 4800, 9600

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckReceiveDataComplete 方法 

[原文連結](http://api.hslcommunication.cn/html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperCheckReceiveDataComplete 方法 |

判断DLT645的报文是否是完整的

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool CheckReceiveDataComplete(
	MemoryStream ms
)
```

```
Public Shared Function CheckReceiveDataComplete ( 
	ms As MemoryStream
) As Boolean
```

```
public:
static bool CheckReceiveDataComplete(
	MemoryStream^ ms
)
```

```
static member CheckReceiveDataComplete : 
        ms : MemoryStream -> bool 
```

#### 参数

ms
:   类型：System.IOMemoryStream  
    内存数据信息

#### 返回值

类型：Boolean  
是否完整的

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckResponse 方法 

[原文連結](http://api.hslcommunication.cn/html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperCheckResponse 方法 |

检查当前的DLT仪表设备反馈数据信息是否正确

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
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
	IDlt645 dlt,
	byte[] send,
	byte[] response
)
```

```
Public Shared Function CheckResponse ( 
	dlt As IDlt645,
	send As Byte(),
	response As Byte()
) As OperateResult
```

```
public:
static OperateResult^ CheckResponse(
	IDlt645^ dlt, 
	array<unsigned char>^ send, 
	array<unsigned char>^ response
)
```

```
static member CheckResponse : 
        dlt : IDlt645 * 
        send : byte[] * 
        response : byte[] -> OperateResult 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信设备

send
:   类型：SystemByte  
    发送到DLT仪表的报文信息

response
:   类型：SystemByte  
    从仪表反馈的数据信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否校验成功

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckResponseCS 方法 

[原文連結](http://api.hslcommunication.cn/html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperCheckResponseCS 方法 |

检查设备返回的报文信息，是否校验码确认通过

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult CheckResponseCS(
	byte[] response,
	int index
)
```

```
Public Shared Function CheckResponseCS ( 
	response As Byte(),
	index As Integer
) As OperateResult
```

```
public:
static OperateResult^ CheckResponseCS(
	array<unsigned char>^ response, 
	int index
)
```

```
static member CheckResponseCS : 
        response : byte[] * 
        index : int -> OperateResult 
```

#### 参数

response
:   类型：SystemByte  
    设备返回的报文

index
:   类型：SystemInt32  
    起始校验的索引

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否校验成功

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FindHeadCode68H 方法 

[原文連結](http://api.hslcommunication.cn/html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperFindHeadCode68H 方法 |

寻找0x68字节开头的位置信息

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static int FindHeadCode68H(
	byte[] buffer
)
```

```
Public Shared Function FindHeadCode68H ( 
	buffer As Byte()
) As Integer
```

```
public:
static int FindHeadCode68H(
	array<unsigned char>^ buffer
)
```

```
static member FindHeadCode68H : 
        buffer : byte[] -> int 
```

#### 参数

buffer
:   类型：SystemByte  
    缓存数据

#### 返回值

类型：Int32  
如果有则为索引位置，如果没有则为空

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FreezeCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/0d285564-5729-5426-87ce-467f2cd7da85.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperFreezeCommand 方法 |

对设备发送冻结命令，默认点对点操作，地址域为 99999999999999 时为广播，数据域格式说明：MMDDhhmm(月日时分)，
99DDhhmm表示月为周期定时冻结，9999hhmm表示日为周期定时冻结，999999mm表示以小时为周期定时冻结，99999999表示瞬时冻结  
Send a freeze command to the device, the default point-to-point operation, when the address field is 9999999999999,
it is broadcast, and the data field format description: MMDDhhmm (month, day, hour and minute),
99DDhhmm means the month is the periodic fixed freeze, 9999hhmm means the day is the periodic periodic freeze,
and 999999mm means the hour It is periodic timed freezing, 99999999 means instantaneous freezing

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult FreezeCommand(
	IDlt645 dlt,
	string dataArea
)
```

```
Public Shared Function FreezeCommand ( 
	dlt As IDlt645,
	dataArea As String
) As OperateResult
```

```
public:
static OperateResult^ FreezeCommand(
	IDlt645^ dlt, 
	String^ dataArea
)
```

```
static member FreezeCommand : 
        dlt : IDlt645 * 
        dataArea : string -> OperateResult 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

dataArea
:   类型：SystemString  
    数据域信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功冻结

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FreezeCommandAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperFreezeCommandAsync 方法 |

对设备发送冻结命令，默认点对点操作，地址域为 99999999999999 时为广播，数据域格式说明：MMDDhhmm(月日时分)，
99DDhhmm表示月为周期定时冻结，9999hhmm表示日为周期定时冻结，999999mm表示以小时为周期定时冻结，99999999表示瞬时冻结  
Send a freeze command to the device, the default point-to-point operation, when the address field is 9999999999999,
it is broadcast, and the data field format description: MMDDhhmm (month, day, hour and minute),
99DDhhmm means the month is the periodic fixed freeze, 9999hhmm means the day is the periodic periodic freeze,
and 999999mm means the hour It is periodic timed freezing, 99999999 means instantaneous freezing

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> FreezeCommandAsync(
	DLT645OverTcp dlt,
	string dataArea
)
```

```
Public Shared Function FreezeCommandAsync ( 
	dlt As DLT645OverTcp,
	dataArea As String
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ FreezeCommandAsync(
	DLT645OverTcp^ dlt, 
	String^ dataArea
)
```

```
static member FreezeCommandAsync : 
        dlt : DLT645OverTcp * 
        dataArea : string -> Task<OperateResult> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLTDLT645OverTcp](4f76bd69-c6bb-adad-7580-6d8a82d243e4.htm)  
    DLT通信对象

dataArea
:   类型：SystemString  
    数据域信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功冻结

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Function1C 方法 

[原文連結](http://api.hslcommunication.cn/html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperFunction1C 方法 |

功能码1C的操作，主要用来控制跳闸（控制类型1A），合闸允许（控制类型1B）

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Function1C(
	IDlt645 dlt,
	string password,
	string opCode,
	string station,
	byte controlType,
	DateTime validTime
)
```

```
Public Shared Function Function1C ( 
	dlt As IDlt645,
	password As String,
	opCode As String,
	station As String,
	controlType As Byte,
	validTime As DateTime
) As OperateResult
```

```
public:
static OperateResult^ Function1C(
	IDlt645^ dlt, 
	String^ password, 
	String^ opCode, 
	String^ station, 
	unsigned char controlType, 
	DateTime validTime
)
```

```
static member Function1C : 
        dlt : IDlt645 * 
        password : string * 
        opCode : string * 
        station : string * 
        controlType : byte * 
        validTime : DateTime -> OperateResult 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

password
:   类型：SystemString  
    密钥信息

opCode
:   类型：SystemString  
    操作者代码

station
:   类型：SystemString  
    站号信息

controlType
:   类型：SystemByte  
    控制类型

validTime
:   类型：SystemDateTime  
    有效截止时间

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否操作成功

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetAddressByteFromString 方法 

[原文連結](http://api.hslcommunication.cn/html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperGetAddressByteFromString 方法 |

将地址解析成BCD码的地址，并且扩充到12位，不够的补0操作

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
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

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperRead 方法 |

根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  
Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual,
from high to position, such as 00-00-00-00. The separator can be any special character or no separator.

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
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
	IDlt645 dlt,
	string address,
	ushort length
)
```

```
Public Shared Function Read ( 
	dlt As IDlt645,
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	IDlt645^ dlt, 
	String^ address, 
	unsigned short length
)
```

```
static member Read : 
        dlt : IDlt645 * 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

length
:   类型：SystemUInt16  
    数据长度信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
结果信息

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;00-00-00-00" 或是 "s=100000;00-00-02-00"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperReadAddress 方法 |

读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：149100007290  
Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
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
	IDlt645 dlt
)
```

```
Public Shared Function ReadAddress ( 
	dlt As IDlt645
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ ReadAddress(
	IDlt645^ dlt
)
```

```
static member ReadAddress : 
        dlt : IDlt645 -> OperateResult<string> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
设备的通信地址

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAddressAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperReadAddressAsync 方法 |

读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：149100007290  
Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<string>> ReadAddressAsync(
	IDlt645 dlt
)
```

```
Public Shared Function ReadAddressAsync ( 
	dlt As IDlt645
) As Task(Of OperateResult(Of String))
```

```
public:
static Task<OperateResult<String^>^>^ ReadAddressAsync(
	IDlt645^ dlt
)
```

```
static member ReadAddressAsync : 
        dlt : IDlt645 -> Task<OperateResult<string>> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
设备的通信地址

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperReadAsync 方法 |

根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  
Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual,
from high to position, such as 00-00-00-00. The separator can be any special character or no separator.

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
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
	IDlt645 dlt,
	string address,
	ushort length
)
```

```
Public Shared Function ReadAsync ( 
	dlt As IDlt645,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	IDlt645^ dlt, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadAsync : 
        dlt : IDlt645 * 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

length
:   类型：SystemUInt16  
    数据长度信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
结果信息

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;00-00-00-00" 或是 "s=100000;00-00-02-00"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDouble 方法 

[原文連結](http://api.hslcommunication.cn/html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperReadDouble 方法 |

读取指定地址的所有的double数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，然后全部转换为double数据信息  
Read all the double data information of the specified address, in general, an address has only one data, but a small number of addresses exist multiple data,
and then all converted to double data information

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<double[]> ReadDouble(
	IDlt645 dlt,
	string address,
	ushort length
)
```

```
Public Shared Function ReadDouble ( 
	dlt As IDlt645,
	address As String,
	length As UShort
) As OperateResult(Of Double())
```

```
public:
static OperateResult<array<double>^>^ ReadDouble(
	IDlt645^ dlt, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadDouble : 
        dlt : IDlt645 * 
        address : string * 
        length : uint16 -> OperateResult<float[]> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

length
:   类型：SystemUInt16  
    读取的数据长度信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Double  

[缺少 "M:HslCommunication.Instrument.DLT.Helper.DLT645Helper.ReadDouble(HslCommunication.Instrument.DLT.Helper.IDlt645,System.String,System.UInt16)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;00-00-00-00" 或是 "s=100000;00-00-02-00"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能  
地址也可以携带是否数据翻转的标记，例如 "reverse=false;00-00-00-00" 解析数据的时候就不发生反转的操作

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDoubleAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/0b933f21-d4c3-790a-8748-f4504694afc0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperReadDoubleAsync 方法 |

读取指定地址的所有的double数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，然后全部转换为double数据信息  
Read all the double data information of the specified address, in general, an address has only one data, but a small number of addresses exist multiple data,
and then all converted to double data information

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<double[]>> ReadDoubleAsync(
	IDlt645 dlt,
	string address,
	ushort length
)
```

```
Public Shared Function ReadDoubleAsync ( 
	dlt As IDlt645,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Double()))
```

```
public:
static Task<OperateResult<array<double>^>^>^ ReadDoubleAsync(
	IDlt645^ dlt, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadDoubleAsync : 
        dlt : IDlt645 * 
        address : string * 
        length : uint16 -> Task<OperateResult<float[]>> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

length
:   类型：SystemUInt16  
    读取的数据长度信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Double  

[缺少 "M:HslCommunication.Instrument.DLT.Helper.DLT645Helper.ReadDoubleAsync(HslCommunication.Instrument.DLT.Helper.IDlt645,System.String,System.UInt16)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;00-00-00-00" 或是 "s=100000;00-00-02-00"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能  
地址也可以携带是否数据翻转的标记，例如 "reverse=false;00-00-00-00" 解析数据的时候就不发生反转的操作

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringArray 方法 

[原文連結](http://api.hslcommunication.cn/html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperReadStringArray 方法 |

读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，例如 01-01-00-00 正向有功总需求及发生时间  
Read all the string data information of the specified address, in general, there is only one data for one address, but there are multiple data for a small number of addresses,
such as 01-01-00-00 Forward active total demand and occurrence time

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
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
	IDlt645 dlt,
	string address
)
```

```
Public Shared Function ReadStringArray ( 
	dlt As IDlt645,
	address As String
) As OperateResult(Of String())
```

```
public:
static OperateResult<array<String^>^>^ ReadStringArray(
	IDlt645^ dlt, 
	String^ address
)
```

```
static member ReadStringArray : 
        dlt : IDlt645 * 
        address : string -> OperateResult<string[]> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
字符串数组信息

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;00-00-00-00" 或是 "s=100000;00-00-02-00"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能  
地址也可以携带是否数据翻转的标记，例如 "reverse=false;00-00-00-00" 解析数据的时候就不发生反转的操作

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringArrayAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[AnalysisBytesAddress 方法](../html/ccedd268-3156-f220-6f80-2cead5b8d003.htm "AnalysisBytesAddress 方法 ")

[AnalysisIntegerAddress 方法](../html/cd8d1447-03ec-eeef-d72d-ffec67c65c55.htm "AnalysisIntegerAddress 方法 ")

[BroadcastTime 方法](../html/dd48e0b5-7957-1d0a-442c-768473767de9.htm "BroadcastTime 方法 ")

[BroadcastTimeAsync 方法](../html/0b1242af-65e6-2c4d-9760-9a405d382ca7.htm "BroadcastTimeAsync 方法 ")

[BuildDlt645EntireCommand 方法](../html/2f3b71bb-0138-0674-64b5-81c1f8617812.htm "BuildDlt645EntireCommand 方法 ")

[ChangeBaudRate 方法](../html/cb47bf4a-1aeb-2bb5-4efb-9697bdf1a694.htm "ChangeBaudRate 方法 ")

[ChangeBaudRateAsync 方法](../html/90f1d452-ac4c-23f4-fb70-ad1849840043.htm "ChangeBaudRateAsync 方法 ")

[CheckReceiveDataComplete 方法](../html/992c926d-30fb-9b3c-5d5a-bd51380a63e3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/2587a315-7949-2cae-44ac-f99ec0ceba53.htm "CheckResponse 方法 ")

[CheckResponseCS 方法](../html/d06c60d7-9201-e615-0780-6d4f00a7efbb.htm "CheckResponseCS 方法 ")

[FindHeadCode68H 方法](../html/1444533e-f9d8-ef99-ebab-67f62e99c929.htm "FindHeadCode68H 方法 ")

[FreezeCommand 方法](../html/0d285564-5729-5426-87ce-467f2cd7da85.htm "FreezeCommand 方法 ")

[FreezeCommandAsync 方法](../html/4b7116ed-6aad-b74d-1284-6df9c44e8329.htm "FreezeCommandAsync 方法 ")

[Function1C 方法](../html/0847bf29-8b25-b11f-0ba2-815c53c07597.htm "Function1C 方法 ")

[GetAddressByteFromString 方法](../html/2629bc4f-8a26-336c-0080-3440d4f37c0f.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/a2ee3fb2-11b8-5a1e-ec06-6ff07ce67efc.htm "Read 方法 ")

[ReadAddress 方法](../html/6ad8eea8-df4d-bd00-a827-3fe496cfb1b4.htm "ReadAddress 方法 ")

[ReadAddressAsync 方法](../html/65c718ae-fe54-3f0b-1293-1f64b8d52c01.htm "ReadAddressAsync 方法 ")

[ReadAsync 方法](../html/4835b409-3524-b766-6d33-8df4aa5f6c60.htm "ReadAsync 方法 ")

[ReadDouble 方法](../html/ffe25163-129a-0a2f-b55e-e2d2e7193ddf.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/0b933f21-d4c3-790a-8748-f4504694afc0.htm "ReadDoubleAsync 方法 ")

[ReadStringArray 方法](../html/c4c395b8-e5c0-e974-c611-56ac6d76442f.htm "ReadStringArray 方法 ")

[ReadStringArrayAsync 方法](../html/cb9e6c9b-5039-972e-b855-a7999eaf8f26.htm "ReadStringArrayAsync 方法 ")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[WriteAddress 方法](../html/9c49030b-7831-e97f-a1b3-e9b4a8e2b498.htm "WriteAddress 方法 ")

[WriteAddressAsync 方法](../html/a499dd5a-2edc-296f-2513-fc67d8fdd5a1.htm "WriteAddressAsync 方法 ")

[WriteAsync 方法](../html/cf1ff974-94ce-89d6-27f0-b2684762cbdd.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperReadStringArrayAsync 方法 |

读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，例如 01-01-00-00 正向有功总需求及发生时间  
Read all the string data information of the specified address, in general, there is only one data for one address, but there are multiple data for a small number of addresses,
such as 01-01-00-00 Forward active total demand and occurrence time

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<string[]>> ReadStringArrayAsync(
	IDlt645 dlt,
	string address
)
```

```
Public Shared Function ReadStringArrayAsync ( 
	dlt As IDlt645,
	address As String
) As Task(Of OperateResult(Of String()))
```

```
public:
static Task<OperateResult<array<String^>^>^>^ ReadStringArrayAsync(
	IDlt645^ dlt, 
	String^ address
)
```

```
static member ReadStringArrayAsync : 
        dlt : IDlt645 * 
        address : string -> Task<OperateResult<string[]>> 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
字符串数组信息

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;00-00-00-00" 或是 "s=100000;00-00-02-00"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能  
地址也可以携带是否数据翻转的标记，例如 "reverse=false;00-00-00-00" 解析数据的时候就不发生反转的操作

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/30670cec-c188-6c46-5d20-71a1605cdd33.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[Write 方法 (IDlt645, String, String, String, Byte[])](../html/7ff3b579-6723-2ef5-b6ac-b646cb493846.htm "Write 方法 (IDlt645, String, String, String, Byte[])")

[Write 方法 (IDlt645, String, String, String, String[])](../html/72105f34-878f-1335-1c37-cf2f1d888a77.htm "Write 方法 (IDlt645, String, String, String, String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [Write(IDlt645, String, String, String, Byte)](7ff3b579-6723-2ef5-b6ac-b646cb493846.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. |
| 公共方法静态成员 | [Write(IDlt645, String, String, String, String)](72105f34-878f-1335-1c37-cf2f1d888a77.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IDlt645, String, String, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/7ff3b579-6723-2ef5-b6ac-b646cb493846.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[Write 方法 (IDlt645, String, String, String, Byte[])](../html/7ff3b579-6723-2ef5-b6ac-b646cb493846.htm "Write 方法 (IDlt645, String, String, String, Byte[])")

[Write 方法 (IDlt645, String, String, String, String[])](../html/72105f34-878f-1335-1c37-cf2f1d888a77.htm "Write 方法 (IDlt645, String, String, String, String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperWrite 方法 (IDlt645, String, String, String, Byte) |

根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  
Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual,
from high to position, such as 00-00-00-00. The separator can be any special character or no separator.

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
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
	IDlt645 dlt,
	string password,
	string opCode,
	string address,
	byte[] value
)
```

```
Public Shared Function Write ( 
	dlt As IDlt645,
	password As String,
	opCode As String,
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IDlt645^ dlt, 
	String^ password, 
	String^ opCode, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member Write : 
        dlt : IDlt645 * 
        password : string * 
        opCode : string * 
        address : string * 
        value : byte[] -> OperateResult 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

password
:   类型：SystemString  
    密钥信息

opCode
:   类型：SystemString  
    操作者代码

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

地址可以携带地址域信息，例如 "s=2;00-00-00-00" 或是 "s=100000;00-00-02-00"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能  
注意：本命令必须与编程键配合使用

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[Write 重载](30670cec-c188-6c46-5d20-71a1605cdd33.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IDlt645, String, String, String, String[])

[原文連結](http://api.hslcommunication.cn/html/72105f34-878f-1335-1c37-cf2f1d888a77.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT.Helper](../html/655bc9a8-c0d2-de47-e43c-8a07038a8377.htm "HslCommunication.Instrument.DLT.Helper")

[DLT645Helper 类](../html/7200043e-ae86-4dd9-df01-8a0a80832ee9.htm "DLT645Helper 类")

[DLT645Helper 方法](../html/ca652825-862f-2923-5313-d1b5b7340723.htm "DLT645Helper 方法")

[Write 方法](../html/30670cec-c188-6c46-5d20-71a1605cdd33.htm "Write 方法 ")

[Write 方法 (IDlt645, String, String, String, Byte[])](../html/7ff3b579-6723-2ef5-b6ac-b646cb493846.htm "Write 方法 (IDlt645, String, String, String, Byte[])")

[Write 方法 (IDlt645, String, String, String, String[])](../html/72105f34-878f-1335-1c37-cf2f1d888a77.htm "Write 方法 (IDlt645, String, String, String, String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645HelperWrite 方法 (IDlt645, String, String, String, String) |

将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  
Write the data to the gauge, address identification according to the manual, from high bit to position,
such as 00-00-00-00, the separator can be any special character or no delimiter.

**命名空间：**
 [HslCommunication.Instrument.DLT.Helper](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)  
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
	IDlt645 dlt,
	string password,
	string opCode,
	string address,
	string[] value
)
```

```
Public Shared Function Write ( 
	dlt As IDlt645,
	password As String,
	opCode As String,
	address As String,
	value As String()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IDlt645^ dlt, 
	String^ password, 
	String^ opCode, 
	String^ address, 
	array<String^>^ value
)
```

```
static member Write : 
        dlt : IDlt645 * 
        password : string * 
        opCode : string * 
        address : string * 
        value : string[] -> OperateResult 
```

#### 参数

dlt
:   类型：[HslCommunication.Instrument.DLT.HelperIDlt645](3aaea7ea-ccf9-0290-7ec5-641c25729089.htm)  
    DLT通信对象

password
:   类型：SystemString  
    密钥信息

opCode
:   类型：SystemString  
    操作者代码

address
:   类型：SystemString  
    地址信息

value
:   类型：SystemString  
    写入的数据值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645Helper 类](7200043e-ae86-4dd9-df01-8a0a80832ee9.htm)

[Write 重载](30670cec-c188-6c46-5d20-71a1605cdd33.htm)

[HslCommunication.Instrument.DLT.Helper 命名空间](655bc9a8-c0d2-de47-e43c-8a07038a8377.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)