# HslCommunication - HslCommunication.Profinet.Omron

> 分類頁數: 30



---
## HslCommunication.Profinet.Omron

[原文連結](http://api.hslcommunication.cn/html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipServer 类](../html/ba4a1901-a9e5-9cbb-eb52-741b40ce6038.htm "OmronCipServer 类")

[OmronConnectedCipNet 类](../html/a1761eed-c88f-301f-ea86-7a8e1cb0a477.htm "OmronConnectedCipNet 类")

[OmronCpuUnitData 类](../html/8b70e014-ce2a-a83e-f0d9-90b1109b1661.htm "OmronCpuUnitData 类")

[OmronCpuUnitStatus 类](../html/3728b500-7fe6-71ee-fc3a-e003e54d6238.htm "OmronCpuUnitStatus 类")

[OmronFinsDataType 类](../html/8970f2b3-91d5-18d2-f6a9-9857c544161d.htm "OmronFinsDataType 类")

[OmronFinsNet 类](../html/e8e99732-9b55-e5e6-7009-51b107a8bd8b.htm "OmronFinsNet 类")

[OmronFinsNetHelper 类](../html/c41e6e70-c5e8-0a1e-2f26-cf45f77140af.htm "OmronFinsNetHelper 类")

[OmronFinsServer 类](../html/51a26f20-e69d-95c5-be5e-f4f64188c1e1.htm "OmronFinsServer 类")

[OmronFinsUdp 类](../html/93245a23-e751-2e3c-9092-36be61442090.htm "OmronFinsUdp 类")

[OmronFinsUdpServer 类](../html/8e82e59b-322b-ae7d-34b8-d6609314a446.htm "OmronFinsUdpServer 类")

[OmronHostLink 类](../html/2462cf2c-e8d2-d4c2-e2e5-64e5f5c2f76c.htm "OmronHostLink 类")

[OmronHostLinkCMode 类](../html/b6e22bbb-3c82-741b-c808-d0c005869ad9.htm "OmronHostLinkCMode 类")

[OmronHostLinkCModeOverTcp 类](../html/6e47e632-2c03-07c7-8066-9d4ea3b965b5.htm "OmronHostLinkCModeOverTcp 类")

[OmronHostLinkCModeServer 类](../html/9cac46b2-a2aa-7ad4-9a7d-d4de58dd89c8.htm "OmronHostLinkCModeServer 类")

[OmronHostLinkOverTcp 类](../html/dc9f53e0-b653-5b42-074f-57b99e4eea7e.htm "OmronHostLinkOverTcp 类")

[OmronHostLinkServer 类](../html/8d6e8be0-1c47-0849-8df2-dbd65125ef46.htm "OmronHostLinkServer 类")

[OmronPlcType 枚举](../html/89c0ba8d-a44f-aad6-7bb0-3b1c8a18183e.htm "OmronPlcType 枚举")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Omron 命名空间 |

[缺少 "N:HslCommunication.Profinet.Omron" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [OmronCipNet](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm) | 欧姆龙PLC的CIP协议的类，支持NJ,NX,NY系列PLC，支持tag名的方式读写数据，假设你读取的是局部变量，那么使用 Program:MainProgram.变量名  Omron PLC's CIP protocol class, support NJ, NX, NY series PLC, support tag name read and write data, assuming you read local variables, then use Program: MainProgram. Variable name |
| 公共类 | [OmronCipServer](ba4a1901-a9e5-9cbb-eb52-741b40ce6038.htm) | 欧姆龙的CIP虚拟服务器 |
| 公共类代码示例 | [OmronConnectedCipNet](a1761eed-c88f-301f-ea86-7a8e1cb0a477.htm) | 基于连接的对象访问的CIP协议的实现，用于对Omron PLC进行标签的数据读写，对数组，多维数组进行读写操作，支持的数据类型请参照API文档手册。  The implementation of the CIP protocol based on connected object access is used to read and write tag data to Omron PLC, and read and write arrays and multidimensional arrays. For the supported data types, please refer to the API documentation manual. |
| 公共类 | [OmronCpuUnitData](8b70e014-ce2a-a83e-f0d9-90b1109b1661.htm) | 欧姆的Cpu的单元信息数据类 |
| 公共类 | [OmronCpuUnitStatus](3728b500-7fe6-71ee-fc3a-e003e54d6238.htm) | 欧姆龙Cpu的状态信息 |
| 公共类 | [OmronFinsDataType](8970f2b3-91d5-18d2-f6a9-9857c544161d.htm) | 欧姆龙的Fins协议的数据类型 |
| 公共类代码示例 | [OmronFinsNet](e8e99732-9b55-e5e6-7009-51b107a8bd8b.htm) | 欧姆龙PLC通讯类，采用Fins-Tcp通信协议实现，支持的地址信息参见api文档信息。本协议下PLC默认的端口号为 9600，也可以手动更改，重启PLC更改生效。  Omron PLC communication class is implemented using Fins-Tcp communication protocol. For the supported address information, please refer to the api document information. The default port number of the PLC under this protocol is 9600, and it can also be changed manually. Restart the PLC to make the changes take effect. |
| 公共类 | [OmronFinsNetHelper](c41e6e70-c5e8-0a1e-2f26-cf45f77140af.htm) | Omron PLC的FINS协议相关的辅助类，主要是一些地址解析，读写的指令生成。  The auxiliary classes related to the FINS protocol of Omron PLC are mainly some address resolution and the generation of read and write instructions. |
| 公共类 | [OmronFinsServer](51a26f20-e69d-95c5-be5e-f4f64188c1e1.htm) | 欧姆龙的虚拟服务器，支持DM区，CIO区，Work区，Hold区，Auxiliary区，可以方便的进行测试  Omron's virtual server supports DM area, CIO area, Work area, Hold area, and Auxiliary area, which can be easily tested |
| 公共类 | [OmronFinsUdp](93245a23-e751-2e3c-9092-36be61442090.htm) | 欧姆龙的Udp协议的实现类，地址类型和Fins-TCP一致，无连接的实现，可靠性不如[OmronFinsNet](e8e99732-9b55-e5e6-7009-51b107a8bd8b.htm)  Omron's Udp protocol implementation class, the address type is the same as Fins-TCP, and the connectionless implementation is not as reliable as [OmronFinsNet](e8e99732-9b55-e5e6-7009-51b107a8bd8b.htm) |
| 公共类 | [OmronFinsUdpServer](8e82e59b-322b-ae7d-34b8-d6609314a446.htm) | 欧姆龙的虚拟服务器，支持DM区，CIO区，Work区，Hold区，Auxiliary区，可以方便的进行测试  Omron's virtual server supports DM area, CIO area, Work area, Hold area, and Auxiliary area, which can be easily tested |
| 公共类代码示例 | [OmronHostLink](2462cf2c-e8d2-d4c2-e2e5-64e5f5c2f76c.htm) | 欧姆龙的HostLink协议的实现，地址支持示例 DM区:D100; CIO区:C100; Work区:W100; Holding区:H100; Auxiliary区: A100  Implementation of Omron's HostLink protocol, address support example DM area: D100; CIO area: C100; Work area: W100; Holding area: H100; Auxiliary area: A100 |
| 公共类 | [OmronHostLinkCMode](b6e22bbb-3c82-741b-c808-d0c005869ad9.htm) | 欧姆龙的HostLink的C-Mode实现形式，地址支持携带站号信息，例如：s=2;D100  Omron's HostLink C-Mode implementation form, the address supports carrying station number information, for example: s=2;D100 |
| 公共类 | [OmronHostLinkCModeOverTcp](6e47e632-2c03-07c7-8066-9d4ea3b965b5.htm) | 欧姆龙的HostLink的C-Mode实现形式，当前的类是通过以太网透传实现。地址支持携带站号信息，例如：s=2;D100  The C-Mode implementation form of Omron’s HostLink, the current class is realized through Ethernet transparent transmission. Address supports carrying station number information, for example: s=2;D100 |
| 公共类 | [OmronHostLinkCModeServer](9cac46b2-a2aa-7ad4-9a7d-d4de58dd89c8.htm) | 欧姆龙的HostLinkCMode协议的虚拟服务器 |
| 公共类 | [OmronHostLinkOverTcp](dc9f53e0-b653-5b42-074f-57b99e4eea7e.htm) | 欧姆龙的HostLink协议的实现，基于Tcp实现，地址支持示例 DM区:D100; CIO区:C100; Work区:W100; Holding区:H100; Auxiliary区: A100  Implementation of Omron's HostLink protocol, based on tcp protocol, address support example DM area: D100; CIO area: C100; Work area: W100; Holding area: H100; Auxiliary area: A100 |
| 公共类 | [OmronHostLinkServer](8d6e8be0-1c47-0849-8df2-dbd65125ef46.htm) | 欧姆龙的HostLink虚拟服务器，支持DM区，CIO区，Work区，Hold区，Auxiliary区，可以方便的进行测试  Omron's HostLink virtual server supports DM area, CIO area, Work area, Hold area, and Auxiliary area, which can be easily tested |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [OmronPlcType](89c0ba8d-a44f-aad6-7bb0-3b1c8a18183e.htm) | 欧姆龙PLC的类型 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OmronCipNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/bda49fa2-ff51-4a37-f7e1-9a83ee8eb428.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 构造函数](../html/bda49fa2-ff51-4a37-f7e1-9a83ee8eb428.htm "OmronCipNet 构造函数 ")

[OmronCipNet 构造函数](../html/2bdc7a11-ca27-f4b4-fb38-f89bcf3b2fff.htm "OmronCipNet 构造函数 ")

[OmronCipNet 构造函数 (String, Int32)](../html/f2b7b1ad-cac4-b9cb-9cf8-eae3547ed52d.htm "OmronCipNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNet 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [OmronCipNet](2bdc7a11-ca27-f4b4-fb38-f89bcf3b2fff.htm) | Instantiate a communication object for a OmronCipNet PLC protocol |
| 公共方法 | [OmronCipNet(String, Int32)](f2b7b1ad-cac4-b9cb-9cf8-eae3547ed52d.htm) | Specify the IP address and port to instantiate a communication object for a OmronCipNet PLC protocol |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OmronCipNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/2bdc7a11-ca27-f4b4-fb38-f89bcf3b2fff.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 构造函数](../html/bda49fa2-ff51-4a37-f7e1-9a83ee8eb428.htm "OmronCipNet 构造函数 ")

[OmronCipNet 构造函数](../html/2bdc7a11-ca27-f4b4-fb38-f89bcf3b2fff.htm "OmronCipNet 构造函数 ")

[OmronCipNet 构造函数 (String, Int32)](../html/f2b7b1ad-cac4-b9cb-9cf8-eae3547ed52d.htm "OmronCipNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNet 构造函数 |

Instantiate a communication object for a OmronCipNet PLC protocol

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OmronCipNet()
```

```
Public Sub New
```

```
public:
OmronCipNet()
```

```
new : unit -> OmronCipNet
```

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[OmronCipNet 重载](bda49fa2-ff51-4a37-f7e1-9a83ee8eb428.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OmronCipNet 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/f2b7b1ad-cac4-b9cb-9cf8-eae3547ed52d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 构造函数](../html/bda49fa2-ff51-4a37-f7e1-9a83ee8eb428.htm "OmronCipNet 构造函数 ")

[OmronCipNet 构造函数](../html/2bdc7a11-ca27-f4b4-fb38-f89bcf3b2fff.htm "OmronCipNet 构造函数 ")

[OmronCipNet 构造函数 (String, Int32)](../html/f2b7b1ad-cac4-b9cb-9cf8-eae3547ed52d.htm "OmronCipNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNet 构造函数 (String, Int32) |

Specify the IP address and port to instantiate a communication object for a OmronCipNet PLC protocol

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OmronCipNet(
	string ipAddress,
	int port = 44818
)
```

```
Public Sub New ( 
	ipAddress As String,
	Optional port As Integer = 44818
)
```

```
public:
OmronCipNet(
	String^ ipAddress, 
	int port = 44818
)
```

```
new : 
        ipAddress : string * 
        ?port : int 
(* Defaults:
        let _port = defaultArg port 44818
*)
-> OmronCipNet
```

#### 参数

ipAddress
:   类型：SystemString  
    PLC IpAddress

port (Optional)
:   类型：SystemInt32  
    PLC Port

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[OmronCipNet 重载](bda49fa2-ff51-4a37-f7e1-9a83ee8eb428.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OmronCipNet 属性

[原文連結](http://api.hslcommunication.cn/html/488cfe9f-3319-2ccc-ad2a-ba3435024419.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 构造函数](../html/bda49fa2-ff51-4a37-f7e1-9a83ee8eb428.htm "OmronCipNet 构造函数 ")

[OmronCipNet 属性](../html/488cfe9f-3319-2ccc-ad2a-ba3435024419.htm "OmronCipNet 属性")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[OmronCipNet 字段](../html/37dabe94-9b03-1948-8f79-4bbdfad15f11.htm "OmronCipNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNet 属性 |

[OmronCipNet](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CipCommand](eddd2d91-8cdd-d1ae-9eff-a40cb561cf95.htm) | 获取或设置整个交互指令的控制码，默认为0x6F，通常不需要修改  Gets or sets the control code of the entire interactive instruction. The default is 0x6F, and usually does not need to be modified. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [ContextCheck](f03df045-7245-ec0e-7d46-c4755856ab45.htm) | Gets or sets a value indicating whether the current context meets the required conditions. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [MessageRouter](65036b90-4e01-b3ca-2349-237cc61dc17b.htm) | 获取或设置当前的通信的消息路由信息，可以实现一些复杂情况的通信，数据包含背板号，路由参数，slot，例如：1.15.2.18.1.1  Get or set the message routing information of the current communication, which can realize some complicated communication. The data includes the backplane number, routing parameters, and slot, for example: 1.15.2.18.1.1 (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [PortSlot](f7e49c4c-4c89-b699-ff9a-44967f95e548.htm) | port and slot information (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共属性 | [ReadArrayUseSegment](d816280d-94a3-1349-a0ef-a6a2640a7a09.htm) | 获取或设置是否在读取数组的时候使用片段读取功能，默认是开启的。  Get or set whether to use the fragment reading function when reading the array. By default, it is enabled. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SessionHandle](c64b9588-9b29-1079-84d9-7e38f9317ee8.htm) | The current session handle, which is determined by the PLC when communicating with the PLC handshake (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Slot](ab3e7758-dd89-4760-cf98-c4ab157dbfbf.htm) | Gets or sets the slot number information for the current plc, which should be set before connections (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [TagValueAndDataBits](0e1970dc-cbb8-b71b-aaf0-97cce53ba136.htm) | 获取当前缓存的Bool变量的类型信息，key为标签地址，value为数据位长度信息 (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OmronCipNet 方法

[原文連結](http://api.hslcommunication.cn/html/a6223914-afee-ab41-ba21-5b30b468808d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[GetBoolWritePadding 方法](../html/41ba1de6-9932-39eb-2e3e-7b3867726987.htm "GetBoolWritePadding 方法 ")

[GetWriteValueLength 方法](../html/3fc102d6-0067-c2c7-f4e7-f67879d2825e.htm "GetWriteValueLength 方法 ")

[Read 方法](../html/c3bc4e40-cb07-a183-5cc3-37dfca2727ae.htm "Read 方法 ")

[ReadAsync 方法](../html/e3167e9d-439a-6ef2-e03f-9f4898ea526c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/08b48880-ba18-76eb-009f-5bf3320e6b47.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/696b89a0-1a79-589c-e14e-75a7a9b26af9.htm "ReadBoolAsync 方法 ")

[ReadString 方法](../html/c72f72db-c199-9b51-d3bf-3d93e87ed17a.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/67d00847-5a0e-781f-3f2b-7f3f9752b8b5.htm "ReadStringAsync 方法 ")

[ReadStruct(T) 方法](../html/03fd2a05-8efd-8508-cdd7-564432b48c75.htm "ReadStruct(T) 方法 ")

[ToString 方法](../html/96331804-b4e1-e30a-7493-2f9fe10ee376.htm "ToString 方法 ")

[Write 方法](../html/9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm "Write 方法 ")

[WriteAsync 方法](../html/558f1912-35b3-64b1-eca5-833d440e5b2f.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNet 方法 |

[OmronCipNet](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [BuildReadCommand(String)](3311d7b3-f1fe-956e-d4c6-2a3947e2f97c.htm) | 创建一个读取多标签的报文  Build a read command bytes (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [BuildReadCommand(String, UInt16)](03006976-dcf4-28ac-6b89-7ac3d9b4e9a1.htm) | 创建一个读取标签的报文指定，标签地址可以手动动态指定slot编号，例如 slot=2;AAA  Build a read command bytes, The label address can manually specify the slot number dynamically, for example slot=2;AAA (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [BuildWriteCommand(String, Boolean)](21f58acf-ba63-86d9-eb96-68cf4d707524.htm) | Create a written message instruction (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 受保护的方法 | [BuildWriteCommand(String, UInt16, Byte, Int32)](7e6c86c0-beb6-1c72-e94d-b39d740d0d11.htm) | Create a written message instruction (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ConnectClose](569a9c0e-7628-b8d5-e124-0853452747b8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](76be8150-ff22-95b9-16b1-7c3440c54ae8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServer](6299cc1f-2fab-b754-4597-1eddc2b21308.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](e7eff5c0-8e2a-3b6c-5027-16f4d8bcf5ae.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](e55ae5e7-d0d9-2b64-7e74-62170d54feab.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](240a4b48-1173-f335-7939-ec2124d3cf1b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 受保护的方法 | [GetBoolWritePadding](41ba1de6-9932-39eb-2e3e-7b3867726987.htm) | 当进行bool写入的时候，是否需要补齐到字节长度，默认不需要  When performing bool writing, whether it is necessary to complement to the byte length is not required by default (重写 [AllenBradleyNetGetBoolWritePadding](5341096d-0498-1607-19ca-ee6dfb9c8279.htm).) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](ef0b2f59-ff51-1a4b-f2f9-60ed23e85826.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [GetWriteValueLength](3fc102d6-0067-c2c7-f4e7-f67879d2825e.htm) | 获取写入数据的长度信息，此处直接返回数组的长度信息 (重写 [AllenBradleyNetGetWriteValueLength(String, Int32)](74c8fcd4-f6f4-04e1-2ec7-24ba19974f9a.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnect](78b23eef-7249-a62f-8807-688647efbf2e.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](a66dbf16-f5f0-8d8c-5696-4424bc933931.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 受保护的方法 | [PackCommandService](f4d5fa10-bcdc-b01c-d6f1-7116039916c0.htm) | 将所有的cip指定进行打包操作。 (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [PackCommandWithHeader](c2b5345a-da00-9916-93ad-8bba94d2a11b.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Read(String)](682e8e26-45f6-f437-bbe2-a345db989fb8.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"，每个地址的数据长度为1，表示一个数据，最终读取返回的是一整个的字节数组，需要自行解析 **[Authorization]** Batch read data information of multiple addresses, for example, I can read the data of two tags "A", "B[0]", the data length of each address is 1, which means one data, and the final read returns a The entire byte array, which needs to be parsed by itself (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Read(String, UInt16)](a6216184-0009-e4a9-96cd-a2158390aa44.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"， 长度为 [1, 5]，返回的是一整个的字节数组，需要自行解析 **[Authorization]** Read the data information of multiple addresses in batches. For example, I can read the data "A", "B[0]" of two tags, the length is [1, 5], and the return is an entire byte array, and I need to do it myself Parsing (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Read(String, UInt16)](2954efdc-df90-5d0f-71ca-fc1f5d94a0c1.htm) | 读取指定地址的二进制数据内容，长度为地址长度，一般都是1，除非读取数组时，如果需要强制使用 片段读取功能码，则地址里携带 x=0x52; 或是 x=82; 则强制使用片段读取。  Read the binary data content of the specified address, the length is the address length, generally 1, unless the array is read, if you need to force the fragment reading function code, the address carries x=0x52; or x=82; then the fragment read is forced. (重写 [AllenBradleyNetRead(String, UInt16)](cc6bd266-1e7b-5052-cdf3-b5e45abf6b48.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String)](d666e005-8c96-6381-fa24-60fd58545d95.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"，每个地址的数据长度为1，表示一个数据，最终读取返回的是一整个的字节数组，需要自行解析 **[Authorization]** Batch read data information of multiple addresses, for example, I can read the data of two tags "A", "B[0]", the data length of each address is 1, which means one data, and the final read returns a The entire byte array, which needs to be parsed by itself (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](856714b2-0b15-fe6e-253f-201b6bd0b010.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"， 长度为 [1, 5]，返回的是一整个的字节数组，需要自行解析 **[Authorization]** Read the data information of multiple addresses in batches. For example, I can read the data "A", "B[0]" of two tags, the length is [1, 5], and the return is an entire byte array, and I need to do it myself Parsing (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](ada9961d-130c-7c7b-e883-5a229b1fbbed.htm) | 读取指定地址的二进制数据内容，长度为地址长度，一般都是1，除非读取数组时，如果需要强制使用 片段读取功能码，则地址里携带 x=0x52; 或是 x=82; 则强制使用片段读取。  Read the binary data content of the specified address, the length is the address length, generally 1, unless the array is read, if you need to force the fragment reading function code, the address carries x=0x52; or x=82; then the fragment read is forced. (重写 [AllenBradleyNetReadAsync(String, UInt16)](12fbde37-b7d0-54e4-4196-fa53d1976d16.htm).) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](3d12b2f9-19b8-28ea-a1e2-5d8d7c18bf2d.htm) | 读取单个的bool数据信息，如果读取的是单bool变量，就直接写变量名，如果是由int组成的bool数组的一个值，一律带"i="开头访问，例如"i=A[0]"   Read a single bool data information, if it is a single bool variable, write the variable name directly, if it is a value of a bool array composed of int, it is always accessed with "i=" at the beginning, for example, "i=A[0]" (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](f6c36d10-c3a9-24cf-e0a1-9eac0b791c04.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (重写 [AllenBradleyNetReadBool(String, UInt16)](bae5fa6c-770b-f495-a578-15bda70b6449.htm).) |
| 公共方法 | [ReadBoolArray](12218b60-0b18-9502-e1fb-138fcc5f754b.htm) | 批量读取的bool数组信息，如果你有个Bool数组变量名为 A, 那么读第0个位，可以通过 ReadBool("A")，但是第二个位需要使用 ReadBoolArray("A[0]") // 返回32个bool长度，0-31的索引，如果我想读取32-63的位索引，就需要 ReadBoolArray("A[1]") ，以此类推。  For batch read bool array information, if you have a Bool array variable named A, then you can read the 0th bit through ReadBool("A"), but the second bit needs to use ReadBoolArray("A[0]" ) // Returns the length of 32 bools, the index is 0-31, if I want to read the bit index of 32-63, I need ReadBoolArray("A[1]"), and so on. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadBoolArrayAsync](20f6eb3c-4f89-f1a7-0a3e-4436bdac324e.htm) | 批量读取的bool数组信息，如果你有个Bool数组变量名为 A, 那么读第0个位，可以通过 ReadBool("A")，但是第二个位需要使用 ReadBoolArray("A[0]") // 返回32个bool长度，0-31的索引，如果我想读取32-63的位索引，就需要 ReadBoolArray("A[1]") ，以此类推。  For batch read bool array information, if you have a Bool array variable named A, then you can read the 0th bit through ReadBool("A"), but the second bit needs to use ReadBoolArray("A[0]" ) // Returns the length of 32 bools, the index is 0-31, if I want to read the bit index of 32-63, I need ReadBoolArray("A[1]"), and so on. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](28957b24-aa96-b3f4-6187-4873cbd9c23b.htm) | 读取单个的bool数据信息，如果读取的是单bool变量，就直接写变量名，如果是由int组成的bool数组的一个值，一律带"i="开头访问，例如"i=A[0]"   Read a single bool data information, if it is a single bool variable, write the variable name directly, if it is a value of a bool array composed of int, it is always accessed with "i=" at the beginning, for example, "i=A[0]" (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](2ff1bdf7-1827-1e13-adbc-54c939c9758b.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (重写 [AllenBradleyNetReadBoolAsync(String, UInt16)](e0d9e8f4-4c9f-f7a7-8d7c-9fada2ce8e66.htm).) |
| 公共方法 | [ReadByte](cd4de42a-dcbc-e370-3573-c2736dc879e4.htm) | 读取PLC的byte类型的数据  Read the byte type of PLC data (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadByteAsync](d1c262df-2d65-8981-41c3-fc7a96756d26.htm) | 读取PLC的byte类型的数据  Read the byte type of PLC data (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadCipFromServer](f71285cd-0319-f128-f17d-1ea6f9270dc5.htm) | 使用CIP报文和服务器进行核心的数据交换 (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadCipFromServerAsync](a4b09c13-8b71-b87b-fe4f-d05a4cc82182.htm) | 使用CIP报文和服务器进行核心的数据交换 (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDate](36230862-44c2-e14d-c708-1496413f5c31.htm) | 读取指定地址的日期数据，最小日期为 1970年1月1日，当PLC的变量类型为 "Date" 和 "TimeAndDate" 时，都可以用本方法读取。  Read the date data of the specified address. The minimum date is January 1, 1970. When the PLC variable type is "Date" and "TimeAndDate", this method can be used to read. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadDateAsync](c398eee6-568d-511d-4478-46fc1421aca3.htm) | 读取指定地址的日期数据，最小日期为 1970年1月1日，当PLC的变量类型为 "Date" 和 "TimeAndDate" 时，都可以用本方法读取。  Read the date data of the specified address. The minimum date is January 1, 1970. When the PLC variable type is "Date" and "TimeAndDate", this method can be used to read. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](ffec9173-531c-53de-b444-db5d981f4e1f.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](de8a1b5c-d5cb-bf6c-5c6c-0b533d93bf56.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadEipFromServer](b006fe69-a6bb-a1ac-7163-6faf2bbe7e33.htm) | 使用EIP报文和服务器进行核心的数据交换 (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadEipFromServerAsync](ebc2f077-871e-e0fb-47c6-79f432c1ac3f.htm) | 使用EIP报文和服务器进行核心的数据交换 (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](d9e89cf3-f22a-3586-8bce-292d0761d507.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](c53a71bd-fbbe-b53b-ae34-4563ab2ea4a5.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](677a10df-eade-a88f-0447-c60826aa0d30.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](18dab656-d028-17ce-293c-e12bc36b133d.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](196a3391-ee8c-1014-09b2-2201174b913f.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](5e33e0a3-5428-3db6-28a8-3140bef5ee61.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](2c015ba8-8290-5757-eab9-d9d8bd637fc7.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](7cbcde11-b758-7812-3dd0-6d1f17e08f3d.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadPlcType](c0b8e0e7-deb5-3117-61df-58d3424c30b0.htm) | 从PLC里读取当前PLC的型号信息  Read the current PLC model information from the PLC (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadPlcTypeAsync](027468c4-60f1-8b07-94ae-e152902e19b5.htm) | 从PLC里读取当前PLC的型号信息  Read the current PLC model information from the PLC (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadSegment](e7f0c7ac-d304-31a7-ab2c-7a4d8ac00716.htm) | Read Segment Data Array form plc, use address tag name (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadSegmentAsync](3c490d83-a966-793e-a990-13085f00b15d.htm) | Read Segment Data Array form plc, use address tag name (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadString(String)](bd50aa1d-e5d1-df14-61e3-6341d1ce151b.htm) | (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](bf38a088-492e-4c7c-2f29-dd56d4230bbc.htm) | 读取字符串数据，默认为UTF8编码  Read string data, default is the UTF8 encoding (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](d07259dd-28d2-2ad5-f266-7d06eec03bbf.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [AllenBradleyNetReadString(String, UInt16, Encoding)](1bfd1b97-413e-55ef-6f0f-f16d8f6c843e.htm).) |
| 公共方法 | [ReadStringAsync(String)](9ae33180-776f-24ed-dc81-760e052d3cd4.htm) | (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](22b8c6d1-466e-3a4a-3bd3-28967e60345b.htm) | 读取字符串数据，默认为UTF8编码  Read string data, default is the UTF8 encoding (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](361b03e5-0b7a-381c-b9ac-e199303aae5d.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [AllenBradleyNetReadStringAsync(String, UInt16, Encoding)](4716e7f5-fc4b-7ea4-0871-169113220e8f.htm).) |
| 公共方法代码示例 | [ReadStructT](03fd2a05-8efd-8508-cdd7-564432b48c75.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (重写 [DeviceCommunicationReadStructT(String, UInt16)](49853082-fbad-52e5-d756-0615cedb4b83.htm).) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadTag](7ca6b7a4-ed19-fff4-9b55-48417d917aa2.htm) | 从PLC里读取一个指定标签名的原始数据信息及其数据类型信息  Read the original data information of a specified tag name and its data type information from the PLC (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadTagAsync](fb95e370-484e-c84e-d2a4-6b2e5e6ba44e.htm) | 从PLC里读取一个指定标签名的原始数据信息及其数据类型信息  Read the original data information of a specified tag name and its data type information from the PLC (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadTags](5ed6fae1-fe05-4913-6128-74d121b88fd6.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"， 长度为 [1, 5]，返回的是一整个的字节数组，需要自行解析 **[Authorization]** Read the data information of multiple addresses in batches. For example, I can read the data "A", "B[0]" of two tags, the length is [1, 5], and the return is an entire byte array, and I need to do it myself Parsing (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadTagsAsync](c01104dc-ea9d-b038-e377-8f446b28fa2f.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"， 长度为 [1, 5]，返回的是一整个的字节数组，需要自行解析 **[Authorization]** Read the data information of multiple addresses in batches. For example, I can read the data "A", "B[0]" of two tags, the length is [1, 5], and the return is an entire byte array, and I need to do it myself Parsing (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadTime](238554b9-f191-784d-6337-31b1c37176ce.htm) | 读取指定地址的时间数据，最小时间为 0，如果获取秒，可以访问 TotalSeconds，当PLC的变量类型为 "Time" 和 "TimeOfDate" 时，都可以用本方法读取。  Read the time data of the specified address. The minimum time is 0. If you get seconds, you can access TotalSeconds. When the PLC variable type is "Time" and "TimeOfDate", you can use this Method to read. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadTimeAsync](aa1f494a-b03a-3b7b-665f-febe3c229a8d.htm) | 读取指定地址的时间数据，最小时间为 0，如果获取秒，可以访问 TotalSeconds，当PLC的变量类型为 "Time" 和 "TimeOfDate" 时，都可以用本方法读取。  Read the time data of the specified address. The minimum time is 0. If you get seconds, you can access TotalSeconds. When the PLC variable type is "Time" and "TimeOfDate", you can use this Method to read. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](b649307b-d4cb-848c-e33d-06381bde093b.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](b1ec9b05-632c-f107-2536-10e04d4d1196.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](bb6af48c-cd4b-344e-72f1-589b844507dc.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](eef95827-5524-901d-64b9-928b08937f98.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](db1c90f9-6351-109f-7cf8-74c89983991e.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](c8edfb08-ddd5-2649-a47f-4fb1be933ed8.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [StructTagEnumerator](20fd92b3-9b3b-2166-f201-66d0ea183d60.htm) | 枚举结构体的方法，传入结构体的标签对象，返回结构体子属性标签列表信息，子属性有可能是标量数据，也可能是另一个结构体。  The method of enumerating the structure, passing in the tag object of the structure, and returning the tag list information of the sub-attributes of the structure. The sub-attributes may be scalar data or another structure. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [TagEnumerator](e5e08dcb-2ebc-d5ae-9e8b-3f4e6162e003.htm) | 枚举当前的所有的变量名字，包含结构体信息，除去系统自带的名称数据信息  Enumerate all the current variable names, including structure information, except the name data information that comes with the system (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [TagEnumeratorAsync](ac93bb16-8676-59ae-712e-7741917ab113.htm) | 枚举当前的所有的变量名字，包含结构体信息，除去系统自带的名称数据信息  Enumerate all the current variable names, including structure information, except the name data information that comes with the system (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ToString](96331804-b4e1-e30a-7493-2f9fe10ee376.htm) | (重写 [AllenBradleyNetToString](e02711ff-3ccb-66de-bbdb-90736220b140.htm).) |
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
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](47e90d71-05a2-8faa-c780-634036f04f16.htm) | 写入单个Bool的数据信息。如果读取的是单bool变量，就直接写变量名，如果是bool数组的一个值，一律带下标访问，例如a[0]  Write the data information of a single Bool. If the read is a single bool variable, write the variable name directly, if it is a value of the bool array, it will always be accessed with a subscript, such as a[0] (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Write(String, Boolean)](31242775-2e0c-2956-b095-d924dfa8ab09.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Write(String, Byte)](41c67fa5-e688-244b-8386-7e8264c67db5.htm) | 当前写入字节数组使用数据类型 0xD1 写入，如果其他的字节类型需要调用 [WriteTag(String, UInt16, Byte, Int32)](173f78d2-bc3f-bc03-96bf-b12c64e8f00e.htm) 方法来实现。  The currently written byte array is written using the data type 0xD1. If other byte types need to be called [WriteTag(String, UInt16, Byte, Int32)](173f78d2-bc3f-bc03-96bf-b12c64e8f00e.htm) Method to achieve. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](eb656f38-9fbe-6e0e-b9f7-e9b547b029d6.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](3f4c41f1-46c6-1451-066e-86e07f3bcdf6.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](4f3756f9-ea50-65d0-09a1-5748213e5482.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](66e66a64-a488-8499-e33b-e825a4a696e4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](32b3207a-a832-044c-1636-336be8c7caca.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](ef9780b5-0708-cac3-25bd-642bc42282e0.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](99e6113b-ccdc-4f2d-e35c-10ee5e3ce598.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](ebb019b9-0154-cac2-37fe-bfa15e77c071.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Write(String, Byte)](5e9c9bd6-abdf-eb6f-2512-672ae2787f4f.htm) | 写入Byte数据，返回是否写入成功，默认使用类型 0xC2, 如果PLC的变量类型不一样，则需要指定实际的变量类型，例如PLC的变量 A 是0xD1类型，那么地址需要携带类型信息，type=0xD1;A   Write Byte data and return whether the writing is successful. The default type is 0xC2. If the variable types of the PLC are different, you need to specify the actual variable type. For example, the variable A of the PLC is of type 0xD1, then the address needs to carry the type information, type= 0xD1;A (重写 [AllenBradleyNetWrite(String, Byte)](aaece82a-eb65-d5cf-163d-972c8cd293b0.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](bdee2591-29af-2e3e-e886-1e12e6ec0156.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (重写 [AllenBradleyNetWrite(String, String, Encoding)](c5c97f35-37a6-c1c2-33c9-ede5f40fef9e.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](a82d8ada-5c2d-046f-ae90-1d5dacdb636f.htm) | 写入单个Bool的数据信息。如果读取的是单bool变量，就直接写变量名，如果是bool数组的一个值，一律带下标访问，例如a[0]  Write the data information of a single Bool. If the read is a single bool variable, write the variable name directly, if it is a value of the bool array, it will always be accessed with a subscript, such as a[0] (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](19627866-fc0d-11a9-b7f6-614e93f7561f.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteAsync(String, Byte)](f772978a-f980-0b80-4bc6-20e012b15503.htm) | 当前写入字节数组使用数据类型 0xD1 写入，如果其他的字节类型需要调用 [WriteTag(String, UInt16, Byte, Int32)](173f78d2-bc3f-bc03-96bf-b12c64e8f00e.htm) 方法来实现。  The currently written byte array is written using the data type 0xD1. If other byte types need to be called [WriteTag(String, UInt16, Byte, Int32)](173f78d2-bc3f-bc03-96bf-b12c64e8f00e.htm) Method to achieve. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](c51b4f41-8ad9-3233-a814-3089a3703472.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](610a8de5-9fc1-c59c-4a98-16f7a8b9a4c8.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](418f77b5-475c-ac6a-a679-39bd2cbdf874.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4406a22e-239d-bbd1-a48a-b2b28919ac15.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](7201b04f-3748-74ce-a0fc-1d4b12ae1235.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](8325bb7f-59bf-9fb1-8fd9-59105a3abc50.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](acd95b54-10a6-2c08-faec-81170d1f78f0.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](505f97f3-fdd6-1428-cf87-a6bc00b5a9b9.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteAsync(String, Byte)](29742e84-ecc1-0569-edb1-339042b3fe41.htm) | 写入Byte数据，返回是否写入成功，默认使用类型 0xC2, 如果PLC的变量类型不一样，则需要指定实际的变量类型，例如PLC的变量 A 是0xD1类型，那么地址需要携带类型信息，type=0xD1;A   Write Byte data and return whether the writing is successful. The default type is 0xC2. If the variable types of the PLC are different, you need to specify the actual variable type. For example, the variable A of the PLC is of type 0xD1, then the address needs to carry the type information, type= 0xD1;A (重写 [AllenBradleyNetWriteAsync(String, Byte)](b7187f26-4770-2958-0e77-7050179e599e.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](3672d794-6b78-d944-b5df-3f4c64c038e3.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (重写 [AllenBradleyNetWriteAsync(String, String, Encoding)](b67045d7-f5f4-0f5b-bead-9a229e0b920d.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteDate](63de3b96-3b79-49a1-bd42-ee5b89a6918a.htm) | 使用日期格式（Date）将指定的数据写入到指定的地址里，PLC的地址类型变量必须为 "Date"，否则写入失败。  Use the date format (Date) to write the specified data to the specified address. The PLC address type variable must be "Date", otherwise the writing will fail. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteDateAsync](d96547ba-383a-90e5-4d55-1224084dac1f.htm) | 使用日期格式（Date）将指定的数据写入到指定的地址里，PLC的地址类型变量必须为 "Date"，否则写入失败。  Use the date format (Date) to write the specified data to the specified address. The PLC address type variable must be "Date", otherwise the writing will fail. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteTag](173f78d2-bc3f-bc03-96bf-b12c64e8f00e.htm) | 使用指定的类型写入指定的节点数据，类型信息参考API文档，地址支持协议类型代号信息，例如 "type=0xD1;A"  Use the specified type to write the specified node data. For type information, refer to the API documentation. The address supports protocol type code information, such as "type=0xD1;A" (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteTagAsync](1dd946b0-877f-2757-f82a-481289c7adfc.htm) | 使用指定的类型写入指定的节点数据，类型信息参考API文档，地址支持协议类型代号信息，例如 "type=0xD1;A"  Use the specified type to write the specified node data. For type information, refer to the API documentation. The address supports protocol type code information, such as "type=0xD1;A" (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteTime](37994d5b-c25c-fc59-d87f-dfa45d9979aa.htm) | 使用时间格式（TIME）将时间数据写入到PLC中指定的地址里去，PLC的地址类型变量必须为 "TIME"，否则写入失败。  Use the time format (TIME) to write the time data to the address specified in the PLC. The PLC address type variable must be "TIME", otherwise the writing will fail. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteTimeAndDate](50c0073b-efb2-acf5-b7ef-495f2f69dc13.htm) | 使用日期格式（Date）将指定的数据写入到指定的地址里，PLC的地址类型变量必须为 "Date"，否则写入失败。  Use the date format (Date) to write the specified data to the specified address. The PLC address type variable must be "Date", otherwise the writing will fail. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteTimeAndDateAsync](8b18984a-04c6-60c7-26c6-fbf5105528c9.htm) | 使用日期格式（Date）将指定的数据写入到指定的地址里，PLC的地址类型变量必须为 "Date"，否则写入失败。  Use the date format (Date) to write the specified data to the specified address. The PLC address type variable must be "Date", otherwise the writing will fail. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteTimeAsync](09fdb4d9-0989-5e8a-a1ac-f1b841ca3fee.htm) | 使用时间格式（TIME）将时间数据写入到PLC中指定的地址里去，PLC的地址类型变量必须为 "TIME"，否则写入失败。  Use the time format (TIME) to write the time data to the address specified in the PLC. The PLC address type variable must be "TIME", otherwise the writing will fail. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteTimeOfDate](ecbdec9f-7ed2-6fef-f443-b200be912592.htm) | 使用时间格式（TimeOfDate）将时间数据写入到PLC中指定的地址里去，PLC的地址类型变量必须为 "TimeOfDate"，否则写入失败。  Use the time format (TimeOfDate) to write the time data to the address specified in the PLC. The PLC address type variable must be "TimeOfDate", otherwise the writing will fail. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteTimeOfDateAsync](7d75bc08-c835-a669-f471-dc3bd8945c4d.htm) | 使用时间格式（TimeOfDate）将时间数据写入到PLC中指定的地址里去，PLC的地址类型变量必须为 "TimeOfDate"，否则写入失败。  Use the time format (TimeOfDate) to write the time data to the address specified in the PLC. The PLC address type variable must be "TimeOfDate", otherwise the writing will fail. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolWritePadding 方法 

[原文連結](http://api.hslcommunication.cn/html/41ba1de6-9932-39eb-2e3e-7b3867726987.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[GetBoolWritePadding 方法](../html/41ba1de6-9932-39eb-2e3e-7b3867726987.htm "GetBoolWritePadding 方法 ")

[GetWriteValueLength 方法](../html/3fc102d6-0067-c2c7-f4e7-f67879d2825e.htm "GetWriteValueLength 方法 ")

[Read 方法](../html/c3bc4e40-cb07-a183-5cc3-37dfca2727ae.htm "Read 方法 ")

[ReadAsync 方法](../html/e3167e9d-439a-6ef2-e03f-9f4898ea526c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/08b48880-ba18-76eb-009f-5bf3320e6b47.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/696b89a0-1a79-589c-e14e-75a7a9b26af9.htm "ReadBoolAsync 方法 ")

[ReadString 方法](../html/c72f72db-c199-9b51-d3bf-3d93e87ed17a.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/67d00847-5a0e-781f-3f2b-7f3f9752b8b5.htm "ReadStringAsync 方法 ")

[ReadStruct(T) 方法](../html/03fd2a05-8efd-8508-cdd7-564432b48c75.htm "ReadStruct(T) 方法 ")

[ToString 方法](../html/96331804-b4e1-e30a-7493-2f9fe10ee376.htm "ToString 方法 ")

[Write 方法](../html/9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm "Write 方法 ")

[WriteAsync 方法](../html/558f1912-35b3-64b1-eca5-833d440e5b2f.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetGetBoolWritePadding 方法 |

当进行bool写入的时候，是否需要补齐到字节长度，默认不需要  
When performing bool writing, whether it is necessary to complement to the byte length is not required by default

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override bool GetBoolWritePadding()
```

```
Protected Overrides Function GetBoolWritePadding As Boolean
```

```
protected:
virtual bool GetBoolWritePadding() override
```

```
abstract GetBoolWritePadding : unit -> bool 
override GetBoolWritePadding : unit -> bool
```

#### 返回值

类型：Boolean  
是否需要字节对齐

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetWriteValueLength 方法 

[原文連結](http://api.hslcommunication.cn/html/3fc102d6-0067-c2c7-f4e7-f67879d2825e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[GetBoolWritePadding 方法](../html/41ba1de6-9932-39eb-2e3e-7b3867726987.htm "GetBoolWritePadding 方法 ")

[GetWriteValueLength 方法](../html/3fc102d6-0067-c2c7-f4e7-f67879d2825e.htm "GetWriteValueLength 方法 ")

[Read 方法](../html/c3bc4e40-cb07-a183-5cc3-37dfca2727ae.htm "Read 方法 ")

[ReadAsync 方法](../html/e3167e9d-439a-6ef2-e03f-9f4898ea526c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/08b48880-ba18-76eb-009f-5bf3320e6b47.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/696b89a0-1a79-589c-e14e-75a7a9b26af9.htm "ReadBoolAsync 方法 ")

[ReadString 方法](../html/c72f72db-c199-9b51-d3bf-3d93e87ed17a.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/67d00847-5a0e-781f-3f2b-7f3f9752b8b5.htm "ReadStringAsync 方法 ")

[ReadStruct(T) 方法](../html/03fd2a05-8efd-8508-cdd7-564432b48c75.htm "ReadStruct(T) 方法 ")

[ToString 方法](../html/96331804-b4e1-e30a-7493-2f9fe10ee376.htm "ToString 方法 ")

[Write 方法](../html/9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm "Write 方法 ")

[WriteAsync 方法](../html/558f1912-35b3-64b1-eca5-833d440e5b2f.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetGetWriteValueLength 方法 |

获取写入数据的长度信息，此处直接返回数组的长度信息

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override int GetWriteValueLength(
	string address,
	int length
)
```

```
Protected Overrides Function GetWriteValueLength ( 
	address As String,
	length As Integer
) As Integer
```

```
protected:
virtual int GetWriteValueLength(
	String^ address, 
	int length
) override
```

```
abstract GetWriteValueLength : 
        address : string * 
        length : int -> int 
override GetWriteValueLength : 
        address : string * 
        length : int -> int
```

#### 参数

address
:   类型：SystemString  
    地址信息

length
:   类型：SystemInt32  
    数组长度信息

#### 返回值

类型：Int32  
实际的写入长度信息

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/c3bc4e40-cb07-a183-5cc3-37dfca2727ae.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[Read 方法](../html/c3bc4e40-cb07-a183-5cc3-37dfca2727ae.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/2954efdc-df90-5d0f-71ca-fc1f5d94a0c1.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetRead 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Read(String)](682e8e26-45f6-f437-bbe2-a345db989fb8.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"，每个地址的数据长度为1，表示一个数据，最终读取返回的是一整个的字节数组，需要自行解析 **[Authorization]** Batch read data information of multiple addresses, for example, I can read the data of two tags "A", "B[0]", the data length of each address is 1, which means one data, and the final read returns a The entire byte array, which needs to be parsed by itself (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Read(String, UInt16)](a6216184-0009-e4a9-96cd-a2158390aa44.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"， 长度为 [1, 5]，返回的是一整个的字节数组，需要自行解析 **[Authorization]** Read the data information of multiple addresses in batches. For example, I can read the data "A", "B[0]" of two tags, the length is [1, 5], and the return is an entire byte array, and I need to do it myself Parsing (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Read(String, UInt16)](2954efdc-df90-5d0f-71ca-fc1f5d94a0c1.htm) | 读取指定地址的二进制数据内容，长度为地址长度，一般都是1，除非读取数组时，如果需要强制使用 片段读取功能码，则地址里携带 x=0x52; 或是 x=82; 则强制使用片段读取。  Read the binary data content of the specified address, the length is the address length, generally 1, unless the array is read, if you need to force the fragment reading function code, the address carries x=0x52; or x=82; then the fragment read is forced. (重写 [AllenBradleyNetRead(String, UInt16)](cc6bd266-1e7b-5052-cdf3-b5e45abf6b48.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/2954efdc-df90-5d0f-71ca-fc1f5d94a0c1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[Read 方法](../html/c3bc4e40-cb07-a183-5cc3-37dfca2727ae.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/2954efdc-df90-5d0f-71ca-fc1f5d94a0c1.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetRead 方法 (String, UInt16) |

读取指定地址的二进制数据内容，长度为地址长度，一般都是1，除非读取数组时，如果需要强制使用 片段读取功能码，则地址里携带 x=0x52; 或是 x=82; 则强制使用片段读取。  
Read the binary data content of the specified address, the length is the address length, generally 1, unless the array is read, if you need to force the fragment reading function code,
the address carries x=0x52; or x=82; then the fragment read is forced.

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
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
    Address format of the node

length
:   类型：SystemUInt16  
    In the case of arrays, the length of the array

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
Result data with result object

#### 实现

[IReadWriteNetRead(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm)  
[IReadWriteNetRead(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm)

![](../icons/SectionExpanded.png)备注

使用片段读取的时候，可以读取一些数量量非常大的地址，例如一个结构体标签有100\_000个字节长度的时候。  
When using fragment reading, you can read a very large number of addresses, such as when a struct tag is 100\_000 bytes long.

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[Read 重载](c3bc4e40-cb07-a183-5cc3-37dfca2727ae.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/e3167e9d-439a-6ef2-e03f-9f4898ea526c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadAsync 方法](../html/e3167e9d-439a-6ef2-e03f-9f4898ea526c.htm "ReadAsync 方法 ")

[ReadAsync 方法 (String, UInt16)](../html/ada9961d-130c-7c7b-e883-5a229b1fbbed.htm "ReadAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String)](d666e005-8c96-6381-fa24-60fd58545d95.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"，每个地址的数据长度为1，表示一个数据，最终读取返回的是一整个的字节数组，需要自行解析 **[Authorization]** Batch read data information of multiple addresses, for example, I can read the data of two tags "A", "B[0]", the data length of each address is 1, which means one data, and the final read returns a The entire byte array, which needs to be parsed by itself (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](856714b2-0b15-fe6e-253f-201b6bd0b010.htm) | **[商业授权]** 批量读取多地址的数据信息，例如我可以读取两个标签的数据 "A","B[0]"， 长度为 [1, 5]，返回的是一整个的字节数组，需要自行解析 **[Authorization]** Read the data information of multiple addresses in batches. For example, I can read the data "A", "B[0]" of two tags, the length is [1, 5], and the return is an entire byte array, and I need to do it myself Parsing (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](ada9961d-130c-7c7b-e883-5a229b1fbbed.htm) | 读取指定地址的二进制数据内容，长度为地址长度，一般都是1，除非读取数组时，如果需要强制使用 片段读取功能码，则地址里携带 x=0x52; 或是 x=82; 则强制使用片段读取。  Read the binary data content of the specified address, the length is the address length, generally 1, unless the array is read, if you need to force the fragment reading function code, the address carries x=0x52; or x=82; then the fragment read is forced. (重写 [AllenBradleyNetReadAsync(String, UInt16)](12fbde37-b7d0-54e4-4196-fa53d1976d16.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/ada9961d-130c-7c7b-e883-5a229b1fbbed.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadAsync 方法](../html/e3167e9d-439a-6ef2-e03f-9f4898ea526c.htm "ReadAsync 方法 ")

[ReadAsync 方法 (String, UInt16)](../html/ada9961d-130c-7c7b-e883-5a229b1fbbed.htm "ReadAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadAsync 方法 (String, UInt16) |

读取指定地址的二进制数据内容，长度为地址长度，一般都是1，除非读取数组时，如果需要强制使用 片段读取功能码，则地址里携带 x=0x52; 或是 x=82; 则强制使用片段读取。  
Read the binary data content of the specified address, the length is the address length, generally 1, unless the array is read, if you need to force the fragment reading function code,
the address carries x=0x52; or x=82; then the fragment read is forced.

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
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
    Address format of the node

length
:   类型：SystemUInt16  
    In the case of arrays, the length of the array

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
Result data with result object

#### 实现

[IReadWriteNetReadAsync(String, UInt16)](59839c27-5a76-2f53-730c-c9d2bb41c2c4.htm)  
[IReadWriteNetReadAsync(String, UInt16)](59839c27-5a76-2f53-730c-c9d2bb41c2c4.htm)

![](../icons/SectionExpanded.png)备注

使用片段读取的时候，可以读取一些数量量非常大的地址，例如一个结构体标签有100\_000个字节长度的时候。  
When using fragment reading, you can read a very large number of addresses, such as when a struct tag is 100\_000 bytes long.

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[ReadAsync 重载](e3167e9d-439a-6ef2-e03f-9f4898ea526c.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/08b48880-ba18-76eb-009f-5bf3320e6b47.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadBool 方法](../html/08b48880-ba18-76eb-009f-5bf3320e6b47.htm "ReadBool 方法 ")

[ReadBool 方法 (String, UInt16)](../html/f6c36d10-c3a9-24cf-e0a1-9eac0b791c04.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadBool 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadBool(String)](3d12b2f9-19b8-28ea-a1e2-5d8d7c18bf2d.htm) | 读取单个的bool数据信息，如果读取的是单bool变量，就直接写变量名，如果是由int组成的bool数组的一个值，一律带"i="开头访问，例如"i=A[0]"   Read a single bool data information, if it is a single bool variable, write the variable name directly, if it is a value of a bool array composed of int, it is always accessed with "i=" at the beginning, for example, "i=A[0]" (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](f6c36d10-c3a9-24cf-e0a1-9eac0b791c04.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (重写 [AllenBradleyNetReadBool(String, UInt16)](bae5fa6c-770b-f495-a578-15bda70b6449.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/f6c36d10-c3a9-24cf-e0a1-9eac0b791c04.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadBool 方法](../html/08b48880-ba18-76eb-009f-5bf3320e6b47.htm "ReadBool 方法 ")

[ReadBool 方法 (String, UInt16)](../html/f6c36d10-c3a9-24cf-e0a1-9eac0b791c04.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadBool 方法 (String, UInt16) |

批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  
Batch read Boolean array information, need to specify the address and length, return Boolean array

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
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
带有成功标识的 bool[] 数组

#### 实现

[IReadWriteNetReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm)  
[IReadWriteNetReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[ReadBool 重载](08b48880-ba18-76eb-009f-5bf3320e6b47.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/696b89a0-1a79-589c-e14e-75a7a9b26af9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadBoolAsync 方法](../html/696b89a0-1a79-589c-e14e-75a7a9b26af9.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (String, UInt16)](../html/2ff1bdf7-1827-1e13-adbc-54c939c9758b.htm "ReadBoolAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadBoolAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadBoolAsync(String)](28957b24-aa96-b3f4-6187-4873cbd9c23b.htm) | 读取单个的bool数据信息，如果读取的是单bool变量，就直接写变量名，如果是由int组成的bool数组的一个值，一律带"i="开头访问，例如"i=A[0]"   Read a single bool data information, if it is a single bool variable, write the variable name directly, if it is a value of a bool array composed of int, it is always accessed with "i=" at the beginning, for example, "i=A[0]" (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](2ff1bdf7-1827-1e13-adbc-54c939c9758b.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (重写 [AllenBradleyNetReadBoolAsync(String, UInt16)](e0d9e8f4-4c9f-f7a7-8d7c-9fada2ce8e66.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/2ff1bdf7-1827-1e13-adbc-54c939c9758b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadBoolAsync 方法](../html/696b89a0-1a79-589c-e14e-75a7a9b26af9.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (String, UInt16)](../html/2ff1bdf7-1827-1e13-adbc-54c939c9758b.htm "ReadBoolAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadBoolAsync 方法 (String, UInt16) |

异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  
Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
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
[IReadWriteNetReadBoolAsync(String, UInt16)](74e40b23-198a-944b-6ed7-3f58ca51da5d.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[ReadBoolAsync 重载](696b89a0-1a79-589c-e14e-75a7a9b26af9.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadString 方法 

[原文連結](http://api.hslcommunication.cn/html/c72f72db-c199-9b51-d3bf-3d93e87ed17a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadString 方法](../html/c72f72db-c199-9b51-d3bf-3d93e87ed17a.htm "ReadString 方法 ")

[ReadString 方法 (String, UInt16, Encoding)](../html/d07259dd-28d2-2ad5-f266-7d06eec03bbf.htm "ReadString 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadString 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadString(String)](bd50aa1d-e5d1-df14-61e3-6341d1ce151b.htm) | (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](bf38a088-492e-4c7c-2f29-dd56d4230bbc.htm) | 读取字符串数据，默认为UTF8编码  Read string data, default is the UTF8 encoding (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](d07259dd-28d2-2ad5-f266-7d06eec03bbf.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [AllenBradleyNetReadString(String, UInt16, Encoding)](1bfd1b97-413e-55ef-6f0f-f16d8f6c843e.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadString 方法 (String, UInt16, Encoding)

[原文連結](http://api.hslcommunication.cn/html/d07259dd-28d2-2ad5-f266-7d06eec03bbf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadString 方法](../html/c72f72db-c199-9b51-d3bf-3d93e87ed17a.htm "ReadString 方法 ")

[ReadString 方法 (String, UInt16, Encoding)](../html/d07259dd-28d2-2ad5-f266-7d06eec03bbf.htm "ReadString 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadString 方法 (String, UInt16, Encoding) |

使用指定的编码，读取字符串数据  
Reads string data using the specified encoding

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<string> ReadString(
	string address,
	ushort length,
	Encoding encoding
)
```

```
Public Overrides Function ReadString ( 
	address As String,
	length As UShort,
	encoding As Encoding
) As OperateResult(Of String)
```

```
public:
virtual OperateResult<String^>^ ReadString(
	String^ address, 
	unsigned short length, 
	Encoding^ encoding
) override
```

```
abstract ReadString : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> OperateResult<string> 
override ReadString : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> OperateResult<string>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数据长度

encoding
:   类型：System.TextEncoding  
    指定的自定义的编码

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有成功标识的string数据

#### 实现

[IReadWriteNetReadString(String, UInt16, Encoding)](fc289def-e53b-b47d-22fe-7e02b5e354e0.htm)  
[IReadWriteNetReadString(String, UInt16, Encoding)](fc289def-e53b-b47d-22fe-7e02b5e354e0.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

String类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
string d100_value = melsec_net.ReadString( "D100", 5 ).Content;

// 如果需要判断是否读取成功，使用 Unicode 编码即可读取中文，如果还是乱码，就需要自己指定编码来实现
OperateResult<string> R_d100_value = melsec_net.ReadString( "D100", 5, Encoding.Unicode );
if (R_d100_value.IsSuccess)
{
    // success
    string value = R_d100_value.Content;
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[ReadString 重载](c72f72db-c199-9b51-d3bf-3d93e87ed17a.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/67d00847-5a0e-781f-3f2b-7f3f9752b8b5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadStringAsync 方法](../html/67d00847-5a0e-781f-3f2b-7f3f9752b8b5.htm "ReadStringAsync 方法 ")

[ReadStringAsync 方法 (String, UInt16, Encoding)](../html/361b03e5-0b7a-381c-b9ac-e199303aae5d.htm "ReadStringAsync 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadStringAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadStringAsync(String)](9ae33180-776f-24ed-dc81-760e052d3cd4.htm) | (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](22b8c6d1-466e-3a4a-3bd3-28967e60345b.htm) | 读取字符串数据，默认为UTF8编码  Read string data, default is the UTF8 encoding (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](361b03e5-0b7a-381c-b9ac-e199303aae5d.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [AllenBradleyNetReadStringAsync(String, UInt16, Encoding)](4716e7f5-fc4b-7ea4-0871-169113220e8f.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringAsync 方法 (String, UInt16, Encoding)

[原文連結](http://api.hslcommunication.cn/html/361b03e5-0b7a-381c-b9ac-e199303aae5d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[ReadStringAsync 方法](../html/67d00847-5a0e-781f-3f2b-7f3f9752b8b5.htm "ReadStringAsync 方法 ")

[ReadStringAsync 方法 (String, UInt16, Encoding)](../html/361b03e5-0b7a-381c-b9ac-e199303aae5d.htm "ReadStringAsync 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadStringAsync 方法 (String, UInt16, Encoding) |

异步使用指定的编码，读取字符串数据  
Asynchronously reads string data using the specified encoding

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<string>> ReadStringAsync(
	string address,
	ushort length,
	Encoding encoding
)
```

```
Public Overrides Function ReadStringAsync ( 
	address As String,
	length As UShort,
	encoding As Encoding
) As Task(Of OperateResult(Of String))
```

```
public:
virtual Task<OperateResult<String^>^>^ ReadStringAsync(
	String^ address, 
	unsigned short length, 
	Encoding^ encoding
) override
```

```
abstract ReadStringAsync : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> Task<OperateResult<string>> 
override ReadStringAsync : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> Task<OperateResult<string>>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数据长度

encoding
:   类型：System.TextEncoding  
    指定的自定义的编码

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有成功标识的string数据

#### 实现

[IReadWriteNetReadStringAsync(String, UInt16, Encoding)](db06baf4-e674-22de-e89e-cb26646413f8.htm)  
[IReadWriteNetReadStringAsync(String, UInt16, Encoding)](db06baf4-e674-22de-e89e-cb26646413f8.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

String类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
string d100_value = melsec_net.ReadString( "D100", 5 ).Content;

// 如果需要判断是否读取成功，使用 Unicode 编码即可读取中文，如果还是乱码，就需要自己指定编码来实现
OperateResult<string> R_d100_value = melsec_net.ReadString( "D100", 5, Encoding.Unicode );
if (R_d100_value.IsSuccess)
{
    // success
    string value = R_d100_value.Content;
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[ReadStringAsync 重载](67d00847-5a0e-781f-3f2b-7f3f9752b8b5.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStruct(T) 方法 

[原文連結](http://api.hslcommunication.cn/html/03fd2a05-8efd-8508-cdd7-564432b48c75.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[GetBoolWritePadding 方法](../html/41ba1de6-9932-39eb-2e3e-7b3867726987.htm "GetBoolWritePadding 方法 ")

[GetWriteValueLength 方法](../html/3fc102d6-0067-c2c7-f4e7-f67879d2825e.htm "GetWriteValueLength 方法 ")

[Read 方法](../html/c3bc4e40-cb07-a183-5cc3-37dfca2727ae.htm "Read 方法 ")

[ReadAsync 方法](../html/e3167e9d-439a-6ef2-e03f-9f4898ea526c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/08b48880-ba18-76eb-009f-5bf3320e6b47.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/696b89a0-1a79-589c-e14e-75a7a9b26af9.htm "ReadBoolAsync 方法 ")

[ReadString 方法](../html/c72f72db-c199-9b51-d3bf-3d93e87ed17a.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/67d00847-5a0e-781f-3f2b-7f3f9752b8b5.htm "ReadStringAsync 方法 ")

[ReadStruct(T) 方法](../html/03fd2a05-8efd-8508-cdd7-564432b48c75.htm "ReadStruct(T) 方法 ")

[ToString 方法](../html/96331804-b4e1-e30a-7493-2f9fe10ee376.htm "ToString 方法 ")

[Write 方法](../html/9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm "Write 方法 ")

[WriteAsync 方法](../html/558f1912-35b3-64b1-eca5-833d440e5b2f.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetReadStructT 方法 |

读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<T> ReadStruct<T>(
	string address,
	ushort length
)
where T : class, new()
```

```
Public Overrides Function ReadStruct(Of T As {Class, New}) ( 
	address As String,
	length As UShort
) As OperateResult(Of T)
```

```
public:
generic<typename T>
where T : ref class, gcnew()
virtual OperateResult<T>^ ReadStruct(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadStruct : 
        address : string * 
        length : uint16 -> OperateResult<'T>  when 'T : not struct, new()
override ReadStruct : 
        address : string * 
        length : uint16 -> OperateResult<'T>  when 'T : not struct, new()
```

#### 参数

address
:   类型：SystemString  
    PLC的地址信息

length
:   类型：SystemUInt16  
    读取的地址长度信息

#### 类型参数

T
:   类型对象信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)T  
如果成功，返回成功的结构体对象

#### 实现

[IReadWriteNetReadStructT(String, UInt16)](50df0eb8-8f9c-ad9c-4dc9-bf75b4438d46.htm)

![](../icons/SectionExpanded.png)示例

此处演示西门子的读取示例，先定义一个类，重点是将需要读取的数据，写入到属性的特性中去。

特性实现示例

[复制](# "复制")

```
// 假设你要读取连续地址数据的情况，我们把需要读取的数据解析成一个个不同的数据类型，则可以使用本方法。
// 我们假设，我们要读取的PLC是西门子PLC，地址数据的假设如下
// 我们假设 设备是否启动是 M0.0
// 产量是 M10 开始的2个地址数据
// 温度信息是 M12开始的4个地址数据
// 报警的IO信息是 M20 开始，5个字节，共计40个IO点信息
// 那么我们可以做如下的定义

public class StructExample
{
    /// <summary>
    /// 设备是否启动
    /// </summary>
    [HslStruct( 0 )]
    public bool Enable { get; set; }

    /// <summary>
    /// 产量信息
    /// </summary>
    [HslStruct( 10 )]
    public short Production { get; set; }

    /// <summary>
    /// 温度信息
    /// </summary>
    [HslStruct( 12 )]
    public float Temperature { get; set; }

    /// <summary>
    /// 连续的位报警信息
    /// </summary>
    [HslStruct( 20 * 8, 40 )]
    public bool[] AlarmStatus { get; set; }
}
```

接下来就可以实现数据的读取了

ReadStruct示例

[复制](# "复制")

```
SiemensS7Net plc = new SiemensS7Net( SiemensPLCS.S1200, "192.168.0.100" );

// 此处需要注意的是，凡是带有 HslStruct 特性的属性都会被读取出来
OperateResult<StructExample> read = plc.ReadStruct<StructExample>( "M0", 30 );
if (read.IsSuccess)
{
    // success
    StructExample data = read.Content;
}
else
{
    // failed
    Console.WriteLine( "读取失败：" + read.Message );
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/96331804-b4e1-e30a-7493-2f9fe10ee376.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[GetBoolWritePadding 方法](../html/41ba1de6-9932-39eb-2e3e-7b3867726987.htm "GetBoolWritePadding 方法 ")

[GetWriteValueLength 方法](../html/3fc102d6-0067-c2c7-f4e7-f67879d2825e.htm "GetWriteValueLength 方法 ")

[Read 方法](../html/c3bc4e40-cb07-a183-5cc3-37dfca2727ae.htm "Read 方法 ")

[ReadAsync 方法](../html/e3167e9d-439a-6ef2-e03f-9f4898ea526c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/08b48880-ba18-76eb-009f-5bf3320e6b47.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/696b89a0-1a79-589c-e14e-75a7a9b26af9.htm "ReadBoolAsync 方法 ")

[ReadString 方法](../html/c72f72db-c199-9b51-d3bf-3d93e87ed17a.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/67d00847-5a0e-781f-3f2b-7f3f9752b8b5.htm "ReadStringAsync 方法 ")

[ReadStruct(T) 方法](../html/03fd2a05-8efd-8508-cdd7-564432b48c75.htm "ReadStruct(T) 方法 ")

[ToString 方法](../html/96331804-b4e1-e30a-7493-2f9fe10ee376.htm "ToString 方法 ")

[Write 方法](../html/9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm "Write 方法 ")

[WriteAsync 方法](../html/558f1912-35b3-64b1-eca5-833d440e5b2f.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetToString 方法 |

[缺少 "M:HslCommunication.Profinet.Omron.OmronCipNet.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
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

[缺少 "M:HslCommunication.Profinet.Omron.OmronCipNet.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[Write 方法](../html/9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm "Write 方法 ")

[Write 方法 (String, Byte)](../html/5e9c9bd6-abdf-eb6f-2512-672ae2787f4f.htm "Write 方法 (String, Byte)")

[Write 方法 (String, String, Encoding)](../html/bdee2591-29af-2e3e-e886-1e12e6ec0156.htm "Write 方法 (String, String, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](47e90d71-05a2-8faa-c780-634036f04f16.htm) | 写入单个Bool的数据信息。如果读取的是单bool变量，就直接写变量名，如果是bool数组的一个值，一律带下标访问，例如a[0]  Write the data information of a single Bool. If the read is a single bool variable, write the variable name directly, if it is a value of the bool array, it will always be accessed with a subscript, such as a[0] (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Write(String, Boolean)](31242775-2e0c-2956-b095-d924dfa8ab09.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Write(String, Byte)](41c67fa5-e688-244b-8386-7e8264c67db5.htm) | 当前写入字节数组使用数据类型 0xD1 写入，如果其他的字节类型需要调用 [WriteTag(String, UInt16, Byte, Int32)](173f78d2-bc3f-bc03-96bf-b12c64e8f00e.htm) 方法来实现。  The currently written byte array is written using the data type 0xD1. If other byte types need to be called [WriteTag(String, UInt16, Byte, Int32)](173f78d2-bc3f-bc03-96bf-b12c64e8f00e.htm) Method to achieve. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](eb656f38-9fbe-6e0e-b9f7-e9b547b029d6.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](3f4c41f1-46c6-1451-066e-86e07f3bcdf6.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](4f3756f9-ea50-65d0-09a1-5748213e5482.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](66e66a64-a488-8499-e33b-e825a4a696e4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](32b3207a-a832-044c-1636-336be8c7caca.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](ef9780b5-0708-cac3-25bd-642bc42282e0.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](99e6113b-ccdc-4f2d-e35c-10ee5e3ce598.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](ebb019b9-0154-cac2-37fe-bfa15e77c071.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [Write(String, Byte)](5e9c9bd6-abdf-eb6f-2512-672ae2787f4f.htm) | 写入Byte数据，返回是否写入成功，默认使用类型 0xC2, 如果PLC的变量类型不一样，则需要指定实际的变量类型，例如PLC的变量 A 是0xD1类型，那么地址需要携带类型信息，type=0xD1;A   Write Byte data and return whether the writing is successful. The default type is 0xC2. If the variable types of the PLC are different, you need to specify the actual variable type. For example, the variable A of the PLC is of type 0xD1, then the address needs to carry the type information, type= 0xD1;A (重写 [AllenBradleyNetWrite(String, Byte)](aaece82a-eb65-d5cf-163d-972c8cd293b0.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](bdee2591-29af-2e3e-e886-1e12e6ec0156.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (重写 [AllenBradleyNetWrite(String, String, Encoding)](c5c97f35-37a6-c1c2-33c9-ede5f40fef9e.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (String, Byte)

[原文連結](http://api.hslcommunication.cn/html/5e9c9bd6-abdf-eb6f-2512-672ae2787f4f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[Write 方法](../html/9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm "Write 方法 ")

[Write 方法 (String, Byte)](../html/5e9c9bd6-abdf-eb6f-2512-672ae2787f4f.htm "Write 方法 (String, Byte)")

[Write 方法 (String, String, Encoding)](../html/bdee2591-29af-2e3e-e886-1e12e6ec0156.htm "Write 方法 (String, String, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetWrite 方法 (String, Byte) |

写入Byte数据，返回是否写入成功，默认使用类型 0xC2, 如果PLC的变量类型不一样，则需要指定实际的变量类型，例如PLC的变量 A 是0xD1类型，那么地址需要携带类型信息，type=0xD1;A   
Write Byte data and return whether the writing is successful. The default type is 0xC2. If the variable types of the PLC are different, you need to specify the actual variable type.
For example, the variable A of the PLC is of type 0xD1, then the address needs to carry the type information, type= 0xD1;A

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
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
	byte value
)
```

```
Public Overrides Function Write ( 
	address As String,
	value As Byte
) As OperateResult
```

```
public:
virtual OperateResult^ Write(
	String^ address, 
	unsigned char value
) override
```

```
abstract Write : 
        address : string * 
        value : byte -> OperateResult 
override Write : 
        address : string * 
        value : byte -> OperateResult
```

#### 参数

address
:   类型：SystemString  
    标签的地址数据

value
:   类型：SystemByte  
    Byte数据

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)备注

如何确认PLC的变量的类型呢？可以在HslCommunicationDemo程序上测试知道，也可以直接调用 ReadWithType(String, UInt16) 来知道类型信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[Write 重载](9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (String, String, Encoding)

[原文連結](http://api.hslcommunication.cn/html/bdee2591-29af-2e3e-e886-1e12e6ec0156.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[Write 方法](../html/9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm "Write 方法 ")

[Write 方法 (String, Byte)](../html/5e9c9bd6-abdf-eb6f-2512-672ae2787f4f.htm "Write 方法 (String, Byte)")

[Write 方法 (String, String, Encoding)](../html/bdee2591-29af-2e3e-e886-1e12e6ec0156.htm "Write 方法 (String, String, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetWrite 方法 (String, String, Encoding) |

写入字符串信息，需要指定的编码信息  
Write string information, need to specify the encoding information

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
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
	string value,
	Encoding encoding
)
```

```
Public Overrides Function Write ( 
	address As String,
	value As String,
	encoding As Encoding
) As OperateResult
```

```
public:
virtual OperateResult^ Write(
	String^ address, 
	String^ value, 
	Encoding^ encoding
) override
```

```
abstract Write : 
        address : string * 
        value : string * 
        encoding : Encoding -> OperateResult 
override Write : 
        address : string * 
        value : string * 
        encoding : Encoding -> OperateResult
```

#### 参数

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemString  
    写入值

encoding
:   类型：System.TextEncoding  
    指定的编码信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
带有成功标识的结果类对象

#### 实现

[IReadWriteNetWrite(String, String, Encoding)](c66a9dfb-840a-3ac4-96ef-48e23fb85ec1.htm)  
[IReadWriteNetWrite(String, String, Encoding)](c66a9dfb-840a-3ac4-96ef-48e23fb85ec1.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

String类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 简单的写入
melsec_net.Write( "D100", "ABCDEFGH" );

// 如果想要判断是否写入成功
OperateResult write = melsec_net.Write( "D100", "ABCDEFGH" );
if (write.IsSuccess)
{
    // success
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[Write 重载](9b749bbf-e11e-8d74-4b59-f3b48cf760da.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/558f1912-35b3-64b1-eca5-833d440e5b2f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[WriteAsync 方法](../html/558f1912-35b3-64b1-eca5-833d440e5b2f.htm "WriteAsync 方法 ")

[WriteAsync 方法 (String, Byte)](../html/29742e84-ecc1-0569-edb1-339042b3fe41.htm "WriteAsync 方法 (String, Byte)")

[WriteAsync 方法 (String, String, Encoding)](../html/3672d794-6b78-d944-b5df-3f4c64c038e3.htm "WriteAsync 方法 (String, String, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetWriteAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](a82d8ada-5c2d-046f-ae90-1d5dacdb636f.htm) | 写入单个Bool的数据信息。如果读取的是单bool变量，就直接写变量名，如果是bool数组的一个值，一律带下标访问，例如a[0]  Write the data information of a single Bool. If the read is a single bool variable, write the variable name directly, if it is a value of the bool array, it will always be accessed with a subscript, such as a[0] (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](19627866-fc0d-11a9-b7f6-614e93f7561f.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteAsync(String, Byte)](f772978a-f980-0b80-4bc6-20e012b15503.htm) | 当前写入字节数组使用数据类型 0xD1 写入，如果其他的字节类型需要调用 [WriteTag(String, UInt16, Byte, Int32)](173f78d2-bc3f-bc03-96bf-b12c64e8f00e.htm) 方法来实现。  The currently written byte array is written using the data type 0xD1. If other byte types need to be called [WriteTag(String, UInt16, Byte, Int32)](173f78d2-bc3f-bc03-96bf-b12c64e8f00e.htm) Method to achieve. (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](c51b4f41-8ad9-3233-a814-3089a3703472.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](610a8de5-9fc1-c59c-4a98-16f7a8b9a4c8.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](418f77b5-475c-ac6a-a679-39bd2cbdf874.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4406a22e-239d-bbd1-a48a-b2b28919ac15.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](7201b04f-3748-74ce-a0fc-1d4b12ae1235.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](8325bb7f-59bf-9fb1-8fd9-59105a3abc50.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](acd95b54-10a6-2c08-faec-81170d1f78f0.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](505f97f3-fdd6-1428-cf87-a6bc00b5a9b9.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)。) |
| 公共方法 | [WriteAsync(String, Byte)](29742e84-ecc1-0569-edb1-339042b3fe41.htm) | 写入Byte数据，返回是否写入成功，默认使用类型 0xC2, 如果PLC的变量类型不一样，则需要指定实际的变量类型，例如PLC的变量 A 是0xD1类型，那么地址需要携带类型信息，type=0xD1;A   Write Byte data and return whether the writing is successful. The default type is 0xC2. If the variable types of the PLC are different, you need to specify the actual variable type. For example, the variable A of the PLC is of type 0xD1, then the address needs to carry the type information, type= 0xD1;A (重写 [AllenBradleyNetWriteAsync(String, Byte)](b7187f26-4770-2958-0e77-7050179e599e.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](3672d794-6b78-d944-b5df-3f4c64c038e3.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (重写 [AllenBradleyNetWriteAsync(String, String, Encoding)](b67045d7-f5f4-0f5b-bead-9a229e0b920d.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (String, Byte)

[原文連結](http://api.hslcommunication.cn/html/29742e84-ecc1-0569-edb1-339042b3fe41.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[WriteAsync 方法](../html/558f1912-35b3-64b1-eca5-833d440e5b2f.htm "WriteAsync 方法 ")

[WriteAsync 方法 (String, Byte)](../html/29742e84-ecc1-0569-edb1-339042b3fe41.htm "WriteAsync 方法 (String, Byte)")

[WriteAsync 方法 (String, String, Encoding)](../html/3672d794-6b78-d944-b5df-3f4c64c038e3.htm "WriteAsync 方法 (String, String, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetWriteAsync 方法 (String, Byte) |

写入Byte数据，返回是否写入成功，默认使用类型 0xC2, 如果PLC的变量类型不一样，则需要指定实际的变量类型，例如PLC的变量 A 是0xD1类型，那么地址需要携带类型信息，type=0xD1;A   
Write Byte data and return whether the writing is successful. The default type is 0xC2. If the variable types of the PLC are different, you need to specify the actual variable type.
For example, the variable A of the PLC is of type 0xD1, then the address needs to carry the type information, type= 0xD1;A

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
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
	byte value
)
```

```
Public Overrides Function WriteAsync ( 
	address As String,
	value As Byte
) As Task(Of OperateResult)
```

```
public:
virtual Task<OperateResult^>^ WriteAsync(
	String^ address, 
	unsigned char value
) override
```

```
abstract WriteAsync : 
        address : string * 
        value : byte -> Task<OperateResult> 
override WriteAsync : 
        address : string * 
        value : byte -> Task<OperateResult>
```

#### 参数

address
:   类型：SystemString  
    标签的地址数据

value
:   类型：SystemByte  
    Byte数据

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)备注

如何确认PLC的变量的类型呢？可以在HslCommunicationDemo程序上测试知道，也可以直接调用 ReadWithType(String, UInt16) 来知道类型信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[WriteAsync 重载](558f1912-35b3-64b1-eca5-833d440e5b2f.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (String, String, Encoding)

[原文連結](http://api.hslcommunication.cn/html/3672d794-6b78-d944-b5df-3f4c64c038e3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[WriteAsync 方法](../html/558f1912-35b3-64b1-eca5-833d440e5b2f.htm "WriteAsync 方法 ")

[WriteAsync 方法 (String, Byte)](../html/29742e84-ecc1-0569-edb1-339042b3fe41.htm "WriteAsync 方法 (String, Byte)")

[WriteAsync 方法 (String, String, Encoding)](../html/3672d794-6b78-d944-b5df-3f4c64c038e3.htm "WriteAsync 方法 (String, String, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNetWriteAsync 方法 (String, String, Encoding) |

异步写入字符串信息，需要指定的编码信息  
Asynchronously write string information, need to specify the encoding information

**命名空间：**
 [HslCommunication.Profinet.Omron](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)  
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
	string value,
	Encoding encoding
)
```

```
Public Overrides Function WriteAsync ( 
	address As String,
	value As String,
	encoding As Encoding
) As Task(Of OperateResult)
```

```
public:
virtual Task<OperateResult^>^ WriteAsync(
	String^ address, 
	String^ value, 
	Encoding^ encoding
) override
```

```
abstract WriteAsync : 
        address : string * 
        value : string * 
        encoding : Encoding -> Task<OperateResult> 
override WriteAsync : 
        address : string * 
        value : string * 
        encoding : Encoding -> Task<OperateResult>
```

#### 参数

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemString  
    写入值

encoding
:   类型：System.TextEncoding  
    指定的编码信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
带有成功标识的结果类对象

#### 实现

[IReadWriteNetWriteAsync(String, String, Encoding)](cda78ec9-ab0c-b6f1-db38-0e51e6aa06d0.htm)  
[IReadWriteNetWriteAsync(String, String, Encoding)](cda78ec9-ab0c-b6f1-db38-0e51e6aa06d0.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

String类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 简单的写入
await melsec_net.WriteAsync( "D100", "ABCDEFGH" );

// 如果想要判断是否写入成功
OperateResult write = await melsec_net.WriteAsync( "D100", "ABCDEFGH" );
if (write.IsSuccess)
{
    // success
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[WriteAsync 重载](558f1912-35b3-64b1-eca5-833d440e5b2f.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OmronCipNet 字段

[原文連結](http://api.hslcommunication.cn/html/37dabe94-9b03-1948-8f79-4bbdfad15f11.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Omron](../html/cc7e1570-0983-ab07-e73a-494145c71ebb.htm "HslCommunication.Profinet.Omron")

[OmronCipNet 类](../html/5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm "OmronCipNet 类")

[OmronCipNet 构造函数](../html/bda49fa2-ff51-4a37-f7e1-9a83ee8eb428.htm "OmronCipNet 构造函数 ")

[OmronCipNet 属性](../html/488cfe9f-3319-2ccc-ad2a-ba3435024419.htm "OmronCipNet 属性")

[OmronCipNet 方法](../html/a6223914-afee-ab41-ba21-5b30b468808d.htm "OmronCipNet 方法")

[OmronCipNet 字段](../html/37dabe94-9b03-1948-8f79-4bbdfad15f11.htm "OmronCipNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OmronCipNet 字段 |

[OmronCipNet](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OmronCipNet 类](5bad2c83-6c75-88fb-a72c-e0232c00fcdf.htm)

[HslCommunication.Profinet.Omron 命名空间](cc7e1570-0983-ab07-e73a-494145c71ebb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)