# HslCommunication - HslCommunication.Robot.ABB

> 分類頁數: 30



---
## HslCommunication.Robot.ABB

[原文連結](http://api.hslcommunication.cn/html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiServer 类](../html/64696c60-86d5-bc02-eea6-eb5f756ca474.htm "ABBWebApiServer 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Robot.ABB 命名空间 |

[缺少 "N:HslCommunication.Robot.ABB" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [ABBWebApiClient](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm) | ABB机器人的web api接口的客户端，可以方便快速的获取到abb机器人的一些数据信息  The client of ABB robot's web API interface can easily and quickly obtain some data information of ABB robot |
| 公共类 | [ABBWebApiServer](64696c60-86d5-bc02-eea6-eb5f756ca474.htm) | **[商业授权]** ABB机器人的虚拟服务器，基于WebApi协议构建，可用于读取一些数据信息 **[Authorization]**The virtual server of ABB robot, built based on the WebApi protocol, can be used to read some data information |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ABBWebApiClient 类

[原文連結](http://api.hslcommunication.cn/html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 构造函数](../html/4c320cc9-6510-60b7-55a5-723d5fd595d4.htm "ABBWebApiClient 构造函数 ")

[ABBWebApiClient 属性](../html/4b1c4570-961b-61f2-39d2-0694b0fb3230.htm "ABBWebApiClient 属性")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClient 类 |

ABB机器人的web api接口的客户端，可以方便快速的获取到abb机器人的一些数据信息  
The client of ABB robot's web API interface can easily and quickly obtain some data information of ABB robot

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetNetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)  
    [HslCommunication.Core.NetNetworkWebApiRobotBase](3303830b-83fb-9aba-9684-41bba031dd90.htm)  
      HslCommunication.Robot.ABBABBWebApiClient

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class ABBWebApiClient : NetworkWebApiRobotBase, 
	IRobotNet
```

```
Public Class ABBWebApiClient
	Inherits NetworkWebApiRobotBase
	Implements IRobotNet
```

```
public ref class ABBWebApiClient : public NetworkWebApiRobotBase, 
	IRobotNet
```

```
type ABBWebApiClient =  
    class
        inherit NetworkWebApiRobotBase
        interface IRobotNet
    end
```

ABBWebApiClient 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ABBWebApiClient(String)](7970e083-23d1-8515-f806-cdcce1bef8c7.htm) | 使用指定的ip地址来初始化对象  Initializes the object using the specified IP address |
| 公共方法 | [ABBWebApiClient(String, Int32)](6cfb437e-abd6-1f0f-6db1-461722366ca2.htm) | 使用指定的ip地址和端口号来初始化对象  Initializes the object with the specified IP address and port number |
| 公共方法 | [ABBWebApiClient(String, Int32, String, String)](cfec4d27-c62f-38ed-faac-8e17528b163c.htm) | 使用指定的ip地址，端口号，用户名，密码来初始化对象  Initialize the object with the specified IP address, port number, username, and password |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddRequestHeadersAction](8802f5ed-2c21-7c00-f866-13e2305b29ad.htm) | 针对请求的头信息进行额外的处理 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [Client](9cd17e98-778d-9fce-82b0-23220be00423.htm) | 获取当前的HttpClinet的客户端  Get the current HttpClinet client (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [DefaultContentType](ac7a98ea-f82d-9616-c8be-3e66b3d6c765.htm) | 默认的内容类型，如果为空，则不进行设置操作。例如设置为 "text/plain", "application/json", "text/html" 等等。  The default content type, if it is empty, no setting operation will be performed. For example, set to "text/plain", "application/json", "text/html" and so on. (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [Host](c35c1951-f4ab-888e-a835-a61ac1fc9055.htm) | 获取当前的远程服务器的地址，可能是ip，也可能是网址。  Get the address of the current remote server, which may be an IP address or a web address. (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [IpAddress](1b964583-7670-bdb2-9d75-d26a31cfc6e4.htm) | 获取或设置远程服务器的IP地址  Gets or sets the IP address of the remote server (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性代码示例 | [LogNet](ea0f7c15-1a07-f1f9-08a1-8f80758e9713.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [Password](e284c994-bd49-b4c5-7a33-7470c9a290b0.htm) | 获取或设置当前的密码  Get or set the current password (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [Port](26bd23af-c78c-31fe-9fc8-e0fabd5efc93.htm) | 获取或设置远程服务器的端口号信息  Gets or sets the port number information for the remote server (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [UseEncodingISO](940f46b4-8176-1d00-e41f-18d2389adf5e.htm) | 获取或设置是否使用ISO的编码信息，默认为 False  Get or set whether to use ISO encoding information, the default is False (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [UseHttps](ce585636-02cb-4e69-4daa-58bfeaa4f7bd.htm) | 是否启用Https的协议访问，对于Https来说，端口号默认为 443  Whether to enable Https protocol access, for Https, the port number defaults to 443 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [UserName](a847a830-ec29-4043-cdb2-1fc6cd6b9bf9.htm) | 获取或设置当前的用户名  Get or set the current username (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [AddRequestHeaders](48e7d402-d12c-b73e-bdcd-3ccd368dc588.htm) | 针对请求的头信息进行额外的处理，可以重写用来实现一些特殊的信息添加到请求头中 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [Get](c7de454e-7f11-4bf3-8820-8416dfff8773.htm) | 使用GET操作从网络中获取到数据信息，地址 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [GetAnIOSignal](d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm) | 获取机器人的IO信号资源  Get an IO signal resource. |
| 公共方法 | [GetAnIOSignalAsync](0e63db29-1826-248c-2c6f-13e158318127.htm) | 获取机器人的IO信号资源  Get an IO signal resource. |
| 公共方法 | [GetAsync](feeda5e6-59eb-bf36-ae9e-45b47fc9f8fe.htm) | 使用GET操作从网络中获取到数据信息，地址 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [GetCtrlState](57e5860a-18bc-4f4a-06b7-0c17555f460d.htm) | 获取当前的控制状态，Content属性就是机器人的控制信息  Get the current control state. The Content attribute is the control information of the robot |
| 公共方法 | [GetCtrlStateAsync](996af0a7-6e58-911a-541b-6dfe417d7bba.htm) | 获取当前的控制状态，Content属性就是机器人的控制信息  Get the current control state. The Content attribute is the control information of the robot |
| 公共方法 | [GetEntireUrl](372891b8-649b-08bd-80ff-61609f6b6c92.htm) | 根据当前的url信息，获取到完整的url地址，如果没有http头信息，则自动添加上  Based on the current url information, the complete url address is obtained. If there is no http header information, it will be automatically added (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [GetErrorState](858227c4-967f-fa0b-4b77-7558db99421c.htm) | 获取当前的错误状态，Content属性就是机器人的状态信息  Gets the current error state. The Content attribute is the state information of the robot |
| 公共方法 | [GetErrorStateAsync](9ea55110-7354-4743-45d7-677395d3b9da.htm) | 获取当前的错误状态，Content属性就是机器人的状态信息  Gets the current error state. The Content attribute is the state information of the robot |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetIO2In](fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm) | 获取当前机器人的本机的输入IO  Gets the input IO of the current robot's native |
| 公共方法 | [GetIO2InAsync](7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm) | 获取当前机器人的本机的输入IO  Gets the input IO of the current robot's native |
| 公共方法 | [GetIO2Out](fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm) | 获取当前机器人的本机的输出IO  Gets the output IO of the current robot's native |
| 公共方法 | [GetIO2OutAsync](8c3a0765-19ba-0676-dc85-22cdbf22498b.htm) | 获取当前机器人的本机的输出IO  Gets the output IO of the current robot's native |
| 公共方法 | [GetIOIn](51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm) | 获取当前机器人的本机的输入IO  Gets the input IO of the current robot's native |
| 公共方法 | [GetIOInAsync](da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm) | 获取当前机器人的本机的输入IO  Gets the input IO of the current robot's native |
| 公共方法 | [GetIOOut](ee3b44e0-a853-ccdd-e475-98552d0966a8.htm) | 获取当前机器人的本机的输出IO  Gets the output IO of the current robot's native |
| 公共方法 | [GetIOOutAsync](7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm) | 获取当前机器人的本机的输出IO  Gets the output IO of the current robot's native |
| 公共方法 | [GetJointTarget](8d619075-fc31-7a20-712e-cb9483e48b09.htm) | 获取当前机器人的物理关节点信息，返回json格式的关节信息  Get the physical node information of the current robot and return the joint information in json format |
| 公共方法 | [GetJointTargetAsync](2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm) | 获取当前机器人的物理关节点信息，返回json格式的关节信息  Get the physical node information of the current robot and return the joint information in json format |
| 公共方法 | [GetLog](16aa6e7b-178b-51ee-8039-400d94447c39.htm) | 获取当前机器人的日志记录，默认记录为10条  Gets the log record for the current robot, which is 10 by default |
| 公共方法 | [GetLogAsync](40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm) | 获取当前机器人的日志记录，默认记录为10条  Gets the log record for the current robot, which is 10 by default |
| 公共方法 | [GetOperationMode](080723d7-a5ca-c1db-c967-22871a231339.htm) | 获取当前机器人的工作模式  Gets the current working mode of the robot |
| 公共方法 | [GetOperationModeAsync](4c37c64d-584f-02b0-912c-d78cf7764b29.htm) | 获取当前机器人的工作模式  Gets the current working mode of the robot |
| 公共方法 | [GetRapidExecution](ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm) | 获取当前机器人的当前程序运行状态  Get the current program running status of the current robot |
| 公共方法 | [GetRapidExecutionAsync](e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm) | 获取当前机器人的当前程序运行状态  Get the current program running status of the current robot |
| 公共方法 | [GetRapidTasks](4fbe4e80-94f1-45ba-6007-92328676e6ea.htm) | 获取当前机器人的任务列表  Get the task list of the current robot |
| 公共方法 | [GetRapidTasksAsync](bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm) | 获取当前机器人的任务列表  Get the task list of the current robot |
| 公共方法 | [GetRobotTarget](78b1b85c-db0f-a511-7e12-746069550e04.htm) | 获取机器人的目标坐标信息  Get the current robot's target information |
| 公共方法 | [GetRobotTargetAsync](17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm) | 获取机器人的目标坐标信息  Get the current robot's target information |
| 公共方法静态成员 | [GetSelectStrings](ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm) | 获取当前支持的读取的地址列表  Gets a list of addresses for currently supported reads |
| 公共方法 | [GetServoEnable](ad00798e-b449-5ac9-1902-e2b8c987003d.htm) | 获取当前机器人的伺服使能状态  Get the current robot servo enable state |
| 公共方法 | [GetServoEnableAsync](7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm) | 获取当前机器人的伺服使能状态  Get the current robot servo enable state |
| 公共方法 | [GetSpeedRatio](51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm) | 获取当前机器人的速度配比信息  Get the speed matching information of the current robot |
| 公共方法 | [GetSpeedRatioAsync](e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm) | 获取当前机器人的速度配比信息  Get the speed matching information of the current robot |
| 公共方法 | [GetSystem](852f33fa-83bf-42e7-9e08-33fc34a82600.htm) | 获取当前机器人的系统信息，版本号，唯一ID等信息  Get the current robot's system information, version number, unique ID and other information |
| 公共方法 | [GetSystemAsync](7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm) | 获取当前机器人的系统信息，版本号，唯一ID等信息  Get the current robot's system information, version number, unique ID and other information |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [GetUserValue](5470b496-462f-c3da-153f-dac92a70066c.htm) | 根据给定的名称，获取当前用户的数据值信息。  According to the given name, gets the data value information of the current user |
| 公共方法 | [GetUserValueAsync](d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm) | 根据给定的名称，获取当前用户的数据值信息。  According to the given name, gets the data value information of the current user |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Post](8be86099-0853-7a12-e310-23383ec0bc7e.htm) | 使用POST命令去提交数据内容，然后返回相关的数据信息 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [PostAsync](e0b8ac89-a0ff-4d7b-1e76-de08e7d31de1.htm) | 使用POST命令去提交数据内容，然后返回相关的数据信息 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [Read](bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm) | 读取对方信息的的数据信息，通常是针对GET的方法信息设计的。如果使用了url=开头，就表示是使用了原生的地址访问  Read the other side of the data information, usually designed for the GET method information.If you start with url=, you are using native address access (重写 [NetworkWebApiRobotBaseRead(String)](e2e4dd7f-efbe-ae6b-abb2-e3bec451182d.htm).) |
| 公共方法 | [ReadAsync](d8667131-2445-aa8f-7c68-865fdc537237.htm) | 读取对方信息的的数据信息，通常是针对GET的方法信息设计的。如果使用了url=开头，就表示是使用了原生的地址访问  Read the other side of the data information, usually designed for the GET method information.If you start with url=, you are using native address access (继承自 [NetworkWebApiRobotBase](3303830b-83fb-9aba-9684-41bba031dd90.htm)。) |
| 受保护的方法 | [ReadByAddress](4ce59148-7bef-049e-a0ca-d190467b1938.htm) | 等待重写的额外的指令信息的支持。除了url的形式之外，还支持基于命令的数据交互  Additional instruction information waiting for rewriting is supported.In addition to the url format, command based data interaction is supported (重写 [NetworkWebApiRobotBaseReadByAddress(String)](106dd981-59b1-b17d-a183-0ef15117631f.htm).) |
| 受保护的方法 | [ReadByAddressAsync](be2f191d-9781-c9a0-5800-ed8c52088aae.htm) | 等待重写的额外的指令信息的支持。除了url的形式之外，还支持基于命令的数据交互  Additional instruction information waiting for rewriting is supported.In addition to the url format, command based data interaction is supported (重写 [NetworkWebApiRobotBaseReadByAddressAsync(String)](e1e23509-a944-bf67-c503-057ddf581cac.htm).) |
| 公共方法 | [ReadString](bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm) | 读取对方信息的的字符串数据信息，通常是针对GET的方法信息设计的。如果使用了url=开头，就表示是使用了原生的地址访问  The string data information that reads the other party information, usually designed for the GET method information.If you start with url=, you are using native address access (重写 [NetworkWebApiRobotBaseReadString(String)](c4dc02be-9ee9-ce47-0b44-b41e08635fb3.htm).) |
| 公共方法 | [ReadStringAsync](03cc49bd-e193-7225-6447-2e62b69ddec1.htm) | 读取对方信息的的字符串数据信息，通常是针对GET的方法信息设计的。如果使用了url=开头，就表示是使用了原生的地址访问  The string data information that reads the other party information, usually designed for the GET method information.If you start with url=, you are using native address access (继承自 [NetworkWebApiRobotBase](3303830b-83fb-9aba-9684-41bba031dd90.htm)。) |
| 公共方法 | [RequestAsync](649024d2-2df8-56e7-14ad-9b884701fab3.htm) | 发起一个请求，然后返回结果，支持任意的http方法  Initiate a request and then return the result, supporting any http method (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [ToString](969cd64f-a556-542a-b5d0-6e82066d44ff.htm) | (重写 [NetworkWebApiRobotBaseToString](bc8a32fc-dd82-78db-875f-de4b14a23dd3.htm).) |
| 公共方法 | [Write(String, Byte)](36744e4b-b339-1a38-ac3c-47ef3dc3b71e.htm) | 使用POST的方式来向对方进行请求数据信息，需要使用url=开头，来表示是使用了原生的地址访问  Using POST to request data information from the other party, we need to start with url= to indicate that we are using native address access (重写 [NetworkWebApiRobotBaseWrite(String, Byte)](74d6da2a-aa64-98ad-3d8c-0ea865a6755f.htm).) |
| 公共方法 | [Write(String, String)](b227fd97-6aeb-90f2-7271-fd676a3e2e54.htm) | 使用POST的方式来向对方进行请求数据信息，需要使用url=开头，来表示是使用了原生的地址访问  Using POST to request data information from the other party, we need to start with url= to indicate that we are using native address access (重写 [NetworkWebApiRobotBaseWrite(String, String)](5654dcab-d566-4511-2024-bc51004c9dad.htm).) |
| 公共方法 | [WriteAsync(String, Byte)](e1222bd2-37fe-3ad9-f20a-15c324142e96.htm) | 使用POST的方式来向对方进行请求数据信息，需要使用url=开头，来表示是使用了原生的地址访问  Using POST to request data information from the other party, we need to start with url= to indicate that we are using native address access (继承自 [NetworkWebApiRobotBase](3303830b-83fb-9aba-9684-41bba031dd90.htm)。) |
| 公共方法 | [WriteAsync(String, String)](3e81cc81-a117-ac9e-8c1a-31b021ab3d3d.htm) | 使用POST的方式来向对方进行请求数据信息，需要使用url=开头，来表示是使用了原生的地址访问  Using POST to request data information from the other party, we need to start with url= to indicate that we are using native address access (继承自 [NetworkWebApiRobotBase](3303830b-83fb-9aba-9684-41bba031dd90.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)备注

参考的界面信息是：http://developercenter.robotstudio.com/webservice/api\_reference
关于额外的地址说明，如果想要查看，可以调用[GetSelectStrings](ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm) 返回字符串列表来看看。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ABBWebApiClient 构造函数 

[原文連結](http://api.hslcommunication.cn/html/4c320cc9-6510-60b7-55a5-723d5fd595d4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 构造函数](../html/4c320cc9-6510-60b7-55a5-723d5fd595d4.htm "ABBWebApiClient 构造函数 ")

[ABBWebApiClient 构造函数 (String)](../html/7970e083-23d1-8515-f806-cdcce1bef8c7.htm "ABBWebApiClient 构造函数 (String)")

[ABBWebApiClient 构造函数 (String, Int32)](../html/6cfb437e-abd6-1f0f-6db1-461722366ca2.htm "ABBWebApiClient 构造函数 (String, Int32)")

[ABBWebApiClient 构造函数 (String, Int32, String, String)](../html/cfec4d27-c62f-38ed-faac-8e17528b163c.htm "ABBWebApiClient 构造函数 (String, Int32, String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClient 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ABBWebApiClient(String)](7970e083-23d1-8515-f806-cdcce1bef8c7.htm) | 使用指定的ip地址来初始化对象  Initializes the object using the specified IP address |
| 公共方法 | [ABBWebApiClient(String, Int32)](6cfb437e-abd6-1f0f-6db1-461722366ca2.htm) | 使用指定的ip地址和端口号来初始化对象  Initializes the object with the specified IP address and port number |
| 公共方法 | [ABBWebApiClient(String, Int32, String, String)](cfec4d27-c62f-38ed-faac-8e17528b163c.htm) | 使用指定的ip地址，端口号，用户名，密码来初始化对象  Initialize the object with the specified IP address, port number, username, and password |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ABBWebApiClient 构造函数 (String)

[原文連結](http://api.hslcommunication.cn/html/7970e083-23d1-8515-f806-cdcce1bef8c7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 构造函数](../html/4c320cc9-6510-60b7-55a5-723d5fd595d4.htm "ABBWebApiClient 构造函数 ")

[ABBWebApiClient 构造函数 (String)](../html/7970e083-23d1-8515-f806-cdcce1bef8c7.htm "ABBWebApiClient 构造函数 (String)")

[ABBWebApiClient 构造函数 (String, Int32)](../html/6cfb437e-abd6-1f0f-6db1-461722366ca2.htm "ABBWebApiClient 构造函数 (String, Int32)")

[ABBWebApiClient 构造函数 (String, Int32, String, String)](../html/cfec4d27-c62f-38ed-faac-8e17528b163c.htm "ABBWebApiClient 构造函数 (String, Int32, String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClient 构造函数 (String) |

使用指定的ip地址来初始化对象  
Initializes the object using the specified IP address

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ABBWebApiClient(
	string ipAddress
)
```

```
Public Sub New ( 
	ipAddress As String
)
```

```
public:
ABBWebApiClient(
	String^ ipAddress
)
```

```
new : 
        ipAddress : string -> ABBWebApiClient
```

#### 参数

ipAddress
:   类型：SystemString  
    Ip地址信息

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[ABBWebApiClient 重载](4c320cc9-6510-60b7-55a5-723d5fd595d4.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ABBWebApiClient 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/6cfb437e-abd6-1f0f-6db1-461722366ca2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 构造函数](../html/4c320cc9-6510-60b7-55a5-723d5fd595d4.htm "ABBWebApiClient 构造函数 ")

[ABBWebApiClient 构造函数 (String)](../html/7970e083-23d1-8515-f806-cdcce1bef8c7.htm "ABBWebApiClient 构造函数 (String)")

[ABBWebApiClient 构造函数 (String, Int32)](../html/6cfb437e-abd6-1f0f-6db1-461722366ca2.htm "ABBWebApiClient 构造函数 (String, Int32)")

[ABBWebApiClient 构造函数 (String, Int32, String, String)](../html/cfec4d27-c62f-38ed-faac-8e17528b163c.htm "ABBWebApiClient 构造函数 (String, Int32, String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClient 构造函数 (String, Int32) |

使用指定的ip地址和端口号来初始化对象  
Initializes the object with the specified IP address and port number

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ABBWebApiClient(
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
ABBWebApiClient(
	String^ ipAddress, 
	int port
)
```

```
new : 
        ipAddress : string * 
        port : int -> ABBWebApiClient
```

#### 参数

ipAddress
:   类型：SystemString  
    Ip地址信息

port
:   类型：SystemInt32  
    端口号信息

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[ABBWebApiClient 重载](4c320cc9-6510-60b7-55a5-723d5fd595d4.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ABBWebApiClient 构造函数 (String, Int32, String, String)

[原文連結](http://api.hslcommunication.cn/html/cfec4d27-c62f-38ed-faac-8e17528b163c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 构造函数](../html/4c320cc9-6510-60b7-55a5-723d5fd595d4.htm "ABBWebApiClient 构造函数 ")

[ABBWebApiClient 构造函数 (String)](../html/7970e083-23d1-8515-f806-cdcce1bef8c7.htm "ABBWebApiClient 构造函数 (String)")

[ABBWebApiClient 构造函数 (String, Int32)](../html/6cfb437e-abd6-1f0f-6db1-461722366ca2.htm "ABBWebApiClient 构造函数 (String, Int32)")

[ABBWebApiClient 构造函数 (String, Int32, String, String)](../html/cfec4d27-c62f-38ed-faac-8e17528b163c.htm "ABBWebApiClient 构造函数 (String, Int32, String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClient 构造函数 (String, Int32, String, String) |

使用指定的ip地址，端口号，用户名，密码来初始化对象  
Initialize the object with the specified IP address, port number, username, and password

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ABBWebApiClient(
	string ipAddress,
	int port,
	string name,
	string password
)
```

```
Public Sub New ( 
	ipAddress As String,
	port As Integer,
	name As String,
	password As String
)
```

```
public:
ABBWebApiClient(
	String^ ipAddress, 
	int port, 
	String^ name, 
	String^ password
)
```

```
new : 
        ipAddress : string * 
        port : int * 
        name : string * 
        password : string -> ABBWebApiClient
```

#### 参数

ipAddress
:   类型：SystemString  
    Ip地址信息

port
:   类型：SystemInt32  
    端口号信息

name
:   类型：SystemString  
    用户名

password
:   类型：SystemString  
    密码

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[ABBWebApiClient 重载](4c320cc9-6510-60b7-55a5-723d5fd595d4.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ABBWebApiClient 属性

[原文連結](http://api.hslcommunication.cn/html/4b1c4570-961b-61f2-39d2-0694b0fb3230.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 构造函数](../html/4c320cc9-6510-60b7-55a5-723d5fd595d4.htm "ABBWebApiClient 构造函数 ")

[ABBWebApiClient 属性](../html/4b1c4570-961b-61f2-39d2-0694b0fb3230.htm "ABBWebApiClient 属性")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClient 属性 |

[ABBWebApiClient](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddRequestHeadersAction](8802f5ed-2c21-7c00-f866-13e2305b29ad.htm) | 针对请求的头信息进行额外的处理 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [Client](9cd17e98-778d-9fce-82b0-23220be00423.htm) | 获取当前的HttpClinet的客户端  Get the current HttpClinet client (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [DefaultContentType](ac7a98ea-f82d-9616-c8be-3e66b3d6c765.htm) | 默认的内容类型，如果为空，则不进行设置操作。例如设置为 "text/plain", "application/json", "text/html" 等等。  The default content type, if it is empty, no setting operation will be performed. For example, set to "text/plain", "application/json", "text/html" and so on. (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [Host](c35c1951-f4ab-888e-a835-a61ac1fc9055.htm) | 获取当前的远程服务器的地址，可能是ip，也可能是网址。  Get the address of the current remote server, which may be an IP address or a web address. (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [IpAddress](1b964583-7670-bdb2-9d75-d26a31cfc6e4.htm) | 获取或设置远程服务器的IP地址  Gets or sets the IP address of the remote server (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性代码示例 | [LogNet](ea0f7c15-1a07-f1f9-08a1-8f80758e9713.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [Password](e284c994-bd49-b4c5-7a33-7470c9a290b0.htm) | 获取或设置当前的密码  Get or set the current password (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [Port](26bd23af-c78c-31fe-9fc8-e0fabd5efc93.htm) | 获取或设置远程服务器的端口号信息  Gets or sets the port number information for the remote server (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [UseEncodingISO](940f46b4-8176-1d00-e41f-18d2389adf5e.htm) | 获取或设置是否使用ISO的编码信息，默认为 False  Get or set whether to use ISO encoding information, the default is False (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [UseHttps](ce585636-02cb-4e69-4daa-58bfeaa4f7bd.htm) | 是否启用Https的协议访问，对于Https来说，端口号默认为 443  Whether to enable Https protocol access, for Https, the port number defaults to 443 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共属性 | [UserName](a847a830-ec29-4043-cdb2-1fc6cd6b9bf9.htm) | 获取或设置当前的用户名  Get or set the current username (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ABBWebApiClient 方法

[原文連結](http://api.hslcommunication.cn/html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClient 方法 |

[ABBWebApiClient](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [AddRequestHeaders](48e7d402-d12c-b73e-bdcd-3ccd368dc588.htm) | 针对请求的头信息进行额外的处理，可以重写用来实现一些特殊的信息添加到请求头中 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [Get](c7de454e-7f11-4bf3-8820-8416dfff8773.htm) | 使用GET操作从网络中获取到数据信息，地址 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [GetAnIOSignal](d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm) | 获取机器人的IO信号资源  Get an IO signal resource. |
| 公共方法 | [GetAnIOSignalAsync](0e63db29-1826-248c-2c6f-13e158318127.htm) | 获取机器人的IO信号资源  Get an IO signal resource. |
| 公共方法 | [GetAsync](feeda5e6-59eb-bf36-ae9e-45b47fc9f8fe.htm) | 使用GET操作从网络中获取到数据信息，地址 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [GetCtrlState](57e5860a-18bc-4f4a-06b7-0c17555f460d.htm) | 获取当前的控制状态，Content属性就是机器人的控制信息  Get the current control state. The Content attribute is the control information of the robot |
| 公共方法 | [GetCtrlStateAsync](996af0a7-6e58-911a-541b-6dfe417d7bba.htm) | 获取当前的控制状态，Content属性就是机器人的控制信息  Get the current control state. The Content attribute is the control information of the robot |
| 公共方法 | [GetEntireUrl](372891b8-649b-08bd-80ff-61609f6b6c92.htm) | 根据当前的url信息，获取到完整的url地址，如果没有http头信息，则自动添加上  Based on the current url information, the complete url address is obtained. If there is no http header information, it will be automatically added (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [GetErrorState](858227c4-967f-fa0b-4b77-7558db99421c.htm) | 获取当前的错误状态，Content属性就是机器人的状态信息  Gets the current error state. The Content attribute is the state information of the robot |
| 公共方法 | [GetErrorStateAsync](9ea55110-7354-4743-45d7-677395d3b9da.htm) | 获取当前的错误状态，Content属性就是机器人的状态信息  Gets the current error state. The Content attribute is the state information of the robot |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetIO2In](fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm) | 获取当前机器人的本机的输入IO  Gets the input IO of the current robot's native |
| 公共方法 | [GetIO2InAsync](7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm) | 获取当前机器人的本机的输入IO  Gets the input IO of the current robot's native |
| 公共方法 | [GetIO2Out](fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm) | 获取当前机器人的本机的输出IO  Gets the output IO of the current robot's native |
| 公共方法 | [GetIO2OutAsync](8c3a0765-19ba-0676-dc85-22cdbf22498b.htm) | 获取当前机器人的本机的输出IO  Gets the output IO of the current robot's native |
| 公共方法 | [GetIOIn](51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm) | 获取当前机器人的本机的输入IO  Gets the input IO of the current robot's native |
| 公共方法 | [GetIOInAsync](da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm) | 获取当前机器人的本机的输入IO  Gets the input IO of the current robot's native |
| 公共方法 | [GetIOOut](ee3b44e0-a853-ccdd-e475-98552d0966a8.htm) | 获取当前机器人的本机的输出IO  Gets the output IO of the current robot's native |
| 公共方法 | [GetIOOutAsync](7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm) | 获取当前机器人的本机的输出IO  Gets the output IO of the current robot's native |
| 公共方法 | [GetJointTarget](8d619075-fc31-7a20-712e-cb9483e48b09.htm) | 获取当前机器人的物理关节点信息，返回json格式的关节信息  Get the physical node information of the current robot and return the joint information in json format |
| 公共方法 | [GetJointTargetAsync](2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm) | 获取当前机器人的物理关节点信息，返回json格式的关节信息  Get the physical node information of the current robot and return the joint information in json format |
| 公共方法 | [GetLog](16aa6e7b-178b-51ee-8039-400d94447c39.htm) | 获取当前机器人的日志记录，默认记录为10条  Gets the log record for the current robot, which is 10 by default |
| 公共方法 | [GetLogAsync](40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm) | 获取当前机器人的日志记录，默认记录为10条  Gets the log record for the current robot, which is 10 by default |
| 公共方法 | [GetOperationMode](080723d7-a5ca-c1db-c967-22871a231339.htm) | 获取当前机器人的工作模式  Gets the current working mode of the robot |
| 公共方法 | [GetOperationModeAsync](4c37c64d-584f-02b0-912c-d78cf7764b29.htm) | 获取当前机器人的工作模式  Gets the current working mode of the robot |
| 公共方法 | [GetRapidExecution](ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm) | 获取当前机器人的当前程序运行状态  Get the current program running status of the current robot |
| 公共方法 | [GetRapidExecutionAsync](e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm) | 获取当前机器人的当前程序运行状态  Get the current program running status of the current robot |
| 公共方法 | [GetRapidTasks](4fbe4e80-94f1-45ba-6007-92328676e6ea.htm) | 获取当前机器人的任务列表  Get the task list of the current robot |
| 公共方法 | [GetRapidTasksAsync](bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm) | 获取当前机器人的任务列表  Get the task list of the current robot |
| 公共方法 | [GetRobotTarget](78b1b85c-db0f-a511-7e12-746069550e04.htm) | 获取机器人的目标坐标信息  Get the current robot's target information |
| 公共方法 | [GetRobotTargetAsync](17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm) | 获取机器人的目标坐标信息  Get the current robot's target information |
| 公共方法静态成员 | [GetSelectStrings](ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm) | 获取当前支持的读取的地址列表  Gets a list of addresses for currently supported reads |
| 公共方法 | [GetServoEnable](ad00798e-b449-5ac9-1902-e2b8c987003d.htm) | 获取当前机器人的伺服使能状态  Get the current robot servo enable state |
| 公共方法 | [GetServoEnableAsync](7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm) | 获取当前机器人的伺服使能状态  Get the current robot servo enable state |
| 公共方法 | [GetSpeedRatio](51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm) | 获取当前机器人的速度配比信息  Get the speed matching information of the current robot |
| 公共方法 | [GetSpeedRatioAsync](e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm) | 获取当前机器人的速度配比信息  Get the speed matching information of the current robot |
| 公共方法 | [GetSystem](852f33fa-83bf-42e7-9e08-33fc34a82600.htm) | 获取当前机器人的系统信息，版本号，唯一ID等信息  Get the current robot's system information, version number, unique ID and other information |
| 公共方法 | [GetSystemAsync](7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm) | 获取当前机器人的系统信息，版本号，唯一ID等信息  Get the current robot's system information, version number, unique ID and other information |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [GetUserValue](5470b496-462f-c3da-153f-dac92a70066c.htm) | 根据给定的名称，获取当前用户的数据值信息。  According to the given name, gets the data value information of the current user |
| 公共方法 | [GetUserValueAsync](d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm) | 根据给定的名称，获取当前用户的数据值信息。  According to the given name, gets the data value information of the current user |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Post](8be86099-0853-7a12-e310-23383ec0bc7e.htm) | 使用POST命令去提交数据内容，然后返回相关的数据信息 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [PostAsync](e0b8ac89-a0ff-4d7b-1e76-de08e7d31de1.htm) | 使用POST命令去提交数据内容，然后返回相关的数据信息 (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [Read](bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm) | 读取对方信息的的数据信息，通常是针对GET的方法信息设计的。如果使用了url=开头，就表示是使用了原生的地址访问  Read the other side of the data information, usually designed for the GET method information.If you start with url=, you are using native address access (重写 [NetworkWebApiRobotBaseRead(String)](e2e4dd7f-efbe-ae6b-abb2-e3bec451182d.htm).) |
| 公共方法 | [ReadAsync](d8667131-2445-aa8f-7c68-865fdc537237.htm) | 读取对方信息的的数据信息，通常是针对GET的方法信息设计的。如果使用了url=开头，就表示是使用了原生的地址访问  Read the other side of the data information, usually designed for the GET method information.If you start with url=, you are using native address access (继承自 [NetworkWebApiRobotBase](3303830b-83fb-9aba-9684-41bba031dd90.htm)。) |
| 受保护的方法 | [ReadByAddress](4ce59148-7bef-049e-a0ca-d190467b1938.htm) | 等待重写的额外的指令信息的支持。除了url的形式之外，还支持基于命令的数据交互  Additional instruction information waiting for rewriting is supported.In addition to the url format, command based data interaction is supported (重写 [NetworkWebApiRobotBaseReadByAddress(String)](106dd981-59b1-b17d-a183-0ef15117631f.htm).) |
| 受保护的方法 | [ReadByAddressAsync](be2f191d-9781-c9a0-5800-ed8c52088aae.htm) | 等待重写的额外的指令信息的支持。除了url的形式之外，还支持基于命令的数据交互  Additional instruction information waiting for rewriting is supported.In addition to the url format, command based data interaction is supported (重写 [NetworkWebApiRobotBaseReadByAddressAsync(String)](e1e23509-a944-bf67-c503-057ddf581cac.htm).) |
| 公共方法 | [ReadString](bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm) | 读取对方信息的的字符串数据信息，通常是针对GET的方法信息设计的。如果使用了url=开头，就表示是使用了原生的地址访问  The string data information that reads the other party information, usually designed for the GET method information.If you start with url=, you are using native address access (重写 [NetworkWebApiRobotBaseReadString(String)](c4dc02be-9ee9-ce47-0b44-b41e08635fb3.htm).) |
| 公共方法 | [ReadStringAsync](03cc49bd-e193-7225-6447-2e62b69ddec1.htm) | 读取对方信息的的字符串数据信息，通常是针对GET的方法信息设计的。如果使用了url=开头，就表示是使用了原生的地址访问  The string data information that reads the other party information, usually designed for the GET method information.If you start with url=, you are using native address access (继承自 [NetworkWebApiRobotBase](3303830b-83fb-9aba-9684-41bba031dd90.htm)。) |
| 公共方法 | [RequestAsync](649024d2-2df8-56e7-14ad-9b884701fab3.htm) | 发起一个请求，然后返回结果，支持任意的http方法  Initiate a request and then return the result, supporting any http method (继承自 [NetworkWebApiBase](b145c7c9-1925-8e28-2bc0-d8604384b4ae.htm)。) |
| 公共方法 | [ToString](969cd64f-a556-542a-b5d0-6e82066d44ff.htm) | (重写 [NetworkWebApiRobotBaseToString](bc8a32fc-dd82-78db-875f-de4b14a23dd3.htm).) |
| 公共方法 | [Write(String, Byte)](36744e4b-b339-1a38-ac3c-47ef3dc3b71e.htm) | 使用POST的方式来向对方进行请求数据信息，需要使用url=开头，来表示是使用了原生的地址访问  Using POST to request data information from the other party, we need to start with url= to indicate that we are using native address access (重写 [NetworkWebApiRobotBaseWrite(String, Byte)](74d6da2a-aa64-98ad-3d8c-0ea865a6755f.htm).) |
| 公共方法 | [Write(String, String)](b227fd97-6aeb-90f2-7271-fd676a3e2e54.htm) | 使用POST的方式来向对方进行请求数据信息，需要使用url=开头，来表示是使用了原生的地址访问  Using POST to request data information from the other party, we need to start with url= to indicate that we are using native address access (重写 [NetworkWebApiRobotBaseWrite(String, String)](5654dcab-d566-4511-2024-bc51004c9dad.htm).) |
| 公共方法 | [WriteAsync(String, Byte)](e1222bd2-37fe-3ad9-f20a-15c324142e96.htm) | 使用POST的方式来向对方进行请求数据信息，需要使用url=开头，来表示是使用了原生的地址访问  Using POST to request data information from the other party, we need to start with url= to indicate that we are using native address access (继承自 [NetworkWebApiRobotBase](3303830b-83fb-9aba-9684-41bba031dd90.htm)。) |
| 公共方法 | [WriteAsync(String, String)](3e81cc81-a117-ac9e-8c1a-31b021ab3d3d.htm) | 使用POST的方式来向对方进行请求数据信息，需要使用url=开头，来表示是使用了原生的地址访问  Using POST to request data information from the other party, we need to start with url= to indicate that we are using native address access (继承自 [NetworkWebApiRobotBase](3303830b-83fb-9aba-9684-41bba031dd90.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetAnIOSignal 方法 

[原文連結](http://api.hslcommunication.cn/html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetAnIOSignal 方法 |

获取机器人的IO信号资源  
Get an IO signal resource.

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetAnIOSignal(
	string network = "Local",
	string unit = "DRV_1",
	string signal = "DRV1K1"
)
```

```
Public Function GetAnIOSignal ( 
	Optional network As String = "Local",
	Optional unit As String = "DRV_1",
	Optional signal As String = "DRV1K1"
) As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetAnIOSignal(
	String^ network = L"Local", 
	String^ unit = L"DRV_1", 
	String^ signal = L"DRV1K1"
)
```

```
member GetAnIOSignal : 
        ?network : string * 
        ?unit : string * 
        ?signal : string 
(* Defaults:
        let _network = defaultArg network "Local"
        let _unit = defaultArg unit "DRV_1"
        let _signal = defaultArg signal "DRV1K1"
*)
-> OperateResult<string> 
```

#### 参数

network (Optional)
:   类型：SystemString  

    [缺少 "M:HslCommunication.Robot.ABB.ABBWebApiClient.GetAnIOSignal(System.String,System.String,System.String)" 的 <param name="network"/> 文档]

unit (Optional)
:   类型：SystemString  

    [缺少 "M:HslCommunication.Robot.ABB.ABBWebApiClient.GetAnIOSignal(System.String,System.String,System.String)" 的 <param name="unit"/> 文档]

signal (Optional)
:   类型：SystemString  

    [缺少 "M:HslCommunication.Robot.ABB.ABBWebApiClient.GetAnIOSignal(System.String,System.String,System.String)" 的 <param name="signal"/> 文档]

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
系统的基本信息

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetAnIOSignalAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/0e63db29-1826-248c-2c6f-13e158318127.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetAnIOSignalAsync 方法 |

获取机器人的IO信号资源  
Get an IO signal resource.

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetAnIOSignalAsync(
	string network = "Local",
	string unit = "DRV_1",
	string signal = "DRV1K1"
)
```

```
Public Function GetAnIOSignalAsync ( 
	Optional network As String = "Local",
	Optional unit As String = "DRV_1",
	Optional signal As String = "DRV1K1"
) As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetAnIOSignalAsync(
	String^ network = L"Local", 
	String^ unit = L"DRV_1", 
	String^ signal = L"DRV1K1"
)
```

```
member GetAnIOSignalAsync : 
        ?network : string * 
        ?unit : string * 
        ?signal : string 
(* Defaults:
        let _network = defaultArg network "Local"
        let _unit = defaultArg unit "DRV_1"
        let _signal = defaultArg signal "DRV1K1"
*)
-> Task<OperateResult<string>> 
```

#### 参数

network (Optional)
:   类型：SystemString  

    [缺少 "M:HslCommunication.Robot.ABB.ABBWebApiClient.GetAnIOSignalAsync(System.String,System.String,System.String)" 的 <param name="network"/> 文档]

unit (Optional)
:   类型：SystemString  

    [缺少 "M:HslCommunication.Robot.ABB.ABBWebApiClient.GetAnIOSignalAsync(System.String,System.String,System.String)" 的 <param name="unit"/> 文档]

signal (Optional)
:   类型：SystemString  

    [缺少 "M:HslCommunication.Robot.ABB.ABBWebApiClient.GetAnIOSignalAsync(System.String,System.String,System.String)" 的 <param name="signal"/> 文档]

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
系统的基本信息

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetCtrlState 方法 

[原文連結](http://api.hslcommunication.cn/html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetCtrlState 方法 |

获取当前的控制状态，Content属性就是机器人的控制信息  
Get the current control state. The Content attribute is the control information of the robot

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetCtrlState()
```

```
Public Function GetCtrlState As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetCtrlState()
```

```
member GetCtrlState : unit -> OperateResult<string> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有状态信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetCtrlStateAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetCtrlStateAsync 方法 |

获取当前的控制状态，Content属性就是机器人的控制信息  
Get the current control state. The Content attribute is the control information of the robot

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetCtrlStateAsync()
```

```
Public Function GetCtrlStateAsync As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetCtrlStateAsync()
```

```
member GetCtrlStateAsync : unit -> Task<OperateResult<string>> 
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有状态信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorState 方法 

[原文連結](http://api.hslcommunication.cn/html/858227c4-967f-fa0b-4b77-7558db99421c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetErrorState 方法 |

获取当前的错误状态，Content属性就是机器人的状态信息  
Gets the current error state. The Content attribute is the state information of the robot

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetErrorState()
```

```
Public Function GetErrorState As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetErrorState()
```

```
member GetErrorState : unit -> OperateResult<string> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有状态信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorStateAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/9ea55110-7354-4743-45d7-677395d3b9da.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetErrorStateAsync 方法 |

获取当前的错误状态，Content属性就是机器人的状态信息  
Gets the current error state. The Content attribute is the state information of the robot

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetErrorStateAsync()
```

```
Public Function GetErrorStateAsync As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetErrorStateAsync()
```

```
member GetErrorStateAsync : unit -> Task<OperateResult<string>> 
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有状态信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetIO2In 方法 

[原文連結](http://api.hslcommunication.cn/html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetIO2In 方法 |

获取当前机器人的本机的输入IO  
Gets the input IO of the current robot's native

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetIO2In()
```

```
Public Function GetIO2In As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetIO2In()
```

```
member GetIO2In : unit -> OperateResult<string> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetIO2InAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetIO2InAsync 方法 |

获取当前机器人的本机的输入IO  
Gets the input IO of the current robot's native

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetIO2InAsync()
```

```
Public Function GetIO2InAsync As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetIO2InAsync()
```

```
member GetIO2InAsync : unit -> Task<OperateResult<string>> 
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetIO2Out 方法 

[原文連結](http://api.hslcommunication.cn/html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetIO2Out 方法 |

获取当前机器人的本机的输出IO  
Gets the output IO of the current robot's native

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetIO2Out()
```

```
Public Function GetIO2Out As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetIO2Out()
```

```
member GetIO2Out : unit -> OperateResult<string> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetIO2OutAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetIO2OutAsync 方法 |

获取当前机器人的本机的输出IO  
Gets the output IO of the current robot's native

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetIO2OutAsync()
```

```
Public Function GetIO2OutAsync As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetIO2OutAsync()
```

```
member GetIO2OutAsync : unit -> Task<OperateResult<string>> 
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetIOIn 方法 

[原文連結](http://api.hslcommunication.cn/html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetIOIn 方法 |

获取当前机器人的本机的输入IO  
Gets the input IO of the current robot's native

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetIOIn()
```

```
Public Function GetIOIn As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetIOIn()
```

```
member GetIOIn : unit -> OperateResult<string> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetIOInAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetIOInAsync 方法 |

获取当前机器人的本机的输入IO  
Gets the input IO of the current robot's native

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetIOInAsync()
```

```
Public Function GetIOInAsync As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetIOInAsync()
```

```
member GetIOInAsync : unit -> Task<OperateResult<string>> 
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetIOOut 方法 

[原文連結](http://api.hslcommunication.cn/html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetIOOut 方法 |

获取当前机器人的本机的输出IO  
Gets the output IO of the current robot's native

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetIOOut()
```

```
Public Function GetIOOut As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetIOOut()
```

```
member GetIOOut : unit -> OperateResult<string> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetIOOutAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetIOOutAsync 方法 |

获取当前机器人的本机的输出IO  
Gets the output IO of the current robot's native

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetIOOutAsync()
```

```
Public Function GetIOOutAsync As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetIOOutAsync()
```

```
member GetIOOutAsync : unit -> Task<OperateResult<string>> 
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetJointTarget 方法 

[原文連結](http://api.hslcommunication.cn/html/8d619075-fc31-7a20-712e-cb9483e48b09.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetJointTarget 方法 |

获取当前机器人的物理关节点信息，返回json格式的关节信息  
Get the physical node information of the current robot and return the joint information in json format

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetJointTarget(
	string mechunit = "ROB_1"
)
```

```
Public Function GetJointTarget ( 
	Optional mechunit As String = "ROB_1"
) As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetJointTarget(
	String^ mechunit = L"ROB_1"
)
```

```
member GetJointTarget : 
        ?mechunit : string 
(* Defaults:
        let _mechunit = defaultArg mechunit "ROB_1"
*)
-> OperateResult<string> 
```

#### 参数

mechunit (Optional)
:   类型：SystemString  
    操作单元，默认为 ROB\_1

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有关节信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetJointTargetAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetJointTargetAsync 方法 |

获取当前机器人的物理关节点信息，返回json格式的关节信息  
Get the physical node information of the current robot and return the joint information in json format

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetJointTargetAsync(
	string mechunit = "ROB_1"
)
```

```
Public Function GetJointTargetAsync ( 
	Optional mechunit As String = "ROB_1"
) As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetJointTargetAsync(
	String^ mechunit = L"ROB_1"
)
```

```
member GetJointTargetAsync : 
        ?mechunit : string 
(* Defaults:
        let _mechunit = defaultArg mechunit "ROB_1"
*)
-> Task<OperateResult<string>> 
```

#### 参数

mechunit (Optional)
:   类型：SystemString  
    操作单元，默认为 ROB\_1

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有关节信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetLog 方法 

[原文連結](http://api.hslcommunication.cn/html/16aa6e7b-178b-51ee-8039-400d94447c39.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetLog 方法 |

获取当前机器人的日志记录，默认记录为10条  
Gets the log record for the current robot, which is 10 by default

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetLog(
	int logCount = 10
)
```

```
Public Function GetLog ( 
	Optional logCount As Integer = 10
) As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetLog(
	int logCount = 10
)
```

```
member GetLog : 
        ?logCount : int 
(* Defaults:
        let _logCount = defaultArg logCount 10
*)
-> OperateResult<string> 
```

#### 参数

logCount (Optional)
:   类型：SystemInt32  
    读取的最大的日志总数

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetLogAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetLogAsync 方法 |

获取当前机器人的日志记录，默认记录为10条  
Gets the log record for the current robot, which is 10 by default

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetLogAsync(
	int logCount = 10
)
```

```
Public Function GetLogAsync ( 
	Optional logCount As Integer = 10
) As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetLogAsync(
	int logCount = 10
)
```

```
member GetLogAsync : 
        ?logCount : int 
(* Defaults:
        let _logCount = defaultArg logCount 10
*)
-> Task<OperateResult<string>> 
```

#### 参数

logCount (Optional)
:   类型：SystemInt32  
    读取的最大的日志总数

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有IO信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetOperationMode 方法 

[原文連結](http://api.hslcommunication.cn/html/080723d7-a5ca-c1db-c967-22871a231339.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetOperationMode 方法 |

获取当前机器人的工作模式  
Gets the current working mode of the robot

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetOperationMode()
```

```
Public Function GetOperationMode As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetOperationMode()
```

```
member GetOperationMode : unit -> OperateResult<string> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有工作模式信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetOperationModeAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetOperationModeAsync 方法 |

获取当前机器人的工作模式  
Gets the current working mode of the robot

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetOperationModeAsync()
```

```
Public Function GetOperationModeAsync As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetOperationModeAsync()
```

```
member GetOperationModeAsync : unit -> Task<OperateResult<string>> 
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有工作模式信息的结果类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetRapidExecution 方法 

[原文連結](http://api.hslcommunication.cn/html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetRapidExecution 方法 |

获取当前机器人的当前程序运行状态  
Get the current program running status of the current robot

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> GetRapidExecution()
```

```
Public Function GetRapidExecution As OperateResult(Of String)
```

```
public:
OperateResult<String^>^ GetRapidExecution()
```

```
member GetRapidExecution : unit -> OperateResult<string> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
机器人的当前的程序运行状态

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetRapidExecutionAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.ABB](../html/54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm "HslCommunication.Robot.ABB")

[ABBWebApiClient 类](../html/75e13fcb-0640-7483-7582-6e0e6bc9e864.htm "ABBWebApiClient 类")

[ABBWebApiClient 方法](../html/5dc4f38c-b444-74f4-c27e-51bd68faa824.htm "ABBWebApiClient 方法")

[GetAnIOSignal 方法](../html/d7d77ccd-7e59-7b3b-83de-90fcd80d273a.htm "GetAnIOSignal 方法 ")

[GetAnIOSignalAsync 方法](../html/0e63db29-1826-248c-2c6f-13e158318127.htm "GetAnIOSignalAsync 方法 ")

[GetCtrlState 方法](../html/57e5860a-18bc-4f4a-06b7-0c17555f460d.htm "GetCtrlState 方法 ")

[GetCtrlStateAsync 方法](../html/996af0a7-6e58-911a-541b-6dfe417d7bba.htm "GetCtrlStateAsync 方法 ")

[GetErrorState 方法](../html/858227c4-967f-fa0b-4b77-7558db99421c.htm "GetErrorState 方法 ")

[GetErrorStateAsync 方法](../html/9ea55110-7354-4743-45d7-677395d3b9da.htm "GetErrorStateAsync 方法 ")

[GetIO2In 方法](../html/fc99d52e-f3e1-563f-94c6-dc144f30a6d1.htm "GetIO2In 方法 ")

[GetIO2InAsync 方法](../html/7cd45d57-f74f-5e98-28c0-b7375bd87eb0.htm "GetIO2InAsync 方法 ")

[GetIO2Out 方法](../html/fd58c42c-1363-fbb3-49ec-dcbd38cd7198.htm "GetIO2Out 方法 ")

[GetIO2OutAsync 方法](../html/8c3a0765-19ba-0676-dc85-22cdbf22498b.htm "GetIO2OutAsync 方法 ")

[GetIOIn 方法](../html/51dbd9f4-0ccf-46e9-41b6-4f304bce5b66.htm "GetIOIn 方法 ")

[GetIOInAsync 方法](../html/da8f283b-b7ca-9388-cd0c-0255aad3fc70.htm "GetIOInAsync 方法 ")

[GetIOOut 方法](../html/ee3b44e0-a853-ccdd-e475-98552d0966a8.htm "GetIOOut 方法 ")

[GetIOOutAsync 方法](../html/7d3bbbe6-67e4-49ac-7be5-7e4d46886e61.htm "GetIOOutAsync 方法 ")

[GetJointTarget 方法](../html/8d619075-fc31-7a20-712e-cb9483e48b09.htm "GetJointTarget 方法 ")

[GetJointTargetAsync 方法](../html/2fdbc6cd-c54a-c74b-1bb9-159405de748d.htm "GetJointTargetAsync 方法 ")

[GetLog 方法](../html/16aa6e7b-178b-51ee-8039-400d94447c39.htm "GetLog 方法 ")

[GetLogAsync 方法](../html/40ca6ebf-c341-1bdc-1dde-3e0cd94da553.htm "GetLogAsync 方法 ")

[GetOperationMode 方法](../html/080723d7-a5ca-c1db-c967-22871a231339.htm "GetOperationMode 方法 ")

[GetOperationModeAsync 方法](../html/4c37c64d-584f-02b0-912c-d78cf7764b29.htm "GetOperationModeAsync 方法 ")

[GetRapidExecution 方法](../html/ee99bc27-e0b3-2b92-c6a8-9244ffbc6e66.htm "GetRapidExecution 方法 ")

[GetRapidExecutionAsync 方法](../html/e2bd208a-f150-9325-e85e-fa18d70d8f3a.htm "GetRapidExecutionAsync 方法 ")

[GetRapidTasks 方法](../html/4fbe4e80-94f1-45ba-6007-92328676e6ea.htm "GetRapidTasks 方法 ")

[GetRapidTasksAsync 方法](../html/bc1ec905-e4d0-5f2f-6a5d-a7f0a3eb8b3d.htm "GetRapidTasksAsync 方法 ")

[GetRobotTarget 方法](../html/78b1b85c-db0f-a511-7e12-746069550e04.htm "GetRobotTarget 方法 ")

[GetRobotTargetAsync 方法](../html/17ecddf9-5371-d3f4-a0a2-b08e2cb8552d.htm "GetRobotTargetAsync 方法 ")

[GetSelectStrings 方法](../html/ed5f2af0-d225-1f65-7b00-c0634baf91e3.htm "GetSelectStrings 方法 ")

[GetServoEnable 方法](../html/ad00798e-b449-5ac9-1902-e2b8c987003d.htm "GetServoEnable 方法 ")

[GetServoEnableAsync 方法](../html/7d6b3944-2a09-4d89-c59d-620417eb4f2c.htm "GetServoEnableAsync 方法 ")

[GetSpeedRatio 方法](../html/51e3d72d-82b8-cbec-64a2-0d2273b34b6b.htm "GetSpeedRatio 方法 ")

[GetSpeedRatioAsync 方法](../html/e8eb9e5e-48d2-c000-cd15-df020d88d2d9.htm "GetSpeedRatioAsync 方法 ")

[GetSystem 方法](../html/852f33fa-83bf-42e7-9e08-33fc34a82600.htm "GetSystem 方法 ")

[GetSystemAsync 方法](../html/7a8c56ac-18f7-7a9f-c22b-4862dd90a46a.htm "GetSystemAsync 方法 ")

[GetUserValue 方法](../html/5470b496-462f-c3da-153f-dac92a70066c.htm "GetUserValue 方法 ")

[GetUserValueAsync 方法](../html/d01798ca-8654-bf6c-9cff-62d0a20ccd22.htm "GetUserValueAsync 方法 ")

[Read 方法](../html/bc860f72-862b-3e4b-2cef-be33aeaf17ed.htm "Read 方法 ")

[ReadByAddress 方法](../html/4ce59148-7bef-049e-a0ca-d190467b1938.htm "ReadByAddress 方法 ")

[ReadByAddressAsync 方法](../html/be2f191d-9781-c9a0-5800-ed8c52088aae.htm "ReadByAddressAsync 方法 ")

[ReadString 方法](../html/bcf8c72a-0661-6c30-9ff1-dfca3d93f342.htm "ReadString 方法 ")

[ToString 方法](../html/969cd64f-a556-542a-b5d0-6e82066d44ff.htm "ToString 方法 ")

[Write 方法](../html/e57de68b-910a-4f48-3d36-99cae73beb15.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ABBWebApiClientGetRapidExecutionAsync 方法 |

获取当前机器人的当前程序运行状态  
Get the current program running status of the current robot

**命名空间：**
 [HslCommunication.Robot.ABB](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> GetRapidExecutionAsync()
```

```
Public Function GetRapidExecutionAsync As Task(Of OperateResult(Of String))
```

```
public:
Task<OperateResult<String^>^>^ GetRapidExecutionAsync()
```

```
member GetRapidExecutionAsync : unit -> Task<OperateResult<string>> 
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
机器人的当前的程序运行状态

![](../icons/SectionExpanded.png)参见

#### 引用

[ABBWebApiClient 类](75e13fcb-0640-7483-7582-6e0e6bc9e864.htm)

[HslCommunication.Robot.ABB 命名空间](54aafae2-8262-7ae4-8ec9-199f5f8d33a2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)