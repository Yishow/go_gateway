# HslCommunication - HslCommunication.DCS

> 分類頁數: 15



---
## HslCommunication.DCS

[原文連結](http://api.hslcommunication.cn/html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.DCS 命名空间 |

[缺少 "N:HslCommunication.DCS" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [DcsNanJingAuto](f155e806-610d-93e3-0132-4e8e93e1b606.htm) | 南京自动化研究所的DCS系统，基于modbus实现，但是不是标准的实现 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DcsNanJingAuto 类

[原文連結](http://api.hslcommunication.cn/html/f155e806-610d-93e3-0132-4e8e93e1b606.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 构造函数](../html/bffa6e3f-99a6-2b12-48b1-5c79e539a80f.htm "DcsNanJingAuto 构造函数 ")

[DcsNanJingAuto 属性](../html/725c0127-21e7-27e0-0157-6b41076ab254.htm "DcsNanJingAuto 属性")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[DcsNanJingAuto 字段](../html/a5c2aaaf-fa3f-5fa2-d5d8-feac7eaf20fd.htm "DcsNanJingAuto 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAuto 类 |

南京自动化研究所的DCS系统，基于modbus实现，但是不是标准的实现

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
        [HslCommunication.ModBusModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)  
          HslCommunication.DCSDcsNanJingAuto

**命名空间：**
 [HslCommunication.DCS](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DcsNanJingAuto : ModbusTcpNet
```

```
Public Class DcsNanJingAuto
	Inherits ModbusTcpNet
```

```
public ref class DcsNanJingAuto : public ModbusTcpNet
```

```
type DcsNanJingAuto =  
    class
        inherit ModbusTcpNet
    end
```

DcsNanJingAuto 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DcsNanJingAuto](d0ae20e9-e5b3-a0ec-8818-364f34ae6e6f.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [DcsNanJingAuto(String, Int32, Byte)](1989e033-4f21-6d66-4df5-8bf67e3198e8.htm) | 指定服务器地址，端口号，客户端自己的站号来初始化  Specify the server address, port number, and client's own station number to initialize |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](a85f8f84-fa3e-9c41-d47f-87f95b203ac2.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [BroadcastStation](bd7ccbb3-3007-701e-c126-e12266643d1a.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [DataFormat](dd96d9f6-d6c6-92fc-a6eb-bfd3a3f5178f.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [DisableFunctionCode06](5e4bc90e-c1bc-2eab-60a5-983b4aa22a83.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [EnableWriteMaskCode](b08f214b-9f66-1900-c47e-60864753c5bf.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [IsCheckMessageId](7c1fc388-5bd8-85b9-8b20-fb25196e8702.htm) | 获取或设置是否进行检查返回的消息ID和发送的消息ID是否一致，默认为true，也就是检查  Get or set whether to check whether the returned message ID is consistent with the sent message ID, the default is true, that is, check (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [IsStringReverse](e015c6dc-b822-8435-d690-a0f523a40c09.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [MessageId](c11777df-d12d-a962-4993-ac2d6921cdf1.htm) | 获取modbus协议自增的消息号，你可以自定义modbus的消息号的规则，详细参见[ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)说明，也可以查找[SoftIncrementCount](fa4c7c56-11c0-0ef9-d17b-2106a420c805.htm)说明。  Get the message number incremented by the modbus protocol. You can customize the rules of the message number of the modbus. For details, please refer to the description of [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm), or you can find the description of [SoftIncrementCount](fa4c7c56-11c0-0ef9-d17b-2106a420c805.htm) (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [Station](c944fb53-6ac9-023e-e820-b2b9e6685844.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [StationCheckMatch](f0fd9f69-16ae-d779-7c53-47b7d5205ffd.htm) | 获取或设置当前是否启用站号检查的功能，默认启用，读写数据时将进行站号确认操作，如果需要忽略站号，则设置为 false 即可。  Gets or sets whether the station ID check function is enabled. It is enabled by default. When reading or writing data, the station ID is confirmed. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [WordReadBatchLength](426fa5d2-03f7-4706-f57b-9af91648238d.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |

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
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](43e46836-7567-0e1b-d5bc-7dabff0cd5a6.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [ModbusTcpNetGetNewNetMessage](c15715fd-c770-77f9-93e1-abde9e33b17d.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](677a8f96-49ba-f62d-0dc5-cae214d7fec9.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](03cfc2ec-fbcf-2db8-8272-185db7b41548.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm).) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](e118259e-9d2b-9360-0519-ccb1933d2ed4.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Read(String, UInt16)](e12f5a0b-9e09-e4a2-3f89-ac70e87078d0.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsync(String, UInt16)](8ed17c90-fb4f-d8b3-48ad-6068e2661618.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](fa64c9b8-2a76-f8d3-268c-a8f20cb20d62.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](29aa40bb-9bcf-989a-2343-0597632112f3.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadCoil(String)](444ee1ea-37de-a67c-c236-c3d7c9f4ef9d.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadCoil(String, UInt16)](cbc2bd1d-8e5c-2138-a529-e3f986c3c321.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadCoilAsync(String)](c68f256f-3df8-6db3-9a36-644abb2b610d.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadCoilAsync(String, UInt16)](fd8c738c-6806-9e41-0ae2-d2a536671c79.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDiscrete(String)](a0ff1b3e-e126-2e6f-444f-973c21259248.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadDiscrete(String, UInt16)](8a9f14ea-e712-cce8-af39-84f0b4ec748e.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String)](6ed73bfc-4374-a156-d440-b83043cf637e.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String, UInt16)](574acebf-b7be-0e3d-1729-84072f2c8f95.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](e3b58b7c-1f3d-f1a7-1556-38867e1aa12b.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](7b0826b6-fc51-1909-d541-5cc9c629e0ec.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadFile](12d6a0f7-30b0-d6e4-cdfc-b45131628241.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadFileAsync](511d6d06-3200-cef9-1756-2da9a7ce185d.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](baa66c96-8aa8-b52a-8822-7c8c4c32a481.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](565d2388-bbc5-09c3-3fed-dc91ff7a9bf8.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](dc91004f-2a78-e371-c080-b27480978b5d.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](188d8a1e-5cb4-b5e1-7e24-b2a27b0a7304.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (重写 [BinaryCommunicationReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm).) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](6fe6fbe9-f780-8858-2c68-e95533970ace.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](e1a7abde-8168-f333-123c-47909d99e6d1.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (重写 [BinaryCommunicationReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm).) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](847954ad-b0a4-a6f8-9876-fb529010082c.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](13bbe04b-c051-b98b-f83f-a314e2730d7a.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](1e4adede-3769-4737-4b2c-32d51873ea09.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](244dce7b-8182-5364-30d6-a93f761cddbc.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
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
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](2acda4c4-4eb0-7acd-06fe-8e30db736498.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](40d740f5-bb28-14e0-6617-e73e3ccefd74.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](fb0f05d7-9802-c5c3-261f-acfae71470b3.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](b5aba72e-019b-0b36-fd12-073383191732.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadWrite](527ffb76-8114-c81b-dcdb-034fafda0e01.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadWriteAsync](1a08f368-946f-030a-6b19-30d53b5af920.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [RegisteredAddressMapping](9061d9ac-f161-d35d-e2a0-10c5e278314e.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](14f25d02-723d-25f8-5839-3d2ed9ff1c97.htm) | (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [TranslateToModbusAddress](0d37943f-91c8-c0f8-df7a-950a75e39cda.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [UnpackResponseContent](08a4c087-2bcd-0510-bb6d-79ef5369e4ef.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
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
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](57a8dbe3-36c0-5239-2212-554a7efb0f39.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [Write(String, Boolean)](fae89069-685e-395f-8169-21e45299aba9.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Byte)](92518aa9-8cb0-48b7-b395-d4f649d697eb.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](25f40a7d-b863-0596-4b39-ab912007e771.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [Write(String, Int16)](9412d207-b2ff-5ee2-e18a-0e1d4b2c6a9a.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](0bb8bd8c-883a-bcde-6241-b6562bda2a41.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](a9ba289e-07a1-aa02-36f0-640d46d3dd34.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](91d90d52-be67-7ece-2097-a8a864d09be3.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [Write(String, UInt16)](e279b310-922c-e915-880d-4c478c6ac65b.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](4f0878bf-7249-c96e-b3e8-e6d276217557.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](65902c1c-281d-6a8b-15f9-39d3cc5e4efe.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](b87ae863-084e-a336-6695-97e85428fc8a.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](f0433324-86a6-4061-bc4a-1b2232637e0e.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](67476307-1a85-714a-71f2-b0d96ab90da2.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](bd48ad19-f5c7-b79c-1dfe-8705aadb1b74.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteAsync(String, Int16)](67d3e3dd-dfa8-02a5-160e-5d68c1918785.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](ddcf0c15-22c9-3f9a-6211-4c5c3f878f4a.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](5d000817-67fa-7661-e746-75c0b594e6d4.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](a78d8b3f-541f-5085-5fe9-5eb47bb69030.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteAsync(String, UInt16)](74cef334-03f8-9e04-ce97-f7899e6a9332.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](aef43300-2ef0-1628-bf60-3510ced7591a.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](f3046409-59c7-6b7b-8534-00dbfda164c6.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteFile](e7b3755b-54f5-3de5-f14d-9723e3d75c90.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteMask](e231f3a4-6715-3342-97b5-ef37f3488c08.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteMaskAsync](015e01d1-ea45-888b-da85-80d31e63e06a.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteOneRegister(String, Int16)](210f3f3d-6af2-548d-c44c-151a0bcf120e.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteOneRegister(String, UInt16)](d01b6d0d-1958-ca4f-3318-8fb3d6422e29.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, Int16)](a91508e1-cbbc-dcc2-157c-379f4dd1f789.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, UInt16)](7ac37e93-3ea3-b62b-0ee5-37cf6ef31c21.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |

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

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DcsNanJingAuto 构造函数 

[原文連結](http://api.hslcommunication.cn/html/bffa6e3f-99a6-2b12-48b1-5c79e539a80f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 构造函数](../html/bffa6e3f-99a6-2b12-48b1-5c79e539a80f.htm "DcsNanJingAuto 构造函数 ")

[DcsNanJingAuto 构造函数](../html/d0ae20e9-e5b3-a0ec-8818-364f34ae6e6f.htm "DcsNanJingAuto 构造函数 ")

[DcsNanJingAuto 构造函数 (String, Int32, Byte)](../html/1989e033-4f21-6d66-4df5-8bf67e3198e8.htm "DcsNanJingAuto 构造函数 (String, Int32, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAuto 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DcsNanJingAuto](d0ae20e9-e5b3-a0ec-8818-364f34ae6e6f.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [DcsNanJingAuto(String, Int32, Byte)](1989e033-4f21-6d66-4df5-8bf67e3198e8.htm) | 指定服务器地址，端口号，客户端自己的站号来初始化  Specify the server address, port number, and client's own station number to initialize |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DcsNanJingAuto 构造函数 

[原文連結](http://api.hslcommunication.cn/html/d0ae20e9-e5b3-a0ec-8818-364f34ae6e6f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 构造函数](../html/bffa6e3f-99a6-2b12-48b1-5c79e539a80f.htm "DcsNanJingAuto 构造函数 ")

[DcsNanJingAuto 构造函数](../html/d0ae20e9-e5b3-a0ec-8818-364f34ae6e6f.htm "DcsNanJingAuto 构造函数 ")

[DcsNanJingAuto 构造函数 (String, Int32, Byte)](../html/1989e033-4f21-6d66-4df5-8bf67e3198e8.htm "DcsNanJingAuto 构造函数 (String, Int32, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAuto 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.DCS](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DcsNanJingAuto()
```

```
Public Sub New
```

```
public:
DcsNanJingAuto()
```

```
new : unit -> DcsNanJingAuto
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[DcsNanJingAuto 重载](bffa6e3f-99a6-2b12-48b1-5c79e539a80f.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DcsNanJingAuto 构造函数 (String, Int32, Byte)

[原文連結](http://api.hslcommunication.cn/html/1989e033-4f21-6d66-4df5-8bf67e3198e8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 构造函数](../html/bffa6e3f-99a6-2b12-48b1-5c79e539a80f.htm "DcsNanJingAuto 构造函数 ")

[DcsNanJingAuto 构造函数](../html/d0ae20e9-e5b3-a0ec-8818-364f34ae6e6f.htm "DcsNanJingAuto 构造函数 ")

[DcsNanJingAuto 构造函数 (String, Int32, Byte)](../html/1989e033-4f21-6d66-4df5-8bf67e3198e8.htm "DcsNanJingAuto 构造函数 (String, Int32, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAuto 构造函数 (String, Int32, Byte) |

指定服务器地址，端口号，客户端自己的站号来初始化  
Specify the server address, port number, and client's own station number to initialize

**命名空间：**
 [HslCommunication.DCS](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DcsNanJingAuto(
	string ipAddress,
	int port = 502,
	byte station = 1
)
```

```
Public Sub New ( 
	ipAddress As String,
	Optional port As Integer = 502,
	Optional station As Byte = 1
)
```

```
public:
DcsNanJingAuto(
	String^ ipAddress, 
	int port = 502, 
	unsigned char station = 1
)
```

```
new : 
        ipAddress : string * 
        ?port : int * 
        ?station : byte 
(* Defaults:
        let _port = defaultArg port 502
        let _station = defaultArg station 1
*)
-> DcsNanJingAuto
```

#### 参数

ipAddress
:   类型：SystemString  
    服务器的Ip地址

port (Optional)
:   类型：SystemInt32  
    服务器的端口号

station (Optional)
:   类型：SystemByte  
    客户端自身的站号

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[DcsNanJingAuto 重载](bffa6e3f-99a6-2b12-48b1-5c79e539a80f.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DcsNanJingAuto 属性

[原文連結](http://api.hslcommunication.cn/html/725c0127-21e7-27e0-0157-6b41076ab254.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 构造函数](../html/bffa6e3f-99a6-2b12-48b1-5c79e539a80f.htm "DcsNanJingAuto 构造函数 ")

[DcsNanJingAuto 属性](../html/725c0127-21e7-27e0-0157-6b41076ab254.htm "DcsNanJingAuto 属性")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[DcsNanJingAuto 字段](../html/a5c2aaaf-fa3f-5fa2-d5d8-feac7eaf20fd.htm "DcsNanJingAuto 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAuto 属性 |

[DcsNanJingAuto](f155e806-610d-93e3-0132-4e8e93e1b606.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](a85f8f84-fa3e-9c41-d47f-87f95b203ac2.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [BroadcastStation](bd7ccbb3-3007-701e-c126-e12266643d1a.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [DataFormat](dd96d9f6-d6c6-92fc-a6eb-bfd3a3f5178f.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [DisableFunctionCode06](5e4bc90e-c1bc-2eab-60a5-983b4aa22a83.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [EnableWriteMaskCode](b08f214b-9f66-1900-c47e-60864753c5bf.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [IsCheckMessageId](7c1fc388-5bd8-85b9-8b20-fb25196e8702.htm) | 获取或设置是否进行检查返回的消息ID和发送的消息ID是否一致，默认为true，也就是检查  Get or set whether to check whether the returned message ID is consistent with the sent message ID, the default is true, that is, check (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [IsStringReverse](e015c6dc-b822-8435-d690-a0f523a40c09.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [MessageId](c11777df-d12d-a962-4993-ac2d6921cdf1.htm) | 获取modbus协议自增的消息号，你可以自定义modbus的消息号的规则，详细参见[ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)说明，也可以查找[SoftIncrementCount](fa4c7c56-11c0-0ef9-d17b-2106a420c805.htm)说明。  Get the message number incremented by the modbus protocol. You can customize the rules of the message number of the modbus. For details, please refer to the description of [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm), or you can find the description of [SoftIncrementCount](fa4c7c56-11c0-0ef9-d17b-2106a420c805.htm) (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [Station](c944fb53-6ac9-023e-e820-b2b9e6685844.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共属性 | [StationCheckMatch](f0fd9f69-16ae-d779-7c53-47b7d5205ffd.htm) | 获取或设置当前是否启用站号检查的功能，默认启用，读写数据时将进行站号确认操作，如果需要忽略站号，则设置为 false 即可。  Gets or sets whether the station ID check function is enabled. It is enabled by default. When reading or writing data, the station ID is confirmed. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [WordReadBatchLength](426fa5d2-03f7-4706-f57b-9af91648238d.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DcsNanJingAuto 方法

[原文連結](http://api.hslcommunication.cn/html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[GetNewNetMessage 方法](../html/43e46836-7567-0e1b-d5bc-7dabff0cd5a6.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/677a8f96-49ba-f62d-0dc5-cae214d7fec9.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/03cfc2ec-fbcf-2db8-8272-185db7b41548.htm "InitializationOnConnectAsync 方法 ")

[ReadFromCoreServer 方法](../html/d94a03bc-789e-1215-74a1-b6cb11665a9e.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/8ddb68e4-e088-b2e3-b4c4-8f02fc8b0284.htm "ReadFromCoreServerAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAuto 方法 |

[DcsNanJingAuto](f155e806-610d-93e3-0132-4e8e93e1b606.htm) 类型公开以下成员。

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
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](43e46836-7567-0e1b-d5bc-7dabff0cd5a6.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [ModbusTcpNetGetNewNetMessage](c15715fd-c770-77f9-93e1-abde9e33b17d.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](677a8f96-49ba-f62d-0dc5-cae214d7fec9.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](03cfc2ec-fbcf-2db8-8272-185db7b41548.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm).) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](e118259e-9d2b-9360-0519-ccb1933d2ed4.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Read(String, UInt16)](e12f5a0b-9e09-e4a2-3f89-ac70e87078d0.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsync(String, UInt16)](8ed17c90-fb4f-d8b3-48ad-6068e2661618.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](fa64c9b8-2a76-f8d3-268c-a8f20cb20d62.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](29aa40bb-9bcf-989a-2343-0597632112f3.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadCoil(String)](444ee1ea-37de-a67c-c236-c3d7c9f4ef9d.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadCoil(String, UInt16)](cbc2bd1d-8e5c-2138-a529-e3f986c3c321.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadCoilAsync(String)](c68f256f-3df8-6db3-9a36-644abb2b610d.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadCoilAsync(String, UInt16)](fd8c738c-6806-9e41-0ae2-d2a536671c79.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDiscrete(String)](a0ff1b3e-e126-2e6f-444f-973c21259248.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadDiscrete(String, UInt16)](8a9f14ea-e712-cce8-af39-84f0b4ec748e.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String)](6ed73bfc-4374-a156-d440-b83043cf637e.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String, UInt16)](574acebf-b7be-0e3d-1729-84072f2c8f95.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](e3b58b7c-1f3d-f1a7-1556-38867e1aa12b.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](7b0826b6-fc51-1909-d541-5cc9c629e0ec.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadFile](12d6a0f7-30b0-d6e4-cdfc-b45131628241.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadFileAsync](511d6d06-3200-cef9-1756-2da9a7ce185d.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](baa66c96-8aa8-b52a-8822-7c8c4c32a481.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](565d2388-bbc5-09c3-3fed-dc91ff7a9bf8.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](dc91004f-2a78-e371-c080-b27480978b5d.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](188d8a1e-5cb4-b5e1-7e24-b2a27b0a7304.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (重写 [BinaryCommunicationReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm).) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](6fe6fbe9-f780-8858-2c68-e95533970ace.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](e1a7abde-8168-f333-123c-47909d99e6d1.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (重写 [BinaryCommunicationReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm).) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](847954ad-b0a4-a6f8-9876-fb529010082c.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](13bbe04b-c051-b98b-f83f-a314e2730d7a.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](1e4adede-3769-4737-4b2c-32d51873ea09.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](244dce7b-8182-5364-30d6-a93f761cddbc.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
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
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](2acda4c4-4eb0-7acd-06fe-8e30db736498.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](40d740f5-bb28-14e0-6617-e73e3ccefd74.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](fb0f05d7-9802-c5c3-261f-acfae71470b3.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](b5aba72e-019b-0b36-fd12-073383191732.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadWrite](527ffb76-8114-c81b-dcdb-034fafda0e01.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [ReadWriteAsync](1a08f368-946f-030a-6b19-30d53b5af920.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [RegisteredAddressMapping](9061d9ac-f161-d35d-e2a0-10c5e278314e.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](14f25d02-723d-25f8-5839-3d2ed9ff1c97.htm) | (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [TranslateToModbusAddress](0d37943f-91c8-c0f8-df7a-950a75e39cda.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [UnpackResponseContent](08a4c087-2bcd-0510-bb6d-79ef5369e4ef.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
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
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](57a8dbe3-36c0-5239-2212-554a7efb0f39.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [Write(String, Boolean)](fae89069-685e-395f-8169-21e45299aba9.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Byte)](92518aa9-8cb0-48b7-b395-d4f649d697eb.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](25f40a7d-b863-0596-4b39-ab912007e771.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [Write(String, Int16)](9412d207-b2ff-5ee2-e18a-0e1d4b2c6a9a.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](0bb8bd8c-883a-bcde-6241-b6562bda2a41.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](a9ba289e-07a1-aa02-36f0-640d46d3dd34.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](91d90d52-be67-7ece-2097-a8a864d09be3.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [Write(String, UInt16)](e279b310-922c-e915-880d-4c478c6ac65b.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](4f0878bf-7249-c96e-b3e8-e6d276217557.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](65902c1c-281d-6a8b-15f9-39d3cc5e4efe.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](b87ae863-084e-a336-6695-97e85428fc8a.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](f0433324-86a6-4061-bc4a-1b2232637e0e.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](67476307-1a85-714a-71f2-b0d96ab90da2.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](bd48ad19-f5c7-b79c-1dfe-8705aadb1b74.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteAsync(String, Int16)](67d3e3dd-dfa8-02a5-160e-5d68c1918785.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](ddcf0c15-22c9-3f9a-6211-4c5c3f878f4a.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](5d000817-67fa-7661-e746-75c0b594e6d4.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](a78d8b3f-541f-5085-5fe9-5eb47bb69030.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteAsync(String, UInt16)](74cef334-03f8-9e04-ce97-f7899e6a9332.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](aef43300-2ef0-1628-bf60-3510ced7591a.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](f3046409-59c7-6b7b-8534-00dbfda164c6.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteFile](e7b3755b-54f5-3de5-f14d-9723e3d75c90.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteMask](e231f3a4-6715-3342-97b5-ef37f3488c08.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteMaskAsync](015e01d1-ea45-888b-da85-80d31e63e06a.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteOneRegister(String, Int16)](210f3f3d-6af2-548d-c44c-151a0bcf120e.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteOneRegister(String, UInt16)](d01b6d0d-1958-ca4f-3318-8fb3d6422e29.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, Int16)](a91508e1-cbbc-dcc2-157c-379f4dd1f789.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, UInt16)](7ac37e93-3ea3-b62b-0ee5-37cf6ef31c21.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/43e46836-7567-0e1b-d5bc-7dabff0cd5a6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[GetNewNetMessage 方法](../html/43e46836-7567-0e1b-d5bc-7dabff0cd5a6.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/677a8f96-49ba-f62d-0dc5-cae214d7fec9.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/03cfc2ec-fbcf-2db8-8272-185db7b41548.htm "InitializationOnConnectAsync 方法 ")

[ReadFromCoreServer 方法](../html/d94a03bc-789e-1215-74a1-b6cb11665a9e.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/8ddb68e4-e088-b2e3-b4c4-8f02fc8b0284.htm "ReadFromCoreServerAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAutoGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.DCS](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)  
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

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnect 方法 

[原文連結](http://api.hslcommunication.cn/html/677a8f96-49ba-f62d-0dc5-cae214d7fec9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[GetNewNetMessage 方法](../html/43e46836-7567-0e1b-d5bc-7dabff0cd5a6.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/677a8f96-49ba-f62d-0dc5-cae214d7fec9.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/03cfc2ec-fbcf-2db8-8272-185db7b41548.htm "InitializationOnConnectAsync 方法 ")

[ReadFromCoreServer 方法](../html/d94a03bc-789e-1215-74a1-b6cb11665a9e.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/8ddb68e4-e088-b2e3-b4c4-8f02fc8b0284.htm "ReadFromCoreServerAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAutoInitializationOnConnect 方法 |

根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  
Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created.

**命名空间：**
 [HslCommunication.DCS](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override OperateResult InitializationOnConnect()
```

```
Protected Overrides Function InitializationOnConnect As OperateResult
```

```
protected:
virtual OperateResult^ InitializationOnConnect() override
```

```
abstract InitializationOnConnect : unit -> OperateResult 
override InitializationOnConnect : unit -> OperateResult
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否初始化成功，依据具体的协议进行重写

![](../icons/SectionExpanded.png)示例

有些协议不需要握手信号，比如三菱的MC协议，Modbus协议，西门子和欧姆龙就存在握手信息，此处的例子是继承本类后重写的西门子的协议示例

西门子重连示例

[复制](# "复制")

```
        /// <inheritdoc/>
        public override OperateResult<byte[]> ReadFromCoreServer( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = base.ReadFromCoreServer( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override OperateResult InitializationOnConnect( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = ReadFromCoreServer( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = ReadFromCoreServer( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }

#if !NET35 && !NET20
        /// <inheritdoc/>
        public async override Task<OperateResult<byte[]>> ReadFromCoreServerAsync( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = await base.ReadFromCoreServerAsync( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override async Task<OperateResult> InitializationOnConnectAsync( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }
#endif
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnectAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/03cfc2ec-fbcf-2db8-8272-185db7b41548.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[GetNewNetMessage 方法](../html/43e46836-7567-0e1b-d5bc-7dabff0cd5a6.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/677a8f96-49ba-f62d-0dc5-cae214d7fec9.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/03cfc2ec-fbcf-2db8-8272-185db7b41548.htm "InitializationOnConnectAsync 方法 ")

[ReadFromCoreServer 方法](../html/d94a03bc-789e-1215-74a1-b6cb11665a9e.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/8ddb68e4-e088-b2e3-b4c4-8f02fc8b0284.htm "ReadFromCoreServerAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAutoInitializationOnConnectAsync 方法 |

根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  
Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created.

**命名空间：**
 [HslCommunication.DCS](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override Task<OperateResult> InitializationOnConnectAsync()
```

```
Protected Overrides Function InitializationOnConnectAsync As Task(Of OperateResult)
```

```
protected:
virtual Task<OperateResult^>^ InitializationOnConnectAsync() override
```

```
abstract InitializationOnConnectAsync : unit -> Task<OperateResult> 
override InitializationOnConnectAsync : unit -> Task<OperateResult>
```

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否初始化成功，依据具体的协议进行重写

![](../icons/SectionExpanded.png)示例

有些协议不需要握手信号，比如三菱的MC协议，Modbus协议，西门子和欧姆龙就存在握手信息，此处的例子是继承本类后重写的西门子的协议示例

西门子重连示例

[复制](# "复制")

```
        /// <inheritdoc/>
        public override OperateResult<byte[]> ReadFromCoreServer( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = base.ReadFromCoreServer( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override OperateResult InitializationOnConnect( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = ReadFromCoreServer( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = ReadFromCoreServer( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }

#if !NET35 && !NET20
        /// <inheritdoc/>
        public async override Task<OperateResult<byte[]>> ReadFromCoreServerAsync( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = await base.ReadFromCoreServerAsync( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override async Task<OperateResult> InitializationOnConnectAsync( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }
#endif
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServer 方法 

[原文連結](http://api.hslcommunication.cn/html/d94a03bc-789e-1215-74a1-b6cb11665a9e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[ReadFromCoreServer 方法](../html/d94a03bc-789e-1215-74a1-b6cb11665a9e.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServer 方法 (CommunicationPipe, Byte[], Boolean, Boolean)](../html/188d8a1e-5cb4-b5e1-7e24-b2a27b0a7304.htm "ReadFromCoreServer 方法 (CommunicationPipe, Byte[], Boolean, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAutoReadFromCoreServer 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](dc91004f-2a78-e371-c080-b27480978b5d.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](188d8a1e-5cb4-b5e1-7e24-b2a27b0a7304.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (重写 [BinaryCommunicationReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServer 方法 (CommunicationPipe, Byte[], Boolean, Boolean)

[原文連結](http://api.hslcommunication.cn/html/188d8a1e-5cb4-b5e1-7e24-b2a27b0a7304.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[ReadFromCoreServer 方法](../html/d94a03bc-789e-1215-74a1-b6cb11665a9e.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServer 方法 (CommunicationPipe, Byte[], Boolean, Boolean)](../html/188d8a1e-5cb4-b5e1-7e24-b2a27b0a7304.htm "ReadFromCoreServer 方法 (CommunicationPipe, Byte[], Boolean, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAutoReadFromCoreServer 方法 (CommunicationPipe, Byte, Boolean, Boolean) |

使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁

**命名空间：**
 [HslCommunication.DCS](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<byte[]> ReadFromCoreServer(
	CommunicationPipe pipe,
	byte[] send,
	bool hasResponseData = true,
	bool usePackHeader = true
)
```

```
Public Overrides Function ReadFromCoreServer ( 
	pipe As CommunicationPipe,
	send As Byte(),
	Optional hasResponseData As Boolean = true,
	Optional usePackHeader As Boolean = true
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ ReadFromCoreServer(
	CommunicationPipe^ pipe, 
	array<unsigned char>^ send, 
	bool hasResponseData = true, 
	bool usePackHeader = true
) override
```

```
abstract ReadFromCoreServer : 
        pipe : CommunicationPipe * 
        send : byte[] * 
        ?hasResponseData : bool * 
        ?usePackHeader : bool 
(* Defaults:
        let _hasResponseData = defaultArg hasResponseData true
        let _usePackHeader = defaultArg usePackHeader true
*)
-> OperateResult<byte[]> 
override ReadFromCoreServer : 
        pipe : CommunicationPipe * 
        send : byte[] * 
        ?hasResponseData : bool * 
        ?usePackHeader : bool 
(* Defaults:
        let _hasResponseData = defaultArg hasResponseData true
        let _usePackHeader = defaultArg usePackHeader true
*)
-> OperateResult<byte[]>
```

#### 参数

pipe
:   类型：[HslCommunication.Core.PipeCommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)  
    管道信息

send
:   类型：SystemByte  
    等待发送的数据

hasResponseData (Optional)
:   类型：SystemBoolean  
    是否需要返回的数据

usePackHeader (Optional)
:   类型：SystemBoolean  

    [缺少 "M:HslCommunication.DCS.DcsNanJingAuto.ReadFromCoreServer(HslCommunication.Core.Pipe.CommunicationPipe,System.Byte[],System.Boolean,System.Boolean)" 的 <param name="usePackHeader"/> 文档]

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[ReadFromCoreServer 重载](d94a03bc-789e-1215-74a1-b6cb11665a9e.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServerAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/8ddb68e4-e088-b2e3-b4c4-8f02fc8b0284.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[ReadFromCoreServerAsync 方法](../html/8ddb68e4-e088-b2e3-b4c4-8f02fc8b0284.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerAsync 方法 (CommunicationPipe, Byte[], Boolean, Boolean)](../html/e1a7abde-8168-f333-123c-47909d99e6d1.htm "ReadFromCoreServerAsync 方法 (CommunicationPipe, Byte[], Boolean, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAutoReadFromCoreServerAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](6fe6fbe9-f780-8858-2c68-e95533970ace.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](e1a7abde-8168-f333-123c-47909d99e6d1.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (重写 [BinaryCommunicationReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServerAsync 方法 (CommunicationPipe, Byte[], Boolean, Boolean)

[原文連結](http://api.hslcommunication.cn/html/e1a7abde-8168-f333-123c-47909d99e6d1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[ReadFromCoreServerAsync 方法](../html/8ddb68e4-e088-b2e3-b4c4-8f02fc8b0284.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerAsync 方法 (CommunicationPipe, Byte[], Boolean, Boolean)](../html/e1a7abde-8168-f333-123c-47909d99e6d1.htm "ReadFromCoreServerAsync 方法 (CommunicationPipe, Byte[], Boolean, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAutoReadFromCoreServerAsync 方法 (CommunicationPipe, Byte, Boolean, Boolean) |

使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁

**命名空间：**
 [HslCommunication.DCS](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<byte[]>> ReadFromCoreServerAsync(
	CommunicationPipe pipe,
	byte[] send,
	bool hasResponseData = true,
	bool usePackHeader = true
)
```

```
Public Overrides Function ReadFromCoreServerAsync ( 
	pipe As CommunicationPipe,
	send As Byte(),
	Optional hasResponseData As Boolean = true,
	Optional usePackHeader As Boolean = true
) As Task(Of OperateResult(Of Byte()))
```

```
public:
virtual Task<OperateResult<array<unsigned char>^>^>^ ReadFromCoreServerAsync(
	CommunicationPipe^ pipe, 
	array<unsigned char>^ send, 
	bool hasResponseData = true, 
	bool usePackHeader = true
) override
```

```
abstract ReadFromCoreServerAsync : 
        pipe : CommunicationPipe * 
        send : byte[] * 
        ?hasResponseData : bool * 
        ?usePackHeader : bool 
(* Defaults:
        let _hasResponseData = defaultArg hasResponseData true
        let _usePackHeader = defaultArg usePackHeader true
*)
-> Task<OperateResult<byte[]>> 
override ReadFromCoreServerAsync : 
        pipe : CommunicationPipe * 
        send : byte[] * 
        ?hasResponseData : bool * 
        ?usePackHeader : bool 
(* Defaults:
        let _hasResponseData = defaultArg hasResponseData true
        let _usePackHeader = defaultArg usePackHeader true
*)
-> Task<OperateResult<byte[]>>
```

#### 参数

pipe
:   类型：[HslCommunication.Core.PipeCommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)  
    管道信息

send
:   类型：SystemByte  
    等待发送的数据

hasResponseData (Optional)
:   类型：SystemBoolean  
    是否需要返回的数据

usePackHeader (Optional)
:   类型：SystemBoolean  

    [缺少 "M:HslCommunication.DCS.DcsNanJingAuto.ReadFromCoreServerAsync(HslCommunication.Core.Pipe.CommunicationPipe,System.Byte[],System.Boolean,System.Boolean)" 的 <param name="usePackHeader"/> 文档]

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[ReadFromCoreServerAsync 重载](8ddb68e4-e088-b2e3-b4c4-8f02fc8b0284.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DcsNanJingAuto 字段

[原文連結](http://api.hslcommunication.cn/html/a5c2aaaf-fa3f-5fa2-d5d8-feac7eaf20fd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DCS](../html/fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm "HslCommunication.DCS")

[DcsNanJingAuto 类](../html/f155e806-610d-93e3-0132-4e8e93e1b606.htm "DcsNanJingAuto 类")

[DcsNanJingAuto 构造函数](../html/bffa6e3f-99a6-2b12-48b1-5c79e539a80f.htm "DcsNanJingAuto 构造函数 ")

[DcsNanJingAuto 属性](../html/725c0127-21e7-27e0-0157-6b41076ab254.htm "DcsNanJingAuto 属性")

[DcsNanJingAuto 方法](../html/29b63c19-7c77-f33c-f099-cb99efcedcda.htm "DcsNanJingAuto 方法")

[DcsNanJingAuto 字段](../html/a5c2aaaf-fa3f-5fa2-d5d8-feac7eaf20fd.htm "DcsNanJingAuto 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DcsNanJingAuto 字段 |

[DcsNanJingAuto](f155e806-610d-93e3-0132-4e8e93e1b606.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DcsNanJingAuto 类](f155e806-610d-93e3-0132-4e8e93e1b606.htm)

[HslCommunication.DCS 命名空间](fa9dcd26-6740-82c2-e1e1-09a7c04b6c53.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)