# HslCommunication - HslCommunication.Instrument.DLT

> 分類頁數: 30



---
## HslCommunication.Instrument.DLT

[原文連結](http://api.hslcommunication.cn/html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645OverTcp 类](../html/4f76bd69-c6bb-adad-7580-6d8a82d243e4.htm "DLT645OverTcp 类")

[DLT645Server 类](../html/3a23493f-ca42-c0a5-5b24-3a675c68107e.htm "DLT645Server 类")

[DLT645Server.DLTAddress 类](../html/39a170fa-c725-97e4-e8ae-e427c9964aa2.htm "DLT645Server.DLTAddress 类")

[DLT645With1997 类](../html/efa18af1-bb26-81d8-5005-fc768ec0cf26.htm "DLT645With1997 类")

[DLT645With1997OverTcp 类](../html/2bf87f0c-b13a-3c26-3171-35b6102d7a4f.htm "DLT645With1997OverTcp 类")

[DLT645With1997Server 类](../html/dbdbddc3-c561-5afa-a5d3-fb617f8e2a7d.htm "DLT645With1997Server 类")

[DLT698 类](../html/7435d2a0-3fdd-2442-8dd1-0bd354dbfa4f.htm "DLT698 类")

[DLT698OverTcp 类](../html/7f63a374-0a57-6ed3-35b6-6957ca4af380.htm "DLT698OverTcp 类")

[DLT698OverTcp.DataReportingDelegate 委托](../html/9bdb3a09-bfed-fd7b-d749-5f34d84d47a6.htm "DLT698OverTcp.DataReportingDelegate 委托")

[DLT698Server 类](../html/1b930240-6344-f3f4-f0e9-563be3abbe76.htm "DLT698Server 类")

[DLT698TcpNet 类](../html/5c555bdc-1220-93ff-20f2-5b886e790b46.htm "DLT698TcpNet 类")

[DLTControl 类](../html/65f784c4-6493-2ab7-1011-7e25717c0a92.htm "DLTControl 类")

[DLTTransform 类](../html/dacea70e-ec57-4c0e-8f93-66c637beab23.htm "DLTTransform 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Instrument.DLT 命名空间 |

[缺少 "N:HslCommunication.Instrument.DLT" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类代码示例 | [DLT645](dc3f4d48-9b60-0203-88c1-3543260dee17.htm) | 基于多功能电能表通信协议实现的通讯类，参考的文档是DLT645-2007，主要实现了对电表数据的读取和一些功能方法， 在点对点模式下，需要在打开串口后调用 [ReadAddress](16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm) 方法，数据标识格式为 00-00-00-00，具体参照文档手册。  The communication type based on the communication protocol of the multifunctional electric energy meter. The reference document is DLT645-2007, which mainly realizes the reading of the electric meter data and some functional methods. In the point-to-point mode, you need to call [ReadAddress](16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm) method after opening the serial port. the data identification format is 00-00-00-00, refer to the documentation manual for details. |
| 公共类代码示例 | [DLT645OverTcp](4f76bd69-c6bb-adad-7580-6d8a82d243e4.htm) | 基于多功能电能表通信协议实现的通讯类，参考的文档是DLT645-2007，主要实现了对电表数据的读取和一些功能方法， 在点对点模式下，需要在连接后调用 [ReadAddress](c01d2aa5-86f0-a56a-9fa2-4e4c163bfdd9.htm) 方法，数据标识格式为 00-00-00-00，具体参照文档手册。  The communication type based on the communication protocol of the multifunctional electric energy meter. The reference document is DLT645-2007, which mainly realizes the reading of the electric meter data and some functional methods. In the point-to-point mode, you need to call [ReadAddress](c01d2aa5-86f0-a56a-9fa2-4e4c163bfdd9.htm) method after connect the device. the data identification format is 00-00-00-00, refer to the documentation manual for details. |
| 公共类 | [DLT645Server](3a23493f-ca42-c0a5-5b24-3a675c68107e.htm) | DLT645协议的虚拟服务器 |
| 受保护的类 | [DLT645ServerDLTAddress](39a170fa-c725-97e4-e8ae-e427c9964aa2.htm) | DLT的地址信息 |
| 公共类代码示例 | [DLT645With1997](efa18af1-bb26-81d8-5005-fc768ec0cf26.htm) | 基于多功能电能表通信协议实现的通讯类，参考的文档是DLT645-1997，主要实现了对电表数据的读取和一些功能方法，数据标识格式为 B6-11，具体参照文档手册。  Based on the communication class implemented by the multi-function energy meter communication protocol, the reference document is DLT645-1997, which mainly implements the reading of meter data and some functional methods, the data identification format is B6-11, please refer to the document manual for details. |
| 公共类代码示例 | [DLT645With1997OverTcp](2bf87f0c-b13a-3c26-3171-35b6102d7a4f.htm) | 基于多功能电能表通信协议实现的通讯类，参考的文档是DLT645-1997，主要实现了对电表数据的读取和一些功能方法，数据标识格式为 B6-11，具体参照文档手册。  Based on the communication class implemented by the multi-function energy meter communication protocol, the reference document is DLT645-1997, which mainly implements the reading of meter data and some functional methods, the data identification format is B6-11, please refer to the document manual for details. |
| 公共类 | [DLT645With1997Server](dbdbddc3-c561-5afa-a5d3-fb617f8e2a7d.htm) | DLT645-1997协议的服务端实现类，主要用于与DLT645-1997设备进行通信。 |
| 公共类代码示例 | [DLT698](7435d2a0-3fdd-2442-8dd1-0bd354dbfa4f.htm) | 698.45协议的串口通信类，面向对象的用电信息数据交换协议，使用明文的通信方式。支持读取功率，总功，电压，电流，频率，功率因数等数据。  The serial communication class of the 698.45 protocol, an object-oriented power consumption information data exchange protocol, uses the communication method of clear text. Support reading power, total power, voltage, current, frequency, power factor and other data. |
| 公共类代码示例 | [DLT698OverTcp](7f63a374-0a57-6ed3-35b6-6957ca4af380.htm) | 698.45协议的串口转网口透传通信类(不是TCP通信)，面向对象的用电信息数据交换协议，使用明文的通信方式。支持读取功率，总功，电压，电流，频率，功率因数等数据。  698.45 protocol serial port to network port transparent transmission communication (not TCP communication), object-oriented power consumption information data exchange protocol, using plaintext communication. Support reading power, total power, voltage, current, frequency, power factor and other data. |
| 公共类 | [DLT698Server](1b930240-6344-f3f4-f0e9-563be3abbe76.htm) | DLT698的虚拟服务器实现 |
| 公共类代码示例 | [DLT698TcpNet](5c555bdc-1220-93ff-20f2-5b886e790b46.htm) | 698.45协议的TCP通信类(不是串口透传通信)，面向对象的用电信息数据交换协议，使用明文的通信方式。支持读取功率，总功，电压，电流，频率，功率因数等数据。  The TCP communication class of the 698.45 protocol (not the serial port transparent transmission communication), the object-oriented power consumption information data exchange protocol, uses the clear text communication method. Support reading power, total power, voltage, current, frequency, power factor and other data. |
| 公共类 | [DLTControl](65f784c4-6493-2ab7-1011-7e25717c0a92.htm) | 基本的控制码信息 |
| 公共类 | [DLTTransform](dacea70e-ec57-4c0e-8f93-66c637beab23.htm) | DTL数据转换 |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [DLT698OverTcpDataReportingDelegate](9bdb3a09-bfed-fd7b-d749-5f34d84d47a6.htm) | 数据上报时的方法委托 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLT645 类

[原文連結](http://api.hslcommunication.cn/html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 构造函数](../html/bd811bd3-984b-994a-2067-fe21138812df.htm "DLT645 构造函数 ")

[DLT645 属性](../html/a5c9f3aa-a151-48c4-b5e5-5b631f09f81c.htm "DLT645 属性")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[DLT645 字段](../html/8cfbb9b7-e286-a44b-c68d-163c41652dc6.htm "DLT645 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645 类 |

基于多功能电能表通信协议实现的通讯类，参考的文档是DLT645-2007，主要实现了对电表数据的读取和一些功能方法，
在点对点模式下，需要在打开串口后调用 [ReadAddress](16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm) 方法，数据标识格式为 00-00-00-00，具体参照文档手册。  
The communication type based on the communication protocol of the multifunctional electric energy meter.
The reference document is DLT645-2007, which mainly realizes the reading of the electric meter data and some functional methods.
In the point-to-point mode, you need to call [ReadAddress](16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm) method after opening the serial port.
the data identification format is 00-00-00-00, refer to the documentation manual for details.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)  
        HslCommunication.Instrument.DLTDLT645

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DLT645 : DeviceSerialPort, IDlt645, 
	IReadWriteDevice, IReadWriteNet
```

```
Public Class DLT645
	Inherits DeviceSerialPort
	Implements IDlt645, IReadWriteDevice, IReadWriteNet
```

```
public ref class DLT645 : public DeviceSerialPort, 
	IDlt645, IReadWriteDevice, IReadWriteNet
```

```
type DLT645 =  
    class
        inherit DeviceSerialPort
        interface IDlt645
        interface IReadWriteDevice
        interface IReadWriteNet
    end
```

DLT645 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DLT645](6b823654-6dd1-b96e-f5f2-b877c02e62bf.htm) | 默认的无参构造函数   Default no-parameter constructor |
| 公共方法 | [DLT645(String, String, String)](45e443e7-6233-ad71-7858-d4fe8716ef8d.htm) | 指定地址域，密码，操作者代码来实例化一个对象，密码及操作者代码在写入操作的时候进行验证  Specify the address field, password, and operator code to instantiate an object, and the password and operator code are validated during write operations, which address field is a 12-character BCD code, for example: 149100007290 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](5554f563-5e29-699c-435b-c8e7ad5fa0e4.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CheckDataId](ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm) | 获取或设置当前的DLT645是否检查数据标识，默认为true  Get or set whether the current DLT645 checks data identifiers, with the default being true |
| 公共属性 | [CommunicationPipe](a4fcb477-64d2-8da7-c712-687f0261755b.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [DLTType](ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm) | 获取当前的DLT645的类型信息  Gets the type information of the current DLT645 |
| 公共属性 | [EnableCodeFE](47d2a4e7-008e-7ae8-5114-539626522168.htm) | 获取或设置是否在每一次的报文通信时，增加"FE FE FE FE"的命令头  Get or set whether to add the command header of "FE FE FE FE" in each message communication |
| 公共属性 | [IsClearCacheBeforeRead](81dc971f-1f86-8535-ae9a-586b6590a584.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [OpCode](d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm) | 获取或设置当前DLT645的操作者代码，当进行写入数据操作的时候，需要指定正确的值  Obtain or set the operator code of the current DLT645, and specify the correct value when writing data |
| 公共属性 | [Password](562b4556-472d-98d6-eec3-f7057da65e08.htm) | 获取或设置当前DLT645的密码，当进行写入数据操作的时候，需要正确的密码才能写入  Obtain or set the password of the current DLT645, and the correct password is required to write data operations |
| 公共属性 | [PortName](f2b8e4e3-7079-5240-4ff1-14c6c9b376e6.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](f800fd56-7d49-87c4-985b-103a38330078.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](dd4057ef-88fe-f7d0-d2c4-f6e15488fe78.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Station](b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm) | 获取或设置当前的地址域信息，是一个12个字符的BCD码，例如：149100007290  Get or set the current address domain information, which is a 12-character BCD code, for example: 149100007290 |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ActiveDeveice](b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm) | 激活设备的命令，只发送数据到设备，不等待设备数据返回  The command to activate the device, only send data to the device, do not wait for the device data to return |
| 公共方法 | [BroadcastTime](1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm) | 广播指定的时间，强制从站与主站时间同步，传入DateTime时间对象，没有数据返回。  Broadcast the specified time, force the slave station to synchronize with the master station time, pass in the DateTime time object, and no data will be returned. |
| 公共方法 | [ChangeBaudRate](5d5ebb18-3450-516c-71ad-e122c04c7541.htm) | 更改通信速率，波特率可选 600,1200,2400,4800,9600,19200，其他值无效，可以携带地址域信息，s=1;9600   Change the communication rate, the baud rate can be 600, 1200, 2400, 4800, 9600, 19200, other values are invalid, you can carry address domain information, s=1;9600 |
| 公共方法 | [Close](d3a2ebfb-4e6b-eeda-8bc8-c0972e21b5ff.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [FreezeCommand](c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm) | 对设备发送冻结命令，默认点对点操作，地址域为 99999999999999 时为广播，数据域格式说明：MMDDhhmm(月日时分)， 99DDhhmm表示月为周期定时冻结，9999hhmm表示日为周期定时冻结，999999mm表示以小时为周期定时冻结，99999999表示瞬时冻结  Send a freeze command to the device, the default point-to-point operation, when the address field is 9999999999999, it is broadcast, and the data field format description: MMDDhhmm (month, day, hour and minute), 99DDhhmm means the month is the periodic fixed freeze, 9999hhmm means the day is the periodic periodic freeze, and 999999mm means the hour It is periodic timed freezing, 99999999 means instantaneous freezing |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](0a775ede-e8a0-d411-d0a3-885147c91128.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
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
| 公共方法 | [PackCommandWithHeader](3de72c9d-d491-5f63-9301-3a87bdff5d03.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (重写 [BinaryCommunicationPackCommandWithHeader(Byte)](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm).) |
| 公共方法 | [Read(String, UInt16)](8f3c5175-3e33-bbc6-01c4-17341746eb78.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAddress](16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：149100007290  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](ed7ea8c5-39da-19b6-6c3f-6482560d445e.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm).) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](9b02db80-d5a4-133f-b384-32ab306370c3.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm).) |
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
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](909dd3db-f091-c133-d10c-b83dd484e56a.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |
| 公共方法 | [ReadStringArray](52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，例如 01-01-00-00 正向有功总需求及发生时间  Read all the string data information of the specified address, in general, there is only one data for one address, but there are multiple data for a small number of addresses, such as 01-01-00-00 Forward active total demand and occurrence time |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](8163a33b-8f07-a792-0382-cdf3651e8849.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |
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
| 公共方法 | [SwitchingOn(DateTime)](72dae241-09c1-00e1-e499-fa5fb6db9d68.htm) | 合闸允许功能，需要指定有效截止时间，如果有需要可以指定其他的站号信息 |
| 公共方法 | [SwitchingOn(String, DateTime)](550fee4c-05ff-cc45-16cc-537b49f0c807.htm) | 合闸允许功能，需要指定有效截止时间，如果有需要可以指定其他的站号信息 |
| 公共方法 | [ToString](ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm) | (重写 [DeviceSerialPortToString](774e7804-97ce-fb07-3764-39fb1ca2e1e7.htm).) |
| 公共方法 | [Trip(DateTime)](807bec7f-0910-f422-fcb3-a9973502893f.htm) | 跳闸功能，需要指定有效截止时间，如果有需要可以指定其他的站号信息 |
| 公共方法 | [Trip(String, DateTime)](2ab2cb9e-bfdd-7948-c29e-be87aabb92d2.htm) | 跳闸功能，需要指定有效截止时间，如果有需要可以指定其他的站号信息 |
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
| 公共方法 | [Write(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](5b910f9c-6e6b-bb9d-6842-1b496640c1f4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](b1131bdd-8b51-ec6d-dd0e-fec1f6cb51b6.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Byte)](594ab255-5dfb-7c87-59ae-7f89e498fb2a.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法 | [Write(String, Double)](3095607b-5add-8e91-b080-ccf2d39d428e.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, Double)](df56e6b1-0688-ece7-3bf1-9db32fec9f6e.htm).) |
| 公共方法 | [Write(String, Int16)](d25a4fe6-7fea-2548-8b87-4f3de89b16a2.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm).) |
| 公共方法 | [Write(String, Int32)](84bb15bf-6cea-5597-a557-617c69c69310.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, Int32)](57dac21c-8b28-de7c-c213-77cda465cccb.htm).) |
| 公共方法 | [Write(String, Single)](6da71e57-fc60-45e4-d834-86ef607b4eb9.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, Single)](bb76227e-3e95-8ca2-5dcf-c087d90b8a01.htm).) |
| 公共方法 | [Write(String, UInt16)](db38a073-1e97-6201-adbd-bb35e214e5bc.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm).) |
| 公共方法 | [Write(String, UInt32)](c9c8eb40-52d9-e00b-ebba-c836d973cc60.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, UInt32)](b80b3099-b34d-fde8-6a65-837da3bdc5ab.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](f33e1c6d-9bd0-aa3d-5598-9aa3a2e45e93.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (重写 [DeviceCommunicationWrite(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAddress](5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：149100007290  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 149100007290 |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4b2ab94c-2891-e0a9-f10b-b67c63bd1cf9.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](ae3df8ec-29c0-42ff-3d6a-68f496997b3d.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Double)](c65a321f-c229-7151-c3af-db2d3ee317bb.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, Double)](36199b83-c4ea-b02f-3675-fb6ac48d0c59.htm).) |
| 公共方法 | [WriteAsync(String, Int16)](d172f655-4715-f811-6624-b9f9670299e3.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm).) |
| 公共方法 | [WriteAsync(String, Int32)](e711549f-00cc-db4c-9a8a-239caa4ead10.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, Int32)](99962a08-9624-cea6-6009-5e6f207236b7.htm).) |
| 公共方法 | [WriteAsync(String, Single)](cbc953db-ae6b-ffd8-6ca2-1ff4ddc2f1c0.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, Single)](dbb2359a-fb54-cf93-4413-20d5a4ef9d36.htm).) |
| 公共方法 | [WriteAsync(String, UInt16)](067c6bdd-8b35-a9c9-87cc-70474d3a3587.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm).) |
| 公共方法 | [WriteAsync(String, UInt32)](932da9d0-27c9-e53c-962c-9465531e5106.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, UInt32)](5b655bf3-5ba7-ad31-6897-89c59051ac60.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](86b33f31-c0ab-88e2-46f9-c0bcd7d65f6a.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (重写 [DeviceCommunicationWriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm).) |
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

如果一对多的模式，地址可以携带地址域访问，例如 "s=2;00-00-00-00"，主要使用 [ReadDouble(String, UInt16)](ed7ea8c5-39da-19b6-6c3f-6482560d445e.htm) 方法来读取浮点数，
[ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) 方法来读取字符串

![](../icons/SectionExpanded.png)示例

具体的地址请参考相关的手册内容，如果没有，可以联系HSL作者或者，下面列举一些常用的地址  
对于电能来说，DI0是结算日的信息，现在的就是写0，上一结算日的就写 01，上12结算日就写 0C

| DI3 | DI2 | DI1 | DI0 | 地址示例 | 读取方式 | 数据项名称 | 备注 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 00 | 00 | 00 | 00 | 00-00-00-00 | ReadDouble | （当前）组合有功总电能(kwh) | 00-00-01-00到00-00-3F-00分别是组合有功费率1~63电能 |
| 00 | 01 | 00 | 00 | 00-01-00-00 | ReadDouble | （当前）正向有功总电能(kwh) | 00-01-01-00到00-01-3F-00分别是正向有功费率1~63电能 |
| 00 | 02 | 00 | 00 | 00-02-00-00 | ReadDouble | （当前）反向有功总电能(kwh) | 00-02-01-00到00-02-3F-00分别是反向有功费率1~63电能 |
| 00 | 03 | 00 | 00 | 00-03-00-00 | ReadDouble | （当前）组合无功总电能(kvarh) | 00-03-01-00到00-03-3F-00分别是组合无功费率1~63电能 |
| 00 | 09 | 00 | 00 | 00-09-00-00 | ReadDouble | （当前）正向视在总电能(kvah) | 00-09-01-00到00-09-3F-00分别是正向视在费率1~63电能 |
| 00 | 0A | 00 | 00 | 00-0A-00-00 | ReadDouble | （当前）反向视在总电能(kvah) | 00-0A-01-00到00-0A-3F-00分别是反向视在费率1~63电能 |
| 02 | 01 | 01 | 00 | 02-01-01-00 | ReadDouble | A相电压(V) |  |
| 02 | 01 | 02 | 00 | 02-01-02-00 | ReadDouble | B相电压(V) |  |
| 02 | 01 | 03 | 00 | 02-01-03-00 | ReadDouble | C相电压(V) |  |
| 02 | 02 | 01 | 00 | 02-02-01-00 | ReadDouble | A相电流(A) |  |
| 02 | 02 | 02 | 00 | 02-02-02-00 | ReadDouble | B相电流(A) |  |
| 02 | 02 | 03 | 00 | 02-02-03-00 | ReadDouble | C相电流(A) |  |
| 02 | 03 | 00 | 00 | 02-03-00-00 | ReadDouble | 瞬时总有功功率(kw) | DI1=1时表示A相，2时表示B相，3时表示C相 |
| 02 | 04 | 00 | 00 | 02-04-00-00 | ReadDouble | 瞬时总无功功率(kvar) | DI1=1时表示A相，2时表示B相，3时表示C相 |
| 02 | 05 | 00 | 00 | 02-05-00-00 | ReadDouble | 瞬时总视在功率(kva) | DI1=1时表示A相，2时表示B相，3时表示C相 |
| 02 | 06 | 00 | 00 | 02-06-00-00 | ReadDouble | 总功率因素 | DI1=1时表示A相，2时表示B相，3时表示C相 |
| 02 | 07 | 01 | 00 | 02-07-01-00 | ReadDouble | A相相角(°) | DI1=1时表示A相，2时表示B相，3时表示C相 |
| 02 | 08 | 01 | 00 | 02-08-01-00 | ReadDouble | A相电压波形失真度(%) | DI1=1时表示A相，2时表示B相，3时表示C相 |
| 02 | 80 | 00 | 01 | 02-80-00-01 | ReadDouble | 零线电流(A) |  |
| 02 | 80 | 00 | 02 | 02-80-00-02 | ReadDouble | 电网频率(HZ) |  |
| 02 | 80 | 00 | 03 | 02-80-00-03 | ReadDouble | 一分钟有功总平均功率(kw) |  |
| 02 | 80 | 00 | 04 | 02-80-00-04 | ReadDouble | 当前有功需量(kw) |  |
| 02 | 80 | 00 | 05 | 02-80-00-05 | ReadDouble | 当前无功需量(kvar) |  |
| 02 | 80 | 00 | 06 | 02-80-00-06 | ReadDouble | 当前视在需量(kva) |  |
| 02 | 80 | 00 | 07 | 02-80-00-07 | ReadDouble | 表内温度(摄氏度) |  |
| 02 | 80 | 00 | 08 | 02-80-00-08 | ReadDouble | 时钟电池电压(V) |  |
| 02 | 80 | 00 | 09 | 02-80-00-09 | ReadDouble | 停电抄表电池电压(V) |  |
| 02 | 80 | 00 | 0A | 02-80-00-0A | ReadDouble | 内部电池工作时间(分钟) |  |
| 04 | 00 | 04 | 03 | 04-00-04-03 | ReadString("04-00-04-03", 32) | 资产管理编码 |  |
| 04 | 00 | 04 | 0B | 04-00-04-0B | ReadString("04-00-04-0B", 10) | 电表型号 |  |
| 04 | 00 | 04 | 0C | 04-00-04-0C | ReadString("04-00-04-0C", 10) | 生产日期 |  |

直接串口初始化，打开串口，就可以对数据进行读取了，地址如上图所示。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLT645 构造函数 

[原文連結](http://api.hslcommunication.cn/html/bd811bd3-984b-994a-2067-fe21138812df.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 构造函数](../html/bd811bd3-984b-994a-2067-fe21138812df.htm "DLT645 构造函数 ")

[DLT645 构造函数](../html/6b823654-6dd1-b96e-f5f2-b877c02e62bf.htm "DLT645 构造函数 ")

[DLT645 构造函数 (String, String, String)](../html/45e443e7-6233-ad71-7858-d4fe8716ef8d.htm "DLT645 构造函数 (String, String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DLT645](6b823654-6dd1-b96e-f5f2-b877c02e62bf.htm) | 默认的无参构造函数   Default no-parameter constructor |
| 公共方法 | [DLT645(String, String, String)](45e443e7-6233-ad71-7858-d4fe8716ef8d.htm) | 指定地址域，密码，操作者代码来实例化一个对象，密码及操作者代码在写入操作的时候进行验证  Specify the address field, password, and operator code to instantiate an object, and the password and operator code are validated during write operations, which address field is a 12-character BCD code, for example: 149100007290 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLT645 构造函数 

[原文連結](http://api.hslcommunication.cn/html/6b823654-6dd1-b96e-f5f2-b877c02e62bf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 构造函数](../html/bd811bd3-984b-994a-2067-fe21138812df.htm "DLT645 构造函数 ")

[DLT645 构造函数](../html/6b823654-6dd1-b96e-f5f2-b877c02e62bf.htm "DLT645 构造函数 ")

[DLT645 构造函数 (String, String, String)](../html/45e443e7-6233-ad71-7858-d4fe8716ef8d.htm "DLT645 构造函数 (String, String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645 构造函数 |

默认的无参构造函数   
Default no-parameter constructor

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DLT645()
```

```
Public Sub New
```

```
public:
DLT645()
```

```
new : unit -> DLT645
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[DLT645 重载](bd811bd3-984b-994a-2067-fe21138812df.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLT645 构造函数 (String, String, String)

[原文連結](http://api.hslcommunication.cn/html/45e443e7-6233-ad71-7858-d4fe8716ef8d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 构造函数](../html/bd811bd3-984b-994a-2067-fe21138812df.htm "DLT645 构造函数 ")

[DLT645 构造函数](../html/6b823654-6dd1-b96e-f5f2-b877c02e62bf.htm "DLT645 构造函数 ")

[DLT645 构造函数 (String, String, String)](../html/45e443e7-6233-ad71-7858-d4fe8716ef8d.htm "DLT645 构造函数 (String, String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645 构造函数 (String, String, String) |

指定地址域，密码，操作者代码来实例化一个对象，密码及操作者代码在写入操作的时候进行验证  
Specify the address field, password, and operator code to instantiate an object, and the password and operator code are validated during write operations,
which address field is a 12-character BCD code, for example: 149100007290

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DLT645(
	string station,
	string password = "",
	string opCode = ""
)
```

```
Public Sub New ( 
	station As String,
	Optional password As String = "",
	Optional opCode As String = ""
)
```

```
public:
DLT645(
	String^ station, 
	String^ password = L"", 
	String^ opCode = L""
)
```

```
new : 
        station : string * 
        ?password : string * 
        ?opCode : string 
(* Defaults:
        let _password = defaultArg password ""
        let _opCode = defaultArg opCode ""
*)
-> DLT645
```

#### 参数

station
:   类型：SystemString  
    设备的地址信息，是一个12字符的BCD码

password (Optional)
:   类型：SystemString  
    密码，写入的时候进行验证的信息

opCode (Optional)
:   类型：SystemString  
    操作者代码

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[DLT645 重载](bd811bd3-984b-994a-2067-fe21138812df.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLT645 属性

[原文連結](http://api.hslcommunication.cn/html/a5c9f3aa-a151-48c4-b5e5-5b631f09f81c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 属性](../html/a5c9f3aa-a151-48c4-b5e5-5b631f09f81c.htm "DLT645 属性")

[CheckDataId 属性](../html/ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm "CheckDataId 属性 ")

[DLTType 属性](../html/ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm "DLTType 属性 ")

[EnableCodeFE 属性](../html/47d2a4e7-008e-7ae8-5114-539626522168.htm "EnableCodeFE 属性 ")

[OpCode 属性](../html/d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm "OpCode 属性 ")

[Password 属性](../html/562b4556-472d-98d6-eec3-f7057da65e08.htm "Password 属性 ")

[Station 属性](../html/b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645 属性 |

[DLT645](dc3f4d48-9b60-0203-88c1-3543260dee17.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](5554f563-5e29-699c-435b-c8e7ad5fa0e4.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CheckDataId](ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm) | 获取或设置当前的DLT645是否检查数据标识，默认为true  Get or set whether the current DLT645 checks data identifiers, with the default being true |
| 公共属性 | [CommunicationPipe](a4fcb477-64d2-8da7-c712-687f0261755b.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [DLTType](ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm) | 获取当前的DLT645的类型信息  Gets the type information of the current DLT645 |
| 公共属性 | [EnableCodeFE](47d2a4e7-008e-7ae8-5114-539626522168.htm) | 获取或设置是否在每一次的报文通信时，增加"FE FE FE FE"的命令头  Get or set whether to add the command header of "FE FE FE FE" in each message communication |
| 公共属性 | [IsClearCacheBeforeRead](81dc971f-1f86-8535-ae9a-586b6590a584.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [OpCode](d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm) | 获取或设置当前DLT645的操作者代码，当进行写入数据操作的时候，需要指定正确的值  Obtain or set the operator code of the current DLT645, and specify the correct value when writing data |
| 公共属性 | [Password](562b4556-472d-98d6-eec3-f7057da65e08.htm) | 获取或设置当前DLT645的密码，当进行写入数据操作的时候，需要正确的密码才能写入  Obtain or set the password of the current DLT645, and the correct password is required to write data operations |
| 公共属性 | [PortName](f2b8e4e3-7079-5240-4ff1-14c6c9b376e6.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](f800fd56-7d49-87c4-985b-103a38330078.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](dd4057ef-88fe-f7d0-d2c4-f6e15488fe78.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Station](b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm) | 获取或设置当前的地址域信息，是一个12个字符的BCD码，例如：149100007290  Get or set the current address domain information, which is a 12-character BCD code, for example: 149100007290 |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckDataId 属性 

[原文連結](http://api.hslcommunication.cn/html/ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 属性](../html/a5c9f3aa-a151-48c4-b5e5-5b631f09f81c.htm "DLT645 属性")

[CheckDataId 属性](../html/ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm "CheckDataId 属性 ")

[DLTType 属性](../html/ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm "DLTType 属性 ")

[EnableCodeFE 属性](../html/47d2a4e7-008e-7ae8-5114-539626522168.htm "EnableCodeFE 属性 ")

[OpCode 属性](../html/d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm "OpCode 属性 ")

[Password 属性](../html/562b4556-472d-98d6-eec3-f7057da65e08.htm "Password 属性 ")

[Station 属性](../html/b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645CheckDataId 属性 |

获取或设置当前的DLT645是否检查数据标识，默认为true  
Get or set whether the current DLT645 checks data identifiers, with the default being true

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool CheckDataId { get; set; }
```

```
Public Property CheckDataId As Boolean
	Get
	Set
```

```
public:
virtual property bool CheckDataId {
	bool get () sealed;
	void set (bool value) sealed;
}
```

```
abstract CheckDataId : bool with get, set
override CheckDataId : bool with get, set
```

#### 属性值

类型：Boolean

#### 实现

[IDlt645CheckDataId](f1c44a47-4854-e681-41c0-afd3dfc72b5f.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLTType 属性 

[原文連結](http://api.hslcommunication.cn/html/ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 属性](../html/a5c9f3aa-a151-48c4-b5e5-5b631f09f81c.htm "DLT645 属性")

[CheckDataId 属性](../html/ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm "CheckDataId 属性 ")

[DLTType 属性](../html/ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm "DLTType 属性 ")

[EnableCodeFE 属性](../html/47d2a4e7-008e-7ae8-5114-539626522168.htm "EnableCodeFE 属性 ")

[OpCode 属性](../html/d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm "OpCode 属性 ")

[Password 属性](../html/562b4556-472d-98d6-eec3-f7057da65e08.htm "Password 属性 ")

[Station 属性](../html/b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645DLTType 属性 |

获取当前的DLT645的类型信息  
Gets the type information of the current DLT645

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DLT645Type DLTType { get; }
```

```
Public ReadOnly Property DLTType As DLT645Type
	Get
```

```
public:
virtual property DLT645Type DLTType {
	DLT645Type get () sealed;
}
```

```
abstract DLTType : DLT645Type with get
override DLTType : DLT645Type with get
```

#### 属性值

类型：[DLT645Type](3adc8897-273c-cb21-b98b-4af641b94e56.htm)

#### 实现

[IDlt645DLTType](a5f8ffe8-e138-cf85-eab8-47a564a04aa7.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EnableCodeFE 属性 

[原文連結](http://api.hslcommunication.cn/html/47d2a4e7-008e-7ae8-5114-539626522168.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 属性](../html/a5c9f3aa-a151-48c4-b5e5-5b631f09f81c.htm "DLT645 属性")

[CheckDataId 属性](../html/ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm "CheckDataId 属性 ")

[DLTType 属性](../html/ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm "DLTType 属性 ")

[EnableCodeFE 属性](../html/47d2a4e7-008e-7ae8-5114-539626522168.htm "EnableCodeFE 属性 ")

[OpCode 属性](../html/d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm "OpCode 属性 ")

[Password 属性](../html/562b4556-472d-98d6-eec3-f7057da65e08.htm "Password 属性 ")

[Station 属性](../html/b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645EnableCodeFE 属性 |

获取或设置是否在每一次的报文通信时，增加"FE FE FE FE"的命令头  
Get or set whether to add the command header of "FE FE FE FE" in each message communication

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool EnableCodeFE { get; set; }
```

```
Public Property EnableCodeFE As Boolean
	Get
	Set
```

```
public:
virtual property bool EnableCodeFE {
	bool get () sealed;
	void set (bool value) sealed;
}
```

```
abstract EnableCodeFE : bool with get, set
override EnableCodeFE : bool with get, set
```

#### 属性值

类型：Boolean

#### 实现

[IDlt645EnableCodeFE](cec146a5-4952-c424-0a6c-50c2ede66f96.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OpCode 属性 

[原文連結](http://api.hslcommunication.cn/html/d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 属性](../html/a5c9f3aa-a151-48c4-b5e5-5b631f09f81c.htm "DLT645 属性")

[CheckDataId 属性](../html/ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm "CheckDataId 属性 ")

[DLTType 属性](../html/ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm "DLTType 属性 ")

[EnableCodeFE 属性](../html/47d2a4e7-008e-7ae8-5114-539626522168.htm "EnableCodeFE 属性 ")

[OpCode 属性](../html/d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm "OpCode 属性 ")

[Password 属性](../html/562b4556-472d-98d6-eec3-f7057da65e08.htm "Password 属性 ")

[Station 属性](../html/b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645OpCode 属性 |

获取或设置当前DLT645的操作者代码，当进行写入数据操作的时候，需要指定正确的值  
Obtain or set the operator code of the current DLT645, and specify the correct value when writing data

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string OpCode { get; set; }
```

```
Public Property OpCode As String
	Get
	Set
```

```
public:
virtual property String^ OpCode {
	String^ get () sealed;
	void set (String^ value) sealed;
}
```

```
abstract OpCode : string with get, set
override OpCode : string with get, set
```

#### 属性值

类型：String

#### 实现

[IDlt645OpCode](c2854a1b-94d8-dbb6-48cf-69dcef13e6ec.htm)

![](../icons/SectionExpanded.png)备注

对于 DLT645/1997 协议来说无效

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Password 属性 

[原文連結](http://api.hslcommunication.cn/html/562b4556-472d-98d6-eec3-f7057da65e08.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 属性](../html/a5c9f3aa-a151-48c4-b5e5-5b631f09f81c.htm "DLT645 属性")

[CheckDataId 属性](../html/ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm "CheckDataId 属性 ")

[DLTType 属性](../html/ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm "DLTType 属性 ")

[EnableCodeFE 属性](../html/47d2a4e7-008e-7ae8-5114-539626522168.htm "EnableCodeFE 属性 ")

[OpCode 属性](../html/d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm "OpCode 属性 ")

[Password 属性](../html/562b4556-472d-98d6-eec3-f7057da65e08.htm "Password 属性 ")

[Station 属性](../html/b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645Password 属性 |

获取或设置当前DLT645的密码，当进行写入数据操作的时候，需要正确的密码才能写入  
Obtain or set the password of the current DLT645, and the correct password is required to write data operations

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Password { get; set; }
```

```
Public Property Password As String
	Get
	Set
```

```
public:
virtual property String^ Password {
	String^ get () sealed;
	void set (String^ value) sealed;
}
```

```
abstract Password : string with get, set
override Password : string with get, set
```

#### 属性值

类型：String

#### 实现

[IDlt645Password](4078cd31-5dc8-aecb-b6c8-fac10fc68cd3.htm)

![](../icons/SectionExpanded.png)备注

对于 DLT645/1997 协议来说无效

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Station 属性 

[原文連結](http://api.hslcommunication.cn/html/b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 属性](../html/a5c9f3aa-a151-48c4-b5e5-5b631f09f81c.htm "DLT645 属性")

[CheckDataId 属性](../html/ff8ffb96-f0fd-bcb0-2708-ea8a95c61b37.htm "CheckDataId 属性 ")

[DLTType 属性](../html/ec1c2c43-d454-ec21-d69f-2b72a50d087e.htm "DLTType 属性 ")

[EnableCodeFE 属性](../html/47d2a4e7-008e-7ae8-5114-539626522168.htm "EnableCodeFE 属性 ")

[OpCode 属性](../html/d10fd60b-5032-6b8e-d7d8-7778f9ff28d8.htm "OpCode 属性 ")

[Password 属性](../html/562b4556-472d-98d6-eec3-f7057da65e08.htm "Password 属性 ")

[Station 属性](../html/b7e9bf03-c6a5-b48d-56ea-df9ed8b149e7.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645Station 属性 |

获取或设置当前的地址域信息，是一个12个字符的BCD码，例如：149100007290  
Get or set the current address domain information, which is a 12-character BCD code, for example: 149100007290

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Station { get; set; }
```

```
Public Property Station As String
	Get
	Set
```

```
public:
virtual property String^ Station {
	String^ get () sealed;
	void set (String^ value) sealed;
}
```

```
abstract Station : string with get, set
override Station : string with get, set
```

#### 属性值

类型：String

#### 实现

[IDlt645Station](98627128-d199-1937-5194-c76a330ebb7f.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DLT645 方法

[原文連結](http://api.hslcommunication.cn/html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ActiveDeveice 方法](../html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm "ActiveDeveice 方法 ")

[BroadcastTime 方法](../html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm "BroadcastTime 方法 ")

[ChangeBaudRate 方法](../html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm "ChangeBaudRate 方法 ")

[FreezeCommand 方法](../html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm "FreezeCommand 方法 ")

[GetNewNetMessage 方法](../html/0a775ede-e8a0-d411-d0a3-885147c91128.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[ReadAddress 方法](../html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[SwitchingOn 方法](../html/279d6777-e35a-c17c-40f6-c5fdffe9ae70.htm "SwitchingOn 方法 ")

[ToString 方法](../html/ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm "ToString 方法 ")

[Trip 方法](../html/eee862c9-b6ae-307c-860f-1daf0dec3fd3.htm "Trip 方法 ")

[Write 方法](../html/5e9d2d45-4f08-9074-0428-626196c5c83b.htm "Write 方法 ")

[WriteAddress 方法](../html/5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm "WriteAddress 方法 ")

[WriteAsync 方法](../html/bb712c3c-a331-ed3e-43af-1b03a975fe1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645 方法 |

[DLT645](dc3f4d48-9b60-0203-88c1-3543260dee17.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ActiveDeveice](b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm) | 激活设备的命令，只发送数据到设备，不等待设备数据返回  The command to activate the device, only send data to the device, do not wait for the device data to return |
| 公共方法 | [BroadcastTime](1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm) | 广播指定的时间，强制从站与主站时间同步，传入DateTime时间对象，没有数据返回。  Broadcast the specified time, force the slave station to synchronize with the master station time, pass in the DateTime time object, and no data will be returned. |
| 公共方法 | [ChangeBaudRate](5d5ebb18-3450-516c-71ad-e122c04c7541.htm) | 更改通信速率，波特率可选 600,1200,2400,4800,9600,19200，其他值无效，可以携带地址域信息，s=1;9600   Change the communication rate, the baud rate can be 600, 1200, 2400, 4800, 9600, 19200, other values are invalid, you can carry address domain information, s=1;9600 |
| 公共方法 | [Close](d3a2ebfb-4e6b-eeda-8bc8-c0972e21b5ff.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [FreezeCommand](c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm) | 对设备发送冻结命令，默认点对点操作，地址域为 99999999999999 时为广播，数据域格式说明：MMDDhhmm(月日时分)， 99DDhhmm表示月为周期定时冻结，9999hhmm表示日为周期定时冻结，999999mm表示以小时为周期定时冻结，99999999表示瞬时冻结  Send a freeze command to the device, the default point-to-point operation, when the address field is 9999999999999, it is broadcast, and the data field format description: MMDDhhmm (month, day, hour and minute), 99DDhhmm means the month is the periodic fixed freeze, 9999hhmm means the day is the periodic periodic freeze, and 999999mm means the hour It is periodic timed freezing, 99999999 means instantaneous freezing |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](0a775ede-e8a0-d411-d0a3-885147c91128.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
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
| 公共方法 | [PackCommandWithHeader](3de72c9d-d491-5f63-9301-3a87bdff5d03.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (重写 [BinaryCommunicationPackCommandWithHeader(Byte)](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm).) |
| 公共方法 | [Read(String, UInt16)](8f3c5175-3e33-bbc6-01c4-17341746eb78.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAddress](16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：149100007290  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](ed7ea8c5-39da-19b6-6c3f-6482560d445e.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm).) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](9b02db80-d5a4-133f-b384-32ab306370c3.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm).) |
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
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](909dd3db-f091-c133-d10c-b83dd484e56a.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |
| 公共方法 | [ReadStringArray](52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，例如 01-01-00-00 正向有功总需求及发生时间  Read all the string data information of the specified address, in general, there is only one data for one address, but there are multiple data for a small number of addresses, such as 01-01-00-00 Forward active total demand and occurrence time |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](8163a33b-8f07-a792-0382-cdf3651e8849.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |
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
| 公共方法 | [SwitchingOn(DateTime)](72dae241-09c1-00e1-e499-fa5fb6db9d68.htm) | 合闸允许功能，需要指定有效截止时间，如果有需要可以指定其他的站号信息 |
| 公共方法 | [SwitchingOn(String, DateTime)](550fee4c-05ff-cc45-16cc-537b49f0c807.htm) | 合闸允许功能，需要指定有效截止时间，如果有需要可以指定其他的站号信息 |
| 公共方法 | [ToString](ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm) | (重写 [DeviceSerialPortToString](774e7804-97ce-fb07-3764-39fb1ca2e1e7.htm).) |
| 公共方法 | [Trip(DateTime)](807bec7f-0910-f422-fcb3-a9973502893f.htm) | 跳闸功能，需要指定有效截止时间，如果有需要可以指定其他的站号信息 |
| 公共方法 | [Trip(String, DateTime)](2ab2cb9e-bfdd-7948-c29e-be87aabb92d2.htm) | 跳闸功能，需要指定有效截止时间，如果有需要可以指定其他的站号信息 |
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
| 公共方法 | [Write(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](5b910f9c-6e6b-bb9d-6842-1b496640c1f4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](b1131bdd-8b51-ec6d-dd0e-fec1f6cb51b6.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Byte)](594ab255-5dfb-7c87-59ae-7f89e498fb2a.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法 | [Write(String, Double)](3095607b-5add-8e91-b080-ccf2d39d428e.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, Double)](df56e6b1-0688-ece7-3bf1-9db32fec9f6e.htm).) |
| 公共方法 | [Write(String, Int16)](d25a4fe6-7fea-2548-8b87-4f3de89b16a2.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm).) |
| 公共方法 | [Write(String, Int32)](84bb15bf-6cea-5597-a557-617c69c69310.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, Int32)](57dac21c-8b28-de7c-c213-77cda465cccb.htm).) |
| 公共方法 | [Write(String, Single)](6da71e57-fc60-45e4-d834-86ef607b4eb9.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, Single)](bb76227e-3e95-8ca2-5dcf-c087d90b8a01.htm).) |
| 公共方法 | [Write(String, UInt16)](db38a073-1e97-6201-adbd-bb35e214e5bc.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm).) |
| 公共方法 | [Write(String, UInt32)](c9c8eb40-52d9-e00b-ebba-c836d973cc60.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWrite(String, UInt32)](b80b3099-b34d-fde8-6a65-837da3bdc5ab.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](f33e1c6d-9bd0-aa3d-5598-9aa3a2e45e93.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (重写 [DeviceCommunicationWrite(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAddress](5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：149100007290  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 149100007290 |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4b2ab94c-2891-e0a9-f10b-b67c63bd1cf9.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](ae3df8ec-29c0-42ff-3d6a-68f496997b3d.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Double)](c65a321f-c229-7151-c3af-db2d3ee317bb.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, Double)](36199b83-c4ea-b02f-3675-fb6ac48d0c59.htm).) |
| 公共方法 | [WriteAsync(String, Int16)](d172f655-4715-f811-6624-b9f9670299e3.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm).) |
| 公共方法 | [WriteAsync(String, Int32)](e711549f-00cc-db4c-9a8a-239caa4ead10.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, Int32)](99962a08-9624-cea6-6009-5e6f207236b7.htm).) |
| 公共方法 | [WriteAsync(String, Single)](cbc953db-ae6b-ffd8-6ca2-1ff4ddc2f1c0.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, Single)](dbb2359a-fb54-cf93-4413-20d5a4ef9d36.htm).) |
| 公共方法 | [WriteAsync(String, UInt16)](067c6bdd-8b35-a9c9-87cc-70474d3a3587.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm).) |
| 公共方法 | [WriteAsync(String, UInt32)](932da9d0-27c9-e53c-962c-9465531e5106.htm) | 将指定的数据写入到仪表中，，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Write the data to the gauge, address identification according to the manual, from high bit to position, such as 00-00-00-00, the separator can be any special character or no delimiter. (重写 [DeviceCommunicationWriteAsync(String, UInt32)](5b655bf3-5ba7-ad31-6897-89c59051ac60.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](86b33f31-c0ab-88e2-46f9-c0bcd7d65f6a.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (重写 [DeviceCommunicationWriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm).) |
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

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ActiveDeveice 方法 

[原文連結](http://api.hslcommunication.cn/html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ActiveDeveice 方法](../html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm "ActiveDeveice 方法 ")

[BroadcastTime 方法](../html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm "BroadcastTime 方法 ")

[ChangeBaudRate 方法](../html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm "ChangeBaudRate 方法 ")

[FreezeCommand 方法](../html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm "FreezeCommand 方法 ")

[GetNewNetMessage 方法](../html/0a775ede-e8a0-d411-d0a3-885147c91128.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[ReadAddress 方法](../html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[SwitchingOn 方法](../html/279d6777-e35a-c17c-40f6-c5fdffe9ae70.htm "SwitchingOn 方法 ")

[ToString 方法](../html/ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm "ToString 方法 ")

[Trip 方法](../html/eee862c9-b6ae-307c-860f-1daf0dec3fd3.htm "Trip 方法 ")

[Write 方法](../html/5e9d2d45-4f08-9074-0428-626196c5c83b.htm "Write 方法 ")

[WriteAddress 方法](../html/5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm "WriteAddress 方法 ")

[WriteAsync 方法](../html/bb712c3c-a331-ed3e-43af-1b03a975fe1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ActiveDeveice 方法 |

激活设备的命令，只发送数据到设备，不等待设备数据返回  
The command to activate the device, only send data to the device, do not wait for the device data to return

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult ActiveDeveice()
```

```
Public Function ActiveDeveice As OperateResult
```

```
public:
virtual OperateResult^ ActiveDeveice() sealed
```

```
abstract ActiveDeveice : unit -> OperateResult 
override ActiveDeveice : unit -> OperateResult
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

#### 实现

[IDlt645ActiveDeveice](e09d62f7-750b-0874-f18c-b890acc97ecc.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BroadcastTime 方法 

[原文連結](http://api.hslcommunication.cn/html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ActiveDeveice 方法](../html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm "ActiveDeveice 方法 ")

[BroadcastTime 方法](../html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm "BroadcastTime 方法 ")

[ChangeBaudRate 方法](../html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm "ChangeBaudRate 方法 ")

[FreezeCommand 方法](../html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm "FreezeCommand 方法 ")

[GetNewNetMessage 方法](../html/0a775ede-e8a0-d411-d0a3-885147c91128.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[ReadAddress 方法](../html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[SwitchingOn 方法](../html/279d6777-e35a-c17c-40f6-c5fdffe9ae70.htm "SwitchingOn 方法 ")

[ToString 方法](../html/ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm "ToString 方法 ")

[Trip 方法](../html/eee862c9-b6ae-307c-860f-1daf0dec3fd3.htm "Trip 方法 ")

[Write 方法](../html/5e9d2d45-4f08-9074-0428-626196c5c83b.htm "Write 方法 ")

[WriteAddress 方法](../html/5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm "WriteAddress 方法 ")

[WriteAsync 方法](../html/bb712c3c-a331-ed3e-43af-1b03a975fe1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645BroadcastTime 方法 |

广播指定的时间，强制从站与主站时间同步，传入DateTime时间对象，没有数据返回。  
Broadcast the specified time, force the slave station to synchronize with the master station time,
pass in the DateTime time object, and no data will be returned.

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult BroadcastTime(
	DateTime dateTime
)
```

```
Public Function BroadcastTime ( 
	dateTime As DateTime
) As OperateResult
```

```
public:
virtual OperateResult^ BroadcastTime(
	DateTime dateTime
) sealed
```

```
abstract BroadcastTime : 
        dateTime : DateTime -> OperateResult 
override BroadcastTime : 
        dateTime : DateTime -> OperateResult
```

#### 参数

dateTime
:   类型：SystemDateTime  
    时间对象

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功

#### 实现

[IDlt645BroadcastTime(DateTime)](8c9b251f-040e-989f-18ba-02ac7381592f.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ChangeBaudRate 方法 

[原文連結](http://api.hslcommunication.cn/html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ActiveDeveice 方法](../html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm "ActiveDeveice 方法 ")

[BroadcastTime 方法](../html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm "BroadcastTime 方法 ")

[ChangeBaudRate 方法](../html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm "ChangeBaudRate 方法 ")

[FreezeCommand 方法](../html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm "FreezeCommand 方法 ")

[GetNewNetMessage 方法](../html/0a775ede-e8a0-d411-d0a3-885147c91128.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[ReadAddress 方法](../html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[SwitchingOn 方法](../html/279d6777-e35a-c17c-40f6-c5fdffe9ae70.htm "SwitchingOn 方法 ")

[ToString 方法](../html/ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm "ToString 方法 ")

[Trip 方法](../html/eee862c9-b6ae-307c-860f-1daf0dec3fd3.htm "Trip 方法 ")

[Write 方法](../html/5e9d2d45-4f08-9074-0428-626196c5c83b.htm "Write 方法 ")

[WriteAddress 方法](../html/5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm "WriteAddress 方法 ")

[WriteAsync 方法](../html/bb712c3c-a331-ed3e-43af-1b03a975fe1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ChangeBaudRate 方法 |

更改通信速率，波特率可选 600,1200,2400,4800,9600,19200，其他值无效，可以携带地址域信息，s=1;9600   
Change the communication rate, the baud rate can be 600, 1200, 2400, 4800, 9600, 19200,
other values are invalid, you can carry address domain information, s=1;9600

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult ChangeBaudRate(
	string baudRate
)
```

```
Public Function ChangeBaudRate ( 
	baudRate As String
) As OperateResult
```

```
public:
virtual OperateResult^ ChangeBaudRate(
	String^ baudRate
) sealed
```

```
abstract ChangeBaudRate : 
        baudRate : string -> OperateResult 
override ChangeBaudRate : 
        baudRate : string -> OperateResult
```

#### 参数

baudRate
:   类型：SystemString  
    波特率的信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否更改成功

#### 实现

[IDlt645ChangeBaudRate(String)](bea572b1-5231-2302-faae-d174508f62d8.htm)

![](../icons/SectionExpanded.png)备注

对于DLT1997来说，只支持 300, 600, 2400, 4800, 9600

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FreezeCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ActiveDeveice 方法](../html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm "ActiveDeveice 方法 ")

[BroadcastTime 方法](../html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm "BroadcastTime 方法 ")

[ChangeBaudRate 方法](../html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm "ChangeBaudRate 方法 ")

[FreezeCommand 方法](../html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm "FreezeCommand 方法 ")

[GetNewNetMessage 方法](../html/0a775ede-e8a0-d411-d0a3-885147c91128.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[ReadAddress 方法](../html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[SwitchingOn 方法](../html/279d6777-e35a-c17c-40f6-c5fdffe9ae70.htm "SwitchingOn 方法 ")

[ToString 方法](../html/ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm "ToString 方法 ")

[Trip 方法](../html/eee862c9-b6ae-307c-860f-1daf0dec3fd3.htm "Trip 方法 ")

[Write 方法](../html/5e9d2d45-4f08-9074-0428-626196c5c83b.htm "Write 方法 ")

[WriteAddress 方法](../html/5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm "WriteAddress 方法 ")

[WriteAsync 方法](../html/bb712c3c-a331-ed3e-43af-1b03a975fe1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645FreezeCommand 方法 |

对设备发送冻结命令，默认点对点操作，地址域为 99999999999999 时为广播，数据域格式说明：MMDDhhmm(月日时分)，
99DDhhmm表示月为周期定时冻结，9999hhmm表示日为周期定时冻结，999999mm表示以小时为周期定时冻结，99999999表示瞬时冻结  
Send a freeze command to the device, the default point-to-point operation, when the address field is 9999999999999,
it is broadcast, and the data field format description: MMDDhhmm (month, day, hour and minute),
99DDhhmm means the month is the periodic fixed freeze, 9999hhmm means the day is the periodic periodic freeze,
and 999999mm means the hour It is periodic timed freezing, 99999999 means instantaneous freezing

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult FreezeCommand(
	string dataArea
)
```

```
Public Function FreezeCommand ( 
	dataArea As String
) As OperateResult
```

```
public:
OperateResult^ FreezeCommand(
	String^ dataArea
)
```

```
member FreezeCommand : 
        dataArea : string -> OperateResult 
```

#### 参数

dataArea
:   类型：SystemString  
    数据域信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功冻结

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/0a775ede-e8a0-d411-d0a3-885147c91128.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ActiveDeveice 方法](../html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm "ActiveDeveice 方法 ")

[BroadcastTime 方法](../html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm "BroadcastTime 方法 ")

[ChangeBaudRate 方法](../html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm "ChangeBaudRate 方法 ")

[FreezeCommand 方法](../html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm "FreezeCommand 方法 ")

[GetNewNetMessage 方法](../html/0a775ede-e8a0-d411-d0a3-885147c91128.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[ReadAddress 方法](../html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[SwitchingOn 方法](../html/279d6777-e35a-c17c-40f6-c5fdffe9ae70.htm "SwitchingOn 方法 ")

[ToString 方法](../html/ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm "ToString 方法 ")

[Trip 方法](../html/eee862c9-b6ae-307c-860f-1daf0dec3fd3.htm "Trip 方法 ")

[Write 方法](../html/5e9d2d45-4f08-9074-0428-626196c5c83b.htm "Write 方法 ")

[WriteAddress 方法](../html/5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm "WriteAddress 方法 ")

[WriteAsync 方法](../html/bb712c3c-a331-ed3e-43af-1b03a975fe1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645GetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
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

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackCommandWithHeader 方法 

[原文連結](http://api.hslcommunication.cn/html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ActiveDeveice 方法](../html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm "ActiveDeveice 方法 ")

[BroadcastTime 方法](../html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm "BroadcastTime 方法 ")

[ChangeBaudRate 方法](../html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm "ChangeBaudRate 方法 ")

[FreezeCommand 方法](../html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm "FreezeCommand 方法 ")

[GetNewNetMessage 方法](../html/0a775ede-e8a0-d411-d0a3-885147c91128.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[ReadAddress 方法](../html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[SwitchingOn 方法](../html/279d6777-e35a-c17c-40f6-c5fdffe9ae70.htm "SwitchingOn 方法 ")

[ToString 方法](../html/ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm "ToString 方法 ")

[Trip 方法](../html/eee862c9-b6ae-307c-860f-1daf0dec3fd3.htm "Trip 方法 ")

[Write 方法](../html/5e9d2d45-4f08-9074-0428-626196c5c83b.htm "Write 方法 ")

[WriteAddress 方法](../html/5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm "WriteAddress 方法 ")

[WriteAsync 方法](../html/bb712c3c-a331-ed3e-43af-1b03a975fe1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645PackCommandWithHeader 方法 |

对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  
The current command is packaged, usually carrying the content of the command header, marking the length of the current command,
and it needs to be rewritten, otherwise it is not packaged by default

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override byte[] PackCommandWithHeader(
	byte[] command
)
```

```
Public Overrides Function PackCommandWithHeader ( 
	command As Byte()
) As Byte()
```

```
public:
virtual array<unsigned char>^ PackCommandWithHeader(
	array<unsigned char>^ command
) override
```

```
abstract PackCommandWithHeader : 
        command : byte[] -> byte[] 
override PackCommandWithHeader : 
        command : byte[] -> byte[]
```

#### 参数

command
:   类型：SystemByte  
    发送的数据命令内容

#### 返回值

类型：Byte  
打包之后的数据结果信息

![](../icons/SectionExpanded.png)备注

对发送的命令打包之后，直接发送给真实的对方设备了，例如在AB-PLC里面，就重写了打包方法，将当前的会话ID参数传递给PLC设备  
After packaging the sent command, it is directly sent to the real counterpart device. For example, in AB-PLC,
the packaging method is rewritten and the current session ID parameter is passed to the PLC device.

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/825a948a-685c-b4ec-d9ed-df7766297562.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/8f3c5175-3e33-bbc6-01c4-17341746eb78.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645Read 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Read(String, UInt16)](8f3c5175-3e33-bbc6-01c4-17341746eb78.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 00-00-00-00. The separator can be any special character or no separator. (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/8f3c5175-3e33-bbc6-01c4-17341746eb78.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/8f3c5175-3e33-bbc6-01c4-17341746eb78.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645Read 方法 (String, UInt16) |

根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 00-00-00-00，分割符可以任意特殊字符或是没有分隔符。  
Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual,
from high to position, such as 00-00-00-00. The separator can be any special character or no separator.

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
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
    数据标识，具体需要查找手册来对应

length
:   类型：SystemUInt16  
    数据长度信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
结果信息

#### 实现

[IReadWriteNetRead(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm)  
[IReadWriteNetRead(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm)

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;00-00-00-00" 或是 "s=100000;00-00-02-00"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[Read 重载](825a948a-685c-b4ec-d9ed-df7766297562.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ActiveDeveice 方法](../html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm "ActiveDeveice 方法 ")

[BroadcastTime 方法](../html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm "BroadcastTime 方法 ")

[ChangeBaudRate 方法](../html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm "ChangeBaudRate 方法 ")

[FreezeCommand 方法](../html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm "FreezeCommand 方法 ")

[GetNewNetMessage 方法](../html/0a775ede-e8a0-d411-d0a3-885147c91128.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[ReadAddress 方法](../html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[SwitchingOn 方法](../html/279d6777-e35a-c17c-40f6-c5fdffe9ae70.htm "SwitchingOn 方法 ")

[ToString 方法](../html/ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm "ToString 方法 ")

[Trip 方法](../html/eee862c9-b6ae-307c-860f-1daf0dec3fd3.htm "Trip 方法 ")

[Write 方法](../html/5e9d2d45-4f08-9074-0428-626196c5c83b.htm "Write 方法 ")

[WriteAddress 方法](../html/5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm "WriteAddress 方法 ")

[WriteAsync 方法](../html/bb712c3c-a331-ed3e-43af-1b03a975fe1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ReadAddress 方法 |

读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：149100007290  
Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> ReadAddress()
```

```
Public Function ReadAddress As OperateResult(Of String)
```

```
public:
virtual OperateResult<String^>^ ReadAddress() sealed
```

```
abstract ReadAddress : unit -> OperateResult<string> 
override ReadAddress : unit -> OperateResult<string>
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
设备的通信地址

#### 实现

[IDlt645ReadAddress](66a2d3e2-d387-ed3b-18f7-e9bf90fe9412.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDouble 方法 

[原文連結](http://api.hslcommunication.cn/html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDouble 方法 (String, UInt16)](../html/ed7ea8c5-39da-19b6-6c3f-6482560d445e.htm "ReadDouble 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ReadDouble 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](ed7ea8c5-39da-19b6-6c3f-6482560d445e.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDouble 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/ed7ea8c5-39da-19b6-6c3f-6482560d445e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDouble 方法 (String, UInt16)](../html/ed7ea8c5-39da-19b6-6c3f-6482560d445e.htm "ReadDouble 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ReadDouble 方法 (String, UInt16) |

读取双浮点数据的数组  
Read double floating point data array

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<double[]> ReadDouble(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadDouble ( 
	address As String,
	length As UShort
) As OperateResult(Of Double())
```

```
public:
virtual OperateResult<array<double>^>^ ReadDouble(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadDouble : 
        address : string * 
        length : uint16 -> OperateResult<float[]> 
override ReadDouble : 
        address : string * 
        length : uint16 -> OperateResult<float[]>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数组长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Double  
带有成功标识的double数组

#### 实现

[IReadWriteNetReadDouble(String, UInt16)](a8468555-105b-337a-196f-b0a667856edd.htm)  
[IReadWriteNetReadDouble(String, UInt16)](a8468555-105b-337a-196f-b0a667856edd.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

Double类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
double[] d100_139 = melsec_net.ReadDouble( "D100", 10 ).Content;

// 如果需要判断是否读取成功

OperateResult<double[]> R_d100_139 = melsec_net.ReadDouble( "D100", 10 );
if (R_d100_139.IsSuccess)
{
    double value_d100 = R_d100_139.Content[0];
    double value_d104 = R_d100_139.Content[1];
    double value_d108 = R_d100_139.Content[2];
    double value_d112 = R_d100_139.Content[3];
    double value_d116 = R_d100_139.Content[4];
    double value_d120 = R_d100_139.Content[5];
    double value_d124 = R_d100_139.Content[6];
    double value_d128 = R_d100_139.Content[7];
    double value_d132 = R_d100_139.Content[8];
    double value_d136 = R_d100_139.Content[9];
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[ReadDouble 重载](01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDoubleAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/9c75cf35-06b1-a750-9f58-31368c242364.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadDoubleAsync 方法 (String, UInt16)](../html/9b02db80-d5a4-133f-b384-32ab306370c3.htm "ReadDoubleAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ReadDoubleAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](9b02db80-d5a4-133f-b384-32ab306370c3.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDoubleAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/9b02db80-d5a4-133f-b384-32ab306370c3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadDoubleAsync 方法 (String, UInt16)](../html/9b02db80-d5a4-133f-b384-32ab306370c3.htm "ReadDoubleAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ReadDoubleAsync 方法 (String, UInt16) |

读取双浮点数据的数组  
Read double floating point data array

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<double[]>> ReadDoubleAsync(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadDoubleAsync ( 
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Double()))
```

```
public:
virtual Task<OperateResult<array<double>^>^>^ ReadDoubleAsync(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadDoubleAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<float[]>> 
override ReadDoubleAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<float[]>>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数组长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Double  
带有成功标识的double数组

#### 实现

[IReadWriteNetReadDoubleAsync(String, UInt16)](3a436322-c7f5-bb1d-e97b-dc6c0d6990d2.htm)  
[IReadWriteNetReadDoubleAsync(String, UInt16)](3a436322-c7f5-bb1d-e97b-dc6c0d6990d2.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

Double类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
double[] d100_139 = melsec_net.ReadDouble( "D100", 10 ).Content;

// 如果需要判断是否读取成功

OperateResult<double[]> R_d100_139 = melsec_net.ReadDouble( "D100", 10 );
if (R_d100_139.IsSuccess)
{
    double value_d100 = R_d100_139.Content[0];
    double value_d104 = R_d100_139.Content[1];
    double value_d108 = R_d100_139.Content[2];
    double value_d112 = R_d100_139.Content[3];
    double value_d116 = R_d100_139.Content[4];
    double value_d120 = R_d100_139.Content[5];
    double value_d124 = R_d100_139.Content[6];
    double value_d128 = R_d100_139.Content[7];
    double value_d132 = R_d100_139.Content[8];
    double value_d136 = R_d100_139.Content[9];
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[ReadDoubleAsync 重载](9c75cf35-06b1-a750-9f58-31368c242364.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadString 方法 

[原文連結](http://api.hslcommunication.cn/html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadString 方法 (String, UInt16, Encoding)](../html/909dd3db-f091-c133-d10c-b83dd484e56a.htm "ReadString 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ReadString 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](909dd3db-f091-c133-d10c-b83dd484e56a.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadString 方法 (String, UInt16, Encoding)

[原文連結](http://api.hslcommunication.cn/html/909dd3db-f091-c133-d10c-b83dd484e56a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadString 方法 (String, UInt16, Encoding)](../html/909dd3db-f091-c133-d10c-b83dd484e56a.htm "ReadString 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ReadString 方法 (String, UInt16, Encoding) |

使用指定的编码，读取字符串数据  
Reads string data using the specified encoding

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
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

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[ReadString 重载](d15c2fc1-0806-87df-5d68-23d4710c1e30.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringArray 方法 

[原文連結](http://api.hslcommunication.cn/html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ActiveDeveice 方法](../html/b3e8f555-2f00-5be6-a2da-a99db9b2f0df.htm "ActiveDeveice 方法 ")

[BroadcastTime 方法](../html/1a1dcab0-25bc-c17c-6fc6-5b0223c7218d.htm "BroadcastTime 方法 ")

[ChangeBaudRate 方法](../html/5d5ebb18-3450-516c-71ad-e122c04c7541.htm "ChangeBaudRate 方法 ")

[FreezeCommand 方法](../html/c14d9771-4ef2-3bb4-809b-cfd507cf30e7.htm "FreezeCommand 方法 ")

[GetNewNetMessage 方法](../html/0a775ede-e8a0-d411-d0a3-885147c91128.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/3de72c9d-d491-5f63-9301-3a87bdff5d03.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/825a948a-685c-b4ec-d9ed-df7766297562.htm "Read 方法 ")

[ReadAddress 方法](../html/16df7ed3-d2a5-81ff-f530-25543a56ec4c.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/01b4a3b9-6313-dd1a-b510-44f85b283a8e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/9c75cf35-06b1-a750-9f58-31368c242364.htm "ReadDoubleAsync 方法 ")

[ReadString 方法](../html/d15c2fc1-0806-87df-5d68-23d4710c1e30.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/52c56468-5e8e-2fd9-9aaa-2fd332283a95.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[SwitchingOn 方法](../html/279d6777-e35a-c17c-40f6-c5fdffe9ae70.htm "SwitchingOn 方法 ")

[ToString 方法](../html/ddc87514-c61b-0a4c-065c-20c4bbf792f9.htm "ToString 方法 ")

[Trip 方法](../html/eee862c9-b6ae-307c-860f-1daf0dec3fd3.htm "Trip 方法 ")

[Write 方法](../html/5e9d2d45-4f08-9074-0428-626196c5c83b.htm "Write 方法 ")

[WriteAddress 方法](../html/5f73e5b3-45d4-9113-ba15-cb1b85a82ae4.htm "WriteAddress 方法 ")

[WriteAsync 方法](../html/bb712c3c-a331-ed3e-43af-1b03a975fe1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ReadStringArray 方法 |

读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据，但是少部分的地址存在多个数据，例如 01-01-00-00 正向有功总需求及发生时间  
Read all the string data information of the specified address, in general, there is only one data for one address, but there are multiple data for a small number of addresses,
such as 01-01-00-00 Forward active total demand and occurrence time

**命名空间：**
 [HslCommunication.Instrument.DLT](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string[]> ReadStringArray(
	string address
)
```

```
Public Function ReadStringArray ( 
	address As String
) As OperateResult(Of String())
```

```
public:
virtual OperateResult<array<String^>^>^ ReadStringArray(
	String^ address
) sealed
```

```
abstract ReadStringArray : 
        address : string -> OperateResult<string[]> 
override ReadStringArray : 
        address : string -> OperateResult<string[]>
```

#### 参数

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
字符串数组信息

#### 实现

[IDlt645ReadStringArray(String)](1cb657b1-0871-d56d-57cd-4d9f7fe77e98.htm)

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;00-00-00-00" 或是 "s=100000;00-00-02-00"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能  
地址也可以携带是否数据翻转的标记，例如 "reverse=false;00-00-00-00" 解析数据的时候就不发生反转的操作

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.DLT](../html/29f34251-d8c0-8a17-1598-9c922c0ae55c.htm "HslCommunication.Instrument.DLT")

[DLT645 类](../html/dc3f4d48-9b60-0203-88c1-3543260dee17.htm "DLT645 类")

[DLT645 方法](../html/1f89f78d-4b8c-9b17-fc93-6bf40e87bcb9.htm "DLT645 方法")

[ReadStringAsync 方法](../html/d6c80dcf-6794-5ab3-a8b7-9e1df01ee229.htm "ReadStringAsync 方法 ")

[ReadStringAsync 方法 (String, UInt16, Encoding)](../html/8163a33b-8f07-a792-0382-cdf3651e8849.htm "ReadStringAsync 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DLT645ReadStringAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](8163a33b-8f07-a792-0382-cdf3651e8849.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DLT645 类](dc3f4d48-9b60-0203-88c1-3543260dee17.htm)

[HslCommunication.Instrument.DLT 命名空间](29f34251-d8c0-8a17-1598-9c922c0ae55c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)