# HslCommunication - HslCommunication.Profinet.Melsec

> 分類頁數: 30



---
## HslCommunication.Profinet.Melsec

[原文連結](http://api.hslcommunication.cn/html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EDataType 类](../html/3f8dc77e-edc4-5d34-bd94-bcde56fd0ff8.htm "MelsecA1EDataType 类")

[MelsecA1ENet 类](../html/b0592620-d255-80a0-2b83-fe4e5101adc9.htm "MelsecA1ENet 类")

[MelsecA1EServer 类](../html/38b26ba3-c5b4-79c8-699f-42fa801614b6.htm "MelsecA1EServer 类")

[MelsecA3CNet 类](../html/143c103b-bf33-3faa-e553-f55ea36a4f99.htm "MelsecA3CNet 类")

[MelsecA3CNetOverTcp 类](../html/d3e542cc-7f74-6bfb-0ce7-608d5093d38a.htm "MelsecA3CNetOverTcp 类")

[MelsecA3CServer 类](../html/f0c4c7fa-4b0d-61e7-8575-41538da313f3.htm "MelsecA3CServer 类")

[MelsecCipNet 类](../html/3015316e-1719-7110-1ba8-c171020fc912.htm "MelsecCipNet 类")

[MelsecFxLinks 类](../html/f4a31e75-64a3-d21d-9f35-be4c8c517301.htm "MelsecFxLinks 类")

[MelsecFxLinksOverTcp 类](../html/eb348247-7bc2-806d-8917-b28e085141a8.htm "MelsecFxLinksOverTcp 类")

[MelsecFxLinksServer 类](../html/73bacead-c74d-ce11-1eff-3188cb31afc6.htm "MelsecFxLinksServer 类")

[MelsecFxSerial 类](../html/c9af0583-6b00-0e21-96c8-c3c306641769.htm "MelsecFxSerial 类")

[MelsecFxSerialOverTcp 类](../html/8ffc6a4e-c343-f48f-c788-2eff8f375696.htm "MelsecFxSerialOverTcp 类")

[MelsecFxSerialServer 类](../html/7140d0dd-5802-d6c2-ba5e-9ad9fefcdbcd.htm "MelsecFxSerialServer 类")

[MelsecHelper 类](../html/937b0c17-a9e6-6aa5-d8e7-f27fce0ab086.htm "MelsecHelper 类")

[MelsecMcAsciiNet 类](../html/f7e7593f-3923-662a-b2d7-10d0b5f86d2e.htm "MelsecMcAsciiNet 类")

[MelsecMcAsciiUdp 类](../html/cb0eefa9-b894-80f7-2f94-588f9fd1b718.htm "MelsecMcAsciiUdp 类")

[MelsecMcDataType 类](../html/55a6718a-d454-fbdd-7f76-099538dcffc1.htm "MelsecMcDataType 类")

[MelsecMcNet 类](../html/0767d258-53d6-f038-ef35-477c8877d789.htm "MelsecMcNet 类")

[MelsecMcRNet 类](../html/3b8959b9-807f-7c05-6782-abebd8ec6407.htm "MelsecMcRNet 类")

[MelsecMcServer 类](../html/abea5116-fde0-a1f9-dee3-19f99ae23e25.htm "MelsecMcServer 类")

[MelsecMcUdp 类](../html/69a6438b-6de9-82e4-a610-947d4f855011.htm "MelsecMcUdp 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Melsec 命名空间 |

在三菱的PLC通信的MC协议中，分为串行通信的报文和以太网接口的报文。  
在串口通信中，共有以下几种帧，其中1C,2C,3C帧支持格式1，2，3，4，在C帧里支持格式5通信  

1. 4C帧，QnA系列串行通信模块专用协议（Qna扩展帧）
2. 3C帧，QnA系列串行通信模块专用协议（Qna帧）
3. 2C帧，QnA系列串行通信模块专用协议（Qna简易帧）
4. 1C帧，A系列计算机链接模块专用协议

在以太网通信中，共有以下几种帧，每种帧支持二进制和ASCII格式

1. 4E帧，是3E帧上附加了“序列号”。
2. 3E帧，QnA系列以太网接口模块的报文格式，兼容SLMP的报文格式
3. 1E帧，A系列以太网接口模块的报文格式

在以太网通信里，HSL主要针对1E帧协议和3E帧协议进行实现，大概说一下怎么选择通信类对象，对于三菱PLC而言，需要事先在PLC侧的网络配置中进行
相关的配置操作，具体是配置二进制格式还是ASCII格式，然后配置端口，配置TCP还是UDP协议。  
[MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)，[MelsecMcAsciiNet](f7e7593f-3923-662a-b2d7-10d0b5f86d2e.htm)，[MelsecMcUdp](69a6438b-6de9-82e4-a610-947d4f855011.htm), [MelsecMcAsciiUdp](cb0eefa9-b894-80f7-2f94-588f9fd1b718.htm) 这四个类都是MC协议的Qna兼容3E帧实现，分别
是TCP二进制，TCP的ASCII，UDP的二进制，UDP的ASCI格式。适用Q系列，L系列，FX5U系列，还有以太网模块QJ71E71。  
[MelsecA1ENet](b0592620-d255-80a0-2b83-fe4e5101adc9.htm), [MelsecA1EAsciiNet](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm) 这两个类是MC协议的Qna兼容1E协议实现，
分别是二进制和ASCII格式的实现，主要适用A系列的PLC，Fx3u，以及有些老的PLC，使用了北辰模块实现的通信也是选择的是 A1E 协议来通信。  
[MelsecA3CNet](143c103b-bf33-3faa-e553-f55ea36a4f99.htm) 是MC协议的3C帧的实现，主要支持串行通信的模块的实现。  
[MelsecFxSerial](c9af0583-6b00-0e21-96c8-c3c306641769.htm) 是FX编程口的协议的实现，测试并不在所有的老型号上都很稳定。具体支持的系列需要参照类的说明，
当使用GOT触摸屏连接PLC时，需要使用类[MelsecFxSerialOverTcp](8ffc6a4e-c343-f48f-c788-2eff8f375696.htm)并将属性[UseGOT](9b45af8d-5e39-fc3c-c2c7-350ee43b80b5.htm)设置为 True。  
[MelsecFxLinks](f4a31e75-64a3-d21d-9f35-be4c8c517301.htm) 是三菱的计算机链接协议，通常是三菱的485接口，而不是编程口，该协议支持多种格式及是否和校验，实际参数选择需要根据PLC的配置来决定。  
  
三菱Q系列PLC带以太网模块，使用MC 二进制协议通讯，在网络中断情况下无法正常连接的情况，解决方案如下：  
1. 生存确认选择确认模式；  
2. 初始设置中将"对象目标生存确认开始间隔定时器"从1200改为12

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [MelsecA1EAsciiNet](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm) | 三菱PLC通讯协议，采用A兼容1E帧协议实现，使用ASCII码通讯，请根据实际型号来进行选取  Mitsubishi PLC communication protocol, implemented using A compatible 1E frame protocol, using ascii code communication, please choose according to the actual model |
| 公共类 | [MelsecA1EDataType](3f8dc77e-edc4-5d34-bd94-bcde56fd0ff8.htm) | 三菱PLC的数据类型，此处包含了几个常用的类型 |
| 公共类 | [MelsecA1ENet](b0592620-d255-80a0-2b83-fe4e5101adc9.htm) | 三菱PLC通讯协议，采用A兼容1E帧协议实现，使用二进制码通讯，请根据实际型号来进行选取  Mitsubishi PLC communication protocol, implemented using A compatible 1E frame protocol, using binary code communication, please choose according to the actual model |
| 公共类 | [MelsecA1EServer](38b26ba3-c5b4-79c8-699f-42fa801614b6.htm) | 三菱MC-A1E协议的虚拟服务器，支持M,X,Y,D,W的数据池读写操作，支持二进制及ASCII格式进行读写操作，需要在实例化的时候指定。  The Mitsubishi MC-A1E protocol virtual server supports M, X, Y, D, W data pool read and write operations, and supports binary and ASCII format read and write operations, which need to be specified during instantiation. |
| 公共类代码示例 | [MelsecA3CNet](143c103b-bf33-3faa-e553-f55ea36a4f99.htm) | 基于Qna 兼容3C帧的格式一的通讯，具体的地址需要参照三菱的基本地址  Based on Qna-compatible 3C frame format one communication, the specific address needs to refer to the basic address of Mitsubishi. |
| 公共类 | [MelsecA3CNetOverTcp](d3e542cc-7f74-6bfb-0ce7-608d5093d38a.htm) | 基于Qna 兼容3C帧的格式一的通讯，具体的地址需要参照三菱的基本地址，本类是基于tcp通讯的实现  Based on Qna-compatible 3C frame format one communication, the specific address needs to refer to the basic address of Mitsubishi. This class is based on TCP communication. |
| 公共类 | [MelsecA3CServer](f0c4c7fa-4b0d-61e7-8575-41538da313f3.htm) | 基于MC协议的A3C格式虚拟服务器，可以模拟A3C格式的PLC，支持格式1，2，3，4，具体可以使用[Format](e41b8234-2fc6-af9e-51df-58a25f23af36.htm)来设置， 支持设置是否校验，同时支持网口或是串口的访问  A3C format virtual server based on MC protocol can simulate A3C format PLC, support format 1, 2, 3, 4, specifically you can use [Format](e41b8234-2fc6-af9e-51df-58a25f23af36.htm) to set, support whether to verify the setting, and also support network Port or serial port access |
| 公共类 | [MelsecCipNet](3015316e-1719-7110-1ba8-c171020fc912.htm) | 三菱PLC的EIP协议的实现，当PLC使用了 QJ71EIP71 模块时就需要使用本类来访问 |
| 公共类 | [MelsecFxLinks](f4a31e75-64a3-d21d-9f35-be4c8c517301.htm) | 三菱计算机链接协议，适用FX3U系列，FX3G，FX3S等等系列，通常在PLC侧连接的是485的接线口  Mitsubishi Computer Link Protocol, suitable for FX3U series, FX3G, FX3S, etc., usually the 485 connection port is connected on the PLC side |
| 公共类 | [MelsecFxLinksOverTcp](eb348247-7bc2-806d-8917-b28e085141a8.htm) | 三菱计算机链接协议的网口版本，适用FX3U系列，FX3G，FX3S等等系列，通常在PLC侧连接的是485的接线口  Network port version of Mitsubishi Computer Link Protocol, suitable for FX3U series, FX3G, FX3S, etc., usually the 485 connection port is connected on the PLC side |
| 公共类 | [MelsecFxLinksServer](73bacead-c74d-ce11-1eff-3188cb31afc6.htm) | 三菱的虚拟的FxLinks服务器 |
| 公共类代码示例 | [MelsecFxSerial](c9af0583-6b00-0e21-96c8-c3c306641769.htm) | 三菱的串口通信的对象，适用于读取FX系列的串口数据，支持的类型参考文档说明  Mitsubishi's serial communication object is suitable for reading serial data of the FX series. Refer to the documentation for the supported types. |
| 公共类代码示例 | [MelsecFxSerialOverTcp](8ffc6a4e-c343-f48f-c788-2eff8f375696.htm) | 三菱串口协议的网络版，如果使用的是 FX3U编程口(fx2n) -> GOT1000(RS232)(或是GOT2000) -> 上位机(以太网) 的方式，那么就需要把[UseGOT](9b45af8d-5e39-fc3c-c2c7-350ee43b80b5.htm)设置为 True  The network version of the Mitsubishi serial port protocol, if you use the FX3U programming port (fx2n) -> GOT1000 (RS232) (or GOT2000) -> host computer (Ethernet) method, then you need to put [UseGOT](9b45af8d-5e39-fc3c-c2c7-350ee43b80b5.htm) is set to True |
| 公共类 | [MelsecFxSerialServer](7140d0dd-5802-d6c2-ba5e-9ad9fefcdbcd.htm) | 三菱编程口协议的虚拟服务器，属性[IsNewVersion](401da5fb-7ddf-5ffe-8eae-f1b4269dc83e.htm)的配置和客户端保持一致才能数据一致  Mitsubishi Programming interface protocol virtual server, attribute [IsNewVersion](401da5fb-7ddf-5ffe-8eae-f1b4269dc83e.htm) and the client configuration must be consistent data consistency |
| 公共类 | [MelsecHelper](937b0c17-a9e6-6aa5-d8e7-f27fce0ab086.htm) | 所有三菱通讯类的通用辅助工具类，包含了一些通用的静态方法，可以使用本类来获取一些原始的报文信息。详细的操作参见例子  All general auxiliary tool classes of Mitsubishi communication class include some general static methods. You can use this class to get some primitive message information. See the example for detailed operation |
| 公共类代码示例 | [MelsecMcAsciiNet](f7e7593f-3923-662a-b2d7-10d0b5f86d2e.htm) | 三菱PLC通讯类，采用Qna兼容3E帧协议实现，需要在PLC侧先的以太网模块先进行配置，必须为ASCII通讯格式  Mitsubishi PLC communication class is implemented using Qna compatible 3E frame protocol. The Ethernet module on the PLC side needs to be configured first. It must be ascii communication. |
| 公共类代码示例 | [MelsecMcAsciiUdp](cb0eefa9-b894-80f7-2f94-588f9fd1b718.htm) | 三菱PLC通讯类，采用UDP的协议实现，采用Qna兼容3E帧协议实现，需要在PLC侧先的以太网模块先进行配置，必须为ascii通讯  Mitsubishi PLC communication class is implemented using UDP protocol and Qna compatible 3E frame protocol. The Ethernet module needs to be configured first on the PLC side, and it must be ascii communication. |
| 公共类 | [MelsecMcDataType](55a6718a-d454-fbdd-7f76-099538dcffc1.htm) | 三菱PLC的数据类型，此处包含了几个常用的类型  Data types of Mitsubishi PLC, here contains several commonly used types |
| 公共类代码示例 | [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm) | 三菱PLC通讯类，采用Qna兼容3E帧协议实现，需要在PLC侧先的以太网模块先进行配置，必须为二进制通讯  Mitsubishi PLC communication class is implemented using Qna compatible 3E frame protocol. The Ethernet module on the PLC side needs to be configured first. It must be binary communication. |
| 公共类 | [MelsecMcRNet](3b8959b9-807f-7c05-6782-abebd8ec6407.htm) | 三菱的R系列的MC协议，支持的地址类型和 [MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm) 有区别，详细请查看对应的API文档说明 |
| 公共类 | [MelsecMcServer](abea5116-fde0-a1f9-dee3-19f99ae23e25.htm) | 三菱MC协议的虚拟服务器，支持M,X,Y,D,W的数据池读写操作，支持二进制及ASCII格式进行读写操作，需要在实例化的时候指定。  The Mitsubishi MC protocol virtual server supports M, X, Y, D, W data pool read and write operations, and supports binary and ASCII format read and write operations, which need to be specified during instantiation. |
| 公共类代码示例 | [MelsecMcUdp](69a6438b-6de9-82e4-a610-947d4f855011.htm) | 三菱PLC通讯类，采用UDP的协议实现，采用Qna兼容3E帧协议实现，需要在PLC侧先的以太网模块先进行配置，必须为二进制通讯  Mitsubishi PLC communication class is implemented using UDP protocol and Qna compatible 3E frame protocol. The Ethernet module needs to be configured first on the PLC side, and it must be binary communication. |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MelsecA1EAsciiNet 类

[原文連結](http://api.hslcommunication.cn/html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 构造函数](../html/d531f4a7-510e-2e2a-c8f8-44188bb72a9e.htm "MelsecA1EAsciiNet 构造函数 ")

[MelsecA1EAsciiNet 属性](../html/f50f623e-d269-67db-4417-de0077b9583d.htm "MelsecA1EAsciiNet 属性")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[MelsecA1EAsciiNet 字段](../html/1e1755d0-19bc-dadb-ef41-c1365e3e2a9d.htm "MelsecA1EAsciiNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNet 类 |

三菱PLC通讯协议，采用A兼容1E帧协议实现，使用ASCII码通讯，请根据实际型号来进行选取  
Mitsubishi PLC communication protocol, implemented using A compatible 1E frame protocol, using ascii code communication, please choose according to the actual model

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
        HslCommunication.Profinet.MelsecMelsecA1EAsciiNet

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class MelsecA1EAsciiNet : DeviceTcpNet
```

```
Public Class MelsecA1EAsciiNet
	Inherits DeviceTcpNet
```

```
public ref class MelsecA1EAsciiNet : public DeviceTcpNet
```

```
type MelsecA1EAsciiNet =  
    class
        inherit DeviceTcpNet
    end
```

MelsecA1EAsciiNet 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [MelsecA1EAsciiNet](c7fbc8f4-9ab8-ea33-bcba-a5277b0a050a.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [MelsecA1EAsciiNet(String, Int32)](3fbf6c6a-edc1-ca78-1eb8-246f71069e4a.htm) | 指定ip地址和端口来来实例化一个默认的对象  Specify the IP address and port to instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PLCNumber](2c54ef5d-6f0a-2cdc-6015-0f548fc49994.htm) | PLC编号，默认为0xFF  PLC number, default is 0xFF |
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
| 公共方法静态成员 | [BuildReadCommand](5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm) | 根据类型地址长度确认需要读取的指令头 |
| 公共方法静态成员 | [BuildWriteBoolCommand](e8f22d3b-2114-f744-daba-f087099077b7.htm) | 根据类型地址以及需要写入的数据来生成指令头 |
| 公共方法静态成员 | [BuildWriteWordCommand](6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm) | 根据类型地址以及需要写入的数据来生成指令头 |
| 公共方法静态成员 | [CheckResponseLegal](19df45fe-ecef-ad2d-2ce7-d2927673383d.htm) | 检测反馈的消息是否合法 |
| 公共方法代码示例 | [ConnectClose](569a9c0e-7628-b8d5-e124-0853452747b8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](76be8150-ff22-95b9-16b1-7c3440c54ae8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServer](6299cc1f-2fab-b754-4597-1eddc2b21308.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](e7eff5c0-8e2a-3b6c-5027-16f4d8bcf5ae.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [ExtractActualData](2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm) | 从PLC反馈的数据中提取出实际的数据内容，需要传入反馈数据，是否位读取 |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](2a7d440d-2475-926e-bf5f-93b73b188939.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
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
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](55e3ae9c-4e04-2673-edb6-3d7e0e6add5d.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](ffd358d4-af01-dd85-da93-141616d87c61.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm).) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](3c9dee5d-08a3-2882-f874-19a70acff835.htm) | 批量读取Boolean数组信息，需要指定地址和长度，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch read Boolean array information, need to specify the address and length, return Boolean array. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](0df9b5c8-ae62-00d1-4f14-0a660772ce98.htm) | 批量读取Boolean数组信息，需要指定地址和长度，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch read Boolean array information, need to specify the address and length, return Boolean array. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm).) |
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
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](fa489c4e-7431-a042-a7dc-2d79584de577.htm) | (重写 [DeviceTcpNetToString](209a196b-90ea-2b73-e915-ce4b11f1263d.htm).) |
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
| 公共方法 | [Write(String, Boolean)](15d948d6-ab6f-9d25-6ef0-3d1a975de405.htm) | 批量写入Boolean数组数据，返回是否成功，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch write Boolean array data, return whether the write was successful. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](49c72530-9bee-f030-99ec-e73774b17de1.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
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
| 公共方法 | [WriteAsync(String, Boolean)](8703ce64-be44-f53a-4666-da5f104cf3c2.htm) | 批量写入Boolean数组数据，返回是否成功，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch write Boolean array data, return whether the write was successful. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationWriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm).) |
| 公共方法 | [WriteAsync(String, Byte)](60031b6f-fb61-9b88-045b-a201ab534826.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm).) |
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

本类适用于的PLC列表

1. FX3U(C) PLC 测试人sandy\_liao

| 重要事项 **重要事项** |
| --- |
| 本通讯类由CKernal推送，感谢 |

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MelsecA1EAsciiNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/d531f4a7-510e-2e2a-c8f8-44188bb72a9e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 构造函数](../html/d531f4a7-510e-2e2a-c8f8-44188bb72a9e.htm "MelsecA1EAsciiNet 构造函数 ")

[MelsecA1EAsciiNet 构造函数](../html/c7fbc8f4-9ab8-ea33-bcba-a5277b0a050a.htm "MelsecA1EAsciiNet 构造函数 ")

[MelsecA1EAsciiNet 构造函数 (String, Int32)](../html/3fbf6c6a-edc1-ca78-1eb8-246f71069e4a.htm "MelsecA1EAsciiNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNet 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [MelsecA1EAsciiNet](c7fbc8f4-9ab8-ea33-bcba-a5277b0a050a.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [MelsecA1EAsciiNet(String, Int32)](3fbf6c6a-edc1-ca78-1eb8-246f71069e4a.htm) | 指定ip地址和端口来来实例化一个默认的对象  Specify the IP address and port to instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MelsecA1EAsciiNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c7fbc8f4-9ab8-ea33-bcba-a5277b0a050a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 构造函数](../html/d531f4a7-510e-2e2a-c8f8-44188bb72a9e.htm "MelsecA1EAsciiNet 构造函数 ")

[MelsecA1EAsciiNet 构造函数](../html/c7fbc8f4-9ab8-ea33-bcba-a5277b0a050a.htm "MelsecA1EAsciiNet 构造函数 ")

[MelsecA1EAsciiNet 构造函数 (String, Int32)](../html/3fbf6c6a-edc1-ca78-1eb8-246f71069e4a.htm "MelsecA1EAsciiNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNet 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public MelsecA1EAsciiNet()
```

```
Public Sub New
```

```
public:
MelsecA1EAsciiNet()
```

```
new : unit -> MelsecA1EAsciiNet
```

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[MelsecA1EAsciiNet 重载](d531f4a7-510e-2e2a-c8f8-44188bb72a9e.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MelsecA1EAsciiNet 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/3fbf6c6a-edc1-ca78-1eb8-246f71069e4a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 构造函数](../html/d531f4a7-510e-2e2a-c8f8-44188bb72a9e.htm "MelsecA1EAsciiNet 构造函数 ")

[MelsecA1EAsciiNet 构造函数](../html/c7fbc8f4-9ab8-ea33-bcba-a5277b0a050a.htm "MelsecA1EAsciiNet 构造函数 ")

[MelsecA1EAsciiNet 构造函数 (String, Int32)](../html/3fbf6c6a-edc1-ca78-1eb8-246f71069e4a.htm "MelsecA1EAsciiNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNet 构造函数 (String, Int32) |

指定ip地址和端口来来实例化一个默认的对象  
Specify the IP address and port to instantiate a default object

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public MelsecA1EAsciiNet(
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
MelsecA1EAsciiNet(
	String^ ipAddress, 
	int port
)
```

```
new : 
        ipAddress : string * 
        port : int -> MelsecA1EAsciiNet
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

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[MelsecA1EAsciiNet 重载](d531f4a7-510e-2e2a-c8f8-44188bb72a9e.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MelsecA1EAsciiNet 属性

[原文連結](http://api.hslcommunication.cn/html/f50f623e-d269-67db-4417-de0077b9583d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 属性](../html/f50f623e-d269-67db-4417-de0077b9583d.htm "MelsecA1EAsciiNet 属性")

[PLCNumber 属性](../html/2c54ef5d-6f0a-2cdc-6015-0f548fc49994.htm "PLCNumber 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNet 属性 |

[MelsecA1EAsciiNet](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PLCNumber](2c54ef5d-6f0a-2cdc-6015-0f548fc49994.htm) | PLC编号，默认为0xFF  PLC number, default is 0xFF |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PLCNumber 属性 

[原文連結](http://api.hslcommunication.cn/html/2c54ef5d-6f0a-2cdc-6015-0f548fc49994.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 属性](../html/f50f623e-d269-67db-4417-de0077b9583d.htm "MelsecA1EAsciiNet 属性")

[PLCNumber 属性](../html/2c54ef5d-6f0a-2cdc-6015-0f548fc49994.htm "PLCNumber 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetPLCNumber 属性 |

PLC编号，默认为0xFF  
PLC number, default is 0xFF

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte PLCNumber { get; set; }
```

```
Public Property PLCNumber As Byte
	Get
	Set
```

```
public:
property unsigned char PLCNumber {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member PLCNumber : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MelsecA1EAsciiNet 方法

[原文連結](http://api.hslcommunication.cn/html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[BuildReadCommand 方法](../html/5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/e8f22d3b-2114-f744-daba-f087099077b7.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm "BuildWriteWordCommand 方法 ")

[CheckResponseLegal 方法](../html/19df45fe-ecef-ad2d-2ce7-d2927673383d.htm "CheckResponseLegal 方法 ")

[ExtractActualData 方法](../html/2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm "ExtractActualData 方法 ")

[GetNewNetMessage 方法](../html/2a7d440d-2475-926e-bf5f-93b73b188939.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ToString 方法](../html/fa489c4e-7431-a042-a7dc-2d79584de577.htm "ToString 方法 ")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNet 方法 |

[MelsecA1EAsciiNet](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm) | 根据类型地址长度确认需要读取的指令头 |
| 公共方法静态成员 | [BuildWriteBoolCommand](e8f22d3b-2114-f744-daba-f087099077b7.htm) | 根据类型地址以及需要写入的数据来生成指令头 |
| 公共方法静态成员 | [BuildWriteWordCommand](6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm) | 根据类型地址以及需要写入的数据来生成指令头 |
| 公共方法静态成员 | [CheckResponseLegal](19df45fe-ecef-ad2d-2ce7-d2927673383d.htm) | 检测反馈的消息是否合法 |
| 公共方法代码示例 | [ConnectClose](569a9c0e-7628-b8d5-e124-0853452747b8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](76be8150-ff22-95b9-16b1-7c3440c54ae8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServer](6299cc1f-2fab-b754-4597-1eddc2b21308.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](e7eff5c0-8e2a-3b6c-5027-16f4d8bcf5ae.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [ExtractActualData](2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm) | 从PLC反馈的数据中提取出实际的数据内容，需要传入反馈数据，是否位读取 |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](2a7d440d-2475-926e-bf5f-93b73b188939.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
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
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](55e3ae9c-4e04-2673-edb6-3d7e0e6add5d.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](ffd358d4-af01-dd85-da93-141616d87c61.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm).) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](3c9dee5d-08a3-2882-f874-19a70acff835.htm) | 批量读取Boolean数组信息，需要指定地址和长度，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch read Boolean array information, need to specify the address and length, return Boolean array. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](0df9b5c8-ae62-00d1-4f14-0a660772ce98.htm) | 批量读取Boolean数组信息，需要指定地址和长度，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch read Boolean array information, need to specify the address and length, return Boolean array. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm).) |
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
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](fa489c4e-7431-a042-a7dc-2d79584de577.htm) | (重写 [DeviceTcpNetToString](209a196b-90ea-2b73-e915-ce4b11f1263d.htm).) |
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
| 公共方法 | [Write(String, Boolean)](15d948d6-ab6f-9d25-6ef0-3d1a975de405.htm) | 批量写入Boolean数组数据，返回是否成功，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch write Boolean array data, return whether the write was successful. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](49c72530-9bee-f030-99ec-e73774b17de1.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
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
| 公共方法 | [WriteAsync(String, Boolean)](8703ce64-be44-f53a-4666-da5f104cf3c2.htm) | 批量写入Boolean数组数据，返回是否成功，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch write Boolean array data, return whether the write was successful. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationWriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm).) |
| 公共方法 | [WriteAsync(String, Byte)](60031b6f-fb61-9b88-045b-a201ab534826.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm).) |
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

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[BuildReadCommand 方法](../html/5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/e8f22d3b-2114-f744-daba-f087099077b7.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm "BuildWriteWordCommand 方法 ")

[CheckResponseLegal 方法](../html/19df45fe-ecef-ad2d-2ce7-d2927673383d.htm "CheckResponseLegal 方法 ")

[ExtractActualData 方法](../html/2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm "ExtractActualData 方法 ")

[GetNewNetMessage 方法](../html/2a7d440d-2475-926e-bf5f-93b73b188939.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ToString 方法](../html/fa489c4e-7431-a042-a7dc-2d79584de577.htm "ToString 方法 ")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetBuildReadCommand 方法 |

根据类型地址长度确认需要读取的指令头

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
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
	string address,
	ushort length,
	bool isBit,
	byte plcNumber
)
```

```
Public Shared Function BuildReadCommand ( 
	address As String,
	length As UShort,
	isBit As Boolean,
	plcNumber As Byte
) As OperateResult(Of List(Of Byte()))
```

```
public:
static OperateResult<List<array<unsigned char>^>^>^ BuildReadCommand(
	String^ address, 
	unsigned short length, 
	bool isBit, 
	unsigned char plcNumber
)
```

```
static member BuildReadCommand : 
        address : string * 
        length : uint16 * 
        isBit : bool * 
        plcNumber : byte -> OperateResult<List<byte[]>> 
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    长度

isBit
:   类型：SystemBoolean  
    指示是否按照位成批的读出

plcNumber
:   类型：SystemByte  
    PLC编号

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)ListByte  
带有成功标志的指令数据

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteBoolCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/e8f22d3b-2114-f744-daba-f087099077b7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[BuildReadCommand 方法](../html/5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/e8f22d3b-2114-f744-daba-f087099077b7.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm "BuildWriteWordCommand 方法 ")

[CheckResponseLegal 方法](../html/19df45fe-ecef-ad2d-2ce7-d2927673383d.htm "CheckResponseLegal 方法 ")

[ExtractActualData 方法](../html/2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm "ExtractActualData 方法 ")

[GetNewNetMessage 方法](../html/2a7d440d-2475-926e-bf5f-93b73b188939.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ToString 方法](../html/fa489c4e-7431-a042-a7dc-2d79584de577.htm "ToString 方法 ")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetBuildWriteBoolCommand 方法 |

根据类型地址以及需要写入的数据来生成指令头

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteBoolCommand(
	string address,
	bool[] value,
	byte plcNumber
)
```

```
Public Shared Function BuildWriteBoolCommand ( 
	address As String,
	value As Boolean(),
	plcNumber As Byte
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteBoolCommand(
	String^ address, 
	array<bool>^ value, 
	unsigned char plcNumber
)
```

```
static member BuildWriteBoolCommand : 
        address : string * 
        value : bool[] * 
        plcNumber : byte -> OperateResult<byte[]> 
```

#### 参数

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemBoolean  
    数据值

plcNumber
:   类型：SystemByte  
    PLC编号

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标志的指令数据

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteWordCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[BuildReadCommand 方法](../html/5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/e8f22d3b-2114-f744-daba-f087099077b7.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm "BuildWriteWordCommand 方法 ")

[CheckResponseLegal 方法](../html/19df45fe-ecef-ad2d-2ce7-d2927673383d.htm "CheckResponseLegal 方法 ")

[ExtractActualData 方法](../html/2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm "ExtractActualData 方法 ")

[GetNewNetMessage 方法](../html/2a7d440d-2475-926e-bf5f-93b73b188939.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ToString 方法](../html/fa489c4e-7431-a042-a7dc-2d79584de577.htm "ToString 方法 ")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetBuildWriteWordCommand 方法 |

根据类型地址以及需要写入的数据来生成指令头

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteWordCommand(
	string address,
	byte[] value,
	byte plcNumber
)
```

```
Public Shared Function BuildWriteWordCommand ( 
	address As String,
	value As Byte(),
	plcNumber As Byte
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteWordCommand(
	String^ address, 
	array<unsigned char>^ value, 
	unsigned char plcNumber
)
```

```
static member BuildWriteWordCommand : 
        address : string * 
        value : byte[] * 
        plcNumber : byte -> OperateResult<byte[]> 
```

#### 参数

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemByte  
    数据值

plcNumber
:   类型：SystemByte  
    PLC编号

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标志的指令数据

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckResponseLegal 方法 

[原文連結](http://api.hslcommunication.cn/html/19df45fe-ecef-ad2d-2ce7-d2927673383d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[BuildReadCommand 方法](../html/5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/e8f22d3b-2114-f744-daba-f087099077b7.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm "BuildWriteWordCommand 方法 ")

[CheckResponseLegal 方法](../html/19df45fe-ecef-ad2d-2ce7-d2927673383d.htm "CheckResponseLegal 方法 ")

[ExtractActualData 方法](../html/2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm "ExtractActualData 方法 ")

[GetNewNetMessage 方法](../html/2a7d440d-2475-926e-bf5f-93b73b188939.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ToString 方法](../html/fa489c4e-7431-a042-a7dc-2d79584de577.htm "ToString 方法 ")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetCheckResponseLegal 方法 |

检测反馈的消息是否合法

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult CheckResponseLegal(
	byte[] response
)
```

```
Public Shared Function CheckResponseLegal ( 
	response As Byte()
) As OperateResult
```

```
public:
static OperateResult^ CheckResponseLegal(
	array<unsigned char>^ response
)
```

```
static member CheckResponseLegal : 
        response : byte[] -> OperateResult 
```

#### 参数

response
:   类型：SystemByte  
    接收的报文

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtractActualData 方法 

[原文連結](http://api.hslcommunication.cn/html/2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[BuildReadCommand 方法](../html/5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/e8f22d3b-2114-f744-daba-f087099077b7.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm "BuildWriteWordCommand 方法 ")

[CheckResponseLegal 方法](../html/19df45fe-ecef-ad2d-2ce7-d2927673383d.htm "CheckResponseLegal 方法 ")

[ExtractActualData 方法](../html/2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm "ExtractActualData 方法 ")

[GetNewNetMessage 方法](../html/2a7d440d-2475-926e-bf5f-93b73b188939.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ToString 方法](../html/fa489c4e-7431-a042-a7dc-2d79584de577.htm "ToString 方法 ")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetExtractActualData 方法 |

从PLC反馈的数据中提取出实际的数据内容，需要传入反馈数据，是否位读取

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> ExtractActualData(
	byte[] response,
	bool isBit
)
```

```
Public Shared Function ExtractActualData ( 
	response As Byte(),
	isBit As Boolean
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ ExtractActualData(
	array<unsigned char>^ response, 
	bool isBit
)
```

```
static member ExtractActualData : 
        response : byte[] * 
        isBit : bool -> OperateResult<byte[]> 
```

#### 参数

response
:   类型：SystemByte  
    反馈的数据内容

isBit
:   类型：SystemBoolean  
    是否位读取

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
解析后的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/2a7d440d-2475-926e-bf5f-93b73b188939.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[BuildReadCommand 方法](../html/5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/e8f22d3b-2114-f744-daba-f087099077b7.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm "BuildWriteWordCommand 方法 ")

[CheckResponseLegal 方法](../html/19df45fe-ecef-ad2d-2ce7-d2927673383d.htm "CheckResponseLegal 方法 ")

[ExtractActualData 方法](../html/2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm "ExtractActualData 方法 ")

[GetNewNetMessage 方法](../html/2a7d440d-2475-926e-bf5f-93b73b188939.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ToString 方法](../html/fa489c4e-7431-a042-a7dc-2d79584de577.htm "ToString 方法 ")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
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

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/55e3ae9c-4e04-2673-edb6-3d7e0e6add5d.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetRead 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Read(String, UInt16)](55e3ae9c-4e04-2673-edb6-3d7e0e6add5d.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/55e3ae9c-4e04-2673-edb6-3d7e0e6add5d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/55e3ae9c-4e04-2673-edb6-3d7e0e6add5d.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetRead 方法 (String, UInt16) |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  
Batch read byte array information, need to specify the address and length, return the original byte array

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
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

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[Read 重载](dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadAsync 方法 (String, UInt16)](../html/ffd358d4-af01-dd85-da93-141616d87c61.htm "ReadAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetReadAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](ffd358d4-af01-dd85-da93-141616d87c61.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/ffd358d4-af01-dd85-da93-141616d87c61.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadAsync 方法 (String, UInt16)](../html/ffd358d4-af01-dd85-da93-141616d87c61.htm "ReadAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetReadAsync 方法 (String, UInt16) |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  
Batch read byte array information, need to specify the address and length, return the original byte array

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<byte[]>> ReadAsync(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadAsync ( 
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
virtual Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>> 
override ReadAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>>
```

#### 参数

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的byte[]数组

#### 实现

[IReadWriteNetReadAsync(String, UInt16)](59839c27-5a76-2f53-730c-c9d2bb41c2c4.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[ReadAsync 重载](6a31e8a9-e949-4170-cb44-09145d5191a5.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBool 方法 (String, UInt16)](../html/3c9dee5d-08a3-2882-f874-19a70acff835.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetReadBool 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](3c9dee5d-08a3-2882-f874-19a70acff835.htm) | 批量读取Boolean数组信息，需要指定地址和长度，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch read Boolean array information, need to specify the address and length, return Boolean array. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/3c9dee5d-08a3-2882-f874-19a70acff835.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBool 方法 (String, UInt16)](../html/3c9dee5d-08a3-2882-f874-19a70acff835.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetReadBool 方法 (String, UInt16) |

批量读取Boolean数组信息，需要指定地址和长度，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  
Batch read Boolean array information, need to specify the address and length, return Boolean array.
Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address.

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
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
带有成功标识的byte[]数组

#### 实现

[IReadWriteNetReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm)

![](../icons/SectionExpanded.png)备注

根据协议的规范，最多读取256长度的bool数组信息，如果需要读取更长的bool信息，需要按字为单位进行读取的操作。

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[ReadBool 重载](4fb3508a-3d36-efff-b0db-a159d2aba37f.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (String, UInt16)](../html/0df9b5c8-ae62-00d1-4f14-0a660772ce98.htm "ReadBoolAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetReadBoolAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](0df9b5c8-ae62-00d1-4f14-0a660772ce98.htm) | 批量读取Boolean数组信息，需要指定地址和长度，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch read Boolean array information, need to specify the address and length, return Boolean array. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/0df9b5c8-ae62-00d1-4f14-0a660772ce98.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (String, UInt16)](../html/0df9b5c8-ae62-00d1-4f14-0a660772ce98.htm "ReadBoolAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetReadBoolAsync 方法 (String, UInt16) |

批量读取Boolean数组信息，需要指定地址和长度，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  
Batch read Boolean array information, need to specify the address and length, return Boolean array.
Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address.

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<bool[]>> ReadBoolAsync(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadBoolAsync ( 
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
virtual Task<OperateResult<array<bool>^>^>^ ReadBoolAsync(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadBoolAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>> 
override ReadBoolAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>>
```

#### 参数

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带有成功标识的byte[]数组

#### 实现

[IReadWriteNetReadBoolAsync(String, UInt16)](74e40b23-198a-944b-6ed7-3f58ca51da5d.htm)

![](../icons/SectionExpanded.png)备注

根据协议的规范，最多读取256长度的bool数组信息，如果需要读取更长的bool信息，需要按字为单位进行读取的操作。

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[ReadBoolAsync 重载](d69179c0-10c0-a6d2-3470-96c996d8a190.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/fa489c4e-7431-a042-a7dc-2d79584de577.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[BuildReadCommand 方法](../html/5ed2ba4a-c6d8-2105-59b0-1d8380dea37a.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/e8f22d3b-2114-f744-daba-f087099077b7.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/6cccd7dd-247b-5d0e-1691-6ddaa401aa18.htm "BuildWriteWordCommand 方法 ")

[CheckResponseLegal 方法](../html/19df45fe-ecef-ad2d-2ce7-d2927673383d.htm "CheckResponseLegal 方法 ")

[ExtractActualData 方法](../html/2c5cac0d-8ebc-90b8-8f46-11a595eddc84.htm "ExtractActualData 方法 ")

[GetNewNetMessage 方法](../html/2a7d440d-2475-926e-bf5f-93b73b188939.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/dbfd271e-c1e5-d5ab-753c-9975f0b1e3af.htm "Read 方法 ")

[ReadAsync 方法](../html/6a31e8a9-e949-4170-cb44-09145d5191a5.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/4fb3508a-3d36-efff-b0db-a159d2aba37f.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/d69179c0-10c0-a6d2-3470-96c996d8a190.htm "ReadBoolAsync 方法 ")

[ToString 方法](../html/fa489c4e-7431-a042-a7dc-2d79584de577.htm "ToString 方法 ")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetToString 方法 |

[缺少 "M:HslCommunication.Profinet.Melsec.MelsecA1EAsciiNet.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
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

[缺少 "M:HslCommunication.Profinet.Melsec.MelsecA1EAsciiNet.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/e14fce8e-c661-af15-9681-401be520c99f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[Write 方法 (String, Boolean[])](../html/15d948d6-ab6f-9d25-6ef0-3d1a975de405.htm "Write 方法 (String, Boolean[])")

[Write 方法 (String, Byte[])](../html/49c72530-9bee-f030-99ec-e73774b17de1.htm "Write 方法 (String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetWrite 方法 |

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
| 公共方法 | [Write(String, Boolean)](15d948d6-ab6f-9d25-6ef0-3d1a975de405.htm) | 批量写入Boolean数组数据，返回是否成功，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch write Boolean array data, return whether the write was successful. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](49c72530-9bee-f030-99ec-e73774b17de1.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/15d948d6-ab6f-9d25-6ef0-3d1a975de405.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[Write 方法 (String, Boolean[])](../html/15d948d6-ab6f-9d25-6ef0-3d1a975de405.htm "Write 方法 (String, Boolean[])")

[Write 方法 (String, Byte[])](../html/49c72530-9bee-f030-99ec-e73774b17de1.htm "Write 方法 (String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetWrite 方法 (String, Boolean) |

批量写入Boolean数组数据，返回是否成功，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  
Batch write Boolean array data, return whether the write was successful.
Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address.

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
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
	bool[] values
)
```

```
Public Overrides Function Write ( 
	address As String,
	values As Boolean()
) As OperateResult
```

```
public:
virtual OperateResult^ Write(
	String^ address, 
	array<bool>^ values
) override
```

```
abstract Write : 
        address : string * 
        values : bool[] -> OperateResult 
override Write : 
        address : string * 
        values : bool[] -> OperateResult
```

#### 参数

address
:   类型：SystemString  
    起始地址

values
:   类型：SystemBoolean  

    [缺少 "M:HslCommunication.Profinet.Melsec.MelsecA1EAsciiNet.Write(System.String,System.Boolean[])" 的 <param name="values"/> 文档]

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
带有成功标识的结果类对象

#### 实现

[IReadWriteNetWrite(String, Boolean)](2a99122c-1496-d766-f0aa-44a18918fcd3.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[Write 重载](e14fce8e-c661-af15-9681-401be520c99f.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/49c72530-9bee-f030-99ec-e73774b17de1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[Write 方法](../html/e14fce8e-c661-af15-9681-401be520c99f.htm "Write 方法 ")

[Write 方法 (String, Boolean[])](../html/15d948d6-ab6f-9d25-6ef0-3d1a975de405.htm "Write 方法 (String, Boolean[])")

[Write 方法 (String, Byte[])](../html/49c72530-9bee-f030-99ec-e73774b17de1.htm "Write 方法 (String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetWrite 方法 (String, Byte) |

写入原始的byte数组数据到指定的地址，返回是否写入成功  
Write the original byte array data to the specified address, and return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
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

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[Write 重载](e14fce8e-c661-af15-9681-401be520c99f.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

[WriteAsync 方法 (String, Boolean[])](../html/8703ce64-be44-f53a-4666-da5f104cf3c2.htm "WriteAsync 方法 (String, Boolean[])")

[WriteAsync 方法 (String, Byte[])](../html/60031b6f-fb61-9b88-045b-a201ab534826.htm "WriteAsync 方法 (String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetWriteAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [WriteAsync(String, Boolean)](8703ce64-be44-f53a-4666-da5f104cf3c2.htm) | 批量写入Boolean数组数据，返回是否成功，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  Batch write Boolean array data, return whether the write was successful. Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address. (重写 [DeviceCommunicationWriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm).) |
| 公共方法 | [WriteAsync(String, Byte)](60031b6f-fb61-9b88-045b-a201ab534826.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/8703ce64-be44-f53a-4666-da5f104cf3c2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

[WriteAsync 方法 (String, Boolean[])](../html/8703ce64-be44-f53a-4666-da5f104cf3c2.htm "WriteAsync 方法 (String, Boolean[])")

[WriteAsync 方法 (String, Byte[])](../html/60031b6f-fb61-9b88-045b-a201ab534826.htm "WriteAsync 方法 (String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetWriteAsync 方法 (String, Boolean) |

批量写入Boolean数组数据，返回是否成功，地址示例M100，S100，B1A，如果是X,Y, X017就是8进制地址，Y10就是16进制地址。  
Batch write Boolean array data, return whether the write was successful.
Examples of addresses M100, S100, B1A, if it is X, Y, X017 is an octal address, Y10 is a hexadecimal address.

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult> WriteAsync(
	string address,
	bool[] values
)
```

```
Public Overrides Function WriteAsync ( 
	address As String,
	values As Boolean()
) As Task(Of OperateResult)
```

```
public:
virtual Task<OperateResult^>^ WriteAsync(
	String^ address, 
	array<bool>^ values
) override
```

```
abstract WriteAsync : 
        address : string * 
        values : bool[] -> Task<OperateResult> 
override WriteAsync : 
        address : string * 
        values : bool[] -> Task<OperateResult>
```

#### 参数

address
:   类型：SystemString  
    起始地址

values
:   类型：SystemBoolean  

    [缺少 "M:HslCommunication.Profinet.Melsec.MelsecA1EAsciiNet.WriteAsync(System.String,System.Boolean[])" 的 <param name="values"/> 文档]

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
带有成功标识的结果类对象

#### 实现

[IReadWriteNetWriteAsync(String, Boolean)](cf8f8e56-df3d-607b-1c52-c1a76378db39.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[WriteAsync 重载](6a177b75-707d-9f1a-c37d-2e1eb103d675.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/60031b6f-fb61-9b88-045b-a201ab534826.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[WriteAsync 方法](../html/6a177b75-707d-9f1a-c37d-2e1eb103d675.htm "WriteAsync 方法 ")

[WriteAsync 方法 (String, Boolean[])](../html/8703ce64-be44-f53a-4666-da5f104cf3c2.htm "WriteAsync 方法 (String, Boolean[])")

[WriteAsync 方法 (String, Byte[])](../html/60031b6f-fb61-9b88-045b-a201ab534826.htm "WriteAsync 方法 (String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNetWriteAsync 方法 (String, Byte) |

写入原始的byte数组数据到指定的地址，返回是否写入成功  
Write the original byte array data to the specified address, and return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.Melsec](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult> WriteAsync(
	string address,
	byte[] value
)
```

```
Public Overrides Function WriteAsync ( 
	address As String,
	value As Byte()
) As Task(Of OperateResult)
```

```
public:
virtual Task<OperateResult^>^ WriteAsync(
	String^ address, 
	array<unsigned char>^ value
) override
```

```
abstract WriteAsync : 
        address : string * 
        value : byte[] -> Task<OperateResult> 
override WriteAsync : 
        address : string * 
        value : byte[] -> Task<OperateResult>
```

#### 参数

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemByte  
    写入值

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

#### 实现

[IReadWriteNetWriteAsync(String, Byte)](1a2e83c8-7782-847a-5d89-08d5556ac4ad.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[WriteAsync 重载](6a177b75-707d-9f1a-c37d-2e1eb103d675.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MelsecA1EAsciiNet 字段

[原文連結](http://api.hslcommunication.cn/html/1e1755d0-19bc-dadb-ef41-c1365e3e2a9d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec](../html/a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm "HslCommunication.Profinet.Melsec")

[MelsecA1EAsciiNet 类](../html/9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm "MelsecA1EAsciiNet 类")

[MelsecA1EAsciiNet 构造函数](../html/d531f4a7-510e-2e2a-c8f8-44188bb72a9e.htm "MelsecA1EAsciiNet 构造函数 ")

[MelsecA1EAsciiNet 属性](../html/f50f623e-d269-67db-4417-de0077b9583d.htm "MelsecA1EAsciiNet 属性")

[MelsecA1EAsciiNet 方法](../html/dc186abc-1c6a-9346-cc58-a862cf1573b9.htm "MelsecA1EAsciiNet 方法")

[MelsecA1EAsciiNet 字段](../html/1e1755d0-19bc-dadb-ef41-c1365e3e2a9d.htm "MelsecA1EAsciiNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MelsecA1EAsciiNet 字段 |

[MelsecA1EAsciiNet](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MelsecA1EAsciiNet 类](9d731f70-a0c0-8e6c-44e0-334f4f74bf52.htm)

[HslCommunication.Profinet.Melsec 命名空间](a75ad826-88b7-0dd0-83a5-a89f9c381b22.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)