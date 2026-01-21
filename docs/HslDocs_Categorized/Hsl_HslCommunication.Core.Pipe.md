# HslCommunication - HslCommunication.Core.Pipe

> 分類頁數: 30



---
## HslCommunication.Core.Pipe

[原文連結](http://api.hslcommunication.cn/html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[PipeBase 类](../html/614e4c21-7d64-71cb-0419-722ac613d708.htm "PipeBase 类")

[PipeDebugRemote 类](../html/8022a689-f08c-9d70-13f7-b703c2de579d.htm "PipeDebugRemote 类")

[PipeDtuNet 类](../html/97793361-a796-5f13-5029-1abdb2bff247.htm "PipeDtuNet 类")

[PipeMoxa 类](../html/2060615b-7a71-50a1-a649-15ba83c835fd.htm "PipeMoxa 类")

[PipeMqttClient 类](../html/09015966-2aa8-f439-a219-145a2e0bda62.htm "PipeMqttClient 类")

[PipeSerial 类](../html/06daf5c2-5e8f-bb51-f583-56ee9416b609.htm "PipeSerial 类")

[PipeSerialPort 类](../html/9d924607-3dfc-3916-b532-c5dec48fe347.htm "PipeSerialPort 类")

[PipeSocket 类](../html/cd37211e-6406-db08-5178-91649cc868e9.htm "PipeSocket 类")

[PipeSslNet 类](../html/3a530c2d-0288-e171-c842-6bd8d0465ce8.htm "PipeSslNet 类")

[PipeTcpNet 类](../html/26b94bd5-2695-36e7-6378-c3d41846be8f.htm "PipeTcpNet 类")

[PipeUdpNet 类](../html/5a0b810e-187e-826b-0a13-34fd5576aa0e.htm "PipeUdpNet 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Core.Pipe 命名空间 |

[缺少 "N:HslCommunication.Core.Pipe" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [CommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm) | 用于通信的管道信息，包含基本的通信实现，当前类是抽象类，需要继承实现具体的收发报文才可以使用。  Channel information for communication, including the basic communication implementation, the current class is an abstract class, need to inherit the implementation of specific messages can be used. |
| 公共类 | [PipeBase](614e4c21-7d64-71cb-0419-722ac613d708.htm) | 管道的基础类对象 |
| 公共类 | [PipeDebugRemote](8022a689-f08c-9d70-13f7-b703c2de579d.htm) | 用于调试的远程管道信息  Remote pipeline information for debugging |
| 公共类 | [PipeDtuNet](97793361-a796-5f13-5029-1abdb2bff247.htm) | DTU(数据转换模块)的管道信息  Pipeline information of the DTU (Data Transfer unit) |
| 公共类 | [PipeMoxa](2060615b-7a71-50a1-a649-15ba83c835fd.htm) | 基于MOXA公司提供的串口驱动实现的管道类，目前仅支持windows平台，需要当前目录下存在 PCOMM.DLL 组件  The pipe class based on the serial port driver provided by MOXA currently only supports the windows platform and requires the existence of the PCOMM.DLL component in the current directory |
| 公共类 | [PipeMqttClient](09015966-2aa8-f439-a219-145a2e0bda62.htm) | 基于MQTT通信实现的管道信息 |
| 公共类 | [PipeSerial](06daf5c2-5e8f-bb51-f583-56ee9416b609.htm) | **已过时。** 串口的管道类对象，可以在不同的串口类中使用一个串口的通道信息  The pipe class object of the serial port can use the channel information of a serial port in different serial port classes |
| 公共类 | [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) | 串口管道信息 |
| 公共类 | [PipeSocket](cd37211e-6406-db08-5178-91649cc868e9.htm) | 基于网络通信的管道信息，可以设置额外的一些参数信息，例如连接超时时间，读取超时时间等等。  Based on the pipe information of network communication, some additional parameter information can be set, such as connection timeout time, read timeout time and so on. |
| 公共类 | [PipeSslNet](3a530c2d-0288-e171-c842-6bd8d0465ce8.htm) | 基于SSL/TLS加密的管道信息，内部基于 TCP/IP 通信实现  Pipe information based on SSL/TLS encryption, and internal TCP/IP communication |
| 公共类 | [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) | 用于TCP/IP协议的传输管道信息  Transport pipe information of the IP protocol |
| 公共类 | [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) | 基于UDP/IT通信的管道信息  Pipeline information based on UDP/IT communication |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CommunicationPipe 类

[原文連結](http://api.hslcommunication.cn/html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 构造函数](../html/1b2f92a0-02e1-fd21-72e3-b7c23c515b48.htm "CommunicationPipe 构造函数 ")

[CommunicationPipe 属性](../html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm "CommunicationPipe 属性")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CommunicationPipe 字段](../html/3706e381-f54a-a217-7fd3-ce9a5ab98b62.htm "CommunicationPipe 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipe 类 |

用于通信的管道信息，包含基本的通信实现，当前类是抽象类，需要继承实现具体的收发报文才可以使用。  
Channel information for communication, including the basic communication implementation, the current class is an abstract class, need to inherit the implementation of specific messages can be used.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Core.PipeCommunicationPipe  
    [HslCommunication.Core.PipePipeMoxa](2060615b-7a71-50a1-a649-15ba83c835fd.htm)  
    [HslCommunication.Core.PipePipeMqttClient](09015966-2aa8-f439-a219-145a2e0bda62.htm)  
    [HslCommunication.Core.PipePipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)  
    [HslCommunication.Core.PipePipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public abstract class CommunicationPipe : IDisposable
```

```
Public MustInherit Class CommunicationPipe
	Implements IDisposable
```

```
public ref class CommunicationPipe abstract : IDisposable
```

```
[<AbstractClassAttribute>]
type CommunicationPipe =  
    class
        interface IDisposable
    end
```

CommunicationPipe 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CommunicationPipe](1b2f92a0-02e1-fd21-72e3-b7c23c515b48.htm) | 实例化一个默认的构造对象  Instantiate a default constructor object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CommunicationLock](0d03f5cf-c63c-903c-31d2-9261b91de976.htm) | 获取或设置当前管道的线程锁对象，默认是简单的一个互斥锁  Gets or sets the thread lock object of the current pipeline, which defaults to a simple mutex |
| 公共属性 | [DecideWhetherQAMessageFunction](bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm) | 用来决定当前接收的消息是否是问答服务的消息，可由外界设置委托来决定  It is used to determine whether the currently received message is a message of the question answering service, which can be determined by the external setting delegate |
| 公共属性 | [IsPersistentConnection](a43de804-55cd-3862-9898-efae0fde9936.htm) | 获取或设置当前的管道是否是长连接，仅对于串口及TCP是有效的，默认都是长连接  Gets or sets whether the current pipe is a long connection. This is valid only for serial ports and TCP. The default is a long connection |
| 公共属性代码示例 | [ReceiveTimeOut](b5e80233-0023-db67-90b1-fbbeff60ce3d.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback |
| 公共属性 | [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. |
| 公共属性 | [UseServerActivePush](b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm) | 获取或设置当前的管道是否激活从设备主动推送的功能，设置为 true 时支持主动从设备方接收数据信息  Gets or sets whether the current pipeline activates the function of actively pushing data from the device. If this is set to true, it supports actively receiving data information from the device |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [CheckMessageComplete](eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm) | 根据给定的消息，发送的数据，接收到数据来判断是否接收完成报文  According to the given message, sent data, received data to determine whether to receive the completed message |
| 公共方法 | [CloseCommunication](9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm) | 关闭当前的管道信息，返回是否关闭成功的结果对象  Close the current pipeline information and return the result object whether the closure was successful |
| 公共方法 | [CloseCommunicationAsync](3b73e48c-3a64-eb2c-a827-fedd9926558e.htm) |  |
| 公共方法 | [Dispose](fb51c8a3-b4bc-5c64-b1ba-742464b90b41.htm) | 释放被 CommunicationPipe 使用的所有资源 |
| 受保护的方法 | [Dispose(Boolean)](94cf19e0-0b6e-a57f-7559-27bf48a8e8f9.htm) | 释放被 CommunicationPipe 使用的非托管资源，并且是否托管资源（可选） |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [HasCacheData](3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm) | 当前管道的缓存里是否有数据存在，有就返回 True  If there is any data in the cache of the current pipeline, return True |
| 受保护的方法 | [IncrConnectErrorCount](946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm) | 自增当前的连续错误计数，并且获取自增后的值信息，最大到20亿为止，无法继续增加了。  Increment the current continuous error count, and obtain the value information after increment, up to 2 billion, can not continue to increase. |
| 公共方法 | [IsConnectError](dba07b04-a27a-754d-1f88-053ef60251f3.htm) | 当前的管道连接对象是否发生了错误，这里的错误通常是串口，或是网络错误  Whether there is an error in the current pipe connection object, where the error is usually a serial port, or a network error |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [OpenCommunication](009e7e59-f95f-29ad-7545-b691bf676e07.htm) | 打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  Opens the current pipe information, returns whether the result object was successfully opened, and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true |
| 公共方法 | [OpenCommunicationAsync](4623d60d-297f-edd6-91df-892f51ab4609.htm) | 打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  Opens the current pipe information, returns whether the result object was successfully opened, and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true |
| 公共方法 | [RaisePipeError](32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm) | 主动引发一个管道错误，从而让管道可以重新打开  Actively causes a pipe error so that the pipe can be reopened |
| 公共方法 | [ReadFromCoreServer](a85930ab-0e17-beee-cc5b-1178ce83f78c.htm) | 将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type, the data sent, whether there is a data response, and the sleep time |
| 公共方法 | [ReadFromCoreServerAsync](3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm) | 将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type, the data sent, whether there is a data response, and the sleep time |
| 受保护的方法 | [ReadFromCoreServerHelper](e208a32f-401a-20af-86a7-3972e220b425.htm) | 将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type, the data sent, whether there is a data response, and the sleep time |
| 受保护的方法 | [ReadFromCoreServerHelperAsync](7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm) | 将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type, the data sent, whether there is a data response, and the sleep time |
| 公共方法 | [Receive(Int32, Int32, ActionInt64, Int64)](6a8d9448-097e-4cb0-0742-95327d4431d3.htm) | 从管道里，接收指定长度的报文数据信息，如果长度指定为-1，表示接收不超过2048字节的动态长度。另外可以指定超时时间，进度报告等  Receives the packet data of a specified length from the pipe. If the length is set to -1, it indicates that the dynamic length of the packet is not more than 2048 bytes. You can also specify timeouts, progress reports, etc |
| 公共方法 | [Receive(Byte, Int32, Int32, Int32, ActionInt64, Int64)](9979783b-37d4-3993-b8b8-aed6055124bf.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. |
| 公共方法 | [ReceiveAsync(Int32, Int32, ActionInt64, Int64)](b56c1109-0eed-300f-ec21-6217d823f3cf.htm) | 从管道里，接收指定长度的报文数据信息，如果长度指定为-1，表示接收不超过2048字节的动态长度。另外可以指定超时时间，进度报告等  Receives the packet data of a specified length from the pipe. If the length is set to -1, it indicates that the dynamic length of the packet is not more than 2048 bytes. You can also specify timeouts, progress reports, etc |
| 公共方法 | [ReceiveAsync(Byte, Int32, Int32, Int32, ActionInt64, Int64)](2eb79f82-e65a-51c6-4b57-8713600bb27e.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. |
| 公共方法 | [ReceiveMessage](fa579a87-d719-3fef-af36-a3d37af08200.htm) | 包含了一个复杂的逻辑，从管道里根据当前的消息格式定义，接收报文信息，这个报文可能是来自服务器主动推送的。具体可以通过参数 useActivePush 来特殊控制。  Contains a complex logic from the pipeline, according to the current message format definition, to receive message information, this message may be actively pushed from the server. The parameter useActivePush can be used for special control. |
| 公共方法 | [ReceiveMessageAsync](c2e7239d-8839-c1ed-6993-f4da7a615b78.htm) | 包含了一个复杂的逻辑，从管道里根据当前的消息格式定义，接收报文信息，这个报文可能是来自服务器主动推送的。具体可以通过参数 useActivePush 来特殊控制。  Contains a complex logic from the pipeline, according to the current message format definition, to receive message information, this message may be actively pushed from the server. The parameter useActivePush can be used for special control. |
| 公共方法 | [ResetConnectErrorCount](c3d25cde-83f6-4257-8344-0e458c8bd685.htm) | 重置当前的连续错误计数为0，并且返回重置前时候的值  Resets the current consecutive error count to 0 and returns the value before the reset |
| 公共方法 | [Send(Byte)](f71f0371-1270-df68-9648-4919974c4fef.htm) | 发送数据到当前的管道中去  Send data to the current pipe |
| 公共方法 | [Send(Byte, Int32, Int32)](45f67684-08ea-f88c-2dff-8e6701e70996.htm) | 将一个数据缓存中的指定的部分字段，发送到当前的管道中去  Sends the specified partial field from a data cache to the current pipeline |
| 公共方法 | [SendAsync(Byte)](79e3ebb0-fd3d-11c1-53ca-00970cc89c25.htm) | 发送数据到当前的管道中去  Send data to the current pipe |
| 公共方法 | [SendAsync(Byte, Int32, Int32)](e3cd5a7d-1e61-0d44-c5f2-81e9e9606cd8.htm) | 将一个数据缓存中的指定的部分字段，发送到当前的管道中去  Sends the specified partial field from a data cache to the current pipeline |
| 受保护的方法 | [SetBufferQA](6677295f-e375-da2e-1931-ec881676e83f.htm) | 设置当前的问答状态下的缓存数据  Set the cache data in the current Q&A state |
| 公共方法 | [StartReceiveBackground](9b977dfb-5c8b-d6d9-7c94-443345235b26.htm) | 开始后台接收相关的报文数据，当[UseServerActivePush](b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm)为True时，则使用本方法 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [autoResetEvent](7630d2a7-ab28-637c-400e-d143cae1bd6d.htm) | 当启用设备方主动发送数据时，用于同步访问方法的信号同步功能 |
| 受保护的字段 | [bufferQA](30fd59b1-50a5-0465-dc46-b98c8870a58d.htm) | 当启用设备方主动发送数据时，用于应答服务机制的数据缓存 |
| 受保护的字段 | [isPersistentConn](9d79ff1c-42ed-bc8a-e4e8-615f065e5afe.htm) | 是否是长连接的状态  Whether it is a long connection state |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CommunicationPipe 构造函数 

[原文連結](http://api.hslcommunication.cn/html/1b2f92a0-02e1-fd21-72e3-b7c23c515b48.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 构造函数](../html/1b2f92a0-02e1-fd21-72e3-b7c23c515b48.htm "CommunicationPipe 构造函数 ")

[CommunicationPipe 属性](../html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm "CommunicationPipe 属性")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CommunicationPipe 字段](../html/3706e381-f54a-a217-7fd3-ce9a5ab98b62.htm "CommunicationPipe 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipe 构造函数 |

实例化一个默认的构造对象  
Instantiate a default constructor object

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public CommunicationPipe()
```

```
Public Sub New
```

```
public:
CommunicationPipe()
```

```
new : unit -> CommunicationPipe
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CommunicationPipe 属性

[原文連結](http://api.hslcommunication.cn/html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 属性](../html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm "CommunicationPipe 属性")

[CommunicationLock 属性](../html/0d03f5cf-c63c-903c-31d2-9261b91de976.htm "CommunicationLock 属性 ")

[DecideWhetherQAMessageFunction 属性](../html/bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm "DecideWhetherQAMessageFunction 属性 ")

[IsPersistentConnection 属性](../html/a43de804-55cd-3862-9898-efae0fde9936.htm "IsPersistentConnection 属性 ")

[ReceiveTimeOut 属性](../html/b5e80233-0023-db67-90b1-fbbeff60ce3d.htm "ReceiveTimeOut 属性 ")

[SleepTime 属性](../html/1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm "SleepTime 属性 ")

[UseServerActivePush 属性](../html/b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm "UseServerActivePush 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipe 属性 |

[CommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CommunicationLock](0d03f5cf-c63c-903c-31d2-9261b91de976.htm) | 获取或设置当前管道的线程锁对象，默认是简单的一个互斥锁  Gets or sets the thread lock object of the current pipeline, which defaults to a simple mutex |
| 公共属性 | [DecideWhetherQAMessageFunction](bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm) | 用来决定当前接收的消息是否是问答服务的消息，可由外界设置委托来决定  It is used to determine whether the currently received message is a message of the question answering service, which can be determined by the external setting delegate |
| 公共属性 | [IsPersistentConnection](a43de804-55cd-3862-9898-efae0fde9936.htm) | 获取或设置当前的管道是否是长连接，仅对于串口及TCP是有效的，默认都是长连接  Gets or sets whether the current pipe is a long connection. This is valid only for serial ports and TCP. The default is a long connection |
| 公共属性代码示例 | [ReceiveTimeOut](b5e80233-0023-db67-90b1-fbbeff60ce3d.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback |
| 公共属性 | [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. |
| 公共属性 | [UseServerActivePush](b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm) | 获取或设置当前的管道是否激活从设备主动推送的功能，设置为 true 时支持主动从设备方接收数据信息  Gets or sets whether the current pipeline activates the function of actively pushing data from the device. If this is set to true, it supports actively receiving data information from the device |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CommunicationLock 属性 

[原文連結](http://api.hslcommunication.cn/html/0d03f5cf-c63c-903c-31d2-9261b91de976.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 属性](../html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm "CommunicationPipe 属性")

[CommunicationLock 属性](../html/0d03f5cf-c63c-903c-31d2-9261b91de976.htm "CommunicationLock 属性 ")

[DecideWhetherQAMessageFunction 属性](../html/bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm "DecideWhetherQAMessageFunction 属性 ")

[IsPersistentConnection 属性](../html/a43de804-55cd-3862-9898-efae0fde9936.htm "IsPersistentConnection 属性 ")

[ReceiveTimeOut 属性](../html/b5e80233-0023-db67-90b1-fbbeff60ce3d.htm "ReceiveTimeOut 属性 ")

[SleepTime 属性](../html/1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm "SleepTime 属性 ")

[UseServerActivePush 属性](../html/b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm "UseServerActivePush 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeCommunicationLock 属性 |

获取或设置当前管道的线程锁对象，默认是简单的一个互斥锁  
Gets or sets the thread lock object of the current pipeline, which defaults to a simple mutex

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ICommunicationLock CommunicationLock { get; set; }
```

```
Public Property CommunicationLock As ICommunicationLock
	Get
	Set
```

```
public:
property ICommunicationLock^ CommunicationLock {
	ICommunicationLock^ get ();
	void set (ICommunicationLock^ value);
}
```

```
member CommunicationLock : ICommunicationLock with get, set
```

#### 属性值

类型：[ICommunicationLock](8b1e90d6-b81e-2e12-00b6-c959ab1b7821.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecideWhetherQAMessageFunction 属性 

[原文連結](http://api.hslcommunication.cn/html/bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 属性](../html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm "CommunicationPipe 属性")

[CommunicationLock 属性](../html/0d03f5cf-c63c-903c-31d2-9261b91de976.htm "CommunicationLock 属性 ")

[DecideWhetherQAMessageFunction 属性](../html/bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm "DecideWhetherQAMessageFunction 属性 ")

[IsPersistentConnection 属性](../html/a43de804-55cd-3862-9898-efae0fde9936.htm "IsPersistentConnection 属性 ")

[ReceiveTimeOut 属性](../html/b5e80233-0023-db67-90b1-fbbeff60ce3d.htm "ReceiveTimeOut 属性 ")

[SleepTime 属性](../html/1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm "SleepTime 属性 ")

[UseServerActivePush 属性](../html/b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm "UseServerActivePush 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeDecideWhetherQAMessageFunction 属性 |

用来决定当前接收的消息是否是问答服务的消息，可由外界设置委托来决定  
It is used to determine whether the currently received message is a message of the question answering service, which can be determined by the external setting delegate

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Func<CommunicationPipe, OperateResult<byte[]>, bool> DecideWhetherQAMessageFunction { get; set; }
```

```
Public Property DecideWhetherQAMessageFunction As Func(Of CommunicationPipe, OperateResult(Of Byte()), Boolean)
	Get
	Set
```

```
public:
property Func<CommunicationPipe^, OperateResult<array<unsigned char>^>^, bool>^ DecideWhetherQAMessageFunction {
	Func<CommunicationPipe^, OperateResult<array<unsigned char>^>^, bool>^ get ();
	void set (Func<CommunicationPipe^, OperateResult<array<unsigned char>^>^, bool>^ value);
}
```

```
member DecideWhetherQAMessageFunction : Func<CommunicationPipe, OperateResult<byte[]>, bool> with get, set
```

#### 属性值

类型：Func[CommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm), [OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte, Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsPersistentConnection 属性 

[原文連結](http://api.hslcommunication.cn/html/a43de804-55cd-3862-9898-efae0fde9936.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 属性](../html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm "CommunicationPipe 属性")

[CommunicationLock 属性](../html/0d03f5cf-c63c-903c-31d2-9261b91de976.htm "CommunicationLock 属性 ")

[DecideWhetherQAMessageFunction 属性](../html/bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm "DecideWhetherQAMessageFunction 属性 ")

[IsPersistentConnection 属性](../html/a43de804-55cd-3862-9898-efae0fde9936.htm "IsPersistentConnection 属性 ")

[ReceiveTimeOut 属性](../html/b5e80233-0023-db67-90b1-fbbeff60ce3d.htm "ReceiveTimeOut 属性 ")

[SleepTime 属性](../html/1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm "SleepTime 属性 ")

[UseServerActivePush 属性](../html/b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm "UseServerActivePush 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeIsPersistentConnection 属性 |

获取或设置当前的管道是否是长连接，仅对于串口及TCP是有效的，默认都是长连接  
Gets or sets whether the current pipe is a long connection. This is valid only for serial ports and TCP. The default is a long connection

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsPersistentConnection { get; set; }
```

```
Public Property IsPersistentConnection As Boolean
	Get
	Set
```

```
public:
property bool IsPersistentConnection {
	bool get ();
	void set (bool value);
}
```

```
member IsPersistentConnection : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReceiveTimeOut 属性 

[原文連結](http://api.hslcommunication.cn/html/b5e80233-0023-db67-90b1-fbbeff60ce3d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 属性](../html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm "CommunicationPipe 属性")

[CommunicationLock 属性](../html/0d03f5cf-c63c-903c-31d2-9261b91de976.htm "CommunicationLock 属性 ")

[DecideWhetherQAMessageFunction 属性](../html/bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm "DecideWhetherQAMessageFunction 属性 ")

[IsPersistentConnection 属性](../html/a43de804-55cd-3862-9898-efae0fde9936.htm "IsPersistentConnection 属性 ")

[ReceiveTimeOut 属性](../html/b5e80233-0023-db67-90b1-fbbeff60ce3d.htm "ReceiveTimeOut 属性 ")

[SleepTime 属性](../html/1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm "SleepTime 属性 ")

[UseServerActivePush 属性](../html/b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm "UseServerActivePush 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeReceiveTimeOut 属性 |

获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   
Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ReceiveTimeOut { get; set; }
```

```
Public Property ReceiveTimeOut As Integer
	Get
	Set
```

```
public:
property int ReceiveTimeOut {
	int get ();
	void set (int value);
}
```

```
member ReceiveTimeOut : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)备注

超时的通常原因是服务器端没有配置好，导致访问失败，为了不卡死软件，所以有了这个超时的属性。

![](../icons/SectionExpanded.png)示例

设置1秒的接收超时的示例

ReceiveTimeOut示例

[复制](# "复制")

```
// 设置反馈的超时时间，单位 毫秒
client.ReceiveTimeOut = 1000;
// 1秒没有接收到就自动返回失败，此处的地址示例是modbus的地址，对于读取也是一样的
OperateResult write = client.Write( "100", 123 );
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SleepTime 属性 

[原文連結](http://api.hslcommunication.cn/html/1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 属性](../html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm "CommunicationPipe 属性")

[CommunicationLock 属性](../html/0d03f5cf-c63c-903c-31d2-9261b91de976.htm "CommunicationLock 属性 ")

[DecideWhetherQAMessageFunction 属性](../html/bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm "DecideWhetherQAMessageFunction 属性 ")

[IsPersistentConnection 属性](../html/a43de804-55cd-3862-9898-efae0fde9936.htm "IsPersistentConnection 属性 ")

[ReceiveTimeOut 属性](../html/b5e80233-0023-db67-90b1-fbbeff60ce3d.htm "ReceiveTimeOut 属性 ")

[SleepTime 属性](../html/1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm "SleepTime 属性 ")

[UseServerActivePush 属性](../html/b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm "UseServerActivePush 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeSleepTime 属性 |

获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  
Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required.

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int SleepTime { get; set; }
```

```
Public Property SleepTime As Integer
	Get
	Set
```

```
public:
property int SleepTime {
	int get ();
	void set (int value);
}
```

```
member SleepTime : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UseServerActivePush 属性 

[原文連結](http://api.hslcommunication.cn/html/b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 属性](../html/afa3007a-06ff-2a1c-c14e-a8ac693e44c1.htm "CommunicationPipe 属性")

[CommunicationLock 属性](../html/0d03f5cf-c63c-903c-31d2-9261b91de976.htm "CommunicationLock 属性 ")

[DecideWhetherQAMessageFunction 属性](../html/bec87c3f-bba6-42ec-7056-e3174a8e46b6.htm "DecideWhetherQAMessageFunction 属性 ")

[IsPersistentConnection 属性](../html/a43de804-55cd-3862-9898-efae0fde9936.htm "IsPersistentConnection 属性 ")

[ReceiveTimeOut 属性](../html/b5e80233-0023-db67-90b1-fbbeff60ce3d.htm "ReceiveTimeOut 属性 ")

[SleepTime 属性](../html/1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm "SleepTime 属性 ")

[UseServerActivePush 属性](../html/b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm "UseServerActivePush 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeUseServerActivePush 属性 |

获取或设置当前的管道是否激活从设备主动推送的功能，设置为 true 时支持主动从设备方接收数据信息  
Gets or sets whether the current pipeline activates the function of actively pushing data from the device. If this is set to true, it supports actively receiving data information from the device

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool UseServerActivePush { get; set; }
```

```
Public Property UseServerActivePush As Boolean
	Get
	Set
```

```
public:
property bool UseServerActivePush {
	bool get ();
	void set (bool value);
}
```

```
member UseServerActivePush : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CommunicationPipe 方法

[原文連結](http://api.hslcommunication.cn/html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipe 方法 |

[CommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [CheckMessageComplete](eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm) | 根据给定的消息，发送的数据，接收到数据来判断是否接收完成报文  According to the given message, sent data, received data to determine whether to receive the completed message |
| 公共方法 | [CloseCommunication](9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm) | 关闭当前的管道信息，返回是否关闭成功的结果对象  Close the current pipeline information and return the result object whether the closure was successful |
| 公共方法 | [CloseCommunicationAsync](3b73e48c-3a64-eb2c-a827-fedd9926558e.htm) |  |
| 公共方法 | [Dispose](fb51c8a3-b4bc-5c64-b1ba-742464b90b41.htm) |  |
| 受保护的方法 | [Dispose(Boolean)](94cf19e0-0b6e-a57f-7559-27bf48a8e8f9.htm) |  |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [HasCacheData](3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm) | 当前管道的缓存里是否有数据存在，有就返回 True  If there is any data in the cache of the current pipeline, return True |
| 受保护的方法 | [IncrConnectErrorCount](946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm) | 自增当前的连续错误计数，并且获取自增后的值信息，最大到20亿为止，无法继续增加了。  Increment the current continuous error count, and obtain the value information after increment, up to 2 billion, can not continue to increase. |
| 公共方法 | [IsConnectError](dba07b04-a27a-754d-1f88-053ef60251f3.htm) | 当前的管道连接对象是否发生了错误，这里的错误通常是串口，或是网络错误  Whether there is an error in the current pipe connection object, where the error is usually a serial port, or a network error |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [OpenCommunication](009e7e59-f95f-29ad-7545-b691bf676e07.htm) | 打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  Opens the current pipe information, returns whether the result object was successfully opened, and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true |
| 公共方法 | [OpenCommunicationAsync](4623d60d-297f-edd6-91df-892f51ab4609.htm) | 打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  Opens the current pipe information, returns whether the result object was successfully opened, and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true |
| 公共方法 | [RaisePipeError](32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm) | 主动引发一个管道错误，从而让管道可以重新打开  Actively causes a pipe error so that the pipe can be reopened |
| 公共方法 | [ReadFromCoreServer](a85930ab-0e17-beee-cc5b-1178ce83f78c.htm) | 将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type, the data sent, whether there is a data response, and the sleep time |
| 公共方法 | [ReadFromCoreServerAsync](3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm) | 将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type, the data sent, whether there is a data response, and the sleep time |
| 受保护的方法 | [ReadFromCoreServerHelper](e208a32f-401a-20af-86a7-3972e220b425.htm) | 将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type, the data sent, whether there is a data response, and the sleep time |
| 受保护的方法 | [ReadFromCoreServerHelperAsync](7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm) | 将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type, the data sent, whether there is a data response, and the sleep time |
| 公共方法 | [Receive(Int32, Int32, ActionInt64, Int64)](6a8d9448-097e-4cb0-0742-95327d4431d3.htm) | 从管道里，接收指定长度的报文数据信息，如果长度指定为-1，表示接收不超过2048字节的动态长度。另外可以指定超时时间，进度报告等  Receives the packet data of a specified length from the pipe. If the length is set to -1, it indicates that the dynamic length of the packet is not more than 2048 bytes. You can also specify timeouts, progress reports, etc |
| 公共方法 | [Receive(Byte, Int32, Int32, Int32, ActionInt64, Int64)](9979783b-37d4-3993-b8b8-aed6055124bf.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. |
| 公共方法 | [ReceiveAsync(Int32, Int32, ActionInt64, Int64)](b56c1109-0eed-300f-ec21-6217d823f3cf.htm) | 从管道里，接收指定长度的报文数据信息，如果长度指定为-1，表示接收不超过2048字节的动态长度。另外可以指定超时时间，进度报告等  Receives the packet data of a specified length from the pipe. If the length is set to -1, it indicates that the dynamic length of the packet is not more than 2048 bytes. You can also specify timeouts, progress reports, etc |
| 公共方法 | [ReceiveAsync(Byte, Int32, Int32, Int32, ActionInt64, Int64)](2eb79f82-e65a-51c6-4b57-8713600bb27e.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. |
| 公共方法 | [ReceiveMessage](fa579a87-d719-3fef-af36-a3d37af08200.htm) | 包含了一个复杂的逻辑，从管道里根据当前的消息格式定义，接收报文信息，这个报文可能是来自服务器主动推送的。具体可以通过参数 useActivePush 来特殊控制。  Contains a complex logic from the pipeline, according to the current message format definition, to receive message information, this message may be actively pushed from the server. The parameter useActivePush can be used for special control. |
| 公共方法 | [ReceiveMessageAsync](c2e7239d-8839-c1ed-6993-f4da7a615b78.htm) | 包含了一个复杂的逻辑，从管道里根据当前的消息格式定义，接收报文信息，这个报文可能是来自服务器主动推送的。具体可以通过参数 useActivePush 来特殊控制。  Contains a complex logic from the pipeline, according to the current message format definition, to receive message information, this message may be actively pushed from the server. The parameter useActivePush can be used for special control. |
| 公共方法 | [ResetConnectErrorCount](c3d25cde-83f6-4257-8344-0e458c8bd685.htm) | 重置当前的连续错误计数为0，并且返回重置前时候的值  Resets the current consecutive error count to 0 and returns the value before the reset |
| 公共方法 | [Send(Byte)](f71f0371-1270-df68-9648-4919974c4fef.htm) | 发送数据到当前的管道中去  Send data to the current pipe |
| 公共方法 | [Send(Byte, Int32, Int32)](45f67684-08ea-f88c-2dff-8e6701e70996.htm) | 将一个数据缓存中的指定的部分字段，发送到当前的管道中去  Sends the specified partial field from a data cache to the current pipeline |
| 公共方法 | [SendAsync(Byte)](79e3ebb0-fd3d-11c1-53ca-00970cc89c25.htm) | 发送数据到当前的管道中去  Send data to the current pipe |
| 公共方法 | [SendAsync(Byte, Int32, Int32)](e3cd5a7d-1e61-0d44-c5f2-81e9e9606cd8.htm) | 将一个数据缓存中的指定的部分字段，发送到当前的管道中去  Sends the specified partial field from a data cache to the current pipeline |
| 受保护的方法 | [SetBufferQA](6677295f-e375-da2e-1931-ec881676e83f.htm) | 设置当前的问答状态下的缓存数据  Set the cache data in the current Q&A state |
| 公共方法 | [StartReceiveBackground](9b977dfb-5c8b-d6d9-7c94-443345235b26.htm) | 开始后台接收相关的报文数据，当[UseServerActivePush](b9fdf7aa-5fd9-1247-4238-ddaf720f84ab.htm)为True时，则使用本方法 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckMessageComplete 方法 

[原文連結](http://api.hslcommunication.cn/html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeCheckMessageComplete 方法 |

根据给定的消息，发送的数据，接收到数据来判断是否接收完成报文  
According to the given message, sent data, received data to determine whether to receive the completed message

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected bool CheckMessageComplete(
	INetMessage netMessage,
	byte[] sendValue,
	ref MemoryStream ms
)
```

```
Protected Function CheckMessageComplete ( 
	netMessage As INetMessage,
	sendValue As Byte(),
	ByRef ms As MemoryStream
) As Boolean
```

```
protected:
bool CheckMessageComplete(
	INetMessage^ netMessage, 
	array<unsigned char>^ sendValue, 
	MemoryStream^% ms
)
```

```
member CheckMessageComplete : 
        netMessage : INetMessage * 
        sendValue : byte[] * 
        ms : MemoryStream byref -> bool 
```

#### 参数

netMessage
:   类型：[HslCommunication.Core.IMessageINetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)  
    消息类对象

sendValue
:   类型：SystemByte  
    发送的数据内容

ms
:   类型：System.IOMemoryStream  
    接收数据的流

#### 返回值

类型：Boolean  
是否接收完成数据

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CloseCommunication 方法 

[原文連結](http://api.hslcommunication.cn/html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeCloseCommunication 方法 |

关闭当前的管道信息，返回是否关闭成功的结果对象  
Close the current pipeline information and return the result object whether the closure was successful

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult CloseCommunication()
```

```
Public Overridable Function CloseCommunication As OperateResult
```

```
public:
virtual OperateResult^ CloseCommunication()
```

```
abstract CloseCommunication : unit -> OperateResult 
override CloseCommunication : unit -> OperateResult
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否关闭成功

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CloseCommunicationAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeCloseCommunicationAsync 方法 |

[缺少 "M:HslCommunication.Core.Pipe.CommunicationPipe.CloseCommunicationAsync" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual Task<OperateResult> CloseCommunicationAsync()
```

```
Public Overridable Function CloseCommunicationAsync As Task(Of OperateResult)
```

```
public:
virtual Task<OperateResult^>^ CloseCommunicationAsync()
```

```
abstract CloseCommunicationAsync : unit -> Task<OperateResult> 
override CloseCommunicationAsync : unit -> Task<OperateResult>
```

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  

[缺少 "M:HslCommunication.Core.Pipe.CommunicationPipe.CloseCommunicationAsync" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[Dispose 方法](../html/fb51c8a3-b4bc-5c64-b1ba-742464b90b41.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/94cf19e0-0b6e-a57f-7559-27bf48a8e8f9.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeDispose 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Dispose](fb51c8a3-b4bc-5c64-b1ba-742464b90b41.htm) | 释放被 [CommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm) 使用的所有资源 |
| 受保护的方法 | [Dispose(Boolean)](94cf19e0-0b6e-a57f-7559-27bf48a8e8f9.htm) | 释放被 [CommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm) 使用的非托管资源，并且是否托管资源（可选） |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/fb51c8a3-b4bc-5c64-b1ba-742464b90b41.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[Dispose 方法](../html/fb51c8a3-b4bc-5c64-b1ba-742464b90b41.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/94cf19e0-0b6e-a57f-7559-27bf48a8e8f9.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeDispose 方法 |

释放被 [CommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm) 使用的所有资源

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void Dispose()
```

```
Public Sub Dispose
```

```
public:
virtual void Dispose() sealed
```

```
abstract Dispose : unit -> unit 
override Dispose : unit -> unit
```

#### 实现

IDisposableDispose

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[Dispose 重载](52f51aa9-905b-cd8b-843e-23d53be3dd95.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 (Boolean)

[原文連結](http://api.hslcommunication.cn/html/94cf19e0-0b6e-a57f-7559-27bf48a8e8f9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[Dispose 方法](../html/fb51c8a3-b4bc-5c64-b1ba-742464b90b41.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/94cf19e0-0b6e-a57f-7559-27bf48a8e8f9.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeDispose 方法 (Boolean) |

释放被 [CommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm) 使用的非托管资源，并且是否托管资源（可选）

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected virtual void Dispose(
	bool disposing
)
```

```
Protected Overridable Sub Dispose ( 
	disposing As Boolean
)
```

```
protected:
virtual void Dispose(
	bool disposing
)
```

```
abstract Dispose : 
        disposing : bool -> unit 
override Dispose : 
        disposing : bool -> unit
```

#### 参数

disposing
:   类型：SystemBoolean  
    为 true 则同时释放托管资源和非托管资源；为 false 则只释放非托管资源

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[Dispose 重载](52f51aa9-905b-cd8b-843e-23d53be3dd95.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HasCacheData 方法 

[原文連結](http://api.hslcommunication.cn/html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeHasCacheData 方法 |

当前管道的缓存里是否有数据存在，有就返回 True  
If there is any data in the cache of the current pipeline, return True

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual bool HasCacheData()
```

```
Public Overridable Function HasCacheData As Boolean
```

```
public:
virtual bool HasCacheData()
```

```
abstract HasCacheData : unit -> bool 
override HasCacheData : unit -> bool
```

#### 返回值

类型：Boolean  
是否有缓存数据

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IncrConnectErrorCount 方法 

[原文連結](http://api.hslcommunication.cn/html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeIncrConnectErrorCount 方法 |

自增当前的连续错误计数，并且获取自增后的值信息，最大到20亿为止，无法继续增加了。  
Increment the current continuous error count, and obtain the value information after increment, up to 2 billion, can not continue to increase.

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected int IncrConnectErrorCount()
```

```
Protected Function IncrConnectErrorCount As Integer
```

```
protected:
int IncrConnectErrorCount()
```

```
member IncrConnectErrorCount : unit -> int 
```

#### 返回值

类型：Int32  
自增后的值信息

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsConnectError 方法 

[原文連結](http://api.hslcommunication.cn/html/dba07b04-a27a-754d-1f88-053ef60251f3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeIsConnectError 方法 |

当前的管道连接对象是否发生了错误，这里的错误通常是串口，或是网络错误  
Whether there is an error in the current pipe connection object, where the error is usually a serial port, or a network error

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual bool IsConnectError()
```

```
Public Overridable Function IsConnectError As Boolean
```

```
public:
virtual bool IsConnectError()
```

```
abstract IsConnectError : unit -> bool 
override IsConnectError : unit -> bool
```

#### 返回值

类型：Boolean  
是否发生了通道的异常

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OpenCommunication 方法 

[原文連結](http://api.hslcommunication.cn/html/009e7e59-f95f-29ad-7545-b691bf676e07.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeOpenCommunication 方法 |

打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  
Opens the current pipe information, returns whether the result object was successfully opened,
and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult<bool> OpenCommunication()
```

```
Public Overridable Function OpenCommunication As OperateResult(Of Boolean)
```

```
public:
virtual OperateResult<bool>^ OpenCommunication()
```

```
abstract OpenCommunication : unit -> OperateResult<bool> 
override OpenCommunication : unit -> OperateResult<bool>
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
是否打开成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OpenCommunicationAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/4623d60d-297f-edd6-91df-892f51ab4609.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeOpenCommunicationAsync 方法 |

打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  
Opens the current pipe information, returns whether the result object was successfully opened,
and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual Task<OperateResult<bool>> OpenCommunicationAsync()
```

```
Public Overridable Function OpenCommunicationAsync As Task(Of OperateResult(Of Boolean))
```

```
public:
virtual Task<OperateResult<bool>^>^ OpenCommunicationAsync()
```

```
abstract OpenCommunicationAsync : unit -> Task<OperateResult<bool>> 
override OpenCommunicationAsync : unit -> Task<OperateResult<bool>>
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
是否打开成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RaisePipeError 方法 

[原文連結](http://api.hslcommunication.cn/html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeRaisePipeError 方法 |

主动引发一个管道错误，从而让管道可以重新打开  
Actively causes a pipe error so that the pipe can be reopened

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void RaisePipeError()
```

```
Public Sub RaisePipeError
```

```
public:
void RaisePipeError()
```

```
member RaisePipeError : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServer 方法 

[原文連結](http://api.hslcommunication.cn/html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeReadFromCoreServer 方法 |

将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  
To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type,
the data sent, whether there is a data response, and the sleep time

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult<byte[]> ReadFromCoreServer(
	INetMessage netMessage,
	byte[] sendValue,
	bool hasResponseData,
	Action<byte[]> logMessage = null
)
```

```
Public Overridable Function ReadFromCoreServer ( 
	netMessage As INetMessage,
	sendValue As Byte(),
	hasResponseData As Boolean,
	Optional logMessage As Action(Of Byte()) = Nothing
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ ReadFromCoreServer(
	INetMessage^ netMessage, 
	array<unsigned char>^ sendValue, 
	bool hasResponseData, 
	Action<array<unsigned char>^>^ logMessage = nullptr
)
```

```
abstract ReadFromCoreServer : 
        netMessage : INetMessage * 
        sendValue : byte[] * 
        hasResponseData : bool * 
        ?logMessage : Action<byte[]> 
(* Defaults:
        let _logMessage = defaultArg logMessage null
*)
-> OperateResult<byte[]> 
override ReadFromCoreServer : 
        netMessage : INetMessage * 
        sendValue : byte[] * 
        hasResponseData : bool * 
        ?logMessage : Action<byte[]> 
(* Defaults:
        let _logMessage = defaultArg logMessage null
*)
-> OperateResult<byte[]>
```

#### 参数

netMessage
:   类型：[HslCommunication.Core.IMessageINetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)  
    当前接收的消息体信息

sendValue
:   类型：SystemByte  
    等待发送的数据

hasResponseData
:   类型：SystemBoolean  
    是否有数据返回

logMessage (Optional)
:   类型：SystemActionByte  
    用于消息记录的日志信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServerAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeReadFromCoreServerAsync 方法 |

将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  
To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type,
the data sent, whether there is a data response, and the sleep time

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual Task<OperateResult<byte[]>> ReadFromCoreServerAsync(
	INetMessage netMessage,
	byte[] sendValue,
	bool hasResponseData,
	Action<byte[]> logMessage = null
)
```

```
Public Overridable Function ReadFromCoreServerAsync ( 
	netMessage As INetMessage,
	sendValue As Byte(),
	hasResponseData As Boolean,
	Optional logMessage As Action(Of Byte()) = Nothing
) As Task(Of OperateResult(Of Byte()))
```

```
public:
virtual Task<OperateResult<array<unsigned char>^>^>^ ReadFromCoreServerAsync(
	INetMessage^ netMessage, 
	array<unsigned char>^ sendValue, 
	bool hasResponseData, 
	Action<array<unsigned char>^>^ logMessage = nullptr
)
```

```
abstract ReadFromCoreServerAsync : 
        netMessage : INetMessage * 
        sendValue : byte[] * 
        hasResponseData : bool * 
        ?logMessage : Action<byte[]> 
(* Defaults:
        let _logMessage = defaultArg logMessage null
*)
-> Task<OperateResult<byte[]>> 
override ReadFromCoreServerAsync : 
        netMessage : INetMessage * 
        sendValue : byte[] * 
        hasResponseData : bool * 
        ?logMessage : Action<byte[]> 
(* Defaults:
        let _logMessage = defaultArg logMessage null
*)
-> Task<OperateResult<byte[]>>
```

#### 参数

netMessage
:   类型：[HslCommunication.Core.IMessageINetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)  
    当前接收的消息体信息

sendValue
:   类型：SystemByte  
    等待发送的数据

hasResponseData
:   类型：SystemBoolean  
    是否有数据返回

logMessage (Optional)
:   类型：SystemActionByte  
    用于消息记录的日志信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServerHelper 方法 

[原文連結](http://api.hslcommunication.cn/html/e208a32f-401a-20af-86a7-3972e220b425.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeReadFromCoreServerHelper 方法 |

将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  
To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type,
the data sent, whether there is a data response, and the sleep time

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected OperateResult<byte[]> ReadFromCoreServerHelper(
	INetMessage netMessage,
	byte[] sendValue,
	bool hasResponseData,
	int sleep,
	Action<byte[]> logMessage = null
)
```

```
Protected Function ReadFromCoreServerHelper ( 
	netMessage As INetMessage,
	sendValue As Byte(),
	hasResponseData As Boolean,
	sleep As Integer,
	Optional logMessage As Action(Of Byte()) = Nothing
) As OperateResult(Of Byte())
```

```
protected:
OperateResult<array<unsigned char>^>^ ReadFromCoreServerHelper(
	INetMessage^ netMessage, 
	array<unsigned char>^ sendValue, 
	bool hasResponseData, 
	int sleep, 
	Action<array<unsigned char>^>^ logMessage = nullptr
)
```

```
member ReadFromCoreServerHelper : 
        netMessage : INetMessage * 
        sendValue : byte[] * 
        hasResponseData : bool * 
        sleep : int * 
        ?logMessage : Action<byte[]> 
(* Defaults:
        let _logMessage = defaultArg logMessage null
*)
-> OperateResult<byte[]> 
```

#### 参数

netMessage
:   类型：[HslCommunication.Core.IMessageINetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)  
    当前接收的消息体信息

sendValue
:   类型：SystemByte  
    等待发送的数据

hasResponseData
:   类型：SystemBoolean  
    是否有数据返回

sleep
:   类型：SystemInt32  
    休眠时间

logMessage (Optional)
:   类型：SystemActionByte  
    用于消息记录的日志信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServerHelperAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[CheckMessageComplete 方法](../html/eef95d76-bcfd-3add-5a7e-27c3bfeda0d7.htm "CheckMessageComplete 方法 ")

[CloseCommunication 方法](../html/9d6768b9-7e2e-ef7a-80e9-c5cc51720b57.htm "CloseCommunication 方法 ")

[CloseCommunicationAsync 方法](../html/3b73e48c-3a64-eb2c-a827-fedd9926558e.htm "CloseCommunicationAsync 方法 ")

[Dispose 方法](../html/52f51aa9-905b-cd8b-843e-23d53be3dd95.htm "Dispose 方法 ")

[HasCacheData 方法](../html/3fd86f4a-0afb-430e-daf4-f8f45f7e7b3d.htm "HasCacheData 方法 ")

[IncrConnectErrorCount 方法](../html/946b05e2-0c30-6480-7f4b-27cd5a4952c3.htm "IncrConnectErrorCount 方法 ")

[IsConnectError 方法](../html/dba07b04-a27a-754d-1f88-053ef60251f3.htm "IsConnectError 方法 ")

[OpenCommunication 方法](../html/009e7e59-f95f-29ad-7545-b691bf676e07.htm "OpenCommunication 方法 ")

[OpenCommunicationAsync 方法](../html/4623d60d-297f-edd6-91df-892f51ab4609.htm "OpenCommunicationAsync 方法 ")

[RaisePipeError 方法](../html/32e1ec57-1e8d-9083-b3ac-2ef2dbe63227.htm "RaisePipeError 方法 ")

[ReadFromCoreServer 方法](../html/a85930ab-0e17-beee-cc5b-1178ce83f78c.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/3e19b13c-851a-9f6f-bae6-9c7719e66ec3.htm "ReadFromCoreServerAsync 方法 ")

[ReadFromCoreServerHelper 方法](../html/e208a32f-401a-20af-86a7-3972e220b425.htm "ReadFromCoreServerHelper 方法 ")

[ReadFromCoreServerHelperAsync 方法](../html/7617d42e-8c6d-bf3d-ed62-45e0d3e019f4.htm "ReadFromCoreServerHelperAsync 方法 ")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[ReceiveAsync 方法](../html/c5eebe6c-027f-1719-4a75-8d682bb89173.htm "ReceiveAsync 方法 ")

[ReceiveMessage 方法](../html/fa579a87-d719-3fef-af36-a3d37af08200.htm "ReceiveMessage 方法 ")

[ReceiveMessageAsync 方法](../html/c2e7239d-8839-c1ed-6993-f4da7a615b78.htm "ReceiveMessageAsync 方法 ")

[ResetConnectErrorCount 方法](../html/c3d25cde-83f6-4257-8344-0e458c8bd685.htm "ResetConnectErrorCount 方法 ")

[Send 方法](../html/58d3111c-8e43-c19a-ccfc-d3828b876ccd.htm "Send 方法 ")

[SendAsync 方法](../html/c59e0211-9380-4e61-09c7-cc38a4af8af1.htm "SendAsync 方法 ")

[SetBufferQA 方法](../html/6677295f-e375-da2e-1931-ec881676e83f.htm "SetBufferQA 方法 ")

[StartReceiveBackground 方法](../html/9b977dfb-5c8b-d6d9-7c94-443345235b26.htm "StartReceiveBackground 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeReadFromCoreServerHelperAsync 方法 |

将数据发送到当前的管道里，并从管道接收相关的数据信息，可以指定消息类型，发送数据，是否有数据响应，休眠时间  
To send data to the current pipeline and receive relevant data information from the pipeline, you can specify the message type,
the data sent, whether there is a data response, and the sleep time

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected Task<OperateResult<byte[]>> ReadFromCoreServerHelperAsync(
	INetMessage netMessage,
	byte[] sendValue,
	bool hasResponseData,
	int sleep,
	Action<byte[]> logMessage = null
)
```

```
Protected Function ReadFromCoreServerHelperAsync ( 
	netMessage As INetMessage,
	sendValue As Byte(),
	hasResponseData As Boolean,
	sleep As Integer,
	Optional logMessage As Action(Of Byte()) = Nothing
) As Task(Of OperateResult(Of Byte()))
```

```
protected:
Task<OperateResult<array<unsigned char>^>^>^ ReadFromCoreServerHelperAsync(
	INetMessage^ netMessage, 
	array<unsigned char>^ sendValue, 
	bool hasResponseData, 
	int sleep, 
	Action<array<unsigned char>^>^ logMessage = nullptr
)
```

```
member ReadFromCoreServerHelperAsync : 
        netMessage : INetMessage * 
        sendValue : byte[] * 
        hasResponseData : bool * 
        sleep : int * 
        ?logMessage : Action<byte[]> 
(* Defaults:
        let _logMessage = defaultArg logMessage null
*)
-> Task<OperateResult<byte[]>> 
```

#### 参数

netMessage
:   类型：[HslCommunication.Core.IMessageINetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)  
    当前接收的消息体信息

sendValue
:   类型：SystemByte  
    等待发送的数据

hasResponseData
:   类型：SystemBoolean  
    是否有数据返回

sleep
:   类型：SystemInt32  
    休眠时间

logMessage (Optional)
:   类型：SystemActionByte  
    用于消息记录的日志信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Receive 方法 

[原文連結](http://api.hslcommunication.cn/html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[Receive 方法 (Int32, Int32, Action(Int64, Int64))](../html/6a8d9448-097e-4cb0-0742-95327d4431d3.htm "Receive 方法 (Int32, Int32, Action(Int64, Int64))")

[Receive 方法 (Byte[], Int32, Int32, Int32, Action(Int64, Int64))](../html/9979783b-37d4-3993-b8b8-aed6055124bf.htm "Receive 方法 (Byte[], Int32, Int32, Int32, Action(Int64, Int64))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeReceive 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Receive(Int32, Int32, ActionInt64, Int64)](6a8d9448-097e-4cb0-0742-95327d4431d3.htm) | 从管道里，接收指定长度的报文数据信息，如果长度指定为-1，表示接收不超过2048字节的动态长度。另外可以指定超时时间，进度报告等  Receives the packet data of a specified length from the pipe. If the length is set to -1, it indicates that the dynamic length of the packet is not more than 2048 bytes. You can also specify timeouts, progress reports, etc |
| 公共方法 | [Receive(Byte, Int32, Int32, Int32, ActionInt64, Int64)](9979783b-37d4-3993-b8b8-aed6055124bf.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Receive 方法 (Int32, Int32, Action(Int64, Int64))

[原文連結](http://api.hslcommunication.cn/html/6a8d9448-097e-4cb0-0742-95327d4431d3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[Receive 方法 (Int32, Int32, Action(Int64, Int64))](../html/6a8d9448-097e-4cb0-0742-95327d4431d3.htm "Receive 方法 (Int32, Int32, Action(Int64, Int64))")

[Receive 方法 (Byte[], Int32, Int32, Int32, Action(Int64, Int64))](../html/9979783b-37d4-3993-b8b8-aed6055124bf.htm "Receive 方法 (Byte[], Int32, Int32, Int32, Action(Int64, Int64))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeReceive 方法 (Int32, Int32, ActionInt64, Int64) |

从管道里，接收指定长度的报文数据信息，如果长度指定为-1，表示接收不超过2048字节的动态长度。另外可以指定超时时间，进度报告等  
Receives the packet data of a specified length from the pipe. If the length is set to -1,
it indicates that the dynamic length of the packet is not more than 2048 bytes. You can also specify timeouts, progress reports, etc

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult<byte[]> Receive(
	int length,
	int timeOut,
	Action<long, long> reportProgress = null
)
```

```
Public Overridable Function Receive ( 
	length As Integer,
	timeOut As Integer,
	Optional reportProgress As Action(Of Long, Long) = Nothing
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ Receive(
	int length, 
	int timeOut, 
	Action<long long, long long>^ reportProgress = nullptr
)
```

```
abstract Receive : 
        length : int * 
        timeOut : int * 
        ?reportProgress : Action<int64, int64> 
(* Defaults:
        let _reportProgress = defaultArg reportProgress null
*)
-> OperateResult<byte[]> 
override Receive : 
        length : int * 
        timeOut : int * 
        ?reportProgress : Action<int64, int64> 
(* Defaults:
        let _reportProgress = defaultArg reportProgress null
*)
-> OperateResult<byte[]>
```

#### 参数

length
:   类型：SystemInt32  
    接收的长度信息

timeOut
:   类型：SystemInt32  
    指定的超时时间

reportProgress (Optional)
:   类型：SystemActionInt64, Int64  
    进行进度报告的委托

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
是否接收成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[Receive 重载](4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Receive 方法 (Byte[], Int32, Int32, Int32, Action(Int64, Int64))

[原文連結](http://api.hslcommunication.cn/html/9979783b-37d4-3993-b8b8-aed6055124bf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Pipe](../html/cb25db69-9277-4ca9-66af-e850f7f1ac02.htm "HslCommunication.Core.Pipe")

[CommunicationPipe 类](../html/ab3f9c5c-144b-4209-4999-b382a1d94450.htm "CommunicationPipe 类")

[CommunicationPipe 方法](../html/432bd70d-59ee-4d18-9f13-9eba2bdf1133.htm "CommunicationPipe 方法")

[Receive 方法](../html/4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm "Receive 方法 ")

[Receive 方法 (Int32, Int32, Action(Int64, Int64))](../html/6a8d9448-097e-4cb0-0742-95327d4431d3.htm "Receive 方法 (Int32, Int32, Action(Int64, Int64))")

[Receive 方法 (Byte[], Int32, Int32, Int32, Action(Int64, Int64))](../html/9979783b-37d4-3993-b8b8-aed6055124bf.htm "Receive 方法 (Byte[], Int32, Int32, Int32, Action(Int64, Int64))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CommunicationPipeReceive 方法 (Byte, Int32, Int32, Int32, ActionInt64, Int64) |

接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  
Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0,
fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received.

**命名空间：**
 [HslCommunication.Core.Pipe](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult<int> Receive(
	byte[] buffer,
	int offset,
	int length,
	int timeOut = 60000,
	Action<long, long> reportProgress = null
)
```

```
Public Overridable Function Receive ( 
	buffer As Byte(),
	offset As Integer,
	length As Integer,
	Optional timeOut As Integer = 60000,
	Optional reportProgress As Action(Of Long, Long) = Nothing
) As OperateResult(Of Integer)
```

```
public:
virtual OperateResult<int>^ Receive(
	array<unsigned char>^ buffer, 
	int offset, 
	int length, 
	int timeOut = 60000, 
	Action<long long, long long>^ reportProgress = nullptr
)
```

```
abstract Receive : 
        buffer : byte[] * 
        offset : int * 
        length : int * 
        ?timeOut : int * 
        ?reportProgress : Action<int64, int64> 
(* Defaults:
        let _timeOut = defaultArg timeOut 60000
        let _reportProgress = defaultArg reportProgress null
*)
-> OperateResult<int> 
override Receive : 
        buffer : byte[] * 
        offset : int * 
        length : int * 
        ?timeOut : int * 
        ?reportProgress : Action<int64, int64> 
(* Defaults:
        let _timeOut = defaultArg timeOut 60000
        let _reportProgress = defaultArg reportProgress null
*)
-> OperateResult<int>
```

#### 参数

buffer
:   类型：SystemByte  
    等待接收的数据缓存信息

offset
:   类型：SystemInt32  
    开始接收数据的偏移地址

length
:   类型：SystemInt32  
    准备接收的数据长度，当length大于0时，接收固定长度的数据内容，当length小于0时，接收不大于2048长度的随机数据信息

timeOut (Optional)
:   类型：SystemInt32  
    单位：毫秒，超时时间，默认为60秒，如果设置小于0，则不检查超时时间

reportProgress (Optional)
:   类型：SystemActionInt64, Int64  
    进行进度报告的委托

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int32  
包含了字节数据的结果类

![](../icons/SectionExpanded.png)参见

#### 引用

[CommunicationPipe 类](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

[Receive 重载](4d32a32d-49a0-ad39-47b2-64a774a71bd6.htm)

[HslCommunication.Core.Pipe 命名空间](cb25db69-9277-4ca9-66af-e850f7f1ac02.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)