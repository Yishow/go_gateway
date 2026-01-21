# HslCommunication - HslCommunication

> 分類頁數: 30



---
## HslCommunication

[原文連結](http://api.hslcommunication.cn/html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[Authorization 类](../html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm "Authorization 类")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslTimeOut 类](../html/fdde4f7d-b56a-e0f8-bbc6-67cc153c82a2.htm "HslTimeOut 类")

[IDataTransfer 接口](../html/90fc4cd3-fb3a-043c-2bd6-109124d6a78c.htm "IDataTransfer 接口")

[NetHandle 结构](../html/a27405cc-a0f8-3902-cb23-7a9c2adfcb65.htm "NetHandle 结构")

[OperateResult 类](../html/a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm "OperateResult 类")

[OperateResult(T) 类](../html/cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm "OperateResult(T) 类")

[OperateResult(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) 类](../html/8b225508-4c86-ee34-f1c2-a3b6b7c55e76.htm "OperateResult(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) 类")

[OperateResult(T1, T2) 类](../html/f52f888f-5e8d-b0c4-2302-81e70c230a31.htm "OperateResult(T1, T2) 类")

[OperateResult(T1, T2, T3) 类](../html/af8952ca-51d9-e98f-6bf7-332d27b0f0be.htm "OperateResult(T1, T2, T3) 类")

[OperateResult(T1, T2, T3, T4) 类](../html/7b14ef70-a90d-6a4a-0671-cc70054d8024.htm "OperateResult(T1, T2, T3, T4) 类")

[OperateResult(T1, T2, T3, T4, T5) 类](../html/236a3482-734e-f63b-d8f0-0b35f1b6a2f5.htm "OperateResult(T1, T2, T3, T4, T5) 类")

[OperateResult(T1, T2, T3, T4, T5, T6) 类](../html/3fed2115-ad0f-68ac-1232-81de620ee7dc.htm "OperateResult(T1, T2, T3, T4, T5, T6) 类")

[OperateResult(T1, T2, T3, T4, T5, T6, T7) 类](../html/7bc1e112-5372-a828-f613-108aa110327d.htm "OperateResult(T1, T2, T3, T4, T5, T6, T7) 类")

[OperateResult(T1, T2, T3, T4, T5, T6, T7, T8) 类](../html/5832c4ca-ad2a-5a63-4040-4490c3a5d918.htm "OperateResult(T1, T2, T3, T4, T5, T6, T7, T8) 类")

[OperateResult(T1, T2, T3, T4, T5, T6, T7, T8, T9) 类](../html/c6b79f48-59c1-7cfb-c11e-52f534e4d830.htm "OperateResult(T1, T2, T3, T4, T5, T6, T7, T8, T9) 类")

[StringResources 类](../html/317bf433-cd99-4760-9219-d9962ba25325.htm "StringResources 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication 命名空间 |

一个工业物联网的底层架构框架，专注于底层的技术通信及跨平台，跨语言通信功能，实现各种主流的PLC数据读写，实现modbus，机器人的各种协议读写等等，
支持快速搭建工业上位机软件，组态软件，SCADA软件，工厂MES系统，助力企业工业4.0腾飞，实现智能制造，智慧工厂的目标。
  
  
本组件付费开源，使用之前请认真的阅读本API文档，对于本文档中警告部分的内容务必理解，部署生产之前请详细测试，如果在测试的过程中，
发现了BUG，或是有问题的地方，欢迎联系作者进行修改，或是直接在github上进行提问。未经测试，直接部署，对设备，工厂造成了损失，作者概不负责。
  
  
官方网站：<http://www.hslcommunication.cn/>，包含组件的在线API地址以及一个MES DEMO的项目展示。
  
  

| 重要事项 **重要事项** |
| --- |
| 本组件的目标是集成一个框架，统一所有的设备读写方法，抽象成统一的接口[IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)，对于上层操作只需要关注地址，读取类型即可，另一个目标是使用本框架轻松实现C#后台+C#客户端+web浏览器+android手机的全方位功能实现。 |

本库提供了C#版本和java版本和python版本，java，python版本的使用和C#几乎是一模一样的，都是可以相互通讯的。
  
在使用本通讯库之前，需要学会如何使用nuget来安装当前的通讯库，可以参考如下的博文：<http://www.cnblogs.com/dathlin/p/7705014.html>  
  
先整体介绍下如何使用本组件库的基本思路，基本上是引用库，从nuget安装的库会自动添加引用到项目中的，就可以直接进行using操作了，当然了在使用之前，需要先激活一下，激活的方式如下：

激活示例

[复制](# "复制")

```
// 授权示例 Authorization example
if (!HslCommunication.Authorization.SetAuthorizationCode( "Your Code" ))
{
    Console.WriteLine( "Authorization failed! The current program can only be used for 8 hours!" );
    return;   // 激活失败应该退出系统
}
```

在你的应用程序刚开起来的时候，激活一次即可，后续都不需要再重复激活了。接下来就可以开始写代码了。任何的设备的操作基本是相同的，实例化，配置参数（有些plc默认的参数即可），连接设备，读写操作，关闭
  
关于Hsl的日志功能，贯穿整个hslcommunication的项目，所有的网络类，都包含了[ILogNet](d95d4704-db0d-b379-41e6-03879927a543.htm)日志功能，当然你也可以继承接口实现你自己的日志，在hsl里提供了三种常见的简单实用的日志

单文件实例化

[复制](# "复制")

```
// 实例化一个日志后，就可以使用了
ILogNet logNet = new LogNetSingle( "D:\\123.txt" );
```

日志基本的使用

[复制](# "复制")

```
// 然后我们的代码就可以调用下面的方法来存储日志了，支持下面的5个等级的
logNet.WriteDebug( "Debug log test" );
logNet.WriteInfo( "Info log test" );
logNet.WriteWarn( "Warn log test" );
logNet.WriteError( "Error log test" );
logNet.WriteFatal( "Fatal log test" );

// 还有下面的几种额外的情况
logNet.WriteNewLine( );                       // 追加一行空行
logNet.WriteDescrition( "test" );             // 写入额外的注释的信息
logNet.WriteAnyString( "any string" );        // 写任意的数据，不受格式化影响

// 此处的5个等级有高低之分 debug < info < warn < error < fatal < all
// 如果我们需要屏蔽debug等级的话
logNet.SetMessageDegree( HslMessageDegree.INFO );

// 如果我们需要屏蔽debug及info等级的
logNet.SetMessageDegree( HslMessageDegree.WARN );

// 如果所有的日志在记录之前需要在控制台显示出来
logNet.BeforeSaveToFile += (object sender, HslEventArgs e) =>
{
    Console.WriteLine( e.HslMessage.ToString( ) );
};

// 带关键字的功能
logNet.WriteDebug( "A","Debug log test" );
logNet.WriteInfo( "B", "Info log test" );
logNet.WriteWarn( "C", "Warn log test" );
logNet.WriteError( "A", "Error log test" );
logNet.WriteFatal( "B", "Fatal log test" );

// 有了关键字之后，我们就可以根据关键字过滤了
logNet.FiltrateKeyword( "B" ); // 我们不需要B的关键字的日志
logNet.RemoveFiltrate( "B" );  // 重新需要B的关键字的日志

// 日志记录时默认获取当前系统的时间，如果需要在记录日志时，所有的时间往后加8小时，则可以这么设置（如果往前，就是负数）
logNet.HourDeviation = 8;
```

  
再开始讲解基本的代码通讯之前，先来了解两个基本的概念，长连接，短连接。为了更好的说明当前的通信情况，我把所有的通信拆分为四个部分，连接，发，收，断开。  
短连接：连接，发，收，断开，连接，发，收，断开，连接，发，收，断开，连接，发，收，断开...无限循环  
长连接：连接，发，收，发，收，发，收，发，收，发，收，发，收，发，收，发，收，发，收....断开  
然后我们来看看异常的情况，短连接的异常比较好处理，反正每次请求都是先连接，关键来看长连接的异常  
长连接：连接，发，收，发，收...异常，连接，发，收，发，收，异常，连接，连接，连接...收，发，收，发，收，发，收，发，收....断开  
这里第一个异常发生后，第二次读写立即连接上去并且成功，第二个异常触发后，一直读写失败，说明就是一直连接不上去。  
对于HSL组件来说，不需要重复连接服务器或是plc，无论是短连接还是长连接，都只需要一直读写就OK了，对读写的结果进行判定，即使发生异常了，读写失败了，也要一直坚持，网络好的时候，读写会恢复成功的。  
  
我们以三菱的PLC为示例，其他的plc调用方式基本是一模一样的，就是调用的类不一样，参数配置不一样而已。以下的逻辑都是适用的。

简单的短连接使用

[复制](# "复制")

```
// 实例化对象，指定PLC的ip地址和端口号
MelsecMcNet melsecMc = new MelsecMcNet( "192.168.1.110", 6000 );
// 举例读取D100的值
short D100 = melsecMc.ReadInt16( "D100" ).Content;
```

简单的长连接使用

[复制](# "复制")

```
// 实例化对象，指定PLC的ip地址和端口号
MelsecMcNet melsecMc = new MelsecMcNet( "192.168.1.110", 6000 );

// 如果网络不太理想，配置了两个端口，一个有问题，立即切换另一个的话，可以配置如下的代码
// melsecMc.GetPipeSocket( ).SetMultiPorts( new int[] { 6000, 6001 } );

// 连接对象
OperateResult connect = melsecMc.ConnectServer( );
if (!connect.IsSuccess)
{
    Console.WriteLine( "connect failed:" + connect.Message );
    return;
}

// 举例读取D100的值
short D100 = melsecMc.ReadInt16( "D100" ).Content;

// 实际上所有的读写都是返回是否成功的标记的，在实际的开发中，需要严格的判定，怎么判定呢？如下的代码
// In fact, all reads and writes return the flags of success. In actual development, strict judgment is required. How to judge? The following code
OperateResult<short> readD100 = melsecMc.ReadInt16( "D100" );
if (readD100.IsSuccess)
{
    // 读取成功，这时候获取Content才是正确的值
    // The reading is successful, and the Content is the correct value at this time.
    short value = readD100.Content;
}
else
{
    // 读取失败，如果仍然坚持去获取Content的值，就为0
    // Failed to read, if it still persists to get the value of Content, it is 0
}

// 读写是否成功的情况，应用于几乎所有的读写代码，只要看清楚返回的数据类型即可
// Whether the reading and writing is successful is applied to almost all reading and writing codes, as long as the type of data returned is clear

melsecMc.ConnectClose( );
```

基本的读取示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );
// melsec_net.CommunicationPipe.IsPersistentConnection = false;   // 如果想设置短连接的话

// 此处以D寄存器作为示例，此处没有判断是否读写成功，所以是有风险的，如果通讯失败，所有的值都不是正确的
// The D register is used as an example here. There is no judgment whether the read or write is successful, so it is risky. 
// If the communication fails, all values are incorrect
short short_D1000 = melsec_net.ReadInt16( "D1000" ).Content;                     // 读取D1000的short值 
ushort ushort_D1000 = melsec_net.ReadUInt16( "D1000" ).Content;                  // 读取D1000的ushort值
int int_D1000 = melsec_net.ReadInt32( "D1000" ).Content;                         // 读取D1000-D1001组成的int数据
uint uint_D1000 = melsec_net.ReadUInt32( "D1000" ).Content;                      // 读取D1000-D1001组成的uint数据
float float_D1000 = melsec_net.ReadFloat( "D1000" ).Content;                     // 读取D1000-D1001组成的float数据
long long_D1000 = melsec_net.ReadInt64( "D1000" ).Content;                       // 读取D1000-D1003组成的long数据
ulong ulong_D1000 = melsec_net.ReadUInt64( "D1000" ).Content;                    // 读取D1000-D1003组成的long数据
double double_D1000 = melsec_net.ReadDouble( "D1000" ).Content;                  // 读取D1000-D1003组成的double数据
string str_D1000 = melsec_net.ReadString( "D1000", 10 ).Content;                 // 读取D1000-D1009组成的条码数据

// 读取数组
short[] short_D1000_array = melsec_net.ReadInt16( "D1000", 10 ).Content;         // 读取D1000的short值 
ushort[] ushort_D1000_array = melsec_net.ReadUInt16( "D1000", 10 ).Content;      // 读取D1000的ushort值
int[] int_D1000_array = melsec_net.ReadInt32( "D1000", 10 ).Content;             // 读取D1000-D1001组成的int数据
uint[] uint_D1000_array = melsec_net.ReadUInt32( "D1000", 10 ).Content;          // 读取D1000-D1001组成的uint数据
float[] float_D1000_array = melsec_net.ReadFloat( "D1000", 10 ).Content;         // 读取D1000-D1001组成的float数据
long[] long_D1000_array = melsec_net.ReadInt64( "D1000", 10 ).Content;           // 读取D1000-D1003组成的long数据
ulong[] ulong_D1000_array = melsec_net.ReadUInt64( "D1000", 10 ).Content;        // 读取D1000-D1003组成的long数据
double[] double_D1000_array = melsec_net.ReadDouble( "D1000", 10 ).Content;      // 读取D1000-D1003组成的double数据

// 正常的安全的操作如下，每一行都需要转变下面的代码
// The normal and safe operation is as follows, each line needs to change the following code
OperateResult<short> resultD1000 = melsec_net.ReadInt16( "D1000" );
if (resultD1000.IsSuccess)
{
    Console.WriteLine( resultD1000.Content );
}
```

批量读取示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 下面演示读取连续的地址的数据，实际中连续的地址，存储的数据类型不一样，可能既有short，又有int，又有float，还可能有倍数计算的关系
// 所以此处假设，
// D100存储的short，10倍的温度值，例如值为 1234，实际表示 123.4℃
// D101存储的short，100倍的压力值，例如值为 456, 实际表示 4.56 Mpa
// D102存储的int，生产产量计数
// D104-D108 共计5个D存储了10个字符，表示条码
OperateResult<byte[]> read = melsec_net.Read( "D100", 10 );  // 读出来20个字节长度的结果
if(read.IsSuccess)
{
    float temp = melsec_net.ByteTransform.TransInt16( read.Content, 0 ) / 10f;
    float press = melsec_net.ByteTransform.TransInt16( read.Content, 2 ) / 100f;
    int count = melsec_net.ByteTransform.TransInt32( read.Content, 4 );
    string barcode = melsec_net.ByteTransform.TransString( read.Content, 8, 10, Encoding.ASCII );
    // do something
}
else
{
    // failed
}
```

需要注意的事，在实际的开发中，我们的一个窗体程序（或是控制台，原理都是一样的），会定时或是不定时的去读写PLC的操作（调用Read或是Write方法），这个本身是没有任何问题的，
但是总会有这样的需求，我们需要在界面上，或是系统里实时体现当前的PLC的在线情况，我相信不少小伙伴会有这样的问题的。所以就出现了下面的代码：

检查连接状态

[复制](# "复制")

```
// 定义在外面的
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );
System.Threading.Timer timer;                                          // 定时器

/// <summary>
/// 启动定时器，我们一般会在窗体打开的时候调用，或是系统启动的时候调用，此处示例5秒检测一次
/// </summary>
public void StartTimer( )
{
    timer = new System.Threading.Timer( ThreadCheck, null, 1000, 5000 );
}

public void ThreadCheck(object obj )
{
    OperateResult connect = melsec_net.ConnectServer( );
    if (connect.IsSuccess)
    {
        // 进行相关的操作，显示绿灯啥的
    }
    else
    {
        // 进行相关的操作，显示红灯啥的
    }
}
```

实际这种操作是非常不可取的，为什么这么说，下面说说原因：  

| 警告 **警告：** |
| --- |
| 首先说明下[ConnectServer](ea436d31-7950-42df-a9c0-3c749c47e31d.htm)方法里发生了什么？这个方法首先会关闭连接，然后重新连接，连接成功，就发送初始化指令（有些PLC就需要握手确认），初始化握手成功，才返回真正的成功！  那么这里为什么不行呢？因为Read和Write方法是有混合锁实现互斥操作的，这样的好处就是多线程调用互不影响，但是[ConnectServer](ea436d31-7950-42df-a9c0-3c749c47e31d.htm)方法，并没有互斥锁，如果调用的时候同时在读写，那就会导致异常， 那么为什么没有加互斥锁呢？因为为了实现读写方法的时候，支持自动重连操作，所以连接方法已经在互斥锁了。如果再加互斥锁，会发生死锁，所以综合考虑，就设计成了现在的样子。 |

既然上面的代码不能使用，那么怎么来看当前的连接状态呢？这里有一点需要注意，[ConnectServer](ea436d31-7950-42df-a9c0-3c749c47e31d.htm)只需要调用0次或1次即可。

检查连接状态

[复制](# "复制")

```
// 定义在外面的
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );
System.Threading.Timer timer;                                          // 定时器

/// <summary>
/// 启动定时器，我们一般会在窗体打开的时候调用，或是系统启动的时候调用，此处示例5秒检测一次
/// </summary>
public void StartTimer( )
{
    // melsec_net.CommunicationPipe.IsPersistentConnection = false; // 默认长连接的操作，如果想用短连接，就调用本行代码
    timer = new System.Threading.Timer( ThreadCheck, null, 1000, 5000 );
}

public void ThreadCheck( object obj )
{
    OperateResult<short> connect = melsec_net.ReadInt16( "D100" );
    if (connect.IsSuccess)
    {
        // 进行相关的操作，显示绿灯啥的
    }
    else
    {
        // 进行相关的操作，显示红灯啥的
    }
}
```

如果你本来就在每秒读取PLC的数据信息了，那么连检测的定时器都不用写了，你每次读取数据的时候顺便判定下，结果就出来了。  
  
其他相关的代码示例需要到各自的目录里查找，下面只列举了一些常见的代码示例

* Hsl组件日志相关示例参考[ILogNet](d95d4704-db0d-b379-41e6-03879927a543.htm)
* 三菱mc协议示例参考[MelsecMcNet](0767d258-53d6-f038-ef35-477c8877d789.htm)
* 西门子S7协议示例参考[SiemensS7Net](e255b778-7bab-195f-3f07-8f57bf5d8fb7.htm)
* 欧姆龙协议示例参考[OmronFinsNet](e8e99732-9b55-e5e6-7009-51b107a8bd8b.htm)
* 罗克韦尔协议示例参考[AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm)
* Modbus协议示例参考[ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm)
* MQTT服务器示例参考[MqttServer](41e2aff2-d1ad-d10c-bad0-61b644d686fb.htm)
* MQTT客户端示例参考[MqttClient](8626b3c4-5983-7308-94ba-65e3dd480b58.htm)
* WebSocket服务器示例参考[WebSocketServer](a46e31a5-1011-9c93-ffcf-fa8404df0a65.htm)
* WebSocket客户端示例参考[WebSocketClient](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)。
* WebApi示例参考[HttpServer](13113930-8c6c-2cd6-26e7-36a1be1546a9.htm)。
* Redis示例参考[RedisClient](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)
* 身份证阅读器串口版[SAMSerial](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)
* 身份证阅读器网口版[SAMTcpNet](e33ad4bf-ecaf-8484-5f34-e472f905402f.htm)
* 文件传送服务器[UltimateFileServer](a89d8f70-aa19-617d-3e28-c97a715d56ba.htm)

| 重要事项 **重要事项** |
| --- |
| 相关的代码示例，可以翻阅左侧的命名空间，基本是按照功能来区分的，只要点进去多看看即可 |

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [Authorization](e89d602f-541c-8b98-af2b-02ee7b19618b.htm) | 系统的基本授权类 |
| 公共类 | [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) | 扩展的辅助类方法 |
| 公共类 | [HslTimeOut](fdde4f7d-b56a-e0f8-bbc6-67cc153c82a2.htm) | 超时操作的类  a class use to indicate the time-out of the connection |
| 公共类 | [OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm) | 操作结果的类，只带有成功标志和错误信息  The class that operates the result, with only success flags and error messages |
| 公共类 | [OperateResultT](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm) | 操作结果的泛型类，允许带一个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [OperateResultT1, T2, T3, T4, T5, T6, T7, T8, T9, T10](8b225508-4c86-ee34-f1c2-a3b6b7c55e76.htm) | 操作结果的泛型类，允许带十个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [OperateResultT1, T2](f52f888f-5e8d-b0c4-2302-81e70c230a31.htm) | 操作结果的泛型类，允许带两个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [OperateResultT1, T2, T3](af8952ca-51d9-e98f-6bf7-332d27b0f0be.htm) | 操作结果的泛型类，允许带三个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [OperateResultT1, T2, T3, T4](7b14ef70-a90d-6a4a-0671-cc70054d8024.htm) | 操作结果的泛型类，允许带四个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [OperateResultT1, T2, T3, T4, T5](236a3482-734e-f63b-d8f0-0b35f1b6a2f5.htm) | 操作结果的泛型类，允许带五个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [OperateResultT1, T2, T3, T4, T5, T6](3fed2115-ad0f-68ac-1232-81de620ee7dc.htm) | 操作结果的泛型类，允许带六个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [OperateResultT1, T2, T3, T4, T5, T6, T7](7bc1e112-5372-a828-f613-108aa110327d.htm) | 操作结果的泛型类，允许带七个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [OperateResultT1, T2, T3, T4, T5, T6, T7, T8](5832c4ca-ad2a-5a63-4040-4490c3a5d918.htm) | 操作结果的泛型类，允许带八个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [OperateResultT1, T2, T3, T4, T5, T6, T7, T8, T9](c6b79f48-59c1-7cfb-c11e-52f534e4d830.htm) | 操作结果的泛型类，允许带九个用户自定义的泛型对象，推荐使用这个类 |
| 公共类 | [StringResources](317bf433-cd99-4760-9219-d9962ba25325.htm) | 系统的字符串资源及多语言管理中心  System string resource and multi-language management Center |

![](../icons/SectionExpanded.png)结构

|  | 结构 | 说明 |
| --- | --- | --- |
| 公共结构代码示例 | [NetHandle](a27405cc-a0f8-3902-cb23-7a9c2adfcb65.htm) | 用于网络传递的信息头，使用上等同于int |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口代码示例 | [IDataTransfer](90fc4cd3-fb3a-043c-2bd6-109124d6a78c.htm) | 用于PLC通讯及ModBus自定义数据类型的读写操作 |

![](../icons/SectionExpanded.png)备注

本软件著作权归Richard.Hu所有。
  
博客地址：<https://www.cnblogs.com/dathlin/p/7703805.html>  
授权付费模式：Hslcommunication企业商业授权，一次付费，终身授权，终身开放源码，终身支持后续更新，支持的设备列表请参考最新的DEMO界面，凡是DEMO列举都是统一打包授权的。不限制电脑，项目，plc，开发人员数量。
  

* 企业授权提供专业版通讯库的所有更新版的 HslCommunication 源代码。包含 .Net Java Python 三大平台。
* 企业授权支持对特殊需求而进行修改，更新源代码的服务，配合企业客户修复源代码错误的服务。
* 企业授权用户拥有对通讯库商用的权利，拥有自己修改源代码并商业使用的权利，组件版权仍归属原作者。
* 企业授权用户需要对源代码保密。禁止公开源代码，禁止对源代码的商业用途。
* 企业授权用户可以免费获得官网的 MES DEMO源代码。
* 企业商业授权 费用请联系QQ200962190咨询，公司即可拥有商用版权，支持任意的开发人员数量，项目数量，支持源代码更新，长期支持，商用软件必须冠名公司标识，官网显示合作伙伴logo。
* 支持专业的一对一培训业务，一小时1000 rmb，一天8小时为5000 rmb

个人赞助二维码：  
![](https://raw.githubusercontent.com/dathlin/HslCommunication/master/imgs/support.png)

![](../icons/SectionExpanded.png)版本历史

| 日期 | 版本 | 说明 |
| --- | --- | --- |
| 2017-10-21 | 3.7.10 | * 正式发布库到互联网上去。 |
| 2017-10-21 | 3.7.11 | * 添加xml文档 |
| 2017-10-31 | 3.7.12 | * 重新设计西门子的数据读取机制，提供一个更改类型的方法。 |
| 2017-11-06 | 3.7.13 | * 提供一个ModBus的服务端引擎。 |
| 2017-11-07 | 3.7.14 | * 紧急修复了西门子批量访问时出现的BUG。 |
| 2017-11-12 | 3.7.15 | * 完善CRC16校验码功能，完善数据库辅助类方法。 |
| 2017-11-13 | 3.7.16 | * 西门子访问类，提供一个批量bool数据写入，但该写入存在安全隐患，具体见博客。 |
| 2017-11-21 | 4.0.0 | * 与3.X版本不兼容，谨慎升级。如果要升级，主要涉及的代码包含PLC的数据访问和同步数据通信。 * 删除了2个类，OperateResultBytes和OperateResultString类，提供了更加强大方便的泛型继承类，多达10个泛型参数。地址见http://www.cnblogs.com/dathlin/p/7865682.html * 将部分类从HslCommunication命名空间下移动到HslCommunication.Core下面。 * 提供了一个通用的ModBus TCP的客户端类，方便和服务器交互。 * 完善了HslCommunication.BasicFramework.SoftBaisc下面的辅助用的静态方法，提供了一些方便的数据转化，在上面进行公开。 |
| 2017-11-24 | 4.0.1 | * 更新了三菱的读取接口，提供了一个额外的字符串表示的方式，OperateResult<byte[]> read = melsecNet.ReadFromPLC("M100", 5); * 更新了西门子的数据访问类和modbus tcp类提供双模式运行，按照之前版本的写法是默认模式，每次请求重新创建网络连接，新增模式二，在代码里先进行连接服务器方法，自动切换到模式二，每次请求都共用一个网络连接，内部已经同步处理，加速数据访问，如果访问失败，自动在下次请求是重新连接，如果调用关闭连接服务器，自动切换到模式一。 |
| 2017-11-25 | 4.0.2 | * 修复Modbus tcp批量写入寄存器时，数据解析异常的BUG。 * 三菱访问器新增长连接模式。 * 三菱访问器支持单个M写入，在数组中指定一个就行。 * 三菱访问器提供了float[]数组写入的API。 * 三菱访问器支持F报警器，B链接继电器，S步进继电器，V边沿继电器，R文件寄存器读写，不过还需要大面积测试。 * 三菱访问器的读写地址支持字符串形式传入。 * 其他的细节优化。 * 感谢 hwdq0012 网友的测试和建议。 * 感谢 吃饱睡好 好朋友的测试 |
| 2017-11-27 | 4.0.3 | * 三菱，西门子，Modbus tcp客户端内核优化重构。 * 三菱，西门子，Modbus tcp客户端提供统一的报文测试方法，该方法也是通信核心，所有API都是基于此扩展起来的。 * 三菱，西门子，Modbus tcp客户端提供了一些便捷的读写API，详细参见对应博客。 * 三菱的地址区分十进制和十六进制。 * 优化三菱的位读写操作。 |
| 2017-11-28 | 4.1.0 | * 修复西门子读取的地址偏大会出现异常的BUG。 * 完善统一了所有三菱，西门子，modbus客户端类的读写方法，已经更新到博客。 |
| 2017-12-02 | 4.1.1 | * 完善日志记录，提供关键字记录操作。 * 三菱，西门子，modbus tcp客户端提供自定义数据读写。 * modbus tcp服务端提供数据池功能，并支持数据订阅操作。 * 提供一个纵向的进度控件。 |
| 2017-12-04 | 4.1.2 | * 完善Modbus tcp服务器端的数据订阅功能。 * 进度条控件支持水平方向和垂直方向两个模式。 |
| 2017-12-05 | 4.1.3 | * 进度条控件修复初始颜色为空的BUG。 * 进度条控件文本锯齿修复。 * 按钮控件无法使用灰色按钮精灵破解。 |
| 2017-12-13 | 4.1.4 | * modbus tcp提供读取short数组的和ushort数组方法。 |
| 2017-12-13 | 4.1.5 | * 修复流水号生成器无法生成不带日期格式的流水号BUG。 |
| 2017-12-18 | 4.1.6 | * OperateResult成功时，消息为成功。 * 数据库辅助类API添加，方便的读取聚合函数。 * 日志类分析工具界面，显示文本微调。 |
| 2017-12-25 | 4.1.7 | * 进度条控件新增一个新的属性对象，是否使用动画。 |
| 2017-12-27 | 4.1.8 | * 新增一个饼图控件。 |
| 2017-12-28 | 4.1.9 | * 饼图显示优化，新增是否显示百分比的选择。 |
| 2017-12-31 | 4.2.0 | * 新增一个仪表盘控件。 |
| 2018-01-03 | 4.2.1 | * 饼图控件新增一个是否显示占比很小的信息文本。 * 新增一个旋转开关控件。 * 新增一个信号灯控件。 |
| 2018-01-05 | 4.2.2 | * 修复modbus tcp客户端读取 float, int, long,的BUG。 |
| 2018-01-08 | 4.2.3 | * 修复modbus tcp客户端读取某些特殊设备会读取不到数据的BUG。 |
| 2018-01-15 | 4.2.4 | * 双模式的网络基类中新增一个读取超时的时间设置，如果为负数，那么就不验证返回。 |
| 2018-01-24 | 4.3.0 | * 信号灯控件显示优化。 * Modbus Tcp服务端类修复内存暴涨问题。 * winfrom客户端提供一个曲线控件，方便显示实时数据，多曲线数据。 |
| 2018-02-05 | 4.3.1 | * 优化modbus tcp客户端的访问类，支持服务器返回错误信息。 * 优化曲线控件，支持横轴文本显示，支持辅助线标记，详细见对应博客。 |
| 2018-02-22 | 4.3.2 | * 曲线控件最新时间显示BUG修复。 * Modbus tcp错误码BUG修复。 * 三菱访问类完善long类型读写。 * 西门子访问类支持1500系列，支持读取订货号。 |
| 2018-03-05 | 4.3.3 | * 曲线控件增加一个新的属性，图标标题。 * Modbus tcp服务器端的读写BUG修复。 * 西门子访问类重新支持200smart。 |
| 2018-03-07 | 4.3.4 | * Json组件更新至11.0.1版本。 * 紧急修复日志类的BeforeSaveToFile事件在特殊情况的触发BUG。 |
| 2018-03-19 | 4.3.5 | * 修复Modbus-tcp服务器接收异常的BUG。 * 修复SoftBasic.ByteTo[U]ShortArray两个方法异常。 |
| 2018-04-05 | 5.0.0 | * 网络核心层重新开发，完全的基于异步IO实现。 * 所有双模式客户端类进行代码重构，接口统一。 * 完善并扩充OperateResult对象的类型支持。 * 提炼一些基础的更加通用的接口方法，在SoftBasic里面。 * 支持欧姆龙PLC的数据交互。 * 支持三菱的1E帧数据格式。 * 不兼容升级，谨慎操作。 |
| 2018-04-10 | 5.0.1 | * OperateResult静态方法扩充。 * 文件引擎提升缓存空间到100K，加速文件传输。 * 三菱添加读取单个bool数据。 * Modbus-tcp客户端支持配置起始地址不是0的服务器。 * 其他代码优化。 |
| 2018-04-14 | 5.0.2 | * ComplexNet服务器代码精简优化，移除客户端的在线信息维护代码。 * 西门子访问类第一次握手信号18字节改为0x02。 * 更新JSON组件到11.0.2版本。 * 日志存储类优化，支持过滤存储特殊关键字的日志。 * Demo项目新增控件介绍信息。 |
| 2018-04-20 | 5.0.3 | * 修复Modbus-Tcp服务器的空异常。 * 修复西门子类写入float，double，long数据异常。 * 修复modbus-tcp客户端读写字符串颠倒异常。 * 修复三菱多读取数据字节的问题。 * 双模式客户端新增异形客户端模式，变成了三模式客户端。 * 提供异形modbus服务器和客户端Demo方便测试。 |
| 2018-04-25 | 5.0.4 | * Modbus-tcp服务器同时支持RTU数据交互。 * 异形客户端新增在线监测，自动剔除访问异常设备。 * modbus-tcp支持读取输入点。 * 所有客户端设备的连接超时判断增加休眠，降低CPU负载。 * 西门子批量读取上限为19个数组。 * 其他小幅度的代码优化。 |
| 2018-04-30 | 5.0.5 | * Modbus相关的代码优化。 * 新增Modbus-Rtu客户端模式，配合服务器的串口支持，已经可以实现电脑本机的通讯测试了。 |
| 2018-05-04 | 5.0.6 | * 提炼数据转换基类，优化代码，修复WordReverse类对字符串的BUG，相当于修复modbus和omron读写字符串的异常。 * 新增一个全新的功能类，数据的推送类，轻量级的高效的订阅发布数据信息。具体参照Demo。 |
| 2018-05-07 | 5.0.7 | * Modbus服务器提供在线客户端数量属性。 * 所有服务器基类添加端口缓存。 * 双模式客户端完善连接失败，请求超时的消息提示。 * 修复双模式客户端某些特殊情况下的头子节NULL异常。 * 修复三菱交互类的ASCII协议下的写入数据异常。 |
| 2018-05-12 | 5.0.8 | * 新增一个埃夫特机器人的数据访问类。 * 双模式客户端的长连接支持延迟连接操作，通过一个新方法完成。 |
| 2018-05-21 | 5.0.9 | * 优化ComplexNet客户端的代码。 * 更新埃夫特机器人的读取机制到最新版。 * Modbus Rtu及串口基类支持接收超时时间设置，不会一直卡死。 * Modbus Tcp及Rtu都支持带功能码输入，比如读取100地址，等同于03X100。（注意：该多功能地址仅仅适用于Read及相关的方法 |
| 2018-05-22 | 5.0.10 | * Modbus Tcp及Rtu支持手动更改站号。也就是支持动态站号调整。 * 修复上个版本遗留的Modbus在地址偏移情况下会多减1的BUG。 |
| 2018-06-05 | 5.1.0 | * Modbus服务器支持串口发送数据时也会触发消息接收。 * IReadWriteNet接口新增Read(string address,ushort length)方法。 * 提炼统一的设备基类，支持Read方法及其扩展的子方法。 * 修复埃夫特机器人的读取BUG。 * 三菱PLC支持读取定时器，计数器的值，地址格式为"T100"，"C100"。 * 新增快速离散的傅立叶频谱变换算法，并在Demo中测试三种周期信号。 |
| 2018-06-16 | 5.1.1 | * 修复西门子fetch/write协议对db块，定时器，计数器读写的BUG。 * 埃夫特机器人修复tostring()的方法。 * modbus客户端新增两个属性，指示是否字节颠倒和字符串颠倒，根据不同的服务器配置。 * IReadWriteNet接口补充几个数组读取的方法。 * 新增一个全新的连接池功能类，详细请参见 https://www.cnblogs.com/dathlin/p/9191211.html * 其他的小bug修复，细节优化。 |
| 2018-06-27 | 5.1.2 | * IByteTransform接口新增bool[]数组转换的2个方法。 * Modbus Server类新增离散输入数据池和输入寄存器数据池，可以在服务器端读写，在客户端读。 * Modbus Tcp及Modbus Rtu及java的modbus tcp支持富地址表示，比如"s=2;100"为站号2的地址100信息。 * Modbus Server修复一个偶尔出现多次异常下线的BUG。 * 其他注释修正。 |
| 2018-07-13 | 5.1.3 | * Modbus服务器新增数据大小端配置。 * Modbus服务器支持数据存储本地及从本地加载。 * 修复modbus服务器边界读写bug。 * ByteTransformBase的double转换bug修复。 * 修复ReverseWordTransform批量字节转换时隐藏的一些bug。 * SoftBasic移除2个数据转换的方法。 * 修复modbus写入单个寄存器的高地位倒置的bug。 * 修复串口通信过程中字节接收不完整的异常。包含modbus服务器和modbus-rtu。 * 添加了.net 4.5项目，并且其他项目源代码引用该项目。添加了单元测试，逐步新增测试方法。 |
| 2018-07-27 | 5.2.0 | * 项目新增api文档，提供离线版和在线版，文档提供了一些示例代码。 * modbus-rtu新增批量的数组读取方法。 * modbus-rtu公开ByteTransform属性，方便的进行数据转换。 * SoftMail删除发送失败10次不能继续发送的机制。 * modbus server新增站号属性，站号不对的话，不响应rtu反馈。 * modbus server修复读取65524和65535地址提示越界的bug。 * Demo项目提供了tcp/ip的调试工具。 |
| 2018-08-08 | 5.2.1 | * API文档中西门子FW协议示例代码修复。 * modbus-rtu修复读取线圈和输入线圈的值错误的bug。 |
| 2018-08-23 | 5.2.2 | * Demo中三菱A-1E帧，修复bool读取显示失败的BUG。 * 数据订阅类客户端连接上服务器后，服务器立即推送一次。 * 串口设备基类代码提炼，提供了多种数据类型的读写支持。 * 仪表盘新增属性IsBigSemiCircle，设置为true之后，仪表盘可显示大于半圆的视图。 * 提供了一个新的三菱串口类，用于采集FX系列的PLC，MelsecFxSerial |
| 2018-08-24 | 5.2.3 | * 修复双模式基类的一个bug，支持不接受反馈数据。 * 修复三菱串口类的读写bug，包括写入位，和读取字和位。 * 相关代码重构优化。 |
| 2018-09-08 | 5.3.0 | * 串口基类接收数据优化，保证接收一次完整的数据内容。 * 新增一个容器罐子的控件，可以调整背景颜色。 * OperateResult成功时的错误码调整为0。 * 修复modbus-tcp及modbus-rtu读取coil及discrete的1个位时解析异常的bug。 * 授权类公开一个属性，终极秘钥的属性，感谢 洛阳-LYG 的建议。 * 修复transbool方法在特殊情况下的bug * NetworkDeviceBase 写入的方法设置为了虚方法，允许子类进行重写。 * SoftBasic: 新增三个字节处理的方法，移除前端字节，移除后端字节，移除两端字节。 * 新增串口应用的LRC校验方法。还未实际测试。 * Siemens的s7协议支持V区自动转换，方便数据读取。 * 新增ab plc的类AllenBradleyNet，已测试读写，bool写入仍存在一点问题。 * 新增modbus-Ascii类，该类库还未仔细测试。 * 埃夫特机器人更新，适配最新版本数据采集。 * 其他的代码优化，重构精简 |
| 2018-09-10 | 5.3.1 | * 修复埃夫特机器人读取数据的bug，已测试通过。 * ByteTransform数据转换层新增一个DataFormat属性，可选ABCD,BADC,CDAB,DCBA * 三个modbus协议均适配了ByteTransform并提供了直接修改的属性，默认ABCD * 注意：如果您的旧项目使用的Modbus类，请务必重新测试适配。给你带来的不便，敬请谅解。 |
| 2018-09-21 | 5.3.2 | * 所有显示字符串支持中英文，支持切换，默认为系统语言。 * Json组件依赖设置为不依赖指定版本。 * modbus-ascii类库测试通过。 * 新增松下的plc串口读写类，还未测试。 * 西门子s7类写入byte数组长度不受限制，原先大概250个字节左右。 * demo界面进行了部分的中英文适配。 * OperateResult类新增了一些额外的构造方法。 * SoftBasic新增了几个字节数组操作相关的通用方法。 * 其他大量的细节的代码优化，重构。 |
| 2018-09-27 | 5.3.3 | * DeviceNet层添加异步的API，支持async+await调用。 * java修复西门子的写入成功却提示失败的bug。 * java代码重构，和C#基本保持一致。 * python版本发布，支持三菱，西门子，欧姆龙，modbus，数据订阅，同步访问。 * 其他的代码优化，重构精简。 |
| 2018-10-20 | 5.4.0 | * python和java的代码优化，完善，添加三菱A-1E类。 * 修复仪表盘控件，最大值小于0会产生的特殊Bug。 * NetSimplifyClient: 提供高级.net的异步版本方法。 * serialBase: 新增初始化和结束的保护方法，允许重写实现额外的操作。 * softBuffer: 添加一个线程安全的buffer内存读写。 * 添加西门子ppi协议类，针对s7-200，需要最终测试。 * Panasonic: 修复松下plc的读取读取数据异常。 * 修复fx协议批量读取bool时意外的Bug。 * NetSimplifyClient: 新增带用户int数据返回的读取接口。 |
| 2018-10-24 | 5.4.1 | * 新增一个温度采集模块的类，基于modbus-rtu实现，阿尔泰科技发展有限公司的DAM3601模块。 |
| 2018-10-25 | 5.4.2 | * 三菱的mc协议新增支持读取ZR文件寄存器功能。 |
| 2018-10-30 | 5.4.3 | * 修复AB PLC的bool和byte写入失败的bug，感谢 北京-XLang 提供的思路。 |
| 2018-11-1 | 5.5.0 | * 新增西门子PPI通讯类库，支持200，200smart等串口通信，感谢 合肥-加劲 和 江阴- ∮溪风-⊙\_⌒ 的测试 |
| 2018-11-5 | 5.5.1 | * 新增三菱计算机链接协议通讯库，支持485组网，有效距离达50米，感谢珠海-刀客的测试。 * 串口协议的基类提供了检测当前串口是否处于打开的方法接口。 * 西门子S7协议新增槽号为3的s7-400的PLC选项，等待测试。 |
| 2018-11-9 | 5.5.2 | * 西门子PPI写入bool方法名重载到了Write方法里。 * 松下写入bool方法名重载到了Write方法里。 * 修复CRC16验证码在某些特殊情况下的溢出bug。 * 西门子类添加槽号和机架号属性，只针对400PLC有效，初步测试可读写。 * ab plc支持对数组的读写操作，支持数组长度为0-246，超过246即失败。 * 三菱的编程口协议修复某些特殊情况读取失败，却提示成功的bug。 * 串口基类提高缓存空间到4096，并在数据交互时捕获COM口的异常。 |
| 2018-11-16 | 5.6.0 | * 修复欧姆龙的数据格式错误，修改为CDAB。 * 新增一个瓶子的控件。 * 新增一个管道的控件。 * 初步新增一个redis的类，初步实现了读写关键字。 |
| 2018-11-21 | 5.6.1 | * AB PLC读取数组过长时提示错误信息。 * 正式发布redis客户端，支持一些常用的操作，并提供一个浏览器。博客：https://www.cnblogs.com/dathlin/p/9998013.html |
| 2018-11-24 | 5.6.2 | * 曲线控件的曲线支持隐藏其中的一条或是多条曲线，可以用来实现手动选择显示曲线的功能。 * Redis功能块代码优化，支持通知服务器进行数据快照保存，包括同步异步。 * Redis新增订阅客户端类，可以实现订阅一个或是多个频道数据。 |
| 2018-11-30 | 5.6.3 | * 串口数据接收的底层机制重新设计。 * 串口底层循环验证缓冲区是否有数据的间隔可更改，默认20ms。 * 串口底层新增一个清除缓冲区数据的方法。 * 串口底层新增一个属性，用于配置是否在每次读写前清除缓冲区的脏数据。 * 新增了一个SharpList类，用于超高性能的管理固定长度的数组。博客：https://www.cnblogs.com/dathlin/p/10042801.html |
| 2018-12-3 | 5.6.4 | * Networkbase: 接收方法的一个多余对象删除。 * 修复UserDrum控件的默认的text生成，及复制问题。 * UserDrum修复属性在设计界面没有注释的bug。 |
| 2018-12-5 | 5.6.5 | * 修复Demo程序在某些特殊情况下无法在线更新的bug。 * 修复曲线控件隐藏曲线时在某些特殊情况的不隐藏的bug。 * modbus协议无论读写都支持富地址格式。 * 修复连接池清理资源的一个bug，感谢 泉州-邱蕃金 * 修复java的modbus代码读取线圈异常的操作。 * Demo程序新增免责条款。 |
| 2018-12-11 | 5.6.6 | * 修复redis客户端对键值进行自增自减指令操作时的类型错误bug。 * 修复redis客户端对哈希值进行自增自减指令操作时的类型错误bug。 * 推送的客户端可选委托或是事件的方式，方便labview调用。 * 推送的客户端修复当服务器的关键字不存在时连接未关闭的Bug。 * Demo程序里，欧姆龙测试界面新增数据格式功能。 |
| 2018-12-19 | 5.6.7 | * ByteTransfer数据转换类新增了一个重载的构造方法。 * Redis客户提供了一个写键值并发布订阅的方法。 * AB-PLC支持槽号选择，默认为0。 * PushNet推送服务器新增一个配置，可用于设置是否在客户端刚上线的时候推送缓存数据。 * PushNet推送服务器对客户端的上下限管理的小bug修复。 * 本版本开始，组件将使用强签名。 * 本版本开始，组件的控件库将不再维护更新，所有的控件在新的控件库重新实现和功能增强，VIP群将免费使用控件库。 * VIP群的进入资格调整为赞助200Rmb，谢谢支持。 |
| 2018-12-27 | 5.7.0 | * 修复modbus服务器地址写入的bug，之前写入地址数据后无效，必须带x=3;100才可以。 * 修复极少数情况内核对象申请失败的bug，之前会引发资源耗尽的bug。 * SoftBasic的ByteToBoolArray新增一个转换所有位的重载方法，不需要再传递位数。 * 埃夫特机器人新增旧版的访问类对象，达到兼容的目的。 * Demo程序新增作者简介。 * 修复Demo程序的redis订阅界面在设置密码下无效的bug。 * Demo程序的免责界面新增demo在全球的使用情况。 * VIP群将免费使用全新的控件库，谢谢支持。地址：https://github.com/dathlin/HslControlsDemo |
| 2018-12-31 | 5.7.1 | * 修复modbus服务器地址读取的bug，之前读取地址数据后无效，必须带x=3;100才可以。 * NetPush功能里，当客户端订阅关键字时，服务器即使没有该关键字，也成功。 * 三菱的通讯类支持所有的字读取。例如读取M100的short数据表示M100-M115。 * VIP群将免费使用全新的控件库，谢谢支持。地址：https://github.com/dathlin/HslControlsDemo |
| 2019-1-15 | 5.7.2 | * 修复三菱A-1E协议的读取数据的BUG错误，给大家造成的不便，非常抱歉。 * VIP群将免费使用全新的控件库，谢谢支持。地址：https://github.com/dathlin/HslControlsDemo |
| 2019-2-7 | 5.7.3 | * 欧姆龙读写机制更改，报警的异常不再视为失败，仍然可以解析数据。 * Modbus地址优化，Modbus服务器的地址读写优化。 * 新增一个数据池类，SoftBuffer，主要用来缓存字节数组内存的，支持BCL数据类型读写。 * Modbus服务器的数据池更新，使用了最新的数据池类SoftBuffer。 * SoftBasic类新增一个GetEnumFromString方法，支持从字符串直接生成枚举值，已通过单元测试。 * 新增一个机器人的读取接口信息IRobotNet，统一化所有的机器人的数据读取。 * Demo程序中增加modbus的服务器功能。 * VIP群将免费使用全新的控件库，谢谢支持。地址：https://github.com/dathlin/HslControlsDemo |
| 2019-2-13 | 5.7.4 | * 日志存储的线程号格式化改为D3，也即三位有效数字。 * 日志存储事件BeforeSaveToFile里允许设置日志Cancel属性，强制当前的记录不存储。 * JSON库更新到12.0.1版本。 * SoftBasic新增一个GetTimeSpanDescription方法，用来将时间差转换成文本的方法。 * 调整日志分析控件不随字体变化而变化。 * 其他的代码精简优化。 * VIP群将免费使用全新的控件库，谢谢支持。地址：https://github.com/dathlin/HslControlsDemo |
| 2019-2-21 | 5.8.0 | * SoftBasic修复AddArrayData方法批量添加数据异常的bug，导致曲线控件显示异常。 * 提炼一个公共的欧姆龙辅助类，准备为串口协议做基础的通用支持。 * RedisHelper类代码优化精简，提炼部分的公共逻辑到NetSupport。 * SoftBuffer: 新增读写单个的位操作，通过位的与或非来实现。 * SiemensS7Server：新增一个s7协议的服务器，可以模拟PLC，进行通讯测试或是虚拟开发。 * 其他的代码精简优化。 * VIP群将免费使用全新的控件库，谢谢支持。地址：https://github.com/dathlin/HslControlsDemo |
| 2019-3-4 | 6.0.0 | * 西门子虚拟PLC的ToString()方法重新实现。 * 埃夫特机器人的json格式化修正换行符。 * IReadWriteNet接口添加Write(address, bytes)的方法。 * Modbus虚拟服务器修复写入位操作时影响后面3个位的bug。 * SoftBuffer内存数据池类的SetValue(byte,index)的bug修复。 * 西门子虚拟PLC和Modbus服务器新增客户端管理，关闭时也即断开所有连接。 * 三菱编程口协议的读取结果添加错误说明，显示原始返回信号，便于分析。 * 三菱MC协议新增远程启动，停止，读取PLC型号的接口。 * 新增三菱MC协议的串口的A-3C协议支持，允许读写三菱PLC的数据。 * 新增欧姆龙HostLink协议支持，允许读写PLC数据。 * 新增基恩士PLC的MC协议支持，包括二进制和ASCII格式，支持读写PLC的数据。 * 所有PLC的地址说明重新规划，统一在API文档中查询。 * 注意：三菱PLC的地址升级，有一些地址格式进行了更改，比如定时器和计数器，谨慎更新，详细地址参考最新文档。 * 如果有公司使用了本库并愿意公开logo的，将在官网及git上进行统一显示，有意愿的联系作者。 * VIP群将免费使用全新的控件库，谢谢支持。地址：https://github.com/dathlin/HslControlsDemo |
| 2019-3-10 | 6.0.1 | * 修复代码注释上的一些bug，三菱的注释修复。 * 调整三菱和基恩士D区数据和W区数据的地址范围，原来只支持到65535。 * SoftIncrementCount: 修复不持久化的序号自增类的数据复原的bug，并添加totring方法。 * IRobot接口更改。针对埃夫特机器人进行重新实现。 * RedisClient: 修复redis类在带有密码的情况下锁死的bug。 * 初步添加Kuka机器人的通讯类，等待测试。 * 西门子的s7协议读写字符串重新实现，根据西门子的底层存储规则来操作。 * Demo的绝大多的界面进行重构。更友好的支持英文版的显示风格。 * 如果有公司使用了本库并愿意公开logo的，将在官网及git上进行统一显示，有意愿的联系作者。 |
| 2019-3-21 | 6.0.2 | * 修复西门子s7协议读写200smart字符串的bug。 * 重构优化NetworkBase及NetwordDoubleBase网络类的代码。 * 新增欧姆龙的FinsUdp的实现，DA1【PLC节点号】在配置Ip地址的时候自动赋值，不需要额外配置。 * FinsTcp类的DA1【PLC节点号】在配置Ip地址的时候自动赋值，不需要额外配置。 |
| 2019-3-28 | 6.0.3 | * NetPushServer推送服务器修复某些情况下的推送卡死的bug。 * SoftBuffer内存数据类修复Double转换时出现的错误bug。 * 修复Kuka机器人读写数据错误的bug，已通过测试。 * 修复三菱的MelsecMcAsciiNet类写入bool值及数组会导致异常的bug，已通过单元测试。 * SoftBasic新增从字符串计算MD5码的方法。 |
| 2019-4-4 | 6.0.4 | * 修复java的NetPushClient掉线重复连接的bug。 * 发布java的全新测试Demo。 * Kuka机器人Demo修改帮助链接。 * 西门子新增s200的以太网模块连接对象。 * 修复文件引擎在上传文件时意外失败，服务器仍然识别为成功的bug。 |
| 2019-4-17 | 6.1.0 | * 修复日志存储自身异常时，时间没有初始化的bug。 * NetworkBase: 新增UseSynchronousNet属性，默认为true，通过同步的网络进行读写数据，异步手动设置为false。 * 修复西门子的读写字符串的bug。 * 添加KeyenceNanoSerial以支持基恩士Nano系列串口通信。 * 其他的代码优化。 * 发布一个基于xamarin的安卓测试demo。 * 发布官方论坛： http://bbs.hslcommunication.cn/ |
| 2019-4-24 | 6.1.1 | * 修复基恩士MC协议读取D区数据索引不能大于100000的bug。 * 修复基恩士串口协议读写bool数据的异常bug。 * 修复数据推送服务器在客户端异常断开时的奔溃bug，界面卡死bug。 * SoftNumericalOrder类新增数据重置和，最大数限制 。 * ModbusTcp客户端公开属性SoftIncrementCount，可以强制消息号不变，或是最大值。 * NetworkBase: 异步的方法针对Net451及standard版本重写。 * 串口基类serialbase的初始化方法新增多个重载方法，方便VB和labview调用。 * NetworkBase: 默认的机制任然使用异步实现，UseSynchronousNet=false。 * 发布官方论坛： http://bbs.hslcommunication.cn/ |
| 2019-4-25 | 6.1.2 | * 紧急修复在NET451和Core里的异步读取的bug。 * 紧急修复PushNetServer的发送回调bug。 * 发布官方论坛： http://bbs.hslcommunication.cn/ |
| 2019-5-6 | 6.2.0 | * SoftBuffer缓存类支持bool数据的读写，bool数组的读写，并修复double读写的bug。 * Modbus虚拟服务器代码重构实现，继承自NetworkDataServerBase类。 * 新增韩国品牌LS的Fast Enet协议 * 新增韩国品牌LS的Cnet协议 * 新增三菱mc协议的虚拟服务器，仅支持二进制格式的机制。 * LogNet支持写入任意的字符串格式。 * 其他的注释添加及代码优化。 * 发布官方论坛： http://bbs.hslcommunication.cn/ |
| 2019-5-9 | 6.2.1 | * 修复三菱读写PLC位时的bug。 * 修复Modbus读写线圈及离散的变量bug。 * 强烈建议更新，不能使用6.2.0版本！或是回退更低的版本。 * 有问题先上论坛： http://bbs.hslcommunication.cn/ |
| 2019-5-10 | 6.2.2 | * 修复上个版本modbus的致命bug，已通过单元测试。 * 新增松下的mc协议，demo已经新增，等待测试。 * github源代码里的支持的型号需要大家一起完善。 |
| 2019-5-31 | 6.2.3 | * Ls的Fast Enet协议问题修复，感谢来自埃及朋友。 * Ls的CEnet协议问题修复，感谢来自埃及朋友。 * Ls新增虚拟的PLC服务器，感谢来自埃及朋友。 * 改进了机器码获取的方法，获取实际的硬盘串号。 * 日志的等级为None的情况，不再格式化字符串，原生写入日志。 * IReadWriteNet接口测试西门子的写入，没有问题。 * 松下MC协议修复LD数据库的读写bug。 * Redis的DEMO界面新增删除key功能。 |
| 2019-6-3 | 6.2.4 | * Redis新增读取服务器的时间接口，可用于客户端的时间同步。 |
| 2019-6-6 | 6.2.5 | * 西门子的SiemensS7Net类当读取PLC配置长度的DB块数据时，将提示错误信息。 |
| 2019-6-22 | 7.0.0 | * 新增安川机器人通信类，未测试。 * 西门子的多地址读取的长度不再限制为19个，而是无限制个。 * NetworkDoubleBase: 实现IDispose接口，方便手动释放资源。 * SerialBase: 实现IDispose接口，方便手动释放资源。 * NetSimplifyClient:新增一个async...await方法。 * NetSimplifyClient:新增读取字符串数组。 * ModbusServer:新增支持账户密码登录，用于构建安全的服务器，仅支持hsl组件的modbus安全访问。 * NetSimplifyServer:新增支持账户密码登录。 * 新增永宏PLC的编程口协议。 * 新增富士PLC的串口通信，未测试。 * 新增欧姆龙PLC的CIP协议通讯。 * 初步添加OpenProtocol协议，还未完成，为测试。 * MelsecMcNet:字单位的批量读取长度突破960长度的限制，支持读取任意长度。 * MelsecMcAsciiNet:字单位的批量读取长度突破480长度的限制，支持读取任意长度。 * AllenBradleyNet:读取地址优化，支持读取数组任意起始位置，任意长度，支持结构体嵌套读取。 * 其他大量的代码细节优化。 |
| 2019-6-25 | 7.0.1 | * IReadWriteNet完善几个忘记添加的Write不同类型参数的重载方法。 * IReadWriteNet新增ReadBool方法，Write(string address, bool value)方法，是否支持操作需要看plc是否支持，不支持返回操作不支持的错误。 * OmronFinsNet:新增一个属性，IsChangeSA1AfterReadFailed，当设置为True时，通信失败后，就会自动修改SA1的值，这样就能快速链接上PLC了。 * OmronFinsNet:新增读写E区的能力，地址示例E0.0，EF.100，E12.200。 * 新增HslDeviceAddress特性类，现在支持直接基于对象的读写操作，提供了一种更加便捷的读写数据的机制，详细的关注后续的论坛。 * 本组件的最后一个对个人免费的版本，企业使用一律需要授权，授权流程为签合同，付款，开票。 |
| 2019-9-10 | 8.0.0 | * SimpleHybirdLock: 简单混合锁的性能优化，基元对象采用懒加载的机制实现，同时增加了高级混合锁的类，支持自旋，线程拥有权，在高竞争的情况下性能大幅增加。 * NetSoftUpdateServer: 软件自动更新的服务器端支持传送指定目录下的文件及其子文件夹下的所有文件内容，都将更新到客户端的电脑上去。 * AllenBradleyNet: 修复字符串的读写bug，支持读写任意长度的字符串信息。 * MelsecFxSerial: 三菱编程口协议支持读写D1024以上地址的数据，感谢 厦门-Mr.T 的贡献。 * PIDHelper: 新增一个Pid的辅助类，用于模拟pid的波形情况。 * NetPushClient: 修改一个时间的注释，追加单位信息，时间的单位是毫秒。 * XGBFastEnet: 感谢埃及的朋友，修复了一些bug信息。 * MelsecFxSerialOverTcp: 新增基于网口透传的三菱的编程口通讯类。 * MelsecFxLinksOverTcp: 新增基于网口透传的三菱的计算机链接协议的通讯类。 * MelsecA3CNet1OverTcp: 新增基于网口透传的三菱的A-3C的协议的通讯类。 * OmronHostLinkOverTcp: 新增基于网口透传的欧姆龙的hostLink协议的通讯类。 * PanasonicMewtocolOverTcp: 新增基于网口透传的松下的Mewtocol协议的通讯类。 * SiemensPPIOverTcp: 新增基于网口透传的西门子PPi协议的通讯类。 * XGBCnetOverTcp: 新增基于网口透传的Lsis的XGBCnet协议的通讯类。 * KeyenceNanoSerialOverTcp: 新增基于网口透传的基恩士的NanoSerial串口协议的通讯类。 * FujiSPBOverTcp: 新增基于网口透传的富士的SPB串口协议的通讯类。 * ModbusRtuOverTcp: 新增基于网口透传的Modbus rtu协议的通讯类。 * Modbus相关的功能类进行代码精简，重构，优化，api标准化为ReadBool,WriteBool,Read,Write，移除了一些特殊的方法api，本次升级不兼容。 * FFTFilter: 新增一个基于FFT（快速离散傅立叶变换）的滤波功能，可以作为一个高级的曲线拟合方案，详细参照demo，感谢 北京-monk 网友的支持。 * KnxUdp: 新增一个KnxUdp的数据通讯类，感谢上海-null的支持。 * ABBWebApiClient: 新增ABB机器人的基于web api的访问机制的通讯类。 * SickIcrTcpServer: 新增一个sick的条码读取类，支持被动连接，主动连接，经过测试，同时支持海康，基恩士，DATELOGIC扫码器。 * Demo: Demo工具新增了一个基于tcp的服务器的测试界面。 * 本组件从v8.0.0开始进入付费模式，谨慎升级，未激活的将只能使用8小时，普通vip群发放激活码，仅支持个人用途使用，禁止破解，感谢对正版的支持。 * 今天是2019年9月10日，祝天下所有的教师节日快乐。 |
| 2019-9-17 | 8.0.1 | * 所有网口透传类对象完善实例化的方法，都新增一个指定ip及端口的实例方法。 * ABBWebClient: 完善实例化方法，修改ToString的格式化内容，提炼了webapi的基类，开放ip地址和端口。 * ABBWebClient: 新增提供了机器人自身IO，扩展IO，最新的报警日志的数据读取API。 * NetSimplifyClient: 修复了当ReceiveTimeOut小于0，但是实际接收时会发生奔溃的bug。 * NetPlainSocket: 新增一个基于socket的明文的网络发送和接收类，采用事件驱动的机制。 * LogNet: 日志类对象新增一个特性，当日志的文件名设置为空的时候，将不会创建文件，仅仅触发 BeforSaveToFile 事件，方便日志显示。 * XGBCnet: Lsis的plc的串口类修复一个bug，感谢埃及朋友的贡献。 * SoftIncrementCount: 消息号自增类新增一个方法，重置当前的消息号。 * PanasonicMewtocol: 修复松下的串口类读写单个bool时异常的bug，地址支持字+位的表示方式，R33=R2.1，方便大家输入测试。 * MqttClient: 新增一个Mqtt协议的客户端类，支持用户名密码，支持发布，支持订阅，支持重连，欢迎一起测试。 * 本组件从v8.0.0开始进入付费授权模式，详细参考官方：http://www.hslcommunication.cn 。 |
| 2019-9-19 | 8.0.2 | * ABBWebClient: abb机器人的api读取日志的接口新增一个参数，读取最近的日志数量。默认为10条。 * MQTTClient: 修复mqtt客户端类的消息重复bug，修复发送空订阅的bug。 * SiemensS7Net: 西门子的s7协议的类新增一个api，支持时间的读写，支持异步，时间格式和s7net一致。 * 本组件从v8.0.0开始进入付费授权模式，详细参考官方：http://www.hslcommunication.cn 。 |
| 2019-9-26 | 8.0.3 | * Networkbase: 修复套接字网络授权失败时不关闭网络的bug。 * SoftBasic: 新增一个数组数据格式化的方法信息。 * MqttServer: 新增一个mqtt的服务器，初步支持订阅，发布订阅，强制发布订阅，在线客户端数量功能等等。 * Demo: 所有的PLC的demo和modbus协议的demo，支持批量读取各种类型的数组数据。 * Nuget: 新增本项目的图标，在nuget上搜索时会显示图标。 * 本组件从v8.0.0开始进入付费授权模式，详细参考官方：http://www.hslcommunication.cn 。 |
| 2019-10-7 | 8.1.0 | * ModbusUdp: 新增一个Modbus的基于udp的协议类，使用的tcp的报文的机制。 * HttpServer: 新增一个http的服务器封装类，方便实现基于webapi的后台功能，集成GET，POST的接口操作。 * Serial Ports: standard项目依赖官方串口库，实现所有的设备的串口支持，可应用于跨平台。 * standard: 在nuget上提供.net standard2.1版本的库。 * 本组件从v8.0.0开始进入付费授权模式，详细参考官方：http://www.hslcommunication.cn 。 |
| 2019-10-11 | 8.1.1 | * Lsis: 感谢埃及朋友的支持，修复了一些bug，支持了bool的操作。 * Redis: 新增db块属性设置，修复短连接下切换db块无效的bug，因为db块是跟随连接的。 * MQTT: 修复客户端和服务器的长度计算bug，支持和其他mqtt组件混合使用。 * MQTT Demo: 优化demo功能，支持文本追加或是覆盖选择，文本格式化查看选择。 * Http Server: 支持跨域属性选择，编码统一为utf-8，兼容浏览器和postman，demo中增加返回类型示例。 * Modbus server及Lsis Server: 针对.net standard版本，开放串口。 * 本组件从v8.0.0开始进入付费授权模式，详细参考官方：http://www.hslcommunication.cn 。 |
| 2019-10-16 | 8.1.2 | * Lsis: 感谢埃及朋友的支持，demo增加了bool操作。 * Knx驱动：增加测试demo，完善驱动，测试通过，有需要的朋友可以查看信息。 * IntegrationFileClient: 完善文件的收发类，新增重载的构造方法，传入ip地址及端口即可。 * melsec: 三菱的MC协议部分错误代码关联了文本信息，在测试的时候即可弹出错误信息，方便排查，常见了已经绑定。 * melsec: 新增3e协议的随机字批量读取操作，支持跨地址，跨数据类型混合交叉读取，一次即可读完。 * fileserver: 修复linux下的bug，新增上传文件后的触发事件，将文件的信息都传递给调用者。 * SiemensMpi: 添加MPI协议，并完善demo，等待测试。 * 本组件从v8.0.0开始进入付费授权模式，详细参考官方：http://www.hslcommunication.cn 。 |
| 2019-10-24 | 8.1.3 | * Lsis: 感谢埃及朋友的支持，demo完善了cpu类型的选择。 * LogNet:新增移除关键字的接口方法，修复linux运行路径解析的bug，完善api文档的示例代码。 * 大量的细节优化，变量名称单次拼写错误的修复。 * Modbus: 当地址为x=3;100时，读正常，写入异常的问题修复，功能码自动替换为0x10。 * FileNet: 修复高并发下载时的下载异常的问题，调整指令头的超时时间。 * AB plc: 公开一个新的api接口，运行配置一些比较高级的数据。 * 接下来计划：1.完善hsl的demo，api文档，准备基础的入门视频；2.开始完善java版本的代码，java版本只对超级VIP群开放。 * 本组件从v8.0.0开始进入付费授权模式，详细参考官方：http://www.hslcommunication.cn 。 |
| 2019-12-3 | 8.2.0 | * 三菱的MC协议支持读取SM和SD，特殊连接继电器，特殊寄存器。 * PushNet优化相关代码。 * MelsecMcUdp: 新增三菱的MC协议的UDP通讯类。 * MelsecMcAsciiUdp: 新增三菱的MC协议的ASCII格式的UDP通讯类。 * MelsecMcServer: 三菱的虚拟服务器修复数据存储加载的bug。 * Serial: 串口的基类公开了一个Rts属性，用于某些串口无法读取的设备的情况。 * OmronFinsServer: 新增欧姆龙的虚拟plc，支持和hsl自身的通讯，支持cio，h区，ar区，d区的通信，不支持E区。 * AllenBradleyServer: 新增ab plc的虚拟plc，支持和hsl的自身的通讯，在demo里预设了4个变量值。不支持结构体和二维及以上数组读写。 * Aline: 异形的服务器对象新增一个设置属性，是否反馈注册结果，默认为True。 * SoftBasic: 数组格式化操作新增格式化的字符串说明。 * Modbus: 调整Write( string address, bool value )采用05功能码写入，而参数为bool[]的话，采用0F功能码。 * Modbus: 提供WriteOneRegister方法，写入单个的寄存器，使用06功能码。 * LogNet: 日志在实例化的时候，添加对当前设置的目录的是否存在的检查，如果不存在，则自动创建目录。 * Python: 大量代码更新，新增了一个欧姆龙的fins-tcp通信。 * Java: 大量代码更新，新增了一个AB plc的读写类。 |
| 2019-12-11 | 8.2.1 | * Cip协议：cip协议开放Eip指令自定义输入，优化指令生成算法。 * Cip协议：Write(string address, byte[] data)方法提示使用WriteTag信息。 * NetworkDoubleBase: 修复bool异步读写提示不支持的bug，现在可以使用异步了。 * SAMSerial：新增身份证阅读器的串口协议，支持读取身份证信息，头像信息还未解密。 * SAMTcpNet：新增身份证阅读器的串口透传协议，支持读取身份证信息，头像信息还未解密。 * BeckhoffAdsNet：新增倍福plc的协议，还未通过测试。 |
| 2020-1-3 | 8.2.2 | * lsis的plc优化，感谢埃及朋友的提供的技术支持。 * Panasonic: 松下的Mewtocol协议增加SR区的支持，解析地址的方法修改为Public，方便外面调用。 * Panasonic: 松下的Mewtocol协议批量读取bool方法是按字为单位，读取长度按照位为单位，地址写Y0，Y1，不能写Y0.4。 * ab-plc: 虚拟服务器修复上个版本造成的bug，导致读写数据成功，但是数据实际没有更改。 * ab-plc: 支持超长的数组读取，可以一次性读取任意长度的数组内容，不再需要手动切片。 * ab-plc: 新增一个api接口，可以遍历所有的ab-plc的变量名称。 * beckoff: 倍福的plc通信通过测试，需要设置正确各种网络号才可以，优化了标签缓存。 * java: java版本的ab-plc类优化，支持超长的数组读取。 * python: python版本的代码新增ab-plc的读取类。 * demo: 安卓的demo增加lsis，mqtt协议的界面。 * Melsec: 三菱PLC的多块批量读取目前只支持字地址，后续继续优化。 * 其他的代码优化和重构。 * java版本的源代码及demo，python版本的源代码及demo仅对商业授权用户开放，谢谢支持。 * 作者于2020年1月5日和王女士结婚，地址是浙江省金华市兰溪市马涧镇，欢迎有空的老铁们来坐坐，带红包就行。 |
| 2020-2-13 | 9.0.0 | * 宣布V9版本脱胎换骨，浴火重生，C#版本的组件库底层网络大幅重构，删除一直以来的伪异步，原先的通过改为纯同步，并从底层提供完整的异步方法。 * 注意：不兼容升级，影响范围，MQTT协议的事件，网络的同步设置，西门子的PPI协议取消WriteByte方法，改为和其他一样的Write(string address,byte value)重载了，升级请谨慎测试。 * 所有的PLC通讯类，机器人类通讯类，Modbus通讯类，身份证类，包括 IReadWriteNet 接口都实现了异步的操作，针对NET45以上及Standard平台。 * MQTT协议修改触发的消息事件，返回session信息，支持自定义返回数据信息，支持当前消息的发布拦截操作，服务器主送发布的消息支持是否驻留，默认主题驻留。 * 新增websocket协议的服务器，客户端，问答客户端。支持直接从C#的后台发送数据到网页前端，支持订阅操作。详细见demo的操作。 * ComplexNet,SimplifyNet,PushNet,FileNet这几个网络引擎代码优化，初步测试OK。 * SoftBasic: 新增方法SpliceStringArray，用来合并字符串信息。增加了ByteToHexString的空校验。 * HttpServer: 异步优化，修复读取数据时可能长度不足的bug。 * YRC1000: 安川机器人修改无法读取的bug，目前已经测试通过，感谢网友的支持。 * Java: 修复ab-plc读取失败的错误信息，原因来自一个强制转换失败的错误。 * 本版本改动较多，尽管我已经仔细测试过，但是仍然不可避免存在一些bug，欢迎大家使用，测试，有问题可以报告给我，相信hsl组件会变的更加强大。 * 本版本依然是商业授权的，如果需要测试，可以付费240rmb，加入vip群，可以将hsl用于测试环境和研究学术用途，欢迎大家加我的支付宝好友，hsl200909@163.com * 加油，武汉！加油，中国！疫情之后，无人自动化工厂将会获得更大的关注和发展，我辈当自强。 |
| 2020-2-19 | 9.0.1 | * 底层的网络在对方关闭连接后，不再等待接收，直接返回对方已关闭的错误信息，提供通信的性能。 * 四个服务器类，complexserver, simplifyserver,mqttserver,websocketserver开发关闭客户端连接的方法，调用者可以手动操作关闭。 * MQTT服务器新增一个客户端上线事件，包含客户端的会话参数，方便实现一些特殊的场景需求，在api文档中增加调用示例。 * Websocket服务器新增一个客户端上线事件，包含客户端的会话参数，方便实现一些特殊的场景需求。 * Websocket服务器添加0x0A的心跳返回，用于有些客户端的心跳验证操作。 * RedisClient: redis相关的代码优化，调整，添加了异步api接口，本机性能测试不如同步，有待优化。 * RedisClient: 新增基于特性的读写，自动组合键批量读取，组合哈希键批量读取操作，提升性能，详细参考api文档。写入操作不支持列表相关的特性。 * Demo的写入byte操作的反射代码获取失败的bug修复。 |
| 2020-2-25 | 9.0.2 | * 修复websocket对某些浏览器的请求验证失败的bug，改为正则表达式验证，适用的范围更加广阔。 * 三菱的mc协议的错误信息更加明确化，将提供更加确切的错误描述，方便大家查找错误。 * websocket客户端新增连接服务器成功的事件，方便实现类似订阅的功能。 * Websocket服务器添加心跳检测功能，将会定期（可以自定义）发送心跳包给客户端，在检测客户端是否在线。 * 文件的服务器和客户端开放文件缓存大小的属性，默认100K，越大的话，性能越高，越占内存。 * Modbus协议功能调整，Write(string,short)和Write(string,ushort)功能码调整为06，如果需要0x10功能码，使用Write(string,short[])和Write(string,ushort[]) * 新增汇川PLC的通讯类，基于modbus协议，但是实现了地址的自动解析，输入D100即可自动转为modbus的地址，包含AM系列，H3U系列，H5U系列等 * 在示例文档中，新增大量的代码说明，完善注释，如果有任何的问题，优先参考api文档。 * 官网新增一个来自上海亦仕智能科技有限公司 MES DEMO： http://111.229.255.209 账号SF 密码 123 |
| 2020-3-3 | 9.0.3 | * 修复汇川PLC的地址示例文档写错的bug。 * IReadWriteNet标准化字符串读写操作，新增定制编码的字符串读写。netDeviceBase移除之前的writeunicode的方法。这点如果有使用，谨慎更新。 * 串口基类和UDP基类的数据交互方法新增日志记录，对发送的数据和接收的数据写入debug等级的日志。 * 数据服务器（主要是三菱虚拟plc，西门子虚拟plc，modbus服务器等）实现IReadWriteNet接口。 * 关于ab-plc，新增MicroCip协议，适用于Micro800系列读写操作。 * 关于序号生成器类SoftIncrement，重置最大值的方法名称更新，添加了重置当前值，重置初始值，支持反向序列，跳跃序列的功能，详细参考api文档。 * 文件的服务器类，新增一些日志记录，关于文件何时被读取，何时读取结束的日志信息，等级为debug。 * NuGet组件更新，json组件更新到12.0.3版本，IO.port更新到4.7.0版本。单元测试框架更新。 * Demo的redis示例，支持不同的db块选择，当你选择数据后自动切换，键值类数据增加格式化显示。 * NetworkBase: 网络基类的连接服务器改为如果连接立即失败(500ms内)，将会休眠100ms后，立即再尝试一次，提高连接的成功率。影响范围为所有客户端类。 * 三菱二进制MC协议：地址里面新增标签访问，缓冲存储器访问，扩展的地址访问的方式，目前开放二进制的mc协议，欢迎测试，顺利的话，完善写入和ascii格式的。 * 大量的代码注释添加，主流的常用的代码添加中英文注释，后续逐步全都改为中英文，方便国外客户阅读。 * 240元的普通vip群的激活码时间调整，改为20年，中间软件重启一次，就又是20年，感谢大家的理解和支持。 * http://www.hslcommunication.cn/MesDemo 官网的地址以后作为优秀的MES产品展示平台，欢迎大家关注。 |
| 2020-3-15 | 9.1.0 | * MQTT: 服务器增加定时检测客户端在线情况，超过设置的时间不活跃，强制下线，开放OnlineSession属性，获取在线客户端，查看ip，端口，在线时间等信息。 * WebSocket: 服务器开放OnlineSession属性，获取在线客户端，查看ip，端口，在线时间等信息。 * Language: 组件的语言系统修复设置英文后，无法切换回中文的bug。 * SoftBasic: 添加SpliceByteArray(params byte[][] bytes)方法，用来将任意个byte[]进行拼接成一个byte[]。 * SoftBasic: 添加BoolOnByteIndex方法，用来获取byte数据的指定位的bool值。 * Panasonic: 松下的mc地址和串口地址统一表示方式：R101=R10.1=[10\*16+1]，R10.F=R10.15(这两种表示方式都可以) * 发布基于HSL扩展组件HslCppExtension，将写入的重载方法名按照类型重写一遍，方便C++调用。 * VC++的demo示例，新增写入数据的例子，基于扩展组件HslCppExtension实现，详细参照demo源代码。 * SoftBasic: 针对byte数组的切割，选择头，尾，中间，移除头，尾的方法全部更改成泛型版本，方法名字已经变更，如果有调用，谨慎更新。 * FanucInterfaceNet: 新增读取fanuc机器人的通讯类，支持读写任意地址数据的功能，详细参考api文档，写入操作谨慎测试。 * FanucRobotServer: 新增fanuc机器人的虚拟服务器，方便进行测试，初始数据来自真实机器人，支持D,I,Q,AI,AQ,M数据区。 * Fanuc: 目前测试通过的型号为R-30iB mate plus，其他型号暂时不清楚支持范围。 * 代码注释优化，api文档大量的更新，添加一些示例代码，包含如果检测状态，长短连接，使用前请仔细阅读下面的信息：http://api.hslcommunication.cn * http://www.hslcommunication.cn/MesDemo 官网的地址以后作为优秀的MES产品展示平台，欢迎大家关注。 * 三年磨一剑，直插工业互联网的心脏。软件开发之艰辛，如人饮水冷暖自知。感谢所有支持的朋友。 |
| 2020-3-22 | 9.1.1 | * feat(SAM): 身份证阅读器修复在某些状态下接受数据不完整的bug，将校验数据的完整性。 * feat(ab-plc): 虚拟服务器的地址支持小数点的形式，支持单个的bool读写，支持string的读写操作，和客户端的体验一致。 * feat(softbasic): 方法针对数组切割的方法，增加扩展方法支持，byte[] a; byte[] b= a.RemoveBegin(2);意思就是移除最前面的2个元素。 * feat(softbasic): Hex字符串和byte[]的转化也支持扩展方法。byte[] a.ToHexString()。 * feat(melsec): 三菱的a-1e协议之前的，x,y地址采用8进制，先修改为自定义，如果要八进制，地址前面加0，例如X017，如果不加就是十六进制，例如X17，默认十六进制，升级需注意。 * feat(melsec): 三菱的a-1e协议增加了F报警继电器，B链接继电器，W链接寄存器，定时器和计数器的线圈，触点，当前值的读取，地址参见api文档说明。 * feat(melsec): 添加a-1e协议的ASCII版本，地址格式和二进制版本是一致的，支持的地址类型也是一致的，还未仔细测试，欢迎老铁们测试。 * feat(melsec): 三菱的mc虚拟服务器支持二进制和ascii，实例化的时候选择，支持和HSL组件自身的通讯。 * lsis: cnet和fenet地址的解析bug修复，感谢埃及朋友的支持。 * 代码注释优化，使用前请仔细阅读下面的信息：http://api.hslcommunication.cn * http://www.hslcommunication.cn/MesDemo 官网的地址以后作为优秀的MES产品展示平台，欢迎大家关注。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-3-29 | 9.1.2 | * ModbusAscii: 修复和rtu指令转换时的bug，之前会发生读写失败，目前已经在台达PLC上测试通过。 * MelsecA1EAscii：修复三菱的A1E协议的ascii格式类写入单个bool异常的bug。 * NetworkUdpServerBase：新增基于UDP协议的服务器基类，后台线程循环接收数据实现。 * CipServer: 虚拟的ab-plc服务器新增字符串数组对象的读写操作，demo相关的完善。 * HyundaiUdpNet: 新增现代机器人的姿态跟踪网络通讯类， * 代码注释优化，使用前请仔细阅读下面的信息：http://api.hslcommunication.cn * http://www.hslcommunication.cn/MesDemo 官网的地址以后作为优秀的MES产品展示平台，欢迎大家关注。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-4-6 | 9.1.3 | * HslExtension: 完善一些转化的api，方便数组和字符串转化，完善对象转JSON字符串。 * LogNet：消息格式化文本的消息等级追随HSL的语言设定，如果是中文，就显示调试，信息，警告，错误，致命。 * Redis: 修复ExpireKey，生存时间参数丢失的bug，完善了说明文档。 * OmronCip: 欧姆龙的CIP协议的类库，修复数组读取的bug，修复字符串写入bug，字符串写入还需要测试。 * Toledo：新增托利多电子秤的串口类及网口服务器类，方便接收标准的数据流，等待测试。 * Java：增加了单元测试的内容，对一些已经完成的类添加单元测试。 * Python：实现了python版本的HslCommunication程序，基于pyqt实现，初步添加了一些PLC的调试界面。 * 代码注释优化，使用前请仔细阅读下面的信息：http://api.hslcommunication.cn * http://www.hslcommunication.cn/MesDemo 官网的地址以后作为优秀的MES产品展示平台，欢迎大家关注。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-4-20 | 9.1.4 | * ServerBase: 服务器对象基类完善客户端下线逻辑，精简相关的代码。 * LogNet：设备网络通讯类及串口类在数据收发的时候增加日志的记录，可以设置PLC类的LogNet属性抓取相关的报文信息。 * ModbusServer: Modbus服务器同时支持TCP,RTU,ASCII，其中RTU和ASCII共用一个串口，根据报文头是否为冒号区分。 * ModbusAscii: 修复通讯的bug，已通过单元测试，支持和ModbusServer完美通讯，欢迎网友继续测试。 * MelsecMcNet：三菱MC协议的数据地址新增对SB，SW，特殊链接继电器，寄存器的支持。 * SiemensServer: 西门子S7虚拟服务器的DB块支持DB1.X，DB2.X，DB3.X，3以上的db块都是使用同一个的DB块。 * HttpServer：自定义轻量级的WebApi服务器支持反射对象的方法名，简化定义API时定义大量的if...else...。 * UdpNet：添加ConnectionId属性，使用的方法获取信息。 * MelsecMcRNet：添加三菱R系列的MC协议二进制的实现，和标准的有一点区别，地址支持也不一样，欢迎测试Demo。 * OmronCip：欧姆龙的读写数组已经测试通过，修改了读写字符串的逻辑实现，等待测试。 * 代码注释优化，使用前请仔细阅读下面的信息：http://api.hslcommunication.cn * http://www.hslcommunication.cn/MesDemo 官网的地址以后作为优秀的MES产品展示平台，欢迎大家关注。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-4-28 | 9.2.0 | * HttpServer: 当客户端发起request请求的时候，在日志记录的时候记录当前的请求的方式，GET,POST,OPTION等等。 * MQTT: mqtt的消息等级追加一个新的等级，为OnlyTransfer等级，用来表示只发送服务器，不触发发布操作。 * MqttServer: 配合Qos等级为OnlyTransfer时，进行相关的适配操作，并触发消息接收的事件。 * MqttSyncClient: 新增MQTT的同步访问的客户端，协议头标记为HUSL，向HSL的mqtt服务器进行数据请求并等待反馈。尚未添加心跳程序。 * MqttServer: 适配同步客户端实现功能，当客户端为同步客户端的时候，调试心跳验证。 * 至此，HSL的MQTT协议已经是兼容几大网络功能了，在线客户端管理，消息发布订阅，消息普通收发，同步网络访问。 * IByteTransform接口属性新增IsStringReverseByteWord，相当于从ReverseByWord挪过来了，默认为false，如果为true，在解析字符串的时候将两两字节颠倒。 * Omron: 欧姆的fins-tcp及fins-udp及hostlink的IByteTransform接口IsStringReverseByteWord调整为true默认颠倒。 * SerialBase: 串口基类的打开串口方法调整返回类型OperateResult，在串口数据读取之前增加打开串口的Open方法，串口类也只需要一直读就可以了。 * NetworkDoubleBase, SerialDeviceBase, NetworkUdpDeviceBase及相关的继承类，对所有的泛型进行了擦除，一律采用接口实现，之后将统一java,python代码。 * FreedomTcp,FreedomUdp,FreeSerial: 添加基于自由协议的tcp，udp，串口协议，可以自由配置IByteTransform接口，可用来读取一些不常见的协议。 * Omron-cip: 读写字符串仍然没有测试通过，请暂时不要调用。 * SiemensS7: 单次读取之前是按照200字节进行拆分的，现在根据s7协议返回的报文来自动调整，1200系列是220字节，1500系列是920字节，其他等待测试。 * 官网的备案失效了，重新备案需要点时间，请访问 http://118.24.36.220 然后去顶部的菜单找相应的入口。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-5-6 | 9.2.1 | * Toledo: 托利多电子秤的字节触发的时候，传递出来携带原始的字节数组，方便自行处理，Demo界面优化，显示信息更加完善。 * Lsis: Lsis的PLC通信类修复一些bug，感谢埃及朋友的提供的技术支持。 * MqttSyncClient: 新增ReadString方法，以字符串的形式来和服务器交互，默认编码UTF8，当然也可以自己指定编码，本质还是读取字节数据。 * WebsocketClient: websocket的客户端类，重新设计异常重连，网络异常时触发 OnNetworkError 事件，用户应该捕获事件，然后在事件里重连服务器，直到成功为止。 * MqttClient: Mqtt客户端类，重新设计异常重连，网络异常时触发 OnNetworkError 事件，用户应该捕获事件，然后在事件里重连服务器，直到成功为止。 * MqttSyncClient: 支持读取数据的进度回调功能，支持三种进度报告，数据上传到服务器的进度报告，服务器处理进度报告，数据返回到客户端的进度报告。 * PanasonicMewtocol: 修复注释错误，L区的数据也可以进行L100F，L2.3访问。 * DLT645: 初步添加电力规约协议的串口实现，目前只实现了读取数据，还未测试，等待后续的测试完善。 * Omron-cip: 读写字符串仍然没有测试通过，请暂时不要调用。 * 官网的备案失效了，重新备案需要点时间，请访问 http://118.24.36.220 然后去顶部的菜单找相应的入口。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-5-11 | 9.2.2 | * MqttClient: 上个版本开放的网络错误事件，如果不进行事件绑定，增加默认实现，每隔10秒去连接服务器，直到成功为止。 * WebsocketClient: 上个版本开放的网络错误事件，如果不进行事件绑定，增加默认实现，每隔10秒去连接服务器，直到成功为止。 * DLT645: 电力规约协议完善，等待后续的测试完善。 * SerialBase: ReadBase提供一个重载的方法，ReadBase( byte[] send, bool sendOnly )支持单向发送，不接收数据返回。 * SoftBasic: HexStringToBytes算法优化，性能提升，移除了转大写字母的步骤。 * SiemensS7: 开放获取 pdu 数据长度属性，属性名称：PDULength * HslExtension: 增加IncreaseBy方法，但是测试发现不适用byte类型。 * Omron-cip: 读写字符串仍然没有测试通过，请暂时不要调用。 * 官网的备案失效了，重新备案需要点时间，请访问 http://118.24.36.220 然后去顶部的菜单找相应的入口。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-5-21 | 9.2.3 | * IReadWriteNet接口新增一个属性，ConnectionId，用来表示设备的唯一ID。 * ModbusTcpServer: Modbus的虚拟服务器支持0x16功能码，支持掩码写入操作，适用Tcp,Rtu,Ascii。 * Modbus客户端(tcp+rtu+ascii+rtuovertcp) 新增掩码写入方法，WriteMask，bool写入时，假如Write("100.1", true)就使用掩码写入寄存器100的第1位为真。 * RedisClient: redis的客户端新增Ping方法，DBSize方法获取key数量，FlushDB方法清除数据库所有key。 * DTUServer: 新增一个DTU服务器，可以用来实现对plc的反向连接操作，根据设备的唯一号来识别。 * Omron-cip: 读写字符串不成功的bug修复，已经测试通过。 * WebsocketClient: 实例化时新增url的额外参数传递，("127.0.0.1", 1883, "/A/B?C=123456")，也可以使用"ws://127.0.0.1:1883/A/B?C=123456"。 * WebsocketClient: 修复未连接服务器的时候，调用关闭方法将会引发发送异常的bug。 * MqttServer: 修复NET35版本不支持同步访问的bug，新增一个客户端断开连接的事件，OnClientDisConnected事件。 * VibrationSensor: 新增一个震动传感器的类，型号为苏州捷杰震动传感器VB31，支持获取速度，加速度，位移，温度信息。 * 官网的备案失效了，重新备案需要点时间，请访问 http://118.24.36.220 然后去顶部的菜单找相应的入口。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-5-29 | 9.2.4 | * Mewtocol: 松下的串口协议修复LD寄存器无法访问的bug，输入LD100，如果只输入L100，就是线圈。 * Modbus: 修复写入寄存器指定位bool失败的bug，写入true的掩码改为 FF FE，00 01 * Modbus：在ModbusRtuOverTcp里填写掩码写入的api方法。 * ab-plc：CIP协议解析标签地址的编码从ASCII编码修改为UTF-8编码，支持中文的标签名访问。 * omron-plc：CIP协议解析标签地址的编码从ASCII编码修改为UTF-8编码，支持中文的标签名访问。 * Websocket: 连接的请求标头修改为GET ws://127.0.0.1:8800/ HTTP/1.1 就是带IP地址及端口信息 * Redis：Redis的客户端添加对集合和有序集合操作的相关API方法，基本支持了所有需要的操作信息，单元测试通过。 * Demo: 所有DEMO写入数据操作，新增Hex写入，输入1A 1B等十六进制数据，然后底层调用Write(string, byte[])方法。 * Demo：Redis的功能菜单新增一个测试界面，用来同步两个不同的redis的数据，也可以同一个redis不同的db块数据。 * 官网的备案失效了，重新备案需要点时间，请访问 http://118.24.36.220 然后去顶部的菜单找相应的入口。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-6-10 | 9.2.5 | * CipServer: Cip的虚拟服务器的数据节点编码修改为UTF8编码，增加了一些可读性比较强的增加节点的api，支持赋值初始化数据。 * Demo: Kuka机器人的连接问题，请参考下面地址：http://blog.davidrobot.com/2019/03/hsl\_for\_kuka.html?tdsourcetag=s\_pctim\_aiomsg * Redis: 增加读取TTL的api方法，方便的获取剩余的生存时间。 * HttpServer: 修复Response为空时进行AppendHeader时发生的bug，进行二次校验。 * VibrationSensorClient: 修复deme站号设置失效的bug，站号根据接收的数据动态调整，增加检测长时间未接收传感器数据，就选择重连的功能。 * 官网的备案失效了，重新备案需要点时间，请访问 http://118.24.36.220 然后去顶部的菜单找相应的入口。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-6-28 | 9.2.6 | * NetworkBase: 同步网络通信的超时检查不再开新的线程检查，使用socket自带的ReceiveTimeout来检查。 * NetworkBase: 发送数据时，增加对发送数据的空检查，如果为空，就认为成功。 * RedisClient: 新增修改密码的API接口，可以进行对redis的密码重置操作。 * MqttServer: 当同步客户端 MqttSyncClient连接上来时，不进行触发上下线事件。 * MqttServer：原先支持获取所有的在线客户端，现在新增获取异步客户端列表，获取同步客户端列表。 * MqttSubscribeMessage: 类型拼写错误修复，如果使用这个类，请谨慎升级。 * Keyence: 基恩士的MC协议，支持CC，TC的数据类型读取。 * FanucSeries0i: 新增一个fanuc机床的数据通讯类，支持读取一些简单的数据，目前在Series0i-F上测试通过。 * Cip: 修复ab-plc的标签地址解析为UTF-8编码，但是长度确实字符串的bug，现在支持中文编码。 * 其他的注释优化 * 官网的备案失效了，重新备案需要点时间，请访问 http://118.24.36.220 然后去顶部的菜单找相应的入口。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。企业终身授权费：8000元(不含税)。 |
| 2020-7-8 | 9.2.7 | * MqttServer: 修复MQTT服务器开启时，当用其他的mqtt客户端订阅时，会发生异常的bug，原因在于订阅质量没有回传。 * WebsocketServer: websocket的服务器端新增一个客户端下线的事件，无论是正常关闭还是异常关闭，都会触发事件。 * MqttClient: Mqtt的客户端新增一个连接成功的事件OnClientConnected，重连成功后也会触发。在该事件的订阅topic会在网络恢复后重新订阅。 * NetworkDoubleBase: 当校验指令头失败的时候，返回的错误信息里追加，收发的报文，方便查找问题。 * MelsecA1EAsciiNet: 修复读取bool时，长度为奇数时，会出现交替失败的bug，原因出自数据粘包。 * WebsocketClient: 添加一个IsClosed属性，修复服务器强制断线导致客户端无限重连的bug。 * OmronConnectedCipNet: 添加一个基于连接的CIP的读写类，等待测试。 * 其他的注释优化 * 官网的备案失效了，重新备案需要点时间，请访问 http://118.24.36.220 然后去顶部的菜单找相应的入口。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-7-20 | 9.2.8 | * KeyenceNanoSerial: 基恩士的串口协议重新实现，实现IReadWrite接口，增加了单元测试。支持的地址需要查阅API文档信息。 * OmronHostLinkCMode: 支持了欧姆龙的HOSTLINK协议的Cmode模式的实现，初步单元测试通过，等待测试。 * MC协议：三菱MC协议的ZR区的地址进制从16进制改为10进制。 * NetworkDoubleBase: 添加一个PING的方法IpAddressPing( ), 对设备当前的IP地址进行PING操作。 * NetworkUdpBase: 添加一个PING的方法IpAddressPing( ), 对设备当前的IP地址进行PING操作。 * yamaha: 添加一个雅马哈机器人协议的实现，初步实现了几个api，等待测试，测试通过继续完善。 * DEMO: 主界面增加一个全国使用情况的分布图，统计DEMO的使用次数实现。 * 官网的备案失效了，重新备案需要点时间，请访问 http://118.24.36.220 然后去顶部的菜单找相应的入口。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-8-3 | 9.3.0 | * Networkbase：核心网络底层的错误码调整，当读写操作因为网络问题失败时，返回错误码为负数-1，如果连续读写失败，就一直递减。 * OmronConnectedCipNet: 地址解析修改为全部上29 00 报文。 * FileNet: 两种文件服务器支持删除多个文件和删除文件夹的所有文件功能，客户端同步适配，初步测试通过。 * NetSimplifyClient: 新增一个构造方法，可以传入IPAddress类型的ip地址。 * MqttSyncClient: 新增一个构造方法，可以传入IPAddress类型的ip地址。 * MqttClient: 修复一个连接反馈信号，解析判断服务器状态错误的bug，该bug导致MqttClient连接不是中国移动的OneNet物联网框架。 * FFT: 傅立叶变换FFTValue方法添加一个可选参数，是否二次开放，波形中的毛刺频段会更加明显。 * HttpServer: webapi的服务器完善注释，添加一个端口号的属性，获取当前配置端口号信息。 * Active: 当前库激活失效的时候，返回的错误消息，携带当前的通信对象的实例化个数，方便查找授权失败的原因。 * Abb机器人：abb机器人支持读取程序执行状态，任务列表功能，伺服状态，机器人位置数据。 * ABB虚拟机器人：新增一个abb机器人的虚拟webapi的服务器，可以用来测试和ABB客户端的通信。 * Demo: 数据转换的界面，新增一个显示指定的文件的二进制的内容的功能。当demo激活成功时，不显示时间及授权信息。 * 新增一篇全新的博文，介绍基于HSL的大一统网络架构实现，满足发布订阅，一对多通信，webapi等：https://www.cnblogs.com/dathlin/p/13416030.html。 * 官网备案成功了，地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-8-28 | 9.3.1 | * Beckhoff: 倍福PLC新增读取设备信息和设备状态的api接口。在demo界面添加测试按钮，状态码检查优化，错误时返回报文信息。 * FanucSeries0i: 修复fanuc机床的读取宏变量解析double数据时，为0的时候解析异常的bug。 * ABBWebApiServer：ABB机器人的虚拟服务器支持用户名和密码设置，在客户端请求数据的时间，支持账户验证。 * Demoserver: 优化根据IP地址获取物理地址的方法，获取不到或是奇怪字符将切换线路重新获取。 * KukaTcpNet: 新增KukaTcp通讯类，支持多变量写入的api，在demo界面增加启动，复位，停止程序的操作。 * .Net Framwork 2.0 支持2.0的框架的dll发布，通过nuget安装即可。 * SimpleHybirdLock: 简单混合锁添加一个当前进入锁的次数的静态属性，可以查看当前共有多少锁，等待多少锁。 * NetworkDeviceBase: 核心交互方便增加错误捕获，异常释放锁，再throw, YamahaRCX类完善异步方法 * NetworkBase: 增加一个线程检查超时的次数统计功能。 * InovanceH3U: 修复汇川的3U的PLC地址类型为SM,SD时解析异常的bug。 * Demo: HslCommunication Test Demo支持PLC及一些连接对象的参数保存功能，使用英文冒号可以分类管理。 * WebSocketSession: 新增url属性，如果客户端请求包含url信息，例如：ws://127.0.0.1:1883/A/B?C=123, 那么url就是这个值。 * Demo: 测试的DEMO程序，支持连接参数存储，不用再每次打开程序重新输入IP地址，端口，站号等等信息，可以存储起来，还支持分类存储。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 * HSL的目标是打造成工业互联网的利器，工业大数据的基础，打造边缘计算平台。 |
| 2020-9-27 | 9.3.2 | * KeyenceNanoSerial: 修复读写R寄存器时，提示地址格式异常的BUG，已经测试通过。 * MelsecMcUdpServer: 新增三菱MC协议的UDP虚拟PLC，支持数据读写，支持二进制和ASCII格式。 * OmronFinsUdpServer: 新增欧姆龙Fins协议的UDP的虚拟PLC，支持数据读写操作。 * MqttServer: 修复MQTT服务器在客户端发送批量订阅的时候，服务器会触发BUG的问题。 * ConnectPool<TConnector>类代码注释优化，新增连接次数峰值属性。 * RedisSubscribe: 订阅服务器重新设计，订阅实现事件触发，支持手动订阅，取消订阅操作。 * RedisClient: 支持了订阅的操作，当订阅的时候，创建订阅的实例化对象，应该在连接参数设置之后再进行订阅。 * RedisClientPool：新增Redis连接池类，默认不限制连接数量，使用起来和普通的RedisClient一样，适合一个项目实例化一个对象。 * MqttSyncClientPool: 新增MqttSyncClient的连接池版本类，默认不限制连接数量，用起来和普通的MqttSyncClient一样。 * LogNetFileSize: 根据文件大小的日志类，实例化时支持设置允许存在的文件上限，如果设置为10，只保留最新的10个日志文件。 * LogNetDateTime: 根据日期的日志类，实例化时支持设置允许存在的文件上限，如果设置为按天存储，上限为10，就是保留10天的日志。 * AllenBradleySLCNet: 新增AB PLC的数据访问类，适合比较老的AB PLC，测试通过的是1747系列。地址格式为A9:0 * AllenBradleyNet: 读写bool值的时候，不带下标访问单bool数据，如果需要访问bool数组，就需要带下标访问，例如：A[0]。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2020-10-23 | 9.5.0 | * AllenBradleyNet: ReadBool方法默认读取单bool类型变量，如果要读取int类型的bool数组，采用"i="开头表示，例如"i=A[10]" * NetworkDataServerBase: 新增一个属性ActiveTimeSpan，可以设置激活时间间隔，默认24小时，锁优化，其他的继承实现的服务器都进行了设置。 * NetworkDeviceBase: Read<T>修改为虚方法，支持继承进行重写，基于特性的类注释完善。 * Siemenss7net: ReadString(string address, ushort length)读取字符串时，如果长度为0，就读取西门子格式的字符串。 * OperateResult: 扩充泛型方法，Check, Convert, Then，实现了结果链，简化代码。参考：https://www.cnblogs.com/dathlin/p/13863115.html * FanucSeries0i: 修复数控机床在读取0i-mf状态时导致长度不够的bug。 * IReadWriteNet: 新增wait方法接口，用于等待一些信号到达指定的值，支持扫描频率设置，超时设置。例如 Wait("M100.0", true, 500, 10000)等待这个信号为true为止。 * MqttServer: 支持调用ReportOperateResult返回错误信息及错误码给客户端，MqttSyncClient会自动识别报文，然后IsSuccess自动适应，网络不会断开。 * MqttSyncClient: 支持设置接收超时时间，默认是60秒，之前是5秒，而且不能更改。 * MqttServer: 支持注册远程RPC的API接口，自动解析json参数，自动调用已经注册的接口自动返回是否成功，MqttSyncClient也支持遍历服务器的接口列表。详细：https://www.cnblogs.com/dathlin/p/13864866.html * SiemensS7Net: 通信类实现ReadBool("M100", 10); 批量读bool方法，通过读Byte间接实现。 * OmronHostLinkCModeOverTcp: 新增欧姆龙的通讯类，Cmode模式的以太网透传实现。 * PLC: 所以的PLC实现了HslMqttApi特性支持，从而在MqttServer里可以直接注册，然后对外开放读写接口操作。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2020-11-2 | 9.5.1 | * MQTTServer: 服务器端在注册M-RPC远程调用接口服务的时候，如果你的接口实现了参数默认值，那么就提取到示例参数里。 * KeyenceSR2000SeriesTcp: 新增基恩士的SR2000系列的扫码器的驱动，目前是TCP版本，支持读码，读取记录，停止复位等，支持自定义的命令。 * XGKFastEnet: 新增Lsis的XGK系列的FastEnet实现。 * XGKCnet: 新增Lsis的XGK的cnet的实现。 * Demo: Demo的TCP调试的服务器端优化，错误获取优化，发送数据失败的问题修复。 * NetworkBase: 底层异步的数据接收的超时优化，优化超时线程池实现，更加节省线程调度。 * MqttSyncClient: 客户端支持读取MQTT服务器的驻留主题列表，读取该主题的相关的数据信息。详细见demo。 * MqttServer: 修改ClientVerification事件，增加会话句柄传递，支持动态修改client id，支持设置当前客户端禁止发布任何数据。 * MqttServer上的MRPC的权限控制仅对商业授权用户开放，MqttSyncClientPool连接池以及RedisClientPool连接池仅对商业授权用户开放。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2020-11-11 | 9.5.2 | * XGKFastEnet修复了一点小问题，感谢埃及朋友提供的技术支持。 * KeyenceSR2000Serial: 新增基恩士的SR2000系列的扫码器的串口版驱动，支持读码，读取记录，停止复位等，支持自定义的命令。 * IKeyenceSR2000Series: 新增基恩士的SR2000系列的扫码器接口，优化读码功能，如果读码失败，则发送LOFF命令。 * LogStatistics: 新增统计次数的辅助类，可以实现关于时间的一些统计信息，比如统计每小时，每天，每周，每月等的登录量，使用量信息。 * LogValueLimit: 新增数据极值类，用于统计一个数的开始值，结束值，最大值，最小值，平均值，然后可以按小时，天，周，月等统计。 * OmronHostLink: 修复解析错误码时，如果错误码不为数字的时候会导致奔溃的bug，错误提供内容更加详细。 * ILogNet: 日志记录的方法实现MqttApi特性，可以在MQTTServe里面注册为RPC服务，从而实现远程调用日志写入的方法。 * ILogNet: 日志类方法新增属性LogStatistics，只要实例化就可以统计当前的日志记录情况，可以每分钟，每小时，每天，每周，每月，每季度，每年。 * MqttClient: 新增一个订阅的api接口，支持直接传递MqttSubscribeMessage对象，可以指定消息质量。 * XinJEXCSerial: 新增信捷的XC系列的串口通讯类，底层是modbus-rtu，地址做了封装，按照信捷的地址输入即可，比如X1,Y7,M1000,D100,F100,E100。 * MelsecFxSerial: 三菱编程口类新增 IsNewVersion 属性，如果为false，就是老版本的协议，修复T,C线圈读写的地址不对BUG。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2020-12-2 | 9.5.3 | * ModbusRtuOverTcp: 更改继承，直接从NetworkDeviceBase进行继承，通过单元测试 * YokogawaLinkTcp: 新增横河PLC的二进制通讯类，支持X,Y,I,E,M,T,C,L,D,B,F,R,V,Z,W,TN,CN读写，部分高级API商业授权用户才能使用，例如读取PLC信息。 * YokogawaLinkServer: 新增横河PLC的二进制格式的虚拟PLC，模拟的真实的PLC的通信机制，实现了读写长度的限制，以及错误信号的返回。 * Networkdoublebase: ReadFromCoreServer( byte[] send, bool hasResponseData ) 新增是否等待数据返回的属性，可以用于某些不需要数据返回的命令。 * Networkbase： 修复异步接收数据时，某些情况下长度为0导致连接关闭的bug。 * FetchWriteServer: 新增西门子fetch/Write协议的虚拟PLC，支持虚拟数据的读写，通信。 * MelsecFxSerialOverTcp: 修改继承体系，从NetworkDeviceBase继承，和MelsecFxSerial的IsStringReverseByteWord调整为true; * 文件引擎服务器修复路径大小写导致的bug问题，文件客户端支持检查文件是否存在的方法，检查文件是否存在。 * MqttServer: 远程调用的MRPC的参数支持自定义类型，通过JSON转换，将字符串转换为实体类。还有其他的优化。 * DeltaDvpSerial, DeltaDvpSerialAscii, DeltaDvpTcpNet: 添加台达的通信类，输入台达的地址即可，会自动转换实际的modbus地址。 * 所有的虚拟PLC的服务器均调整为商业授权用户专享，还有一些高级的API，具体看api注释是否带有[商业授权]字样，基本的数据读写功能将一直对个人用户开放。 * Demo: 数据读写示例的界面，写入现在支持批量写入，数据写[1,2,3]，然后写入short，就是写入short数组了。 * 普通VIP的个人使用不再限制100个PLC对象，连续运行时间调整为10年，高级的一些API限制商用，参考注释是否带[商业授权]字样。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2020-12-22 | 9.6.0 | * YokogawaLink: 异步的方法添加和完善，虚拟PLC侧支持Special:100开头的地址，表示特殊模块寄存器，从而支持各种类型的读写。 * LogStatistics, LogValueLimit: 两个数据日志分析类支持获取指定时间段的数据，数据存储文件格式重新设计，为不兼容更新。 * LogStatisticsDict, LogValueLimitDict: 新增数据日志分析的词典类，用来统计多个数据的不同时间段的使用情况。 * 所有的基于tcp的plc，机器人，redis, mqtt, websocket等通讯类的ip地址支持直接输入域名，会自动调用Dns.GetIpAddress来解析。 * MqttRpcApi: MRpc的API特性支持应用在属性上，不需要传递参数，直接获取属性的值，在demo上显示的小图标不一样。PLC的通讯类的基本属性在MRPC公开。 * MelsecFxLinks: 支持读取PLC的型号，读写数据的地址支持了站号指定，地址可以写成[s=2;D100]，方便多站号读取。 * AllenBradleyNet: 地址支持slot参数，例如：slot=2;AAA ，也可以不携带，这个是可选的。 * FatekProgram, FujiSPB, XGBCnet, MelsecA3CNet1, OmronHostLink, OmronHostLinkCMode, PanasonicMewtocol, SiemensPPI，信捷，汇川类及其透传类支持地址携带站号，例如 s=2;D100 * FujiSPBServer: 新增富士PLC的虚拟服务器，支持串口和网口，原先的富士PLC存在bug，不能读取，欢迎网友对富士PLC测试。 * HttpServer: 删除原先的HttpGet和HttpPost特性，改用MRPC的特性，支持注册webapi服务，使用方式和MRPC类似，demo增加httpsclient可浏览接口，https://www.cnblogs.com/dathlin/p/14170802.html * HttpServer, MqttServer: 服务器端支持接口调用的次数统计，支持客户端查询接口调用情况，demo客户端实现mqttclient,方便服务器管理在线客户端信息。 * 其他优化改进，如果有网友发现bug，配合作者测试并修复bug，将根据实际情况给与现金红包奖励。 * 普通VIP的个人使用不再限制100个PLC对象，连续运行时间调整为10年，高级的一些API限制商用，参考注释是否带[商业授权]字样。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-1-12 | 9.6.1 | * Lsis XGK: 修复部分的问题，感谢埃及朋友提供的支持。 * FujiSPB: 修复未在net20, net35, net standard项目里添加的bug。 * MqttServer和HttpServer: 注册API的方法支持对静态方法的注册，注册时传入类型对象即可。 * Modbus: tcp, rtu, ascii, rtu over tcp在读写int,uint,float,double,long,ulong时支持动态指定dataformat，地址示例：format=BADC;100 * MqttServer: 扩展MQTT的子协议FILE，支持文件的上传，下载，删除，查看信息，权限控制操作，支持获取上传下载网速监控。 * MqttSyncClient: 扩展文件的方法接口，支持上传，下载，删除，遍历文件操作，每个操作都是短连接的，使用的全新的socket对象。 * SiemensS7Net: 修复西门子s7协议某些情况数据批量写入失败的bug，原因来自PDU长度信息不对。 * DLT645: 修复一些问题，已经测试通过，新增 DLT645OverTcp，感谢 QQ：542023033 提供的技术支持。 * FanucInterface: 机器人的解析数据时，当shift\_jis编码不存在时，将会引发异常，现在自动替换UTF8 * HslCommunication: 所有的异步通信代码优化，优化超时检测机制，现在大大提升了服务器的高并发的能力，异步通信的性能。 * AllenBradleyNet及OmronCipNet协议支持 UINT, UDINT, ULING类型的写入，对应的C#的类型是 ushort, uint, ulong * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-1-26 | 9.6.2 | * OmronConnectedCipNet: 新增欧姆龙的基于连接的CIP实现，测试读写欧姆龙PLC成功，支持数组读写，详细参考API文档，欧姆龙的CIP请使用本类。 * AllenBradleyNet: 修复当自动重连时，连接的ID 还是上次的 ID 导致读写失败的bug。 * DLT645: 点对点模式下，在读取地址域的时候，新增读取成功后即修改当前的地址域信息，也即是在打开串口后，读取地址域即可通信。 * DLT645: 修复有些数据格式不一致导致读取数据不正常的bug，已经测试可以读取功率，能耗，电压电流，电表基本信息，还支持自定义的解析格式。 * NetworkAlienClient: DTU客户端增加对连接客户端的注册包的数据校验，修复数据意外的情况导致程序奔溃的bug。 * Demo: 在 演示程序里，Modbus的DTU的示例界面，修复 ID 设置时，结果设置到 IP 导致异常的bug。另外增加西门子的DTU演示界面。 * LSisServer: 修复同一地址，数据读写不对的bug，和 XGKFastEnet 客户端读写测试通过，包括bool类型地址，字地址 * GeSRTPNet: 新增 GE-PLC（通用电气） 的SRTP协议实现的客户端，支持I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作。 * GeSRTPServer: 新增 GE 的 SRTP 协议的虚拟PLC，支持和 GeSRTPNet 通信测试。支持类型和客户端支持的一致。 * MqttServer: 在启动文件服务功能时，增加对分类路径，以及文件名的合法性进行校验，防止注入特殊字符攻击及意外bug。 * MqttSession: 新增一个方法，GetTopics() 用于获取当前的会话对象所订阅的主题的副本数据。 * PanasonicMewtocol: 修复 Mewtocol及串口转网口类，在批量读取 bool 数组地址解析不准确的bug。 * MelsecCipNet: 新增三菱的CIP协议功能，PLC使用了 QJ71EIP71 模块时就需要使用本类来访问。 * SickIcrTcpServer: 修复当关闭服务器的时候，现有的连接没有关闭的bug，没有关闭的话，仍然会接收到来自设备发来的条码信息。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-2-15 | 9.6.3 | * SickIcrTcpServer: 修复手动连接扫码设备的网络在关闭服务器后仍然会自动重连的bug。 * SoftBasic: 删除SpliceTwoByteArray及SpliceByteArray方法，改为泛型支持的SpliceArray方法，支持任意类型拼接，添加了扩展方法支持。 * Modbus: 支持 0x16 功能码，用于掩码操作，支持对寄存器地址的位操作，需要设备方支持，该功能仅支持商业授权使用。 * Modbus: 读取线圈和输入线圈的长度支持任意，内部按照2000长度自动切割，读取寄存器和输入寄存器按照120自动切割，该功能商业授权特权，普通的VIP用户存在长度限制。 * MqttSyncClient: 新增ReadRpc<T>(string topic, string payload )方法，专门用来读取注册的RPC接口的，自动json转换类型。 * MqttSyncClientPool: 连接池优化，注释优化，添加了一些缺失的方法。该功能商业授权特权。 * RedisClientPool: 连接池优化，注释优化。该功能商业授权特权。 * LogNet: 日志部分新增一个 ConsoleOutput 属性，如果设置为 true，那么日志就会在控制台进行输出，等级不一样的日志，文字颜色不一样。 * LogNet: 日志部分的记录优化调整，取消了一些底层的重复记录的日志内容，针对 MQTT, Websocket, HTTP 及虚拟PLC相关的日志记录根据信息进行优化。 * 祝大家新年生意滚滚，身体健康，牛年大吉。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-3-11 | 9.6.4 | * ILogNet: 当ConsoleOutput为true时，修复在空字符串的存储路径时，不进行控制台输出的bug。 * SiemensS7Server: 优化握手报文的信息，现在支持sharp7进行数据通信了。 * SiemensS7Net: 开放了 ConnectionType 和 LocalTSAP 属性，方便按照自己的需求修改，具体看属性注释。 * SiemensS7Net: 西门子的PLC读取支持0x0A和0x06的错误码，当读取DB块不存在时，提示错误消息。 * SiemensS7Net: 支持了 WString 类型的读写，使用ReadWString和WriteWString方法，支持中文的读写 * SiemensS7Net: 西门子S7协议的地址解析，DB块地址优化，DB1.DBW1 的 DBX,DBB,DBW,DBD都会自动屏蔽。 * SiemensS7Net: 在核心的报文交互上，自动忽略只有7字节的TPKT和ISO的报文的情况。 * MqttServer, HttpServer:RPC注册的方法，原先只支持一个泛型的结果类 OperateResult<T>， 现在支持任意个泛型的结果类对象。 * FanucSeries0i: 所有的方法实现异步接口，并增加了 RPC 的特性支持，方便直接注册就可以调用。 * SimensFetchWriteServer: 修复在standard项目里没有添加的bug。 * MelsecMcNet: 修复ReadMemory的报文错误，增加读取智能模块的ReadSmartModule方法。 * NetworkDataServerBase: 所有的虚拟PLC的基类，添加获取所在在线客户端信息的属性 GetOnlineSessions * Demo: 所有读写PLC的界面在读写的时候，增加提示耗时的信息，包括最大值，最小值。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-3-15 | 9.6.5 | * SiemensS7Net: 修复上个版本的DB块位地址解析的bug，写入DB1.0.5为True的时候，却写入了DB1.0.0为True。 * AllenBradleyDF1Serial: 初步添加AB-PLC的DF1协议，支持了简单的读写，等待测试，地址示例：N7:0。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-3-25 | 9.6.6 | * SoftCRC16: 计算CRC16的辅助方法，开放预置值的设定，可以自由的指定。 * FanucSeries0i: 新增一个读取机床系统语言的api，读取之后将会自动切换语言，暂不支持根据消息自动匹配编码解析。 * SiemensS7Net: OperateResult<byte[]> Read( string[] address, ushort[] length )接口添加RPC支持 * NetworkDataServerBase: OnDataReceived事件签名修改为DataReceivedDelegate( object sender, object source, byte[] data )，追加一个source参数，可用来获取客户端IP地址，具体看api文档 * NetworkDoubleBase: 增加LocalBinding属性，如果需要绑定本地ip或是端口的，可以设置，所有的网络类PLC都支持绑定本地的ip端口操作了。 * NetworkUdpBase: 增加LocalBinding属性，如果需要绑定本地ip或是端口的，可以设置，所有的网络类PLC都支持绑定本地的ip端口操作了。 * SiemensS7Net: 完善异步的PDU自动长度信息，新增AI,AQ地址的读写，地址格式：AI0,AQ0，欢迎大家测试。 * OmronFinsNet: 欧姆龙FINSTCP协议的SA1机制调整为自动获取，不需要在手动设置，修复错误信息文本和错误码不匹配的bug。 * MqttClient: 修复在网络异常导致正在重连服务器的时候，调用ConnectClose方法后，后台仍然不停的重连服务器的BUG。 * NetworkDeviceSoloBase: 删除这个文件，并优化相关的串口透传类。全部改为继承自：NetworkDeviceBase * NetworkDataServerBase: 所有派生类的虚拟服务器，包括modbus，s7, mc, fins等服务器全部支持设置是否允许远程写入操作，modbus的demo界面添加是否允许的选项。 * WebSocketClient: 修复客户量的Request报文少一个换行信号在某些服务器会连接失败的bug，新增两个发送数据的api，发送数据更加的灵活。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-4-15 | 9.7.0 | * OmronFinsNet, OmronFinsUdp, HostLink: 地址的解析优化，在读取的API方法里，自动按照500长度进行切割，可以由 ReadSplits 更改。 * FanucSeriesOi: Fanuc数控机床支持R数据的读取，参考API： ReadRData(int start, int end) * HslExtension: 新增从字节数组获取位值的扩展方法: GetBoolValue( this byte[] bytes, int bytIndex, int boolIndex ) * FatekProgram: 地址解析优化，修复 RT,RC地址解析不正确的bug，读取的字及位长度自动切割，调用不受长度限制。 * SoftBasic: 添加设置byte数据的某个位的方法SetBoolOnByteIndex，也可以调用byte的扩展方法，byte.SetBoolByIndex(2, true) 就是设置第二位为true * FujiSPHNet: 新增支持富士的SPH以太网协议，支持M1.0, M3.0, M10.0, I0, Q0地址的读写操作，支持位的读写操作。写位需要谨慎，先读字，修改位，再写入。 * net20, net35, net451三个框架版本的项目引用 http.web 组件，用来修复 HttpServer 里url携带中文时，会导致解析乱码的情况，现在支持了中文的api接口注册，中文参数。 * HttpServer: 使用了注册RPC接口时，返回调用方的数据内容格式调整为json格式，方便postman等测试工具识别内容。 * FujiSPHServer: 新增富士SPH协议的虚拟服务器，支持和FujiSPHNet进行测试通信。支持的地址是一致的。 * KeyenceNanoSerial: 基恩士的上位链路协议优化，支持了B，VB的bool读写，W，VM的字读写，新增bool数组写入功能。 * KeyenceNanoSerial: 支持了plc型号读取，状态读取，注释读取，扩展缓存器的读写，错误代码提示携带更详细文本，适用于 KeyenceNanoSerialOverTcp * KeyenceNanoServer: 新增基恩士上位链路协议的虚拟服务器，可以和 KeyenceNanoSerialOverTcp 进行通信测试。 * KeyenceSR2000: 基恩士扫描的协议的错误提示信息新增了英文模式下的注释，原来的只有中文的提示。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-5-13 | 9.8.0 | * MqttSyncClient: ReadRpc<T>方法提示JSON转换错误消息更加友好，注释完善。 * MelsecMcServer: 三菱的虚拟plc新增对B，R，ZR的地址支持。 * FujiSPHNet: 富士PLC实现按照240长度自动切割，支持无限长的数据读取。 * MqttClient(不兼容更新): 接收服务器数据的事件签名修改，MqttMessageReceiveDelegate( MqttClient client, string topic, byte[] payload ); * WebSocketServer: 优化数据发送部分的功能逻辑代码，将数据发送从锁中解脱出来。 * MqttClient: 新增属性 ConnectionOptions 用来获取当前连接参数信息。 * Omron: Hostlink协议的读取字数据时，长度进行切割，按照260字切割，可通过ReadSplits属性修改。 * PanasonicMewtocolServer: 初步添加mewtocol的虚拟plc，初步测试R地址成功。 * 设备通信核心: tcp, udp, 串口三大通信内核添加，封包和解包的虚方法，可以重写实现自定义需求。 * IModbus: 新增IModbus的设备接口，用来描述Modbus相关的设备，包含站号，DataFormat属性等。 * Modbus: 包含TCP,RTU,ASCII,RTU-over-tcp，UDP全部结构优化，重写，完善，最终一套代码实现覆盖以上类，接口无变化。   但是如果用调用了ReadFromCoreServer则不兼容，现在都只需要传核心报文，01 03 00 00 00 01，无论rtu,tcp,ascii * Modbus: 支持对寄存器，输入寄存器的位数据读取，ReadBool("100.1") 就是读取寄存器地址100的第一个位的bool值。 * MqttSubscribeMessage: Identifier默认设置为1，这样可以修复在某些服务器（mosquitto）订阅异常的bug。 * SerialBase(不兼容更新): ReadBase串口基本的交互方法重命名为ReadFromCoreServer，这样与TCP，及UDP的方法标准一致。 * FanucInterfaceNet: 支持 R 地址的读写，支持R1-R10，其中R1-R5为int数据，R6-R10为float数据。SR1-SR6进行字符串读写。 * BeckhoffAdsServer: 新增ADS的虚拟PLC，支持M100, I100, Q100地址格式。暂不支持内存地址，变量名。 * SiemensS7Net: 修复9.6.4版本添加的ConnectionType, LocalTSAP属性对 200，200smart型号的影响。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-5-23 | 10.0.0 | * V10版本：本来上个版本就应该定为V10版本，因为已经内核优化，出现不兼容更新，所以这个版本果断设定为V10版本。 * SerialBase: 修复在串口异常的情况下，会抛出异常的bug，是上个版本的问题，现在会返回失败的OperateResult对象。 * ShineInLightSourceController: 新增昱行智造科技（深圳）有限公司的光源控制器的通信对象，主要用于视觉的打光操作。 * MqttSession：Mqtt的会话信息增加一个object Tag属性，用来自己绑定一些自定义的数据。 * SerialBase: 串口初始化的方法修改为虚方法，允许在继承类里进行重写，修改一些默认参数信息。 * NetworkBase: 修复ReceiveAsync异步方法在length=-1时，对方关闭时返回仍然为成功的bug，只有在极少数情况下会触发。 * ModbusTcpServer: 新增一个属性UseModbusRtuOverTcp，只要设置为True，就可以创建ModbusRtuOverTcp的对应的服务器，使用TCP传送RTU报文。 * HttpServer: 新增SetLoginAccessControl( MqttCredential[] credentials )方法，用于增加默认的账户控制，如果传入null，则不启动账户控制。 * IReadWriteDevice: 新增设备读写接口，继承自IReadWriteNet，然后所有设备实现IReadWriteDevice接口，相关继承关系优化，接口增加ReadFromCoreServer。 * All: 统一所有的设备核心层打包报文方法名为:PackCommandWithHeader 解包的方法名为UnpackResponseContent，允许重写实现自定义操作。 * Omron: 对OmronFinsTcp和OmronFinsUdp的通信层大幅度优化，统一代码规则，新增run，stop，读取cpu数据，cpu状态的高级方法。 * DTSU6606Serial: 新增德力西电表的采集类，基于modbusrtu实现，ReadElectricalParameters方法可以直接获取电表相关参数。 * HslExtension: 有两个获取byte的位的方法，功能重复，删除GetBoolOnIndex方法，使用GetBoolByIndex方法。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-6-11 | 10.0.1 | * FatekProgram: 串口类和串口转网口透传类优化，统一一套代码来读写设备。 * IDisposable: NetworkAlienClient, NetworkAlienClient, LogNetBase, MqttClient, MqttServer, WebSocketClient, WebSocketServer实现释放接口。 * SiemensS7net: 新增DestTSAP属性，优化了LocalTSAP和DestTSAP属性对不同系列plc的值设置，当plc为s200系列时，支持设置自定义的值来访问plc。 * UltimateFileServer: 文件服务器删除目录所有文件调整为直接删除整个目录，新增支持删除指定目录下所有空的子目录的功能。文件客户端新增匹配操作的方法。 * PanasonicMcNet: 地址新增支持SD数据类型，示例SD0，返回的错误代码修改为松下的专用信息，和三菱的不一致。 * IModbus: Modbus接口新增TranlateToModbusAddress( string, byte) 接口，只要继承重写该方法，即可轻松实现自定义地址解析转modbus地址。 * Delta: 台达相关的类根据modbus最新的优化，全部进行优化，每个类只有一点点代码了。 * FujiSPB: 富士的串口协议代码和串口透传代码优化，修复串口类调用异步写bool失败的bug。 * XinJE: XinJEXCSerial重命名为 XinJESerial类，根据modbus的优化进行精简，支持了信捷系列选择，可选XC,XJ,XD，地址支持根据所选型号自动解析。 * XinJE: 新增基于串口透传的XinJESerialOverTcp类，以及modbustcp协议的XinJETcpNet类，DEMO上支持测试。 * Inovance: 汇川的类优化，删除原来的AM,H3U,H5U类，改用InovanceSeries枚举来区分系列，然后解析不同的地址。同时添加InovanceSerialOverTcp串口转网口类。 * OmronFinsServer: 欧姆龙的FinsTCP虚拟服务器端支持E数据块，E0.0-E31.0 都是指同一个数据块。 * IByteTransform: 新增二维数组的解析方法接口，主要是short,ushort,int,uint,long,ulong,float,double类型。 * Demo: MelsecSerialOverTcp的demo界面添加是否新版的选择。 * 如果有用到汇川，信捷的类库，请注意升级时出现不兼容，需要修改下类型，指定PLC的系列，感谢支持。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-6-22 | 10.0.2 | * OmronFinsServer: 修复服务器的内存数据保存本地文件及加载文件时，不支持em区数据的bug。 * ISessionContext: 新增会话信息接口，在MqttServer和Httpserver中，在注册RPC时可以在方法参数里追加ISessionContext接口的上下文信息，用来控制当前api的对不同账户的权限。 * Modbus: TranslateToModbusAddress 单词拼写错误的修复。 * FujiSPBAddress: 地址类的继承改成 DeviceAddressDataBase * ModbusHelper: 在所有modbus及派生类里，当实现地址转换后，修复写bool,bool[]时地址仍然不转换的bug。 * KeyenceNano: 新增 UseStation 属性，用来设定是否开启使用站号的报文功能，有些特殊的情况需要站号。 * KeyenceNano: 串口类和串口转网口透传类优化，统一一套代码来读写设备。 * 其他的注释优化，代码优化 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-7-14 | 10.1.0 | * Melsec: 三菱的MC协议TCP，UDP，二进制，ASCII代码优化，一套代码实现，新增IReadWriteMc接口，针对MC协议的设备通用读写类 * NetworkUdpBase: 新增LogMsgFormatBinary属性，可以指示当前的数据交互报文记录时按照ASCII编码显示，或是二进制显示。 * DTSU6606Serial: 德力西的电表的读取方法ReadElectricalParameters支持HslMqttApi特性，方便MRPC及WEBAPI接口服务注册。 * Demo: Modbus rtu的demo界面的报文读取取消crc封装，因为内部已经集成封装。 * Melsec: 统一 MelsecMcNet, MelsecMcUdp, MelsecMcAsciiNet, MelsecMcAsciiUdp, MelsecMcRNet的代码逻辑结构，修复了ASCII格式类的一些bug。 * MelsecMcServer: 三菱的虚拟服务器限制了bool读取长度7168限制字读取长度960，三菱MC客户端的bool读取支持自动切割。 * OmronHostLink: OmronHostLink及OmronHostLinkOverTcp代码优化，完善错误代码文本提示，增加返回命令和发送命令校验的操作。 * HslExtension: ToStringArray的扩展方法支持对GUID的解析功能，不支持.net20, .net35 * NetworkDataServerBase: 修复数据类服务器在主动关闭引擎时，在线客户端的数量未及时复原的bug，影响范围，所有的虚拟PLC服务器。 * OmronHostLinkServer: 新增欧姆龙HostLink协议的虚拟PLC，支持网口和串口的进行读写操作。优化hostlink协议的客户端错误代码含义展示，优化数据接收机制。 * Demo: httpclient界面支持对https接口测试，在内容请求的header支持添加content-type信息，提供了一些选项。 * SimensWebApi: 新增NetworkWebApiDevice设备类，实现IReadWritteNet接口，新增SimensWebApi类，用于西门子1500的webapi接口，可实现读写标签变量信息。 * FanucSeries0i: fanuc的通信类支持NC程序文件的上传和下载，删除，设置主程序号，启动加工操作。修复刀具信息读取时，某个刀具信息失败导致读取失败的bug。 * AllenBradleyServer: ab-plc的虚拟服务器支持会话id的生成，支持对客户端校验会话id是否一致。 * Melsec: MC协议的类支持对字地址按照位读取，例如读取D100.5 开始的3个位，使用ReadBool("D100.5", 3)即可 * NetworkBase: 优化ReceiveByMessage及异步版本的性能，减少一次内容数据的拷贝操作，提升内存利用效率，提升读写的性能。 * SiemensS7Net: 西门子s7协议的地址支持 VB100, VW100, VD100的写法。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-7-18 | 10.1.1 | * UltimateFileServer: 支持获取指定目录的所有文件大小，客户端IntegrationFileClient实现相应的调用方法GetTotalFileSize * MqttServer: 支持获取指定目录的所有文件大小，客户端MqttSyncClient实现相应的调用方法GetTotalFileSize * MelsecA3CNet1重命名为MelsecA3CNet，MelsecA3CNet1OverTcp重命名为MelsecA3CNetOverTcp，所有引用这两个类的话无法兼容更新。 * MelsecA3CNet,MelsecA3CNetOverTcp修复和校验bug，支持是否和校验，支持格式1,2,3,4四种通信机制，已经通过单元测试。 * MelsecA3CServer: 新增三菱的A3C虚拟plc，支持是否和校验，格式1，2，3，4，支持和MelsecA3CNet,MelsecA3CNetOverTcp通信测试。 * FanucSeries0i: fanuc机床客户端新增ReadCutterNumber读取当前刀具号的API信息。 * MqttClient: 新增属性UseTimerCheckDropped是否启动定时器检测，其他优化。 * SoftBasic: 新增一个静态辅助方法string GetAsciiStringRender( byte[] content ); 用来获取ASCII显示的，如果遇见不可见字符，则用十六进制替代。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-8-3 | 10.1.2 | * feat(UltimateFileServer): 获取文件服务器的指定目录的大小的方法改为获取大小，文件数量，最后一个文件修改时间的方法，客户端同步更新GetGroupFileInfo。 * feat(MqttServer): 获取文件服务器的指定目录的大小的方法改为获取大小，文件数量，最后一个文件修改时间的方法，客户端同步更新GetGroupFileInfo。 * feat(UltimateFileServer): 新增一个获取指定目录所有子目录的基本信息，包含总大小，文件数量，最后一个文件修改时间，GetSubGroupFileInfos * feat(MqttServer): 文件服务器新增一个获取指定目录所有子目录的基本信息，包含总大小，文件数量，最后一个文件修改时间，GetSubGroupFileInfos * SiemensWebApi: 修复在.net standard2.0及2.1中未添加SiemensWebApi类的bug。 * MelsecA3CNetHelper: 优化数据解析时，对错误异常的处理，增加错误捕获。 * LSisServer: 重写串口Cnet协议的数据接收，处理，和返回，支持了单变量数据，和连续数据的读写操作。 * XGBCnet, XGBCnetOverTcp: 重新实现了类，统一了代码，重新解析的数据内容，支持了对错误码的提取和错误消息的解析。 * memobus: 新增memobus通信协议，初步实现了读取的操作，具体还需要测试，如有网友有测试条件，可以联系我测试。 * Yamatake: 新增山武的数字示波器的CPL协议的通信对象和虚拟服务器，分别是DigitronCPL,DigitronCPLServer * MqttServer, HttpServer: 使用HslMqttApi特性注册远程RPC接口时，支持对异步方法(async...await)进行注册，并进行异步调用，返回相关数据。仅支持NET451及以上。 * HslMqttApi: 在注册RPC接口时，增加了对方法签名的解析过程，客户端可以浏览服务器接口的方法签名，清楚的看到返回类型，参数类型信息。 * Delta: DeltaDvp系列的网口，串口协议，修复在跨区域读写M1536及D4096时，地址偏移不正确的bug。现在会自动跳转。 * MqttServer: 修复Mqtt客户端在取消订阅时，传入多个的Topic时导致服务器解析异常的bug。 * XGBCnet: 支持Read(string[]); 读取多个地址的数据信息，自动按照每16进行拆分访问。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-8-23 | 10.1.3 | * JSON: .NET framework的dll对newtonsoft.Json不依赖特定的版本。 * XGBFastEnet: 修复读取单个的bool时报文不正确的bug. * MqttServer: 新增GetMqttSessionsByTopic方法，用来获取订阅某个主题的所有客户端列表。 * HttpServer: 修复httpserver中文编码问题，在谷歌，微软浏览器下显示中文乱码的bug。因为Content-Type传值不正确 * HttpServer: 修复在账户验证模式下，使用ajax跨域请求OPTIONS导致401错误的bug。 * FujiCommandSettingType: 新增富士的CommandSettingType通信协议的实现，基于TCP访问，支持地址见API文档. * FujiCommandSettingTypeServer: 新增富士的CommandSettingType协议的虚拟PLC，支持B,M,K,D,W9,BD,F,A,WL,W21地址 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登字第5219522号，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-9-3 | 10.1.4 | * NetworkAlienClient: DTU(异形)服务器增加对报文的固定头的检查。 * NetworkServerBase: 连接异形DTU(异形)的server的方法ConnectHslAlientClient支持密码输入。 * NetworkDoubleBase: 当设置DTU会话时，修复恰好正在读取导致报文错乱的bug，初始化成功才成功切换DTU。 * McServer: 修复三菱的MC虚拟服务器在ASCII模式下，远程客户端读写B继电器时，地址解析异常的bug。 * LSisServer: 修复LSisServer在客户端读写位时，PX20之类的地址时，写入true不成功的bug。 * TemperatureController: 新增RKC的数字式温度控制器的读写类，地址支持M1,M2,M3,等等 * TemperatureControllerOverTcp: 新增RKC的数字式温度控制器的网口透传读写类，地址支持M1,M2,M3,等等 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-9-12 | 10.2.0 | * NetSoftUpdateServer: 代码优化，新增一个OnlineSessionss属性，用来获取当前正在更新的客户端的数量。 * RSAHelper: 提供了从PEM格式的私钥和公钥创建RSA对象的辅助方法，还提供了RSA对象导出PEM格式的私钥公钥方法。支持加密解密超长数据 * RSACryptoServiceProvider: 增加RSA对象的扩展方法GetPEMPrivateKey及GetPEMPublicKey方便快捷的获取PEM格式的公钥和私钥，扩展加密解密超长数据。 * SerialBase: 串口基类新增虚方法CheckReceiveDataComplete用于检查报文是否接收完成，一旦接收完成，立即返回，增加通信性能。 * ModbusRtu: 重写了CheckReceiveDataComplete方法，根据功能码的不同来判断当前的数据长度是否完整，以此判断报文是否完整。 * ModbusAscii: 重写了CheckReceiveDataComplete方法，根据开头及结尾的固定字符来是否指定值，以此判断报文是否完整。 * ModbusTcpServer: 优化服务器的串口接收功能，现在回复报文很快。1. 先接收串口数据，再Sleep。2. 增加数据完整性校验，一旦成功，立即返回报文。 * ModbusTcpServer: 新增属性RequestDelayTime，设置非0时用来防止有客户端疯狂进行请求而导致服务器的CPU占用率上升问题。 * MelsecA1EServer: 新增MC-A1E协议的虚拟服务器，支持了二进制和ASCII格式，已经配合客户端通过单元测试。 * AesCryptography, DesCryptography: 新增AES及DES加密解密对象，使用密钥实例化即可加密解密操作。 * MQTTServer: 在原有的基础上增加了加密的功能，如果MQTTClient，MQTTSyncClient启用加密功能，那么数据的传输就是加密的，无法抓包破解账户名密码及交互数据。 * AllenBradleyNet: cip协议支持了日期，日期时间的读写操作，这样omron的plc的cip协议也支持了日期时间的读写，在omroncip的demo上添加测试操作。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-9-19 | 10.2.1 | * SiemensS7Net: 修复读写WString字符串时，乱码的情况。 * NetSoftUpdateServer: 修复在某些情况下在线客户端数量新增后不会减少的bug。 * demo: 字节转换工具界面增加字节数组和base64字符串的转换功能。 * MqttServer: 当MQTTClient在加密模式下，订阅一个主题后，修复服务器仍然返回没有加密的bug，导致客户端解密失败。 * MqttSyncClient: 修复加密模式下，使用SetPersistentConnection设置长连接，不进行ConnectServer直接第一次请求失败的bug。 * AllenBradleyNet: 优化读取bool的功能方法，新增读取bool数组的实现。 * NetworkDataServerBase: 所有串口类的服务器（从机），功能代码优化调整，部分的类实现报文完整性判断，实现数据瞬间回复客户端（主机）。 * Serial: 大量的串口类的设备进行了优化，对接收结果是否结束进行判断，提高了串口通信的性能。 * NetSoftUpdateServer: 新增另一种更新机制，更新软件在收到文件信息后，先比对MD5信息来确定是否下载更新，从而提高升级速度，仍然兼容旧的更新模式。 * NetSoftUpdateServer配套的更新客户端，软件自动更新重新命名为 AutoUpdate, 针对差异化更新做了优化。 * DEMO里面所有的读写PLC界面的读写字符串功能支持了可选编码，支持ASCII,UTF8,UNICODE,UNICODE-BIG,GB2312,ANSI * 其他一些细节的优化，和注释的完善。DEMO界面的大量优化，显示调整，DEMO使用了新的更新软件，AutoUpdate.exe * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-9-30 | 10.2.2 | * NetworkDoubleBase: 在异步的网络通信方法中，原来的同步锁会在特殊的情况下导致UI线程卡死，所以改为异步锁，相关继承类也修改。 * HslReflectionHelper: 通过HslDeviceAddressAttribute反射Read，Write的读写自定义对象的功能支持对byte类型的读写操作，需要通信对象本身支持才能成功读写。 * SerialBase: 新增protect级别的AtLeastReceiveLength变量，用来表示从串口中至少接收的字节长度信息，一般为1。 * IReadWriteNet: 新增api支持ReadCustomer( string address, T obj )，允许传入实例对象，对属性进行赋值，方便wpf进行数据绑定操作。 * NetworkWebApiBase: 新增属性UseEncodingISO，在访问某些特殊的API的时候，会发生异常"The character set provided in ContentType is invalid...."，这时候本属性设置为true即可。 * Cip: 欧姆龙cip及rockwell的cip支持读取plc型号的方法ReadPlcType()，omron的cip重新设计了WriteTag，对于0xD1类型数据，支持偶数个写入操作。 * SiemensPPI: 修复writebyte方法的api接口名称还未注册的问题, 和串口透传类的相关代码优化，精简。 * MelsecFxSerial: AtLeastReceiveLength变量设置为2，和串口透传类的相关代码优化，精简。 * MqttServer: 新增属性：TopicWildcard，当设置为true时，支持主题订阅通配符，使用#,+来订阅多个主题的功能。具体参考该属性的API文档。 * demo: 修复demo的HTTPClient界面浏览在linux创建的Webapi服务器的api列表功能失败的bug，使用HttpWebRequest来实现。 * demo: rsa加密解密算法测试界面实现签名和验签操作。签名算法可选SHA1，SHA256, SHA512, MD5等等。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-11-19 | 10.3.0 | * allenbrandlyNet: ExtraOnDisconnect方法增加socket的空检测操作，因为在极少数特殊情况会发生空的异常的bug。 * OmronCipNet: 修复Write数据类型为short[]时长度不是1的bug。同步方法调用Write(string,short[])时会发生写入失败。 * Demo: OmronHostLinkCModeOverTcp的Demo测试界面添加读取型号的功能，可以更加便捷的测试。 * AllenBradleyNet: 重写读取字符串方法ReadString( string address, ushort length )，编码修改为UTF8 * ReverseWordTransform: 字节变换的对象代码优化精简，ReverseWordTransform新增IsInteger16Reverse属性指示short,ushort是否发生翻转，默认翻转。适用极个别特殊的Modbus格式 * YRC1000TcpNet: 安川的机器人的API功能大量添加，支持读取机器人坐标，各种变量信息，启动程序，停止，复位，读报警等待。 * SiemensS7Net: 写入bool数组的方法优化，修改为先读取byte[], 修改中间的位，然后再写回去操作，可以写入一个字节的中间几个位，单元测试通过，仍然有风险，谨慎调用。 * OmronHostLinkCMode: 修复OmronHostLinkCMode读写字节顺序颠倒的问题，新增了读取plc型号，plc状态，修改状态的功能，增加错误码文本解释。 * OmronHostLinkCModeServer: 新增CMode协议对应的虚拟PLC服务器，初步实现了和客户端通信功能，包括串口和网口的支持，单元测试通过。 * 添加在NET2.0，NET3.5及Standard项目缺失的SiemensFetchWriteServer, MelsecA3CServer, MelsecA1EServer, 以及安川的 memobusTcpNet * SiemensS7Net: 在读取单个的bool的值的时候，增加对结果数据的合法性检测，当遇到错误码时，提示错误信息。修复写入bool在某些特殊的情况下，实际成功但是提示失败的bug。 * YRCHighEthernet: 安川机器人添加高速以太网通信，基于UDP实现，在DEMO界面可以直接原生命令测试，支持读取报警，字节，整型，字符串等变量。 * PipeSerial, PipeSocket: 新增串口管道和网口管道，串口基类和网口基类优化重构，允许多个串口设备对象共享一个串口通道，多个网口的设备对象共享一个网口通道。 * ModbusTcpNet:新增属性IsCheckMessageId，用于配置modbustcp协议是否检查设备返回的消息ID是否一致，默认检查，设置为False，也即是不检查ID一致。 * MemobusTcpNet: 安川PLC的memobus协议的修复报文错误的bug，增加了线圈读写，寄存器读写，扩展寄存器的读写功能，例如扩展保持寄存器地址：x=9;100 * Demo：测试界面添加了一个TCP转TCP的界面，可以用于两个网络连接中捕捉通信的数据报文。还有其他的一些细节优化。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-12-2 | 10.4.0 | * Delta:台达系列的PLC类重命名，新增PLC系列，可选DVP系列和AS系列，相当于新增了对AS300的支持，如果用到，谨慎更新。 * WebSocketServer: 原先订阅的操作使用header来完成，新增客户端从url添加订阅主题，例如：ws://127.0.0.1:1883/HslSubscribes=A,B * AllenBradleyPcccNet:新增AB的CIP协议PCCC格式的通信，测试适用ab1400，地址格式为 N7:5, I:0/2, F2:0 等 * AllenBradleyPcccServer: 新增PCCC的虚拟服务器功能，可以配合客户端进行通信的测试，支持基本的地址。 * MC-A3C: 修复A3C的服务器在格式4的情况下，返回写入操作的控制代码错误的bug，应该为06，写成了02。 * 三菱的A3C协议的读取字数据及位数据长度自动切割，可以输入理论最大长度，修复校验是否写入成功确提示失败的bug。 * NetworkDataServerBase:串口接收增加最小接收时间，默认是20，表示20毫秒，如果重写了报文结束检查，则可以设置大一点。 * NetworkServerBase: 基础服务器对象新增属性EnableIPv6，当设置为true的时候，则使用IPv6协议进行数据通信和访问。 * 所有的TCP及UDP的相关的客户端通信类的 IP 地址设置，都支持了 IPv6 地址，所以目前IP地址支持 IPv4 和 IPv6 及 网址。 * Demo:fanuc机器人的测试界面写入数据操作进行提示相关的确认，防止手抖不小心写入导致机器人运行不正常。 * Demo，所有的设备通信类，新增一个支持的设备列表功能，大家可以手动添加支持的列表，方便别人查看支持的型号，请大家务必真实填写。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-12-17 | 10.4.1 | * Delta: 台达的DVP系列，当地址带站号时，并且出现了跨地址切割时，修复站号丢失的bug。跨地址指跨越M1536,D4096地址。 * Delta: 台达AS系列的地址 SR 地址支持输入 SR0 以及 sr0 * 针对FinsTcp协议及ADS协议的剩余字节解析，本来是int类型，现在加了最大长度10000的限制，防止异常攻击时导致申请内存过大而奔溃。 * Inovance: 汇川相关的PLC的数据排列规则调整为 CDAB，DEMO上的选择也默认改为这个。 * MelsecFxSerial: 移除三菱编程口协议最小接收回复报文2字节的限制，实际测试有时候只需要一个字节就可以了，否则会一直处于接收超时。 * ByteTransform: TransBool的方法的偏移和长度均修改为以位为单位，具体调用API时查看注释说明，长度为10就表示10个长度的bool。 * AllenBradleySLCNet,AllenBradleyPcccNet: 地址支持ST10:2这种字符串地址，并在AllenBradleyPcccNet上实现动态读取字符串，长度为0或没有，则自动长度。 * AllenBradleySLCNet,AllenBradleyPcccNet: 以及SLC的协议地址支持了 L17:0 的地址，长整型类型 * VigorSerial,VigorSerialOverTcp: 新增丰炜PLC的编程口协议及透传类，支持VS系列，地址支持动态站号信息，例如 s=2;D100 * VigorServer: 新增丰炜的虚拟PLC，模拟了VS系列的通信，可以和对应的客户端进行数据读写测试，位地址支持 X,Y,M,S，字地址支持 D,R,SD * MelsecFxSerial: 三菱的编程口协议及其透传类添加一个激活PLC的API方法ActivePlc，在某些特殊的PLC读写数据之前，需要先调用一下。 * Modbus:Modbus所有协议的ReadFromCoreServe(byte[])增加特别的注释，只需要传递modbus核心报文即可，不管tcp,还是rtu,还是ascii，都是一样的。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2021-12-22 | 10.4.2 | * NetworkBase: 新增Receive原始字节数据时，也就是申请byte[]缓存时，对可能存在的异常捕获操作，返回失败的结果对象。 * AllenBradleyPcccNet: 修复对于ReadAsync也注册了API接口，导致重复注册的bug。 * WebSocketClient: 修复因为掉线重连服务器的时候，重新携带订阅的主题。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-1-1 | 10.4.3 | * SerialBase: SetPipSerial重，命名为SetPipeSerial，如果有使用串口管道，则需要更改相关的名称。 * MelsecMcDataType: 修复三菱累计定时器当前的值的地址进制转换应该为10进制，结果写成100导致转换失败的bug。 * Keyence: 上位链路协议的串口及网口的通信类，ByteTransform的IsStringReverseByteWord调整为true，读写字符串时将两两颠倒。 * IByteTransform: 转换接口类的注释进行完善，提示更加详细完整，中英文并行提示。 * Vigor: 丰炜PLC的读取位和读取字的功能方法，对读取长度进行内置切割，相当于支持了无限长度的数据读取。 * EstunTcpNet: 新增埃斯顿机器人通信类，内置定时器保持心跳，支持读取机器人的基本信息，详细见DEMO界面。 * FanucInterfaceNet: 修复fanuc机器人的中文编码异常的bug，使用标准的GB编码解析，如果编码获取异常，需要自行nuget安装System.Text.Encoding.CodePages组件，并注册编码。 * Device: 设备类增加ReadStruct{T}方法，根据特性从原始字节里解析出实际的数据对象。影响范围所有的设备类对象。 * Demo: Demo程序支持了手动设置版本更新忽略提醒，忽略之后再菜单栏进行提示新版，以及增加添加激活操作功能，本地保存加密的激活码。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-1-17 | 10.5.0 | * ABBWebApiClient: GetRobotTarget方法，返回的数据信息，增加q4数据信息。 * Modbus: 在modbus的派生类协议中，重写了modbus的地址转换的情况中，修复读写bool操作，因为地址中带小数点导致地址转换异常的bug。例如汇川AM系列读写QX3000.0 的bool值 * MelsecFxSerialHelper: 三菱编程口协议的检查和校验的方法，增加try...catch，在某些特殊的情况下，会导致异常，直接奔溃。 * NetSoftUpdateServer: 针对之前旧版的软件升级功能，增加30分钟的超时时间限制，如果30分钟后仍然没有更新完成，则自动移除会话。 * SiemensS7: S7的地址支持大小写，且都支持带X,B,D,W的地址，比如MD100, MW100, MX100.2 * OmronConnectedCipNet: 注册报文里的不通信超时修改为42分钟，读取short，及ushort数组时，按照994长度进行切割 * AllenBradleyPcccNet: 支持使用0F-AB的掩码写入功能，写入一个bool值到PLC中，PCCC虚拟服务器实现了这个AB功能码，可以虚拟测试。 * FanucInterfaceNet: 修复demo上读取WO数据时，地址偏移不正确的bug。 * Freedom: 串口，网口的自由通信协议增加委托CheckResponseStatus，可以自定义对报文结果进行校验，完善注释。 * DLT645: 优化数据接收部分的代码，如果数据完整，立即返回，数据前面兼容无用的字符数据。 * INetMessage: 消息类新增方法PependedUselesByteLength( byte[] headByte )并在DLT645OverTcp消息类重写，支持前置无效的字符。 * AllenBradleyNet：支持添加消息路由功能，默认不开启，实例化属性MessageRouter, 例如：1.15.2.18.1.1，支持在demo界面进行配置操作。 * AllenBradleyNet: 支持遍历全局变量和局部变量。新增StructTagEnumerator( AbTagItem structTag )方法遍历结构体的成员变量信息。 * Demo:AllenBradleyNet节点浏览的界面支持了查找数据，显示数据，结构体嵌套遍历，还支持显示当前的数据信息。 * ModbusTcpServer的RTU接收时间调整为500ms，如果报文完整立即接收结束。 * 本版本可能是春节前的最后一个版本了，提前祝大家新春快乐。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-2-28 | 10.5.1 | * AllenBradleyNet: Write( string address, bool value )写入bool方法优化，如果是类型代号0xD3的，地址前面需要增加 "i=" 标记。 * PanasonicMewtocolServer: 松下虚拟PLC地址类型支持完善线圈支持X,Y,R,L, 字单位的整型支持 X,Y,R,L,D,LD,F，支持RCP及WCP指令离散读写线圈。 * PanasonicMewtocol：松下协议及其网口类的代码优化提炼共同代码，新增加离散bool地址的读写，传入多个地址，返回一个bool数组。 * AsciiControl: 在命名空间slCommunication.Core下面增加一个AsciiControl类，包含控制类ASCII常量定义，例如ENQ,NAK,STX,ETX等。 * MelsecFxLinksServer: 新增三菱的FxLinks计算机链接协议的虚拟PLC，支持串口和网口透传访问，支持格式1，4切换，支持是否和校验。 * MelsecFxLinks: 1. 同时支持了格式1，格式4，2. 支持了地址大于10000的时候使用QR,QW命令，3. 修复报文创建时数据长度及站号不是16进制的bug。 * KeyenceNanoServer: 修复启动串口时，和上位链路客户端通信时，一直返回通信校验错误的bug，现在针对CC指令和CQ指令都能正确的返回。 * keyenceNano: 基恩士上位链路协议的串口类和网口透传类记录报文的格式调整为ASCII码，这样更加直观。 * NetworDataServerBase: 添加GetNewNetMessage( )及ReadFromCoreServer方法，精简所有继承的子类虚拟PLC的服务器的代码。 * MelsecA1EServer: 修复三菱A1E协议服务器报文接收异常的bug，导致客户端读写数据不正常。 * FatekPrograme: 永宏编程口协议读取字按照64字长度自动切割，支持了RUN,STOP，读取状态接口方法，新增对应的虚拟服务器实现。 * SpecifiedCharacterMessage: 新增基于指定字符结尾的消息类，通过ProtocolHeadBytesLength属性变种而来，NetworkBase的ReceiveByMessage功能适配了SpecifiedCharacterMessage消息。 * Turck: 新增图尔克的Reader协议实现，支持对字节读写，bool读写，实现了虚拟服务器，通过了单元测试，主要用来和RFID进行通信。 * IEC104: 初步添加IEC104协议实现和解析，增加demo测试，但是目前还清楚怎么设计API和使用场景，等待继续优化，欢迎相关需求的人联系测试优化。 * DEMO: RSA加密解密的测试界面，支持对超长的数据进行加密解密操作。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-4-1 | 10.6.0 | * ModbusHelper: 修复ReadBoolHelper在Modbus继承类里进行了自定义地址转换，导致读字地址的位失败，例如汇川AM系列读取M240.1的bool数据失败 * LSisServer: 修复基于Cnet串口协议的通讯时，客户端使用SB连续读取指令读取长度大于9时，服务器对长度解析失败导致读取失败的bug。 * GroupFileContainer: 新增获取所有文件的下载次数总和的属性，名称为：TotalDownloadTimes。 * SerialPort: 新增串口的扩展方法IniSerialByFormatString，支持格式化的初始串口信息，例如 COM3-9600-8-N-1，简单明了。ModbusServer界面支持串口参数配置，支持配置最短接收时间配置。 * MelsecFxLinksServer: 三菱Fxlinks虚拟服务器，修复远程客户端在读取位地址（例如M100）的字数据时，返回错误失败的bug。 * HslReflectionHelper: MRPC及Webapi注册的方法接口功能中，支持参数类型为自定义类型数组的参数，例如Student[]，不支持List{Student}。 * MqttServer: 文件读写功能的引擎中，遍历所有子目录统计信息时，可选参数是否携带返回最新的一个文件的信息，客户端做了相关的适配，参数支持。 * XGBFastEnet: 修复向PLC读取数据到时候，类型指定不正确的bug，之前无论是什么类型，都会设置为0，PLC从而值返回一个字节长度的数据。 * Mewtocol: 松下Mewtocol协议，当读写字单位超出协议限制时，自动升级为扩展协议标识，标记字符为 < (小于号), 服务器读写长度根据实际限制，支持扩展协议标识。 * FanucSeries0i: 修复fanuc机床类读取虚拟机床系统时，显示运行时间为负数的bug。 * NetworkDoubleBase: 新增受保护的属性UseServerActivePush，只要设置为True，就可以将当前的连接切换为既支持设备主动发送，又支持同步访问的客户端，使用信号同步器来支持问答通信。 * SiemensPPIServer: 新增PPI Server类，支持TCP以及串口的通讯，PPI的串口和网口透传类客户端增加读取PLC类型的功能方法。 * SiemensS7Server: 修复客户端读取数据时，返回客户端FF04的数据内容时，位长度信息赋值不正确的bug，这会导致某些客户端发生数据异常。 * SiemensPPI: 西门子PPI协议优化，读取bool[]的功能方法调整为先读取字，再解析出位数据信息，以此来支持批量读取位。 * XinJETcpNet: 信捷的modbustcp协议新增当某些地址超过临界范围时，自动升级modbus协议到信捷内部TCP协议，以此访问到更大范围的数据内容，并支持了临界地址的跨地址访问，自动切割重组。 * XinJEInternalNet: 新增完全的信捷内部协议实现的TCP通信，可以读取更大范围的数据内容，比如D200000地址，支持的地址参考API文档。 * DeltaDvp: 修复台达系列在进行读取D寄存器的位的时候，也即是时 ReadBool("D100.1") 方法时，地址无法识别的bug，现在可以正确的读到D100的位信息。 * FanucSeries0i: 新增WriteRData的接口方法，新增读写PMC数据的方法，ReadPMCData及WritePMCData，与R数据的读写类似，都是字节操作的方法。 * IReadWriteDevice: 扩充方法支持，支持ReadFromCoreServer(List{byte[]})的读取的方法支持，支持多个报文读取，结果打包返回。 * BeckhoffAdsNet: 倍福通信类新增是否使用自动AMS属性，默认为false，设置为true时，网络初始化时从server端加载相关的NETID参数信息，目前在twincat3上成功测试。 * BeckhoffAdsNet: 修复ReadBool数组时，返回的长度都是8的倍数的bug。属性UseServerActivePush调整为True，支持PLC方主动推送数据信息。 * BeckhoffAdsNet: 修改读取bool时地址机制，地址支持小数点，例如M100.2 ，所以M800等于M100.0， 修复读取数组时值解析不正确的bug， 修复写入bool数组失败的bug。 * BeckhoffAdsNet: 修复直接SetPersistentConnection设置长连接无法读取的bug，以及掉线后，继续读写一直失败的bug。 * BeckhoffAdsNet: 修复倍福服务器重启后继续读写标签变量一直失败的bug，原因来自倍福重启导致标签内存地址变更，但是缓存还是一直使用旧的。 * SecsHsms: 新增secs协议的hsms实现，初步实现了主动方，也即是客户端，支持了一个通用的接口和数据，ReadSecsMessage方法。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-4-11 | 10.6.1 | * MelsecMcRNet: 修复McType属性设置不正确的bug，导致依然使用普通的mc协议来进行通信，因为R系列MC协议地址不一样。 * ModbusServer: ModbusServer增加属性StationCheck，默认为True，对请求的客户端的站号进行检查操作，也可以设置false，不检查站号一致。 * BeckhoffAdsNet: 新增支持批量读取的方法Read( string[] address, ushort[] length )，完善了API文档说明，可以一次性读取多个标签信息，需要自行解析。 * BeckhoffAdsNet： 修改默认的端口为851，主要用于TWINCAT3，当没有设置 NETID 时，直接连接PLC的话，自动升级为自动获取NETID信息。 * BeckhoffAdsNet: 重写实现新的Read{T}(), 一次性读取所有的类型数据，然后自动解析赋值到对象的属性，在API文档上增加示例的代码。 * SiemensPPI: 西门子PPI协议修复输入M5这种只有两个字符的地址的时候，地址解析报异常的bug，同样也适用于SiemensPPIOverTcp * KeyenceNano: 基恩士上位链路协议优化，地址重新设计，修复B,VB,W的十六进制地址解析异常的bug，FM地址不正确解析的bug。 * KeyenceNano: 上位链路协议读取支持自动切割，根据不同的地址切割出不同的长度信息，一般都是256长度切割，虚拟服务器也进行了地址修复和长度限制适配。 * OmronFinsUdp: 欧姆龙Finsudp协议的DA1（PLC节点号）调整为0x00，自动获取，在demo测试界面，增加GCT和SID的参数设置，一般默认就好。 * FanucSeries0i: 读取R数据的接口重新设计，修改为ReadPMCData，地址支持G,F,Y,X,A,R,T,K,C,D,E 例如R1200，写入也是一样，具体看demo演示。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-5-25 | 11.0.0 | * TcpForward: 新增TCP转TCP的类，设置本地端口及远程端口，客户端连接之后立即连接远程服务器，进行数据转发，支持多个连接，每个客户端都创建一个远程连接对象，DEMO的转发测试使用这个类。 * BeckhoffAdsServer: 倍福的虚拟PLC修复M100.1位地址读写数据不正确的bug，支持了符号地址，支持创建符号数据，自动生成唯一的句柄地址。支持客户端批量的地址读取，也即功能码为 0xF080。 * AllenBradleyPcccServer: 修复PCCC虚拟PLC读写地址数据解析异常的bug，导致和客户端的地址数据不一致。提炼优化地址信息。 * NetworkDataServerBase: 优化串口接收时候的数据检测的代码，修复modbus虚拟服务器接收串口数据时，功能码为0x10时，SerialReceiveAtleastTime属性设置无效的bug。 * XinJE: 基于modbus协议的信捷类对象，当地址是位地址时(X,Y,M等)，支持了当使用Read（字单位）读取原始字节的方法，返回原始的字节信息。 * FanucSeries0i: Fanuc数控机床通信类，新增读取机床型号信息的api，ReadSysInfo，信息在连接初始时获取生成，后续直接读取内存。 * FanucSeries0i: 读取程序文件及下载程序到机床支持路径参数，例如操作PATH2时只需要输入 //CNC\_MEM/USER/PATH2/ 即可，新增DeleteFile的API用于删除任意路径的程序文件。 * FanucSeries0i: 新增ReadAllDirectoryAndFile方法，用来遍历指定路径下的所有的子路径和文件的信息，文件大小使用字节为单位。 * ISessionContext: 在MRPC的接口及HttpWebapi注册接口的方法里，参数ISessionContext新增Tag对象，方便自定义捆绑一些数据信息。 * MqttServer: 文件服务器相关的功能增加GetGroupFromFilePath方法，从路径信息里获取到文件列表管理容器，进而可以获取单个文件的对象，在权限控制时可以更加的细致。 * Modbus: ModbusHelper针对modbus-rtu的返回接收报文检测，增加对站号一致性的检测，如果发现站号不一样，返回错误信息。 * ICryptography: 添加直接对字符串进行加密解密的方法，并在AesCryptography类和DesCryptography实现方法，方便直接操作。 * OmronConnectedCipNet: 优化批量读取数组的方法，自动根据类型来切割不同的读取长度，分批次读取所有的数据信息。修复握手报文丢失Timeout时间报文的bug，导致v10.5.0以来不能读取的问题 * OmronConnectedCipNet: 优化和PLC持续几秒钟不读写就被PLC强制断线的问题，现在支持很长的时间，修复关闭立即重连（包括意外断线重连）连接不成功，一直报错重复打开会话的异常。 * NetworkConnectedCip: 基于连接的CIP协议里，在OpenForward初始化里，自动获取连接ID，如果失败，并尝试五次自增连接ID重新Open。 * AllenBradleyConnectedCipNet: 新增ab的connected cip协议的实现，继承自欧姆龙的connected cip协议，重写了部分的逻辑实现，缩小了单次读取的字节上限，初步测试通过。 * NetworkServerBase,NetworkDoubleBase: 新增属性SocketKeepAliveTime，用来设置底层的socket的KeepAlive信号，单位毫秒，默认不开启，开启后将会自动间隔的发送KeepAlive信号。 * AllenBradleyServer: cip的虚拟服务器同时支持了基于连接的CIP的访问。标签定义类AllenBradleyItemValue新增TypeCode属性，在返回客户端数据的时候，带类型数据返回。支持bool[]节点创建 * SoftBasic: GetAsciiStringRender 不可见字符判断除了0x20以下，追加0x7e以上的字符。新增GetFromAsciiStringRender方法，把GetAsciiStringRender 结果反向转换回来。 * AllenBradleyNet: 支持写入bool[]功能，支持写入byte[]数组功能（类型代号默认为0xD1），WriteTag方法支持地址携带类型信息，例如地址 type=0xD1;A[0] * MelsecFxSerial: 三菱的编程口协议在新版协议配置为True时，修复D8000地址以上读写不正确的bug，修复X,Y,M,S部分地址读写不正确的bug，都调整为新版的报文，支持的范围更大更准确。 * MelsecFxSerialOverTcp: 新增加通过GOT连接的方式，。PanasonicMewtocol: 松下的Mewtocol协议支持了经过值及目标值。ILogNet: 日志的接口新增属性LogThreadID，用来配置是否记录线程ID的数据信息。 * SecsHsmsServer: 新增Secs协议的虚拟服务器，支持事件自定义处理所有的消息。并添加示例的代码。Demo程序界面。 * Mqtt: MqttClient及MqttServer支持了遗言消息，客户端A在连接参数里设置了遗言主题及数据，当客户端A非正常掉线的时候，服务器即发布该遗言主题信息。 * MqttSyncClient: 新增HslCancelToken对象，在文件上传和下载的过程中支持取消操作，MqttServer适配相应的功能实现，取消示例代码增加。 * HttpServer: 新增Action对象DealWithHttpListenerRequest，用来处理定义的Header数据到Session会话中，实现更加复杂自定义的操作。 * Demo: DEMO写入各种PLC的数据功能里，写入的数据支持byte数组，例如数据输入[1,2,3], 点击写入byte，就自动写入new byte[]{ 0x01, 0x02, 0x03}; * Demo: TCP调试界面优化，新增握手报文设置。新增一个测试正则表达式的界面，增加一些常用的正则表达式的内容。CIP浏览节点界面支持正则表达式筛选节点。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-6-9 | 11.0.1 | * HslExtension: 新增了一些扩展方法，现在可以直接从short，int类型获取位及设置位的功能，例如 bool value = shortValue.GetBoolByIndex(10); 设置位也是同理。 * DLT645: Demo界面优化，无论是串口的，还是网口的，都支持了密码及操作者ID的输入。修复了站号输入后，但是站号初始化不成功的bug。 * LogStatisticsBase: 当传入的长度为-1时，表示不限制长度信息，会一直的新增。 * Upgrade.exe: 用于自动升级更新的项目重新优化，支持了下载时候的网速显示，支持了下载百分比显示，删除了自动创建快捷方式的代码。 * SerialBase: 初始化串口的方法SerialPortInni( string portName )支持格式化的输入，例如COM3-9600-8-N-1，COM5-19200-7-E-2。 * MelsecFxSerialOverTcp: 当启动GOT透传时，修复报文封装不正确导致写入失败，读写值不正常的bug。 * InovanceHelper): 汇川的PLC针对M地址的寄存器，支持使用 MX0.0 来对位进行读取操作。支持了MB100读取byte类型的数据。 * Demo: 在Demo所有的设备通信测试界面，几乎所有的设备都支持了显示报文信息，想要查看通信的报文协议，菜单里点击 报文日志 。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-6-28 | 11.0.2 | * LogNet: 修复日志对象在配置了LogThreadID=False时，控制台输出界面仍然显示线程号信息的bug。 * AllenBradleyServer: 新增属性CreateTagWithWrite，当手动设置为true时，从服务器端写入不存在的标签将根据写入的类型自动创建标签。 * NetworkUdpBase: udp的通信基类修复某些设备通信情况下每9次通信就会跟随一次失败的bug，原因来自不停的创建socket，现在修改为连接正常就不重新创建socket * AllenBradleyServer: ab虚拟plc修复写入short的类型分配不正确的bug，并且增加一个创建byte字节数组的标签的方法接口。 * AllenBradleyNet: 修复同步方法ReadString(string address); 读取字符串返回数据乱码的bug，原因来自调用了基类的字符串解析。 * DLT645: 大面积优化，支持了更多的地址读取，支持了一个地址多个数据读取，修复了部分数据(如电流，功率因数等)不识别正负号的bug。 * PipeSocket: 端口号信息由int型调整为数组，使用SetMultiPorts方法可以设置多个端口号信息，当PLC对象重连时，就会切换端口号，循环反复使用指定的端口号信息。 * PipeSocket: 多端口号使用方法为 plc.GetPipeSocket( ).SetMultiPorts( new int[] {6000, 6001} ); 实例化之后调用一次即可。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-7-22 | 11.0.3 | * MelsecFxSerialOverTcp: 修复在启用GOT透传时，碰到特殊的报文解析异常，最终导致偶尔读写失败的bug。 * InovanceHelper: 修复汇川PLC的AM系列时，读写M区域位地址的数据偏移不正确的bug，例如MX100.0，实际地址应该是MW50的第0位。 * ModbusTcpServer: 新增属性StationDataIsolation，默认为false，表示server只有一个数据区，所有站号用的一个数据区。如果设置为 true，表示server给每个站号开辟一个数据区。 * ModbusTcpServer: 修改Station站号类型为byte，删除属性StationCheck，现在不需要检查站号了。服务器侧读写数据支持了输入站号信息，例如 s=2;100，可以访问不同站号的数据。 * FujiSPBServer: 修复富士SPB虚拟服务器再读写位数据的时候，修复地址偏移解析异常导致bug。 * AllenBradleyServer: CIP协议的虚拟服务器修复当客户端写入bool变量时，无论写入True还是False，服务器都写入True的bug。 * MelsecFxSerialOverTcp: 三菱的编程口协议里，读取的地址长度突破了254字节的长度限制，现在可以读取无限个bool值，或是其他类型的值。 * HslReflectionHelper: 在MRPC及webapi接口的参数数据提取过程中，支持传入JSON对象及JSON对象的字符串，都会自动解析成正确的对象值或是JObject值。 * MelsecA1ENet: 三菱的A1E协议，读取字及位数据时，支持了读取超过64地址长度的数据，内部自动切割重组。 * DLT698: 初步添加DLT698.45的协议实现，使用明文的通信方式。支持读取功率，总功，电压，电流，频率，功率因数等数据。 * XGBCnet: Lsis的XGBCnet协议修复部分的bool读写位置和真实PLC对应不上的bug，bool的读取支持了MW100.2 带小数点表示的方式。 * XGBCnet: 支持了批量的bool数组读取功能（内部自动读取字节，解析出bool数组，不支持写入bool数组）。修复XGBFastEnet在配置XGB型号时读写数据异常的bug。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注，新官网：http://www.hsltechnology.cn/。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-7-27 | 11.0.4 | * NET451: 紧急修复NET451项目因为代码合并时导致签名丢失的问题，否则将导致签名的应用程序无法编译。 * DLT698OverTcp: 新增DLT698的串口透传类，使用串口转网口即可进行通信。 * DLT698TcpNet: 新增基于网口的DLT698协议，存在握手包，还未详细测试。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注，新官网：http://www.hsltechnology.cn/。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-8-12 | 11.0.5 | * ModbusAsciiOverTcp: 新增modbusascii协议的串口转网口透传的通信类。 * DeltaSerialAsciiOverTcp: 新增一个台达的ASCII通信协议的串口转网口透传通信类，和其他的基于modbus的台达类一致的用法，一致的地址格式。 * FujiCommandSettingType: 修复FujiCommandSettingType及其虚拟服务器再解析地址不正确的bug，导致大量的其他地址解析失败，并完善了注释说明。 * MelsecFxSerial: 新增属性AutoChangeBaudRate，默认false，当设置为true时，在串口打开时会和PLC交互改变PLC默认的波特率。 * SerialBase: 新增串口类基类属性ReceiveEmptyDataCount，用来表示接收空白数据次数，然后判断接收完成的标记，单次耗费的时间是 SleepTime。 * Modbus: DEMO相关的界面都默认的DataFormat都调整为 CDAB，因为大多数的设备数据大小端都是这个格式。 * OmronConnectedCipNet: 新增IReadWriteCip接口，欧姆龙的基于连接的CIP协议实现了日期及时间的读写接口。 * MqttServer: MQTT服务器支持向自定义规则的客户端会话发布指定的主题数据，详细查看PublishTopicPayload重载方法 * MqttSyncClient: 上传文件到服务器的功能接口里新增流数据的上传，需要指定服务器保存的文件的名称。 * NetworkConnectedCip: 基于连接的CIP协议的连接id信息初始化时及赋予一个随机数，增加一个新的错误描述信息，连接过多的错误。 * Upgrade.exe: 用于更新的exe修复当文件长度为0时，导致百分比计算异常的bug。 * OmronFinsAddress: 欧姆龙FINS协议的地址，增加支持了定时器TIM,计数器CNT，索引寄存器IR，数据寄存器DR的读写。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注，新官网：http://www.hsltechnology.cn/。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-8-18 | 11.0.6 | * TcpForward: TCP转TCP类新增OnRemoteMessageReceived及OnClientMessageReceive事件，可以捕获所有的通信报文。 * TcpForward: TCP转TCP的对象新增ConnectTimeOut属性，支持设定连接超时 * NetSoftUpdateServer: 更新的服务器新增一个方法，GetDealSizeAndReset每秒调用即可获取当前的下载网速信息 * AllenBradleyNet: 完善了路由信息，支持了一种特殊的路由信息的消息，格式为 1.1.2.130.133.139.61.1.0。 * TcpFoward: TCP转发的类新增一个属性OnlineSessionsCount，获取当前在线会客户端会话数量信息。 * SecsHsms: 新增属性 InitializationS0F0，默认为false，指示是否在初始化的时候，发送功能码 S0F0。 * Lognet: 优化日志存储器在计算过期日志文件并删除的代码，现在触发新增文件时，才进行删除旧的文件信息。 * 官网地址： http://www.hslcommunication.cn 官网的界面全新设计过，感谢浏览关注，新官网：http://www.hsltechnology.cn/。 * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-9-13 | 11.1.0 | * SiemensS7Net: 新增ReadDate，WriteDate的功能方法，实现了对日期类型的变量的读写操作，PLC测标记为D#2022-8-19 * SecsHsms: 客户端和服务器的程序，都支持了字符串编码属性设置，可以支持其他的编码，默认的编码调整为Default，主要是支持了中文的消息。 * FanucInterfaceNet: fanuc机器人的地址支持了直接的地址，Bool地址支持：SDO, SDI, RDI, RDO, UI, UO, SI, SO ,字单位地址支持：GI, GO, D。修复日志重复记录。 * OmronFinsNet: OmronFinsNet，OmronFinsUdp, OmronHostLink, OmronHostLinkOverTcp支持0104功能码，批量读取随机字地址，使用方法 Read(string[]) * OmronFinsServer: fins虚拟服务器支持了0104随机读取地址的功能码，对代码进行了优化，提炼了重复的代码。 * MqttRpcDevice: 新增基于MQTT-RPC接口的设备对象，用来访问使用MQTTServer注册plc对象为RPC接口的设备，让单连接的PLC瞬间支持N连接，参考：https://www.cnblogs.com/dathlin/p/16632767.html * Toledo: 托利多电子秤增加通信报文日志记录，支持了扩展模式下的输出格式，并对数据分析增加异常捕获，日志记录。优化串口的通信接收，增加是否配置和校验的属性。 * MqttClient: 每个订阅的主题升级为SubscribeTopic类对象，携带一个事件，支持每个主题绑定不同的事件内容。方便子窗体使用同一个MqttClient对象订阅不同的主题并触发不同的事件。 * SiemensS7Server: 优化字符串的读写操作，支持WString字符串的读写，自动使用S7格式的字符串，在demo写入字符串数据的时候，客户端也可以正确的读取。 * ILogNet: ILogNet日志接口新增一个属性LogStxAsciiCode用来配置是否在每条日志开头记录0x02的ASCII码字符，默认为true，设置false后也就无法使用hsl自带的日志分析工具。 * RedisClient: redis客户端类支持了集群的服务器情况，当标签A在另一个服务器时，redisclient会自动连接到对应的服务器进行获取。 * FanucSeries0i: fanuc机床的通信类支持了ReadOperatorMessage方法，用来读取机床的操作信息。 * MqttServer: 修复Mqtt服务器在进行订阅结果反馈时，没有将topic主题的Qos也返回的bug，在某些情况下，客户端会引发异常。 * NetworkDataServerBase: 新增加属性ForceSerialReceiveOnce，默认为false，当多个modbusserver使用485总线串到一起时，需要设置为true * Iec104: 修复连接失败的bug，修复在I帧消息号返回接收的id信息数量不正确的bug，现在可以正确的接收设备的数据。 * Demo界面新增一个各种校验码测试的界面，主要用于测试CRC16，LRC，FCS（异或校验）,ACC(和校验)，欧姆龙的Fins协议Demo界面上支持配置字符串是否翻转。 * Json: 依赖的json组件 Newtonsoft.Json 更新到最新的 v13.0.1 版本。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-10-12 | 11.1.1 | * Demo: 修复山武设备串口通讯调试界面，站号设置无效的bug。 * MelsecMcServer： MC协议的虚拟服务器支持L寄存器的读写操作，原先的D寄存器范围为 D0~D65535 扩展到 D0~D131071 * MqttSession: Mqtt会话新增属性DeveloperPermissions，标记开发者权限，服务器可以自由控制会话是否具有遍历MRPC接口的权限。 * DLT: DLT645及DLT698协议及其网口透传类新增属性：EnableCodeFE, 莫认false，当设置为True时每次发报文到设备时都会追加FEFEFEFE的指令头。 * OmronConnectedCipNet: 新增属性ConnectionTimeoutMultiplier用来设置异常超时时间，默认02表示32秒。 * SiemensS7Server: 优化握手报文的数据信息，现在支持了node-red的s7功能节点的读写数据，支持了s7netplus组件的连接及读写。 * SoftBasic: 新增从XML元素获取泛型值的便捷方法，DEMO程序的一些界面优化，串口调试界面优化，webapi服务器接口界面支持接口存储。 * OmronFinsTcp: FinsTcp协议优化，在接收fins消息报文时，即使前面有多个无用字节数据，依然可以正确的接收并处理数据。 * SiemensS7Server: 西门子的s7虚拟服务器的DB块优化，采用词典存储，内置DB1，2，3，其他DB块需要使用AddDbBlock方法来手动增加，否则无法读写。 * YokogawaLinkServer: 修复横河PLC虚拟服务器的读取bool数据时，无论输入什么地址类型，都会读取m寄存器的bug。 * HslExtension: ToHexString 的扩展方法以及SoftBasic.ByteToHexString 方法在转换byte[]到string时，支持格式化的参数，默认 {0:X2} * SerialBase: 修复在特殊情况下读取数据导致没有响应，也不触发超时的bug，优化modbusrtu在错误码报文下的接收完成判断。 * MemobusTcpNet: 安川PLC修复读写扩展寄存器失败的bug，通信细节进行优化，支持对扩展寄存器的随机字读取和写入操作。并增加MemobusUdpNet类，基于UDP的实现。 * MemobusTcpServer: 新增加安川PLC的虚拟服务器实现，支持功能码为 01 02 03 04 05 06 08 09 0A 0B 0D 0E 0F 10 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-10-25 | 11.2.0 | * SoftBasic: 新增对url字符串编码和解码的方法，UrlEncode及UrlDecode方法，主要是中文和%开头的互转，代码参考微软实现。例如 AB好 和 AB%E5%A5%BD 互转 * HttpServer: 新增委托HandleFileUpload，当客户端上传文件时触发，然后获取到文件名和内容，可以自己进行存储，目前仅支持单个文件上传。 * DLT645-2007: 网口串口代码优化，相同的逻辑代码提炼，并增加属性 DLTType，返回枚举：DLT645Type.DLT2007 * DLT645With1997: 新增 DLT645/1997 协议的实现，使用ReadDouble来读取相关的数据，例如地址：B6-11 读取电压A数据。 * DLT645With1997OverTcp: 新增 DLT645/1997 网口透传实现的版本，使用TCP/IP通信并数据交互，使用方法和 DLT645With1997 一致。 * SiemensPPI: 优化西门子PPI协议数据接收时完整性判断，数据一旦接收完成，立即返回。因为在某些特殊情况下，数据接收不完整时会发生异常。 * SiemensPPIOverTcp: 接收数据时循环接收，并判断数据完整性，数据一旦接收完成，立即返回，增加循环超时时间，为ReceiveTimeout。 * CJT188: 新增CJT188-2004的水表（燃气表）协议，使用ReadDouble 和 ReadStringArray通信，地址示例：90-1F读取当前的累积流量，同时添加网口透传的版本。 * FanucSeries0i: fanuc的机床新增属性OperatePath，默认为1，当机床为多路径的时候，可以设置为2，然后对路径2进行相关的操作，例如读取报警，读取程序。 * MqttServer: 优化安全性，当没有指定RSA密钥时，针对每个连接的客户端对象，都创建随机的RSA密钥，来提升安全性。 * ILogNet: 日志类接口新增属性HourDeviation，如果需要日志记录时进行时间偏移操作，可以设置该属性实现，如果为-8，就是所有时间往前挪8小时。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-11-1 | 11.2.1 | * SiemensPPIOverTcp: 修复在通信速率比较慢的情况下，数据包发生多次接收的时候，接收不正确的bug。 * SecsValue: 修复在创建None类型的数据时，引发空对象的异常。 * SecsHsms: 优化Secs协议的客户端和服务器的心跳接收，修复了和其他客户端测试连接不正常的bug，DEMO界面添加大量的测试信息。 * CJT188: CJT188OverTcp 协议新增属性 StationMatch 用来标记是否匹配到站号一致的报文再返回，默认为false，修复读取数据失败时，直接报异常的bug。 * SiemensS7Server: 修复西门子虚拟PLC在使用node-red的s7组件访问时频繁在线，掉线的bug，原因来自消息id不匹配。 * DLTTransform: 当没有遇到内置转换DLT数据的时候，不再返回失败，返回数据的HEX字符串形式，用来支持扩展的情况读取。 * HslCommunication的netstandard版本的库所依赖的System.IO.Port 版本调整为 6.0.0 * Demo: 完善TCP/UDP调试窗体，Server端修复选择文本不显示字节长度的bug，修复关闭时连接不断开的bug，界面位置调整优化。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-11-18 | 11.2.2 | * SharpList: 修复Add一个新的值的时候，每隔指定长度次数时数据丢失一次的bug。新增LastValue( )方法获取最后一个值。 * AllenBradleyNet: ReadBool(string address, short length); 支持读取以int类型为基础的bool数组，支持从中间任意位置开始读，返回任意长度。地址前加i=,例如i=AAA[10] * NetSoftUpdateServer: 修复路径在linux下运行时，路径分隔符识别不正确，导致文件下载失败的bug。 * Modbus: modbus的地址支持自定义写入功能码，例如自定义的服务器有个数据区，读使用07功能码，写使用08功能码，那么地址可以表示为 x=7;w=8;100 依然可以支持站号。 * CJT188Helper: 修复CJT188协议再读取水表数据的时候，数据位颠倒的bug。 * Omron: 在欧姆龙Fins协议里，当读取位数据超过1998时会发生错误，现在根据1998长度自动切割，也就是支持任意长度了。 * ReverseBytesTransform: 修复指定DataFormat对象实例化的时候，指定DataFormat无效的bug。 * NetworkDoubleBase: 基础的网络设备通信基类支持了使用MQTT中转读取的方式，指定mqtt对象，读Topic，写Topic就可以远程读写PLC信息。可以配合我司的边缘网关系统的设备转MQTT使用。使用手册后面推出。 * DLT645With1997: 修复DLT645-2007解析数据异常的bug，修复读取电压倍率不对的bug。 * ModbusHelper: 修复汇川PLC读取bool时，当指定MX500.7能读bool值，但是MX501.7 这种地址时，获取不到正确的值的bug。 * Authorization: 支持一种在线激活方式，使用指定ip，端口，令牌获取远程激活码来激活的操作，后续扩展临时激活授权测试的功能。 * FanucSeries0i: 读取程序的API接口ReadProgram再接收完程序后，再接收一次0x18指令的报文。 * OmronHostLinkHelper: 欧姆龙HostLink协议在解析报文的错误码时，当解析失败时，增加错误捕获，并提示原始报文。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-12-9 | 11.3.0 | * byte[]: 类型byte[]新增加一个扩展方法，ReverseByWord，直接可以调用按照字反转字节数据信息。 * OmronCipNet: 欧姆龙的CIP协议优化，写入数组时，根据地址是否带索引，来自动确认写入的长度信息。 * MelsecMcServer: MC协议的虚拟PLC部分，在二进制下支持了随机字读取（0403）和块读取（0406）的功能码。 * WebSocketServer: 现在服务端新增属性：TopicWildcard，如果客户端使用通配符订阅数据点。规则和MQTT的通配符一致。 * HttpServer: 修复request.ContentType为空的情况下引发的异常，原因来自对文件上传的判断。 * SiemensS7Net: 修复西门子PLC在不连接的情况下，直接同步方法读写PLC相关的值的时候，因为pdu初始值不正常导致第一次读写失败的bug。 * Demo: Tcp的服务器调试界面里，修复点击发送客户端时，发送的数据不显示的bug。 * DLT645With1997OverTcp: 修复使用异步ReadDoubleAsync读取数据的时候，数据结果解析不正确的bug。 * OpenProtocolNet: 完善开放以太网协议，并且基本测试通过，支持自定义的功能码读取，支持订阅操作。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2022-12-28 | 11.3.1 | * Dll: 用于激活的类Authorization新增一个方法SetDllContact，用于设置提醒激活的联系方式，可以设置自定义的联系方式。 * Demo: 修复安川的MemobusUdpNet的测试界面的点击实例化按钮后，界面还是灰的bug。 * AllenBradleyNet: 修复ab-plc的CIP协议，当触发0x64错误码的时候，一直重复读写无法恢复成功的bug，现在会自动重连操作。 * SiemensS7Net: 新增对DTL格式的时间数据读写的API接口，使用ReadDTLDataTime和WriteDTLTime来读写操作 * LogStatisticsBase: 修复当系统的时间往回调的时候引发异常的bug，原因来自索引偏移位置会变成负数。 * websocketServer: Websocket服务器header接收订阅时，支持url转义，以及修复开启通配符模式情况下，消息驻留发送失败的bug。 * Demo: 串口调试界面以及网口调试界面，返回的数据支持勾选自动返回，只要在发送框里提前输入数据。 * Dlt645: DLT645的地址支持reverse标记，用来支持在某些不标准的协议的设备，数据顺序没有颠倒的情况，例如：reverse=false;02-01-01-00 * Inovance: 汇川的AM系列的地址，支持MD的地址格式解析，同时适用于读写操作，MD100=MW200 * SiemensS7Helper: 修复了西门子虚拟PLC及SiemensS7Net的写入字符串的同步方法，字符创长度信息和实际不匹配的bug。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-2-7 | 11.3.2 | * WebsocketServer: 在发布Topic数据的时候，新增一个API，允许指定reatain属性，如果不指定，则用对象自己的设置。 * InovanceHelper:汇川读取byte的方法新增错误异常捕获，返回失败的说明信息。 * SiemensS7Server: 修复虚拟服务器在接收远程客户端的写入定时器，计数器时，地址偏移不正确的bug。 * modbus: 当调用地址为 x=1;100 进行写入bool数据的时候，修复功能码不会自动修正为正确的默认功能码的bug。适用于 rtu, tcp, rtu overTcp * BeckhoffAdsServer: 修复当远程客户端读取M100.0 的bool类型时，指定长度大于1时，返回数据不正确的bug。 * BeckhoffAdsNet: 修复在读写I,Q地址类型时，初始偏移位置不正确的bug，无论是字读写，还是位读写都测试通过。虚拟PLC侧自动识别偏移地址。 * DLT645: 修复DLT645的1997及2007版本串口通信时，在某些特殊情况下通信速度慢导致接收数据长度不够引发异常的bug。 * MegMeet: 新增麦格米特的PLC的通信类，底层使用Modbus协议，支持串口，网口，支持系列：MC80/MC100/MC200/MC280/MC200E * KeyenceDLEN1: 新增基恩士的以太网模块通信协议实现，主要用于各种传感器的数据读写操作，提供DEMO测试。 * NetworkBase: NetworkBase基础通信基类，新增了一些基于SSL加密通信的报文收发的方法支持。 * MqttClient: MqttClient支持了使用证书的方式进行通信，设置MqttConnectionOptions.CertificateFile为证书路径即可实现证书加密及校验，已和EMQX测试通过。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-2-20 | 11.3.3 | * HslCertificate: 新增自定义的证书类，支持持有者颁发证书，对证书验证签名，例如可以使用证书实现对API接口权限验证，详细查看API文档。 * OmronCipNet: 修复异步写入字符串时，编码无法指定ASCII之外的bug，导致无法使用UTF8写入中文。 * ModbusTcpServer: 重新增加属性 StationCheck, ，当服务器只有一个站号的时候，设置为True表示校验客户端请求的站号，反之则不校验。 * DLT645及DLT645OverTcp: 支持写入字符串的方法，字符串将会转为二进制并且颠倒顺序，写入到仪表里去，写int及double数组也会转字符串数组写入。 * WebsocketServer: 服务器初步支持了SSL通信，使用证书通信，需要带私钥的证书才能成功创建，等待进一步的测试。 * LogNetBase: 日志获取相关的存储列表的方法中，增加一个错误捕获，在极端情况下可能会因为错误异常。 * Authorization: 激活系统的方法新增基于证书激活的方式，具体使用方法参考官网文档。 * Demo: Demo主界面的左侧的设备列表信息改为浮动窗体，可以自由隐藏，DEMO也支持了使用证书激活。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-3-25 | 11.5.0 | * DLT698: 698协议新增属性UseSecurityResquest，设置是否使用安全模式来通信，默认true，影响范围包括串口，透传，网口类，特殊的仪表需要设置 false。 * DLT698: 新增写入时间数据到地址地址的接口，WriteDateTime( string address, DateTime time )，新增698协议的接口 IDlt698 * FanucInterfaceNet: Fanuc机器人对象读取机器人数据时，FanucData解析字符串数据的编码可以手动调整指定。 * MelsecFxLinksHelper: 三菱的计算机链接协议里，如果传入的报文以ENQ开头的，则认为传入原始包，不进行封包处理。 * OmronFinsAddress: 欧姆龙的Fins协议支持CF地址，可以读取脉冲信号，地址一秒脉冲地址示例：CF1.2 新增IHostLinkCMode, IOmronFins接口。 * 新增三菱的编程口接口 IMelsecFxSerial，新增西门子PPI接口 ISiemensPPI，永宏编程提炼接口 IFatekProgram * AllenBradleyNet: 批量读取地址数组的方法，传入长度数组参数由 int 改为 ushort 类型，所有的CIP协议实现了统一的IReadWriteCip接口。（可能不兼容） * FanucSeries0i: ReadProgram方法取消再次接收0x18的命令，因为在 32i 的机子上是会触发超时接收。 * YamahaRCX: 雅马哈机器人新增 JogXY 接口，使用微动的方式操作六轴的动作信息，在demo界面上支持JOG的测试。 * AllenBradleyHelper: 修复CIP协议创建命令为0x4D写入数据报文时，数据长度过大导致溢出的bug，在解析返回结果数据时增加错误捕获操作。 * Modbus: 所有Modbus的客户端以及RTU主站类，接口IModbus新增使用功能码0x17的一条报文实现读写方法ReadWrite，虚拟服务器支持该功能码。 * SiemensS7Net: Read( string[] address, ushort[] length )方法不仅仅根据19个长度切割，还根据长度之和是否满足200字节切割。 * FinsTcp: Fins协议读取bool时，发现返回数据长度小于指定长度时，自动偏移地址并读取剩余数据，然后返回全部给调用者。 * SimpleHybirdLock: 通信的混合锁新增当前锁对象等待的计数属性：LockingTick，可以用来监视锁对象的竞争情况。 * Demo: TCP转TCP的支持二进制和ASCII格式的转换显示，完善MegMeet，Lsis的添加，批量读取控件增加提示功能。 * HslCommunication: 大部分PLC通信的消息解析代码增加合法性验证，并且增加异常捕获，防止及其特殊的情况导致系统奔溃。 * Demo: Demo程序大升级，测试界面支持定时读取，定时写入，写入支持动态数据变化。批量读取，报文读取，特殊功能测试优化，更加便捷。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-4-14 | 11.5.1 | * MelsecFxSerial: MelsecFxSerial的报文日志记录格式调整为ASCII格式，当设备的日志实例化时，就会自动当前的通信日志。 * IReadWriteMc,IReadWriteA3C: 新增属性EnableWriteBitToWordRegister，当开启后支持往D寄存器写入位数据，实现方式为读字，修改位，写字操作，具体风险查看该属性说明。 * KeyenceMcNet, KeyenceMcAsciiNet: 基恩士的MC协议的通信，修复ZR地址因为地址进制导致读取不到正确地址的bug，现在使用16进制地址表示。 * KeyenceMcNet, KeyenceMcAsciiNet: 地址兼容了基恩士自己的地址格式，可以直接输入 DM100, FM100, EM100, MR015, 或是直接 MR1.2 等 * HslHelper: 新增属性LockLimit表示竞争锁上限(默认1000)，NetworkDoubleBase,NetworkUdpBase,SerialBase检测当锁竞争达LockLimit次时，直接返回失败，不在增加竞争。 * HslTimeOut: 检查socket超时部分的功能代码，当处于linux系统下时，在确认超时close之前，增加一个 Disconnect操作，解决特殊情况close会卡时间的bug。 * HslHelper: HslHelper新增属性UseAsyncLock, 标记本通信库的单个通信对象在异步通信的时候是否使用异步锁，默认True, 当使用控制台或是纯后台线程采集时，配置 False 更好。 * Modbus: 在写入单个的short及ushort值时，默认使用06功能码，现在支持如果在地址里使用w标记，例如 w=16;100 时，写入时报文就使用16功能码。 * TcpForward: 新增属性CacheSize，表示中转时候的数据缓存大小，默认是2048 * Net2.0, 3.5, Standard项目添加缺失的类：VigorServer, OmronHostLinkCModeServer, MelsecFxLinksServer, KeyenceDLEN1 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-4-19 | 11.5.2 | * OmronCipNet: 修复写入单个Byte数据对象，或是奇数长度的byte[]时，plc报错31错误码的异常，现在可以正常的写入了。 * SiemensS7Net: 西门子批量读取地址数组的方法里，原先总长度按照200字节切割变更为按照自动获取的pdu长度进行切割。 * MqttClient: Mqtt客户端的连接类MqttConnectionOptions新增属性UseSSL指示是否开始SSL/TLS加密功能，验证证书时修复服务器名称输入错误导致有些服务器连接失败的bug。 * NetworkUdpBase: 优化UDP通信基类，增加设置获取PipeSocket方法, 支持设置多个端口，因网络读取失败时自动切换另一个端口。 * Modbus: ModbusRtu设备在提取接收到的报文时，校验数据长度并且如果发现长度太长，则按照标准报文长度切割，即使后面跟0xFF扰码也能正确读取 * Memobus, DigitronCPL, MegMeet, Delta 相关的通信类，补全缺失的 MPRC 接口注册功能。 * Demo: 三菱相关的PLC测试界面在连接PLC失败的时候，原先只提示连接失败，现在提示更加详细的信息。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-5-15 | 11.5.3 | * NetworkBase: 在ReceiveByMessage方法里，新增判断消息报文剩余长度小于0时，直接返回头报文信息，防止在万一接收到错乱数据导致小于0引发bug的异常。 * FFTFilter: 傅立叶变化的Filter滤波方法，修复当数据有负数的时候，只能得到正数结果的bug，现在可以还原出带有负数情况的波形图。 * SiemensS7Server: SiemensS7Server增加在握手报文处理时，对握手报文合法性的检查操作，防止收到一些奇怪数据时导致奔溃。 * DLT698: 修复在UseSecurityResquest为默认True情况下，APDU指令创建不正确的bug，涵盖DLT698, DLT698TcpNet, DLT698OverTcp类。 * FanucSeries0i: 新增根据字符串程序名读取程序的方法接口ReadProgram( string program, string path = "" )，用来支持读取非数字名称的程序号。 * Toyota: 新增丰田工机的计算机链接协议实现，ToyoPuc 类，通过2PORT-EFR模块实现对PLC的读写数据操作，同时增加ToyoPucServer虚拟PLC进行测试。 * DLT698: 修复698的协议中的CheckResponse，当请求类型不支持时，引发空对象的异常。 * SerialBase: 间歇时间SleepTime允许设置0及小于0，小于0就是不进行任何休眠，三菱编程口MelsecFxSerial增加对消息完整长度的校验机制。 * OmronCipServer: 新增欧姆龙自身的CIP服务器，支持客户端使用OmronCipNet类和OmronConnectedCipNet进行通信，支持普通点位，数组，字符串的读写操作。 * CipServer: AllenBradleyServer和OmronCipServer支持所有的标签数据存储到文件和从文件加载的操作，demo里服务器的存储加载支持路径选择。 * MelsecA1EServer: 新增模拟PLC的实际情况，针对bool读取长度超过256返回错误码，针对字单位读取长度超过64即返回错误。 * MelsecA1E: MelsecA1ENet及MelsecA1EAsciiNet类修复读取bool数组时长度超过256报错的bug，MelsecA1EServer完善读取bool超过256返回错误。 * Mewtocol: 松下的Mewtocol的串口及网口类对象新增支持读取PLC型号的方法ReadPlcType，虚拟服务器也支持了该功能码。 * SerialBase, NetworkDoubleBase: 优化通信代码，针对Thread.Sleep方法也进行错误捕获，防止可能的异常。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-6-25 | 11.6.0 | * DLT645: DLT645-2007 新增对设备返回的报文进行和校验码校验的步骤，如果校验失败，返回错误信息。 * MelsecA1EAsciiNet: 修复在字读取长度超过64字时，或是bool读取超过256位时，读取结果不正确的bug。现在支持读取任意的长度数据。 * AllenBradleyServe: 修复AllenBradleyServer及OmronCipServer在配置了写入时创建新标签后，写入字符串数据仍然报错，提示标签不存在的bug。 * IDlt645: DLT645接口新增Password属性及OpCode，方便在通信时动态修改这两个参数，主要针对DLT645/2007的,对于 DLT645/1997 协议来说无效 * MelsecMcServer: 修复MC虚拟服务器在ASCII模式下，写入D之外的寄存器，但是仍然写入D寄存器的bug，影响范围包括A3CServer。 * Melsec: 三菱的MC协议开放属性 PlcNumner，可以自由的设置，表示PLC编号。影响范围:MelsecMcNet, MelsecMcAsciiNet, MelsecMcUdp, MelsecMcAsciiUdp * OmronFinsServer: 欧姆龙的FinTcp协议的虚拟PLC支持了CF数据地址，无论是服务器上，还是通过客户端来读写CF数据都可以正确的读取。 * FanucSeries0i: 客户端demo修复删除程序，只能删除主目录的程序号的bug，现在可以删除任意指定路径的程序文件信息。 * WebsocketServer: websocket服务器运行时，当客户端使用火狐浏览器连接服务器时，修复检测websocket连接失败的bug。 * ToyoPuc: 修复丰田工机PLC中，当地址是 U0 及 H0 的时候，地址输入解析不正确的bug，原先的bug是忽略第一个数字。 * AllenBradleyServer: CIP虚拟服务器优化通信的细节，正确的设置了基于连接模式下的各种连接ID信息，当PLC使用MSG模块读写虚拟PLC时也顺利通过。 * Demo: HslCommunicationDemo程序几乎所有的设备测试界面的线程压力测试界面单独拎出来，显示单次通信的平均耗时。 * Demo: HslCommunicationDemo程序几乎所有的设备测试界面增加设备的地址示例说明，带地址示例，含义，部分PLC提供地址范围说明。 * Demo: HslCommunicationDemo程序几乎所有的设备测试界面增加数据点位表的功能，可以配置点位表，保存，加载，然后多个点位同时刷新。 * .net standard 框架的dll依赖的 System.IO.Ports 版本由 6.0.0 升级到 7.0.0。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-7-18 | 11.6.1 | * Demo: 保存设备的连接参数及点位数据的时候，支持了设置密码，设置密码后，加载该设备参数时，必须输入正确的密码才能打开。 * Demo: 串口调试界面，TCP调试界面，TCP Server调试界面大幅度的调整，优化，更好的区分收发数据，以及数据长度信息。 * Demo: 在demo界面的曲线控件界面，从实时曲线修改为历史曲线控件，新增设置颜色，曲线样式，显示最大值，最小值信息，自动滚动条往右。 * Demo: 调试的程序在连接成功，或是打开串口之后，直接在测试界面上显示当前的连接类的示例代码，初始化参数和demo界面保持一致，方便直接复制。 * ModbusTcpServer: modbus虚拟plc当属性UseModbusRtuOverTcp = true时，同时支持rtu over tcp报文及 ascii over tcp报文，并对 ascii 报文做了更加严格的区分，防止弄乱。 * ModbusRtuOverTcp: 通过调整底层基类的读取空消息的报文时，可能引发数据不完整的bug，针对部分回数据断断续续的设备，现在可以读取正确的值信息。 * AllenBradleyNet: 优化CIP协议的读写字符串功能，对长度判断更加的完善，当写入字符串时，可以指定额外的type值，例如地址为"type=0xD1;AA"，影响范围：AllenBradleyMicroCip * HslTimeOut: 优化超时的代码信息，单独添加超时的日志信息，如需记录，可以实例化 HslTimeOut.TimeoutLogNet 属性，并且可以获取到当前所有的超时判断次数 HslTimeOut.TimeoutDealCount * SecsHsms: 优化secs的代码，增加在解析secs消息的时候，对异常错误的消息捕获并记录日志，demo界面优化显示信息，并提供代码示例。 * AllenBradleyNet: 写入字符串的时候，如果类型为 DA 的情况，则使用另外的写入方式，在microcip上测试成功。 * PanasonicHelper: 修复松下的mewtocol协议使用字的方式写入数据到Y,R,L地址时，引发0x41错误码的bug，原因来自结尾地址偏移异常。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-8-13 | 11.6.2 | * YRCHighEthernet: 新增ReadManagementTimeSpan接口，用来读取持续的时间，返回的是字符串信息。完善demo界面，提供示例代码信息。 * NetworkDoubleBase: 新增一个发送前的报文头设置属性 SendBeforeHex，当使用lora中转通信时，可以设置为lora站号，例如00 00 00 02 * ModbusUdpServer: 新增基于UDP通信的ModbusUdpServer虚拟服务器类，在demo界面上的ModbusTcpServer启动时，可选TCP还是UDP功能，方便测试。 * MelsecFxLinksOverTcp: 修复属性 WaittingTime 可以设置0xF 以上值的bug，实际这个值不允许设置 0x0f 以上。 * XinJEServer: 当服务器站号不一致的时候，并且服务器设置检查站号的情况，直接返回站号不一致的错误消息。 * KeyenceNanoServer: 基恩士的上位链路虚拟PLC的字符串是否反转默认设置 true, 如果需要，手动修改server.ByteTransform.IsStringReverseByteWord * FanucSeries0i: 新增方法 ReadAxisNames 读取轴名称以及 ReadSpindleNames读取主轴名称。 * IEC104: 添加总召唤方法TotalSubscriptions，优化值解析的过程，Demo界面进行了优化，显示报文的时候，显示其对应的功能说明，api文档新增示例。 * AllenBradleyNet: 修复当PLC对象修改了ByteTransform.DataFormat的值后，导致读取PLC时提示会话ID不一致的bug。 * DEMO: 在TCP Server调试的界面上，增加保存当前界面的配置参数的功能，连带配置的报文列表也一并存储。 * DEMO: 修复DebugControl控件在窗体拖动时，重复触发发送数据按钮的bug，多来回拖动后将会导致程序异常卡顿。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-9-25 | 11.6.3 | * ModbusTcpServer: 当客户端使用了0xff的站号（10进制为255）的时候，将自动忽略对站号的校验，直接可以进行读写操作。 * SiemensS7Net: 新增强制位ForceBool方法，以及取消强制CancelAllForce方法，目前针对200smart测试通过，支持对I,Q的点位强制操作。 * OmronFinsServer: 优化返回的报文信息，方便wireshark软件抓包进行分析操作，识别到正确的Fins协议。 * SiemensS7Net: 西门子s7协议的报文中，消息号id修改为自增，修复了在某些特殊的200smart型号读取超时的异常。 * OmronCipNet, OmronConnectedCipNet: 读取字符串的方法，在解析字符串数据内容的时候，增加错误捕获，并返回数据的二进制内容。 * NetworkDeviceBase: 写入数据的几个类型的重载方法，也更改为虚方法，支持子类重写操作。 * SiemensS7Net: 新增方法 Write( string[] address, List<byte[]> data )支持同时写入多个数据块的地址数据，SiemensS7Server也支持多个数据块的写入功能。 * Melsec: 三菱的协议MelsecA1ENet，MelsecA1EAsciiNet, MelsecFxLinks(含串口及以太网), MelsecA3CNet(含串口及以太网), 读取bool时支持字地址的位，例如 D100.0 * CJT188Helper: 水表协议188修复再某些特殊的情况下解析结算累积量异常时，导致读取当前累积量也失败的bug。 * OmronFinsNet: 新增属性 ReceiveUntilEmpty, 默认false，设置true后可以在一些及其特殊的场景里防止数据错误的情况。 * LogNetManagment: 日志管理对象的一个静态方法优化，增加错误捕获，防止极小概率的意外情况。 * SiemensS7Plus: 新增基于s7plus协议的通信对象，使用SSL/TLS加密通信，支持读写数据，不需要开启put/get，支持遍历点位，仅支持.net framework平台，还需要依赖两个openssl的C++库。 * Demo: 支持将PLC数据定时读取，然后导出到单个文件。然后可以在虚拟服务器上进行加载文件，然后定时写入，可以本地模拟现场设备数据。 * Demo: 原始字节的批量读取功能显示效果，新增整型显示，当切换hex，ascii，整型的时候，自动更新显示数据，更加便捷。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-10-26 | 11.6.4 | * ByteTransformHelper: 一个转换的静态辅助方法GetSuccessResultFromOther增加异常的捕获，防止在一些极其特殊的情况下导致异常然后奔溃。 * IReadWriteMc: MC协议的接口新增属性TargetIOStation，用来表示目标模块IO编号，默认0x03FF, 表示管理的CPU，当访问多CPU时，需要设置本值。 * ABBWebApiClient: abb机器人的类新增接口GetUserValue，用来获取用户自定义的一些数据内容，返回double[]，demo增加测试功能。 * ModbusRtu: rtu报文的解析新增加对头部干扰码的支持，就算报文前跟了一个字节的无用报文，现在也可以正确的解析，也适用于rtuovertcp. * MqttServer: 文件引擎新删除整个目录的接口方法，重命名目录的接口方法，MqttSyncClient代码优化，并且添加配套的接口实现。 * HttpServer: 新增Action类型属性ApiCalledAction，当一个接口调用完之后触发，可以获得该接口一些调用信息，方便二次处理。 * BeckhoffAdsNet: 新增加ReadStruct<T>( string address )方法及WriteStruct<T>( string address, T value )针对struct类型的对象读写。 * Demo: 在所有的PLC设备的读写测试的控件里，在点击了byte类型的读取时，支持了设置长度大于1的情况下，批量读取了byte数据，显示为十进制数据。 * Demo: DEMO中的所有的PLC测试界面优化，在没有连接PLC的情况下，也可以切换选项卡查看地址示例。 * Nuget: netstandard框架的项目依赖的System.IO.Ports组件版本从7.0降级为6.0，原因是7.0在某些特殊的ARM硬件里存在问题。 * HslHelper: 方法CalculateBitStartIndex优化匹配代码，使用了正则表达式进行判断，提升程序速度。 * SerialBase: 串口基类，网口基类的一点点代码优化，针对及其特殊情况的一点异常捕获。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2023-11-10 | 11.7.0 | * SiemensS7Net: 西门子的用于计数的消息id在重连plc的时候，自动重置为1，然后重新开始计数。 * DLT645: DLT645协议新增对返回的站号检查是否一致的功能，如果刚好站号倒序，则也认为一致。优化消息接收，即使报文前面跟了几个字节的无用信息。 * ModbusRtuOverTcp: 优化通信机制，对于单次接收的报文支持了更多的容错，如果CRC校验失败，不再关闭当前的socket连接操作。 * MqttServer: 当客户端订阅一个主题之后，如果这个主题属于消息驻留的主题，给当前客户端返回的消息里增加Retain=true标记，方便客户端标记。 * MqttServer: 注册MRPC的接口时，如果之前已经有相关的接口注册，现在改为直接覆盖，然后变量参数名自动适配 value 和 values，这两个没找到就相互替换名称。 * DLT645Server: 新增DLT645-2007协议虚拟仪表实现，支持一些常见的数据的简单测试，例如电压，电流，频率，时间等信息。 * Newtonsoft.Json, NET20项目，NET35项目，NET451项目对JSON库的引用调整SpecificVersion为 False * Demo: 地址示例的表格支持右键弹出上下文菜单，点击复制地址，即可复制当前的地址信息。 * Demo: 数据点位表控件DataTableControl修复虚拟服务器不能刷新数据的bug，新增了双击值表格写入数据的功能。 * Demo: 修复DEMO服务器的部分ip地址获取网址解析不出数据的bug，所有的襄阳地址均改为襄樊，方便旧版地址显示。 * MqttClient: 不兼容更新，收到MQTT的事件签名OnMqttMessageReceived修改为MqttMessageReceiveDelegate( MqttClient client, MqttApplicationMessage message );支持携带Retain信息。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-01-26 | 11.8.0 | * SiemensS7Server: 西门子PLC的虚拟服务器，当客户端连接的时候握手报文也进行输出到日志上面。 * OpenProtocolNet: 新增属性RevisonOnConnected，用来设置连接初始化时的MID0001的revison信息，如果小于0，初始化连接时则不使用MID0001 * AdsHelper: 倍福的ADS协议里的部分错误码提示的消息附带一些建议的解决方法，可以快速查找问题。 * SiemensS7Helper: s7协议里如果执行写入wstring类型的字符串时，当检测到字符串最大长度为0时，自动分配254长度的字符串。 * MC1EServer: MC协议的虚拟服务器及MC1EServer的虚拟服务器支持了S0,F0 地址，也支持远程客户端读写该地址。 * MqttRpc: 在Mqtt服务器注册了RPC服务器接口后，当注册的自定义的接口方法本身发生异常时，增加提供原始异常相关的本文及堆栈信息。 * SiemensS7Net: S7协议新增对地址 SM0.0 的支持，在200系列上测试成功。 新增 PIW0, PQW0 外设地址，在S7-300上测试成功。 * HslReflectionHelper: 解析HslStructAttribute特性时，类型为字符串时，修复编码为GB2312，但是编码设置异常的bug，现在自动忽略大小写。 * OmronFinsTcpServer: 虚拟PLC支持了CF地址及DR地址类型的读写操作。 * OmronCipServer: 当服务器创建了bool[]数组的标签时，修复客户端读写该标签不一致的bug。 * IEC104: 104规约总召唤方法的传送原因默认值改为06，新增属性Station表示单元公共地址，并在DEMO上可设置。 * ABBWebApiClient: abb机器人代码优化精简，部分接口自动解析更多的数据，demo测试界面显示优化，同时显示解析结果和html文本。 * OmronFinsUdp: 报文里的SID字段每次通信将自增1，然后校验返回报文是否一致，虚拟服务器适配客户端。 * OmronFinsNetHelper: Fins协议针对EndCode结束码判断增加网络异常标记位判断，直接返回读写失败及错误消息。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-01-29 | 11.8.1 | * MelsecMcServer: 修复MC虚拟服务器当远程客户端读取D，W，R，ZR字地址数据时回复不支持的bug，影响范围包括Udp,A1E,A3C服务器。 * NetworkDataServerBase: 修复虚拟服务器基类启动串口接收后，触发OnDataReceived事件的报文不是接收到客户端报文的bug。 * NetworkDataServerBase: 修复虚拟PLC服务器基类的事件OnDataSend没有触发的bug，就是服务器回发客户端数据的事件。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-02-14 | 11.8.2 | * OmronFinsNetHelper: 欧姆龙的Fins协议的从错误码获取消息文本的方法修复部分错误码获取文本不匹配的bug。 * MqttServer: MQTT服务器针对每个会话的客户端(针对客户端启用了RSA加密的情况)使用临时生成的AES随机密钥，防止AES密钥被伪装客户端获取。 * NetworkDoubleBase类的ConnectClose方法以及ConnectCloseAsync方法取消了锁操作，一般建议系统退出时再调用一次，demo上接收超时时可以立即关闭连接。 * NetworkBase类以及NetworkDoubleBase类所有的异步操作方法添加.ConfigureAwait( false )操作，防止在部分特殊的情况下使用异步导致死锁卡死UI * MqttServer: 支持MqttSyncClient客户端获取服务器当前的所有的所有在线的会话信息，并增加发布的topic统计功能，也就是客户端可以远程查看服务器在线会话，需要开发者权限帐户。 * DEMO MqttServer服务器的测试界面显示的内容优化，显示了在线客户端的详细信息，显示了客户端发布的主题的基本信息。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-04-28 | 12.0.0 | * Controls: 删除controls控件，删除图标资源，FormSupport窗体移动到DEMO，FormPortraitSelect移动到hslcontrols组件，本库更加纯粹的负责通信及算法。 * SharpList: 更改锁的机制为lock，优化了代码注释以及部分代码逻辑，提升数组长度比较大的时候的性能，提升批量增加数据时候的性能。 * NetSoftUpdateServer: 用于程序升级的服务器修复当传送大文件（120M以上时）传送到一半会卡住，传送非常慢的bug。 * Modbus: 修复ModbusTcp协议当设备有前置干扰码但是数据长度判断不正确的bug，ModbusRtu新增站号不匹配时持续接收匹配站号的数据。 * SecsMessage: 修复secs消息解析的时候，屏蔽了function字段的最高位的bug，该bug会导致大于127的功能码解析异常，现在可以正确的识别。 * FormLogNetView: 日志显示的窗体新增带日志路径的构造函数，当指定路径后，显示页面时立即显示日志内容及统计信息。 * OmronFinsTcp: 增加对属性SID信息功能，每次请求的时候进行自增，然后校验PLC返回的数值信息，直到接收到一致为止。 * ABBWebApiClient: GetJointTarget方法新增mechunit参数，可以指定不同单元。新增GetAnIOSignal方法。 * CommunicationServer: 新增服务器类，可以开启TCP,UDP,SerialPort服务器，绑定事件即可接收到客户端发送的数据。 * FujiCommandSettingType: 修复写入数据到plc，提示写入失败但是实际成功的bug; 修复西门子PPI虚拟服务器对写入结果的报文返回不正确的bug。 * SecsHsms: 针对SendByCommand方法，如果通信异常，则自动连接操作，也就是所有的发送指令都会将管道标记为异常，下次发送时进行重连。 * DLT645Server: 修复主站读取时站号不一致时提示校验失败的bug，修复修改站号命令失败的bug，新增属性站号是否校验。 * HslHelper: 修复GetIpAddressFromInput方法输入ipv6地址时会导致解析失败的bug，现在同时支持输入ipv4地址，ipv6地址，网址。 * Modbus: 新增属性EnableWriteMaskCode，默认支持掩码功能码，如果发现设备不支持，会自动切换false，使用读字修改位写字的方式来实现写位到字寄存器。 * ModbusTcpServer: 新增自由控制是否支持掩码功能码; 台达，麦格米特，汇川，信捷的modbus类地址功能优化，并支持了部分字寄存器的写位操作。 * ByteTransform: ByteTransformBase基类删除，字节转换类优化，两个字节的整型short,ushort将根据ABCD,CDAB时使用大端，其他情况使用小端转换。 * AllenBradleyNet: 地址支持了x=0x52;A1 这种携带52功能码直接使用片段读写的方式，还支持了class=0x6b;0xf68f这种符号实例ID地址，服务器增加配合，修复相关bug。 * LogNet: 存储日志因为文件占用等问题发生异常的时候，不再记录存储失败的信息，存储多次才成功的日志加入[Retry:3]标志，相关代码精简优化。 * SecsValue: 新增加一个从string[]对象初始化的构造函数，直接生成一个list对象，包含多个字符串。 * DLT698: DLT698协议新增客户端地址CA属性，可以配置更改，代码优化，修复读取电能量数据的扩展参数属性时，缩小倍数不正确的bug。 * IOmronFins: FinsTcp,finsUdp,hostlink协议新增属性PlcType，可选CS/CJ 和 CV系列，在CV系列地址码略有不同，fins协议支持ReadCpuTime读取PLC时间。支持了地址 E100 当前数据区 * BeckhoffAdsNet: 倍福后台接收数据的地方针对设备通知消息（08功能码）的情况，进行数据屏蔽操作，demo添加ig=0xF030地址格式的示例。 * KeyenceMcNet: 修复基恩士MC协议在读取R100.5， MR100.5，LR100.5，CR100.5 地址格式转化真实位地址时地址有偏差的bug。 * WebSocketClient: 新增方法 UseSSL( string certificateFile )支持开启证书模式，可以连接 wss 的远程服务器。WebSocketServer也做了代码优化。 * V12: 通信类继承体系重新设计, BinaryCommunication(二进制通信)->DeviceCommunication(设备通信)->TCP,UDP,串口通信->每个设备的通信。 * V12: 重新设计通信的管道，CommunicationPipe(通信管道)->TCP,UDP,串口,MQTT管道，前者的二进制通信类，可以随意的设置不同的管道，具体参考手册。 * V12: 重新设计所有的虚拟PLC类，DeviceCommunication(设备通信)->DeviceServer(服务器)->各个PLC类，同时支持开启TCP，UDP，串口从站的功能。 * Demo: 串口网口的收发数据的控件，为了方便模拟设备，输入框支持分批并且延时发送，单独一行<sleep=100>即可，数字100表示休眠的时间，可以修改其他值。 * Demo: 修复demo的PLC测试界面里，定时读取的时候，PLC发生掉线后不停弹窗的bug，现在代码进行了优化精简。 * 请注意，本次升级可能会造成不兼容的情况，具体需要看使用了哪部分的功能，如果简单的使用PLC类例如ModbusTcpNet可以兼容升级，如果使用了基类NetworkDeviceBase，需要替换为DeviceCommunication。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-05-24 | 12.0.1 | * MqttServer: 当启动文件服务器时，客户端进行文件上传，下载，删除触发的OnFileChangedEvent事件里添加一个属性：映射文件名称。 * PipeSslNet: 新增基于 SSL/TLS 的管道信息类 PipeSslNet, 从PipeTcpNet类继承实现，可以配合 HslCommunication 自身的虚拟服务器实现加密通信，防止抓包。 * CommunicationTcpServer: TCP服务器基类代码优化，添加对SSL/TLS的支持，签名修改：ThreadPoolLogin( PipeTcpNet pipe, IPEndPoint endPoint ) 这里如有使用可能不兼容升级。 * ModbusTcpServer: 修复当DataFormat配置DCAB时，解析客户端的报文的时候，导致地址及长度信息高地位弄反的bug。 * PipeSerialPort: 串口管道新增属性DtrEnable控制，表示串行通信中是否启用数据终端就绪 (Drt) 信号 * WebSocketHelper: websocket消息接收使用了管道的实现，删除了原先部分多余的代码，mqtt管道公开一些额外的参数信息。 * DeviceSerialPort: 当串口管道调用关闭方法Close的时候，修复管道被重新指定其他非串口管道时，连接不会关闭的bug。 * DeviceCommunication: 设备类基类重新添加Dispose接口，然后调用管道的Dispose, 方便代码升级考虑兼容性。 * BeckhoffAdsNet: 修复管道配置为SSL及自动获取AMS的情况下，获取ams不使用ssl的bug。 * DeviceTcpNet: 重新添加回方法SetPersistentConnection()并标记弃用的状态，方便部分用户升级的时候可以保持向后兼容。 * DLT645-2007, DLT645-1997, DLT698添加无参的构造方法，这样DEMO生成的示例代码就可以复制使用，然后优化DEMO上的界面。 * MelsecMcServer: 三菱的虚拟服务器针对字地址支持了bool读写操作，也就是ReadBool("D100.1")，写入也是一样。 * HttpServer: 注册接口的时候，支持输入参数HttpListenerRequest request时，自动传入客户端的请求头，方便进行二次分析处理。 * Wecon: demo新增维控PLC的测试界面，如何实例化相关的对象及支持的地址列表参考demo界面上的地址示例及代码示例。 * KeyenceNanoServer: 针对客户端单次读写的数据长度从256，提升到1000长度，因为上位链路协议在某些型号里是1000长度的。 * Omron: 修复欧姆龙CV系列的情况下，针对CIO数据区执行bool写入操作异常的bug，原因来自PLC不支持bool的写入CIO功能码，使用读字修改位写字来实现。 * Inovance: 修复了Modbus协议地址映射后，针对写入字地址的位情况下可能引发了bug，修复汇川的AM系列读写QX，MX位不匹配的bug。 * PipeSerialPort: 优化串口管道的数据接收机制，当定义了报文消息的完整性校验对象后，就不再以收到空数据为终结点，如果数据不完整，就一直接收到超时为止。 * ModbusRtu: 修复当设备的数据分段延时返回的时候，设备数据的内容刚好符合一些错误规律的时候，会直接返回部分的报文给调用者，然后判断出 crc异常 * MqttClient: 修复调用ConnectServerAsync方法连接的情况下触发的OnClientConnected事件中，获取IsConnected属性还是未连接的bug。 * LSis: 删除原来按照XGB,XGK这种命名方式，改为了LsCnet, LSFastEnet类，具体类型的实例化，可以参考demo程序生成的代码示例。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-6-3 | 12.0.2 | * OmronFinsServer: 客户端握手报文返回时的命令码修改为1，因为在某些情况下第三方客户端会连接失败 * CommunicationServer: 新增管道上线事件OnClientOnline，管道下线的事件OnClientOffline * PipeSerialPort: 修复串口管道当设备方一直不间断发送数据的情况下导致始终不引发超时的bug，影响范围包括所有的串口类设备。 * Lsis: lsis的代码优化，支持了DB100这种字节地址，支持直接读short，int多字节数据，修复地址转换异常的bug。 * PanasonicMewtocol: 修复日志记录的时候不按照ASCII格式记录的bug。 * Demo: 服务器端新增数据模拟的界面，可以使用表达式来动态执行脚本，方便生成一个特殊变化的曲线。 * Modbus: 修复Modbus协议在DataFormat为 BADC 及 DCBA 的情况下，写入short及ushort字节不发生转换的bug，影响范围为所有modbus。 * FanucSeries0i: 初步添加读取诊断信息方法ReadDiagnoss( int number, int length, int axis )，demo上可以直接测试，欢迎反馈信息。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-7-10 | 12.0.3 | * MqttServer: 修复在客户端的掉线事件里如果调用RemoveAndCloseSession移除会话导致内存堆栈溢出的bug，虽然一般都不会这么写代码。 * MqttSession: 属性DeveloperPermissions默认值修改为false，在服务器上会话登录后admin用户名默认设置为Ture，也可以自定义更改。 * NetSupport: SocketSend及SocketSendAsync方法增加对socket空的校验，防止部分plc实例化后直接调用close直接报异常的bug。 * HttpServer: 修改HandleRequest方法及属性HandleRequestFunc返回方法为object, 支持返回string及 byte[]，可以用来传文件，此更新为兼容的。 * MemobusTcpNet: 安川的协议地址支持了M100,G100,I100,Q100,S100这种地址，支持超大地址，支持对位原生操作。 * NetworkConnectedCip: 两个连接时的会话ID属性OTConnectionId及TOConnectionId修改为公开状态。 * WebSocketClient: 新增属性GetCarryHostAndPort用来标记HTTP请求头GET是否协议HOST信息，默认为false，DEMO界面优化，添加了输出示例代码，修复某些特殊服务器连接失败的bugs。 * OperateResult: 创建失败的类型返回方法新增一个重载的字符串参数，可以额外添加信息说明，方便后续查问题。 * SecsHsmsServer: 新增一个PublishSecsMessage重载方法，支持指定是否要求客户端返回，修复demo上勾选要求返回失败的bug。方法SendByCommand新增带wbit是否要求返回信息的重载方法。 * OpenProtocolNet: 修复创建报文参数异常，修复并优化demo界面问题，改从TcpNetCommunication继承，新增OpenProtocolServer配合进行本地测试。 * BeckhoffAdsNet: 修复读取的字节数大于10000时返回数据长度不正常的bug，修复读取bool数组时返回数据数量不正确的bug，虚拟服务器修复位读取时的返回数据不正确的bug。 * PanasonicMewtocol: 松下的Mewtocol协议支持读取 D，LD寄存器的位，地址格式为 D100.0， LD100.0 ， 影响范围包括网口及串口。 * BeckhoffAdsServer: 修复服务器端自己读写数组的某个索引的数据的时候，提示标签不存在的bug，暂不支持客户端读读写某个索引的标签数组。 * ModbusMappingAddress: 新增英威腾PLC的modbus地址转换实现，现在modbus可以注册这个地址转换，直接支持PLC原生地址。 * Demo: 数据表控件支持了曲线实现功能，kuka机器人界面增加更多的地址说明，fanuc的机床添加下载程序到本地。 * Demo: demo程序的PLC侧的数据表监视界面支持对值进行脚本运算，例如缩小100倍，输入 x/100.0 或是 x\*0.01，也可以更加复杂三角函数运算 * Demo: demo程序优化，写入bool支持输入0, 1，然后demo关闭服务器自动停掉所有的定时器及线程操作，否则将会触发异常退出。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-9-3 | 12.1.0 | * FanucRobotServer: 修复客户端读取M,I,Q数据区的时候，服务器从内存数据区取数据时偏移量不正确的bug。 * SecsHsmsServer: 新增属性DeviceId，当调用方法 PublishSecsMessage 时使用该属性指定的设备ID信息来主动发送客户端。 * OmronFinsServer: 虚拟PLC支持了读取plc型号及plc状态的功能码，修复如果遇到不支持的功能码时候直接奔溃的bug，现在返回错误的命令。 * FanucSeries0i: 新增加接口读取刀组寿命ReadToolInfoByGroup，读取正在使用的刀组号ReadUseToolGroupId，以及清除刀组ClearToolGroup。 * HttpServer: 注册API接口时从之前的只支持GET/POST的方式，现在可以随意的支持其他的文本了，例如PUT * DeviceServer: 虚拟服务器基类针对ReadFromCoreServer进行异常捕获，避免一些特殊情况下乱七八糟的数据导致虚拟服务异常奔溃。 * FormLogNetView: 新增属性OpenDialogDefaultPath，设置一个路径后，日志窗体点击打开文件时，将会默认打开该路径。 * SoftAuthorize: 授权相关类的实例化方法新增参数useHDD用来指示是否使用HDD参数信息，设置false可以避免插拔U盘导致机器码变化。 * Dtu: 重新设计的DTU功能，新增DTU管道，支持TCP,UDP,串口的设备直接变为DTU模式的设备，并直接支持从云上访问，原先的TCP管道删除AlienSession属性。 * DEMO界面的管道新增DTU管道选择，可以直接测试，TCP,串口调试界面支持直接转转远程DTU，文档地址：http://www.hsltechnology.cn/Doc/HslCommunication?chapter=HslCommChapter4-18 * CommunicationServer: 新增属性SessionsMax，用来设置服务器当前支持的最大会话数量，默认为uint.MaxValue * Inovance: 汇川的AM系列的modbus地址，支持 MB100 地址，内部将自动转为字地址，所以必须为偶数开头。 * HttpServer: 新增方法UseHttps( )用来启动https，具体文档参考：http://www.hsltechnology.cn/Doc/HslCommunication?chapter=HslCommChapter6-5 * SiemensS7: S7虚拟服务器模拟真实PLC，读取结果数据大于226字节时返回错误码，优化SiemensS7Net的离散读取方法，现在地址传入任意个数或是任意长度，就能正确的分割读取，然后结果合并。 * DeviceSerialPort: 修复串口设备基类里，当管道设置了非PipeSerialPort串口管道时，调用打开方法Open()时，仍然打开原串口的bug。 * PipeMoxa: 新增一个基于MOXA公司开发的串口库的一个串口管道，在某些情况下比微软官方自带的串口好用，具体参考官网文档。 * MqttHelper: 直接捕获ParseMqttClientApplicationMessage解析的异常，防止一些乱七八糟情况下报文不对导致直接程序奔溃的bug。 * DEMO: 西门子的测试界面的时间，字符串等读写按钮增加光标悬停提示，提示调用的方法名称，方便大家使用。 * Demo: 修复创建示例代码功能中，如果碰到PipeTcpNet管道时，如果属性UseServerActivePush为true时，示例代码没有生成该文本的bug。 * Demo: 修复TCP、UDP调试界面下，当使用了定时发送之后，远程连接强制断开后，一直提示会话为空的bug。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-9-24 | 12.1.1 | * DeviceSerialPort: 修复当设置了非PipeSerialPort串口管道的时候，调用 IsOpen( ) 结果返回不正确的bug。 * DeviceServer: 新增方法 ServerStart( int tcpPort, int udpPort ), 用来同时启动TCP服务和UDP服务，三菱的DEMO测试增加启动功能测试。 * CommunicationServer: 修复服务器的类再只启动UDP服务器的情况下，服务器侧进行关闭操作时，引发后台线程异常的bug。 * MqttServer: 当客户端连接上来的时候，指定了keepAlive后，客户端的ActiveTimeSpan调整为1.5倍的keepAlive值，防止其他客户端极其特殊情况下可能误判下线。 * DLT645: DLT645-2007的串口网口类新增跳闸方法Trip( DateTime validTime )及合闸允许方法SwitchingOn( DateTime validTime ) * LSCnet: 进行了优化，全部使用连续读取的功能实现，添加携带站号的接口IReadWriteDeviceStation， LSCnet实现了该接口。 * SecsValue: 新增加ToSourceCode( )方法，获取获取该对象的源代码表示方式，用于DEMO程序自动生成发数据的源代码。 * Demo: 在PLC测试界面上的管道的功能当选择了DTU管道时，设备连接上来的时候，修复ID匹配不正确的bug。 * Demo: PLC测试界面上的读写按钮操作，都会在下面的示例代码里面显示出实际的代码情况，一些特殊功能测试的界面全部添加点击按钮生成示例代码，包括机床，机器人界面。 * SecsGem: Secs的服务器测试界面和客户端测试界面，功能码列表均支持了编辑和添加删除，然后点击保存连接的时候，自动将全部的功能码列表进行了存储，方便后续测试。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-10-30 | 12.1.2 | * FujiSPBServer: DataFormat属性值修改为和客户端一致，DEMO界面上可以调整该值，修复英文界面下显示文本不正确的问题。 * XinJEServer: 支持输入 D，SD, HD 格式的地址，针对客户端使用XinJE专用协议的时候有效，如果用modbus协议的客户端，则还是输入modbus的地址 * SiemensS7Helper: 修复写入字符串的方法里，判断字节长度超出PLC最大设置定值时，字符串长度获取不正确的bug。 * MqttClient: 新增加一个属性 Tag, 可以用来携带任意的自定义对象，方便客户端做一些个性化的处理。 * Melsec: 三菱MC协议增加错误码17165的错误描述及解决办法，解决办法是将用户认证功能设置为禁用。 * LogNet: 日志类方法HslMessageFormat修改为虚方法，重写这个方法即可以实现存储时自定义的格式，DEMO界面在英文模式下显示英文，显示按钮的代码。 * MQTTServer: 新增方法UnRegisterMqttRpcApi用于卸载当前已经注册的指定接口，方便用于动态控制接口，此处的卸载是批量的。 * HttpServer: 新增方法UnRegisterHttpRpcApi用于卸载当前已经注册的某些接口，重复注册的时候将会覆盖旧的接口内容。 * SimpleHybirdLock: 新增进入锁的时间，当获取锁对象字符串的时候，输出当前锁的状态以及占用锁的时间。 * MqttServer: 修复ClientVerification事件签名里密码单词拼写不正确的bug。 * OmronFinsServer: 欧姆龙的FINS虚拟服务器支持了客户端使用 E9 这种地址的数据读写操作，不过所有的EM地址，E地址都使用同一个数据块。 * SecsServer: Secs服务器日志记录改为直接记录原始二进制的报文，修复服务器创建中文字符串格式的SecsValue时，长度显示不正确的bug。 * DLT645Server: 修复当StationMatch属性为true时，无论站号是否匹配都返回的bug，现在只有站号一直才支持读写数据操作。 * HslDeviceAddressAttribute: 新增字符串编码的属性，当类型是字符串有效，支持字符串数组，修复长度为0时应用于西门子读取字符串时读取字符串不正确的bug。 * Inovance: 汇川的PLC在AM系列下，读取MB1005单数地址的字符串时，也支持进行读取操作。 * IModbus: 所有Modbus设备新增属性BroadcastStation，设置0-255后，使用该站号时将不检测返回的报文，直接返回成功，ModbusServer适配站号0不返回数据。 * OmronFinsUdp: 属性SID信息不再进行自增操作，将由用户手动设置，然后demo界面支持手动设置操作。 * Demo: 批量读取的界面，新增是否字反转操作，方便部分的PLC可以看到一些特殊的数据情况。 * Demo: Secs服务器英文界面优化，西门子Smart200测试界面增加随机写入测试，并直接输出示例代码。 * Demo: MQTT测试界面增加正则过滤功能。设备类的测试界面添加数据模拟的功能，可以写入一些脚本实现的规律数据。 * Demo: HttpClient测试界面的ContentType修改为可输入。Secs界面修复英文环境下，标签名称不正确的问题 * Demo: 设备测试界面的读取功能，支持勾选屏蔽重复的值，只有数据变化的时候，才在读取界面里显示，重复就显示重复次数，菜单调试下新增设置显示毫秒时间的功能。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-11-19 | 12.1.3 | * NetworkWebApiBase: 修复webapi基类实例化后，仍然继续设置IpAddress属性时，实际请求地址不更新的bug。 * LogNet: 日志类新增属性EncoderShouldEmitUTF8Identifier可以控制UTF编码的是否包含BOM头，默认为 true，如果不需要BOM头，可以设置为false * LogNetDateTime: 日期时间类型的日志新增一个实例化方法，可以同时指定文件大小，超过指定大小即自动拆分日志文件。 * CommunicationTcpServer: 修复tcp的服务器在.net8等环境里运行的时候，在关闭服务器的时候，弹出未捕获的异常的bug，现在可以正确的关闭。 * IEC104: 修复IEC104协议在接收某些特殊情况的报文数据导致解析异常，直接退出程序的bug，现在会进行日志的输出。 * LogNetAnalysisControl: 日志分析的控件界面新增一个\0转义的选择框，解决日志消息里带有\0的情况导致显示阶段的问题。 * DLT698Helper: 修复DLT698.45协议里，在某些特殊的设备上，时间数据不存在导致解析异常崩溃的bug，现在直接返回异常的时间内容。 * HttpServer: HttpServer及MqttServer服务器新增卸载单个API接口的方法 UnRegisterHttpRpcApiSingle( string apiTopic ) * BinaryCommunication: 所有设备的基类当因通信失败的时候，返回当前管道的基本信息，例如ip端口，或是串口的端口，方便查错。 * HslStructAttribute: 结构体偏移字节特性新增字符串解析模式属性，可以解析普通字符串，西门子等字符串，结构体解析的属性支持了自定义类，会继续解析子类操作。 * Demo: 显示弹窗的信息框优化，使用了自定义的窗体显示消息本文，根据文本的内容大小，来自动调整当前窗体的大小，并且父窗体居中显示。 * Demo: 模拟数据写入的地方，支持使用随机数，使用r表示，例如 (short)r.Next(100,200) 表示100-200的随机数，最后转short类型。 * LogNet: 日志部分的demo界面部分字符串再英文状态时，显示英文信息，优化相关的界面的测试功能。 * Demo: 欧姆龙的虚拟服务器支持配置DataFormat及是否字符串反转的功能，所有设备的数据导出功能增加随时停止的按钮。 * Demo: PLC测试界面读取字符串的时候，如果碰到\0 字符的时候，原先直接截断显示，对大家造成困扰，现在直接转义为 \0 显示，请注意查看 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2024-12-26 | 12.2.0 | * MelsecMcServer: 对客户端的功能码支持了读取型号，远程run，远程stop，远程reset，错误状态恢复，方便客户端进行测试 * FanucSeries0i: 修复ReadFanucAxisLoad在某些型号读取异常的bug，修复读取系统状态在某些型号读取异常bug，新增读取主轴负载及程序号的方法。 * RemoteConnectInfo: 优化了DTU连接类，虚拟服务器可以直接连接远程DTU服务器，然后返回连接对象，支持进一步控制操作，例如断开，DEMO相关界面优化。 * MelsecFxSerial: 三菱的编程口协议支持写入bool数组到指定地址里去，不过前提条件是bool数组长度是16的倍数，且位起始地址也需要为16的倍数。 * SiemensS7Server: 修复S7虚拟服务器读写bool数组失败，提示不支持的bug，现在可以正确的读取和写入bool数组操作了。 * Keyence: 基恩士上位链路协议客户端支持直接字单位读写位线圈的地址，虚拟服务器同步支持这些命令操作，Demo界面优化。 * ModbusInfo： modbusASCII协议提取真实数据的时候，遇到中间有0d0a结束字符的时候，自动截断操作。 * ToyoPuc: 丰田工机协议写入位数据到字地址的时候，直接提示写入失败，防止错误写入的情况，因为没有这样的功能码操作。 * ToyoPuc: 丰田工机协议新增随机字读取的方法ReadRandom(string[])，虚拟服务器ToyoPucServer也支持了该功能码。 * SecsValue: 格式化输出代码样式的时候，支持format属性，设置true后自动换行，更加符合阅读，DEMO测试界面的客户端优化，支持设置自动回发数据，支持可选显示XML和代码样式等。 * OpenProtocolNet: 新增属性KeepAliveMessageEnable控制是否启用心跳，修复设置管道的方式下导致订阅数据接收不到的bug，DEMO界面优化增加MID0106和MID0107的数据解析。 * BeckhoffAdsNet: 修复再直接设置管道的情况下，可能忘记设置管道的UseServerActivePush属性导致数据在某些特殊情况下不正确的bug，现在都自动设置。 * Demo: 测试工具中批量读取原始字节数据的时候，展示的类型新增boolarray类型选择，自动转为bool数组，表示为0,1 方便大家查看。 * Demo: Demo界面的原始字节读取结果数据框使用等宽字体，然后支持手动设置每行显示的数据个数，不同类型的数据做了自动对齐操作。 * Demo: 程序界面显示报文交互日志窗体时，该窗体新增过滤功能，支持日志等级，正则表达式过滤，支持输出到本地的txt文件。 * Demo: 界面优化，点击写入bool数组的功能按钮时，输入bool数组支持0,1表示，简化输入，例如 [1,0,1,0] * Demo: FormKeyenceSR2000: 界面优化，点击按钮自动显示调用的代码，连接失败提示更加详细的内容。 * Demo: 各种类型读写控件UserControlReadWriteOp基本数据的读写界面新增加搜索文本的功能，方便快速定位数据。 * Demo: 点位变量表控件DataTableControl中，在导入数据的时候，如果碰到不支持的数据类型，错误提示进行了优化，更方便的查找问题。 * Demo: 网口、串口的调试界面收发的数据的显示格式，从只能一样的设置，修改为支持设置不一样的格式，例如发送是ASCII，接收是二进制。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2025-03-20 | 12.3.0 | * SecsValue: 新增输出SML格式字符串，优化SECS协议，SecsHsms新增SendControlMessage发送控制码的接口，相关的demo界面显示优化。 * SecsHsms: secs客户端对象新增属性AutoBackS1F1，因为HSL内部会自动处理S1F1，S1F13，S2F17并返回数据，如果需要完全的手动控制，则需要设置为false。 * 三菱MC: 支持了三菱的MC二进制协议的对标签的写入功能，支持ASCII格式对标签的读写功能，标签地址格式为 "s=" 开头。 * Authorization: 激活判定优化，当处于部分linux系统里且配置了开机启动后，修复系统时间仍为初始时间导致激活判定直接失效的bug * CimonHmiProtocol: 新增西蒙PLC通信类，支持X,Y,M,K,L,T,C,D,F等寄存器，添加配套的虚拟服务器 CimonServer * DLT698: 新增读取多个地址的方法ReadStringArray( string[] address )，新增698配套的虚拟服务器：DLT698Server, 支持DEMO界面的大部分地址。 * NetworkAlienClient: DTU服务器新增属性UseRegistrationPackage 默认为 true， 如果实际情况不需要注册包，可以直接设置为 false. * DLT698OverTcp: DLT648.45 电力规约协议新增属性IsServerActivePush，默认为false, 当遇到某些特殊情况，例如设备会主动上报的时候，可以设置为 true * YamahaRCX: 雅马哈机器人代码优化，命令方法不在需要传入行数，修复了部分机器人连接后因为会发送欢迎消息导致消息不正确的bug，Demo界面优化，增加点动操作。 * MelsecFxSerialServer: 新增三菱编程口的虚拟服务器，方便和客户端配合测试，优化MelsecFxSerialOverTcp报文记录，格式调整为ASCII格式，方便查看报文。 * OpenProtocolServer: OpenProtocol服务器支持了mid0008， mid0009功能码，不过直接返回操作成功，修复DEMO界面关闭服务器无效的bug及不能保存连接的bug。 * OpenProtocolNet: 新增属性ExtraSubscribeMID，类型为int[]，可以用来设置额外的用于订阅并且要求返回的mid，DEMO界面上支持配置该属性。 * KeyenceKvOld: 新增基恩士的通信类，针对老型号：Kv-10xx, 16xx, 24xx, 40xx, kv-80, kv-300，其中 xx 表示 AR/AT/DR/DT，仅支持单点读写，读写short，ushort类型 * PipeTcpNet: TCP管道当开启UseServerActivePush=true异步模式下，调整顺序，先处理QA消息，再重新接收后面的数据内容。 * IEC104: 优化IEC104类的发送id，接收id的计算和判断，修复因网络异常再次发送时发送id没有重置的bug，DEMO界面优化，显示数据公共地址，显示32位比特串。 * IEC104Server: 新增IEC104的虚拟服务器，添加DEMO界面，内置了一些连续的地址数据，可以在demo界面上查看，并编辑修改这些值，支持总召唤功能，修改数据时支持突发上传单个数据信息。 * FanucSeries0i: 新增进给倍率和主轴倍率读取的方法接口，分别为：ReadFeedRate 和 ReadSpindleRate * Modbus: Modbus协议支持读写文件，0x14功能码及0x15功能码，地址格式为 "file=0;100", 其中file是文件号，虚拟服务器也相应的配套支持。 * ModbusRtu: 修复ModbusRtu类使用0x17功能码及其他非常见功能码时，判断数据接收完成异常的bug。包括rtu over tcp类。 * MqttServer: 新增表达式属性DownloadFileRedirect，用来在下载文件的时候，可以重新指定实际下载的文件路径，并可以指定下载是否删除。 * DebugRemoteServer: 新增远程调试服务器，让设备通信支持远程透传调用，支持多个设备调用，DEMO界面支持直接测试，网址：http://www.hsltechnology.cn/Doc/HslCommunication?chapter=HslCommChapter7-9 * DEMO: 测试界面右上角的保存的当前连接信息（数据电表，数据模拟都支持存储）支持重命名，支持批量删除操作，支持修改密码。 * DEMO: MQTT服务器界面新增支持显示收到主题消息的统计，接收次数及大小，分段时间的接收次数，可以按每秒，每分，每小时统计次数。 * DEMO: Demo程序界面里每个设备读写测试界面里，在选项卡[点位变量]表，支持设置每个请求和上个请求的时间间歇。 * Demo: Demo程序界面支持了激活hslcontrols组件的功能，菜单界面active可以输入两个激活操作。 * Demo: 调试工具界面上的调试菜单下面新增窗体置顶的功能，方便大家根据需要让DEMO窗体置顶，方便实际调试的时候操作。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2025-05-09 | 12.3.1 | * SecsHsmsServer: 新增方法PublishSecsMessage( byte stream, byte function, SecsValue data, uint messageId, bool wBit )，可以发布的时候手动控制系统数据。 * HslExtension: 优化地址判断方法StartsWithAndNumber，支持了对负数的地址判断，因为台达等PLC存在负数地址的情况，例如 D-3530 * CommunicationTcpServer: 修复TCP服务器在linux系统运行时，某些特殊情况调用 ServerClose 后会异常奔溃的bug。 * WebSocketServer: 服务器端新增方法 public void PublishAllClientPayload( byte[] payload ) 使用opcode = 2来发送二进制的数据消息，DEMO界面优化显示。 * PipeSslNet: 新增属性 SslProtocols 可以自定义控制SSL协议的版本，CommunicationTcpServer新增方法 SetSslPipeAction 可用于设置SSL的自定义属性。 * PipeTcpNet: 新增属性 CloseOnRecvTimeOutTick，当使用同步方法读写PLC数据时，想要实现超时几次内不关闭当前的连接对象，就设置该属性，注意，只对同步方法有效。 * TcpForward: TCP转发类接收到双方数据后，顺序调整为先转发再重新接收数据，修复在某些连续发送的场景下导致数据接收顺序错乱的bug。 * AllenBradleyNet: 当使用了i=的地址格式写入单个的bool的值的时候，修复地址为单个dint类型时，提示数据不存在的bug。 * MelsecFxSerial: 优化三菱编程口协议对报文完整性的判断，就算读取的报文前面出现 3F 3F 3F字节，也可以正确的识别和解析数据。 * HttpServer: 内部日志记录针对所有的请求操作及url信息，耗时进行记录操作，新增属性 LogHttpBody 默认false，设置true可以记录body数据，跨域属性默认为 true。 * HttpServer: 新增属性LogHttpHeader，用于表示日志是否记录请求头的信息，DEMO界面优化，支持日志设置请求头或是Body数据。 * ModbusRtu: 修复当收到CRC16不正确的报文时，却表现为接收超时的bug，现在会直接返回并提示CRC不正确，当站号不一致时，还是会持续接收直到有指定站号数据或是引发超时。 * Modbus: IModbus接口新增属性WordReadBatchLength，表示批量读取时自动切割的长度信息，默认根据120长度分批读取，现在开放修改。 * MqttServer: 修改继承体系，从V12版本的CommunicationTcpServer继承实现，支持了SSL/TLS加密通信功能，优化了Demo界面的部分显示信息，测试了连续几万次的发布接收功能。 * MelsecMcRNet: 修复三菱R系列的协议情况下，命令码中部分内容不正确的bug，此前一直读写不了R系列的PLC，需要使用普通MC协议。 * OmronCipNet: 修复欧姆CIP协议写入单个bool数据到PLC失败的bug，修复bool数组读取解析异常，支持从数组中间位置开始读取。 * SecsHsms: 方法SendByCommand返回携带消息ID信息，方便二次处理，修复SecsValue在某些值信息为空的情况下转成SML字符串时异常的bug。 * SecsHsmsServer: 服务器侧新增一个同步方法ReadSecsMessage读取客户端的返回数据信息，该方法内部有锁，线程安全的。 * OmronCMode: 欧姆龙的HostLinkCMode支持了读写位的操作，实际是通过读写字间接实现的，修复了HostLinkCModeServer当客户端读取EM地址失败的bug。 * Demo: demo界面的模拟数据写入测试界面新增单次测试，如果碰到地址异常，自动选中该单元格，修复设备保存时模拟数据界面不保存的bug。 * Demo: DEMO界面上选择串口管道的时候，额外属性里新增ReceiveTimeOut属性，可以设置超时时间了，参数支持保存加载，支持示例代码生成。 * Demo: 修复CJT188及串口透传类，在DEMO中生成示例代码的时候，实际复制到项目出错的bug，现在可以正确的生成代码。 * Demo: 修复DEMO界面的点位变量表，请求间歇没设置直接启动报错的bug，该信息也存储到XML里，当设备关闭时，点位变量表也根据停止读取刷新。 * Demo: fanuc机床的程序测试界面增加特殊情况下载失败的信息提示，openprotocol界面增加实际代码显示，方便复制。 * Demo: 修复demo碰到激活码不正确的情况下，启动弹窗的bug，现在不再弹窗操作。http服务器和websocket服务器的界面控件支持滚动条。 * Demo: 测试工具界面的写入地址的数据，支持输入连续数组，例如数组1,2,3,...,100，直接输入 [1:100] 即可，然后点击写入的类型。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2025-06-06 | 12.3.2 | * LogNetBase: 优化了日志组件在设置属性ConsoleOutput输出到控制台时的性能，减少颜色重复设置，连续日志批量输出，输出性能增加5倍以上。 * EcFanMachine: 新增一个特殊的协议EcFanMachine: FFU 用EC 风机风机通信协议，支持风机的转速控制，风机的地址设置，风机的状态读取等功能 * EmergencyState: 风机对象优化，读取数据时支持返回应急状态，然后DEMO优化，MQTT服务器界面主题优化，AB-CIP读取显示类型备注。 * ILogNet: 日志接口接口ILogNet新增属性ConsoleColorEnable，默认为true，表示控制台输出时不同等级日志颜色不一样，设置false关闭颜色。 * LogNetManagment: 日志按照时间及大小存储时的文件名的开头字符串信息可以修改，设置属性 LogFileHeadString 即可。 * MelsecFxSerialServer: 新增对05报文的的回复，新增对客户端修改波特率的命令的支持，根据客户端的命令来自动设置波特率。 * InovanceSeries: 汇川的系列新增子系列 Easy，但是实际的modbus地址映射和H5U是一致的，DEMO界面可以选择新添加的系列内容。 * LogNetBase: 日志基础类，修复在存储日志信息到文件的时候，如果文件的目录被删除导致日志存储不会恢复的bug，现在会自动重新创建文件夹。 * InovanceHelper: 汇川AM系列支持bool读取MD地址，例如 MD100.0 等同于 MW200.0 地址 * S7Message: 西门子S7协议在接收PLC返回的报文的时候，新增对报文的消息id检查操作，检查和发送是否一致，虚拟服务器配套修改。 * DEMO: demo上写入数据的测试方法支持了地址重复的情况，例如想写100个short全部都是0的数组，地址里输入 [0\*100] 即可。 * DEMO: demo读写byte,short,ushort,int,uint,long,ulong类型时，支持使用16进制的数据格式读写，使用16进制时，暂时只支持正整数。 * Demo: Demo界面上设备的读写测试界面新增一个小框框，点击可以隐藏读写按钮测试界面，让特殊功能测试界面最大化。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2025-07-30 | 12.3.3 | * DLT645的1997及2007协议新增属性CheckDataId，默认为true，表示检查读取返回报文中的数据ID信息，修复在某些特殊情况下报文接收不正确的bug。 * 新增DLT1997的虚拟服务器类DLT645With1997Server，支持常见地址的写入读取，修复DLT1997写入数据时报文组包不正常的bug，现在可以正确的写入。 * OmronCipServer: 修复欧姆龙CIP虚拟服务器在相应基于连接的CIP协议报文时，针对数组索引访问时，返回数据的bug，之前会全部返回，现在只返回读取的数据。 * CIPServer，包含ab的CIP服务器及欧姆龙的CIP服务器，在返回客户端报文时，复制客户端的context数据头内容，然后直接返回给客户端。 * Omron: 修复OmronHostLink协议，FinsTcpNet，FinsUdp类在读写DR, IR地址的bool数据时，读写数据表现不正确的bug。 * Modbus: 所有modbus测试界面功能里，新增站号搜索功能，自动的发报文匹配对应站号，包括tcp.rtu, ascii,及其透传模式都支持测试，直接搜索站号 0 - 255方便快速查找。 * PanasonicHelper: 松下的PLC地址 DT 区的地址同时支持 D100， DT100 两种表示方式，表示同一个地址。 * AllenBradleyNet: 新增属性ContextCheck，默认false，当设置true的时候，将自增Senderontext属性当为消息id，接收消息匹配处理。 * 优化网络socket在异步接收数据的时候，频繁创建方法委托实例的小问题，在很多PLC同时高频通信的时候，会存在句柄持续增大的问题。 * ToyoPuc: 支持了bool数组的写入，支持了prg=1;M102 开头的地址的bool写入功能，通过读字修改位，写字来实现。 * ModbusTcpNet: 新增属性 StationCheckMacth 表示是否对站号进行检查操作。消息接收优化，支持某些特殊的不规范的消息结构报文。 * AllenBradleyNet: AB的PLC新增属性ReadArrayUseSegment，表示当读取数组的时候，是否使用52片段功能码读取，默认为true，表示使用52功能码。 * KeyenceNano: 针对基恩士的上位链路协议的DM,CM,EM,FM,ZF,VM这种字地址支持了位读写操作，格式为 DM100.0 这种带小数点，然后bool读写操作。 * DTU: DEMO中DTU的默认端口调整为10010，支持了长度信息使用ASCII字符表示，注册包支持输出ASCII字符，方便部分只支持ASCII配置的DTU设备。 * Demo: 虚拟PLC的数据存储加载代码优化，然后DEMO测试上支持定时存储虚拟PLC的数据到文件。 * Demo: DEMO界面的DEBUG界面优化了显示数据，显示了发送数据次数，接收数据次数信息，压力测试控件界面优化。 * Demo: 新增加一个PING ip地址的测试界面，可以设置线程数量来快速确定可以访问的，在线的ip地址。 * Demo：曲线界面优化，数据表界面支持设置曲线颜色，样式，粗细。纯曲线界面支持对多个曲线设置颜色，样式，粗细。 * Demo: demo菜单新增一个选项，是否固定测试界面大小，当程序运行在分辨率特别低的显示器的时候，可以设置该属性来显示滚动条。 * Demo: 所有的通信测试界面优化，当测试通信的窗体在没有正常关闭连接的情况下，直接点击叉号关闭PLC测试界面，原来的连接会被自动的关闭，服务器也会自动的关闭端口。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2025-09-10 | 12.5.0 | * SiemensS7Server: 修复当客户端读取位的方式，读取不存在的DB块的时候，没有返回的异常，现在返回DB块不存在的消息。 * Modbus: Modbus相关的协议类的属性StationCheckMatch，修复单次拼写错误，之前的是StationCheckMacth，如有使用该属性，升级时需要注意。 * MelsecMcNet: 三菱的属性EnableWriteBitToWordRegister默认值调整为 true, DEMO界面的选项也自动为 true，可以写位到字寄存器地址。 * WebSocketServer: 新增方法public void SendClientPayload( WebSocketSession session, byte[] payload ) 用来发送 byte[] 数据。 * ModbusRtu: 优化RTU消息完整性校验的代码，现在支持了报文前面即使有三个干扰字节，也可以正确的识别报文并正确的提炼出数据。 * IModbus: 新增属性DisableFunctionCode06，表示禁用功能码06，禁用之后使用10功能码来写入short, ushort 数据。 * IEC104: 修复写入归一化设定值异常的bug，对应的IEC104虚拟服务器支持了单点遥控，双点遥控，归一化，标度化，短浮点，32位比特串的写入操作，并增加最大地址范围设置。 * AllenBradleyNet: 在握手报文阶段，当碰到设备错误码1的时候，不再直接关闭返回，而是视为成功，继续操作，部分的PLC是可以继续通信的。 * HttpServer: 新增一个启动的重载方法 public void Start( string prefix ), 可以输入 http://127.0.0.1:8000/ 来启动服务器。 * NetPlainSocket: 优化普通客户端类，针对高频接收可能发生的数据错乱做了优化，新增byte[]的收发方法，官网完善相关的文档信息。 * IDlt645: DLT645接口新增修改波特率的方法ChangeBaudRate，相关DEMO界面增加该测试功能，虚拟服务器支持这个功能码测试。 * YRCHighEthernet: 新增加一个读取位置型变量的方法ReadPositionVariable，相关的测试界面添加测试功能。 * DeviceServer: 包括基类CommunicationTcpServer新增属性LocalAddress，用于服务器端绑定本地的ip地址来启动服务，DEMO里所有的服务器界面调整优化。 * AllenBradleyNet: 修复写入空字符串时实际成功，提醒失败的bug，修复CIPServer服务器返回数据报文命令头不匹配的bug。 * LogNetDateTime: 时间日期存储的日志，路径地址支持输入格式化信息，达到动态的效果，具体参考官网。 * Demo: Demo界面的点位变量监视表里，修复右键删除指定行不动作的问题，现在可以批量选择删除。JSON测试界面新增替换测试功能。 * Demo: Demo界面新增企业用户相关的支持证书发布的功能，允许企业用户自己注册及发布证书，该证书可用于激活C#版及Java版本的hslcommunication。 * Demo: 当测试界面写入数组数据到地址里的时候，修复生成的示例代码不正确的bug，涉及所有类型的数组写入操作。 * Demo: 修复CIP协议特殊功能测试的界面，读取显示数据类型的时候，点击写入之后，提示类型错误的小bug。 * Demo: 修复测试界面 DTU报文生成时，选中ASCII格式时，并且DTU的id小于11位数字的时候，创建的ASCII的注册包的显示异常问题。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2025-09-18 | 12.5.1 | * KeyenceKvOld: 修复当设备网络问题掉线的时候，读取结果仍然还是成功的bug，原因来自失败的时候返回了错误的对象。 * SiemensS7Server: 新增DB块数据区时支持指定长度，优化客户端和服务器读写超过地址范围时提示错误，DEMO界面支持存储加载增加的DB块信息，显示优化。 * AllenBradleyNet: 新增方法OperateResult<AbTagValue[]> ReadTags( string[] address, ushort[] length )每个标签返回一个AbTagValue对象，包含类型及二进制数据。 * AllenBradleyServer: AB的虚拟服务器新增加对客户端同时读取多个数据的支持，一条报文实现读取返回结果操作，客户端测试通过。 * Server: SiemensS7Server, OmronFinsServer, GeSRTPServer, AllenBradleyServer, AllenBradleyPcccServer这种有握手包的虚拟服务器新增对串口通信的支持。 * ModbusRtu: 优化ModbusRtu数据包的报文接收完整性的判断，修复极少部分情况下接收部分报文却判断为报文完整的bug。 * Demo: PLC测试界面的数据模拟子界面，支持右键菜单操作，删除行，导出，导入数据，然后优化脚本的提示及复制功能，部分中英文界面优化。 * Demo: 修复PLC测试界面不关闭连接，直接打叉关闭的情况下，远程调试功能处于打开状态时，端口号不释放的bug。 * Demo: Demo界面的HttpClient测试界面，当使用了GET测试功能，参数写到body里时，测试请求的时候会自动解析参数赋值到url中，示例url显示实际带参数的url(调整为可复制)。 * Demo: 修复Demo程序菜单栏的激活界面刚打开就直接关闭，几秒后直接提示闪退的bug。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2025-10-20 | 12.5.2 | * OmronHostLinkCModeServer: 修复CMode协议的虚拟服务器，在客户端执行远程写入操作的时候，导致数据错位的bug。 * OmronCipServer: 修复欧姆龙PLC的CIP虚拟服务器在客户端读取数组索引时，返回数据不正确的bug，优化字符串标签读写时候的数据包大小。 * AllenBradleyNet: 在写入地址 i=A[10] 这类地址时，写入位到PLC里的整数方法现在支持SINT,INT,DINT,LINT类型，同时支持 i=A.10 这种地址格式。 * Modbus: Modbus协议里，设置DataFormat为BADC，DCBA时，读取"100.0"的bool数组时，调整为小端模式，跟随DataFormat变化。 * AllenBradleyHelper: ab-cip的协议支持二维数组的地址格式，同时支持了A[1,2]格式和 A[1][2] 格式，修复原本二维数组地址输入后面格式解析异常的bug。 * PipeTcpNet: 优化TCP管道的异常消息号计数的功能，优化连接及重新连接的部分代码。 * HttpServer: 修复HttpServer服务器启动后，手动关闭服务的时候，不停的反复触发异常日志记录的bug。 * Demo: 倍福的ADS协议新增一个地址示例说明，如果是TWINCAT2的情况，全局变量地址前加一个点，例如 s=.A * Demo: HttpClient的webapi测试界面，当get请求时，参数在body里输入时，修复数组的参数解析不正确的bug。 * Demo: 在所有设备测试读取界面，显示读取到的单个变量值时，在十进制和十六进制之外，支持显示二进制的数据内容。 * Demo: 修复虚拟服务器在生成示例代码的时候，如果绑定了本地的ip地址，结果生成代码异常的bug。 * Demo: 西门子的虚拟服务器手动添加DB块的时候，支持直接生成示例代码，方便显示复制。 * Demo: 修复DEMO已经保存的测试界面，重新保存不一样的名字窗口时，两个菜单指向一致的小问题。 * Demo: 已保存的窗体双击打开时，修复之前的bug，现在改为绑定配置的GUID码，无论树形控件怎么刷新都支持重新打开旧的。 * Demo: 保存列表双击当前的树形图窗体时，优先显示已经打开的窗体界面，如果需要新的窗体，可以鼠标右键菜单操作。 * Demo: 测试界面的"点位变量"和"数据模拟"支持导入导出CSV格式字符串到剪贴板，可以直接复制到Excel，数据模拟界面可以自己设置常用表达式。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |
| 2025-11-19 | 12.5.3 | * WebSocketHelper: websocket消息支持了数据压缩的情况，当收到压缩的payload数据的时候，会自动根据压缩标记位来解压缩数据信息。 * WebSocketClient: 新增属性SupportDeflate, AutoDecompress用来控制是否要求服务器不支持压缩，以及数据压缩的情况下，对压缩数据进行手动解析的功能。 * BeckhoffAdsNet: 新增加一个方法 ClearTagsCache 用来清除当前PLC里所有的缓存的标签内存偏移地址的字典。 * Ftpclient: 添加了Ftp的客户端，支持文件上传，文件下载(支持进度报告)，文件删除，目录查列表，删除文件夹，创建文件夹，重命名文件夹，重命名文件功能。 * NetPlainSocket: 普通TCP类优化接收数据逻辑，并添加CreateNewMessage委托用来绑定消息类，官网对应的文档内容优化。 * AllenBradleyServer: 开放服务器的数据词典属性TagValueDict，DEMO界面的CIP服务器增加查看所有数据词典的功能，保存加载时自动保存数据词典的数据。 * AllenBradleyPcccNet: 实现了读取bool数组的功能方法，实际使用读字，解析出位的方式，同步方法，异步方法都支持该操作。 * Modbus: 修复DataFormat为BADC、DCBA时，写入bool数组到字寄存器的时候，高地位不正确的bug，影响范围包括rtu, rtuovertcp, ascii, modbustcp。 * DLT645With1997: 及DLT645With1997OverTcp支持int32类型的读写方法，实际是读写double方法强制转int实现。 * Demo: fanuc机床的测试界面启动按钮增加是否启动的确认，防止误操作，优化部分信息提示中英文结合。 * Demo: 串口管道设置属性新增读取前是否清除缓存IsClearCacheBeforeRead，默认false，modbusrtu测试界面移除这个属性的配置改到管道里配置。 * Demo: Demo的设备测试界面，当添加或是修改了点位变量表或是数据模拟表，没有保存直接关闭窗体的时候，增加弹窗提示操作。 * Demo: 设备测试界面的数据模拟功能块，当配置写入数据是字符串的时候，字符串写入支持设置编码，脚本例子新增加一个字符串的例子。 * Demo: 读取数据显示的界面，显示16进制及二进制的方式也对数组读取及写入也有效，主要针对short, ushort, int. uint, long, ulong。 * Demo: 主界面的设备列表新增搜索功能，可以快速搜索到需要的协议，然后在程序关闭的时候支持记忆功能，记录打开的是协议列表还是保存列表。 * Demo: 修复DEMO界面使用MQTT管道的时候，反复关闭重新打开，会导致重新创建MQTTClient对象的bug，这个属于DEMO的bug，库本身没问题。 * 新官网：http://www.hsltechnology.cn/，还有全新的使用文档的地址：http://www.hsltechnology.cn/Doc/HslCommunication * 本软件已经申请软件著作权，软著登记号：2020SR0340826，任何盗用软件，破解软件，未经正式合同授权而商业使用均视为侵权。 |

![](../icons/SectionExpanded.png)参见

#### 引用

[SoftBasicGetUniqueStringByGuidAndRandom](2758b3f4-e30c-08e5-fc7f-fd56015cc71c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Authorization 类

[原文連結](http://api.hslcommunication.cn/html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[Authorization 类](../html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm "Authorization 类")

[Authorization 构造函数](../html/e68d527d-1fdd-fffd-63a5-3fadaa948b88.htm "Authorization 构造函数 ")

[Authorization 方法](../html/d40a2af1-6e2f-d0cb-cf48-1096e1009001.htm "Authorization 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| Authorization 类 |

系统的基本授权类

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunicationAuthorization

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class Authorization
```

```
Public Class Authorization
```

```
public ref class Authorization
```

```
type Authorization =  class end
```

Authorization 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Authorization](e68d527d-1fdd-fffd-63a5-3fadaa948b88.htm) | 初始化 Authorization 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [SetAuthorizationCode(String)](a8ac31f8-be69-3288-f661-ced47c9710a2.htm) | 设置本组件系统的授权信息，如果激活失败，只能使用24小时，24小时后所有的网络通信不会成功  Set the authorization information of this component system. If the activation fails, it can only be used for 8 hours. All network communication will not succeed after 8 hours |
| 公共方法静态成员 | [SetAuthorizationCode(String, Int32, String)](ffa5f3ba-639e-cf82-c4e3-36916072b28c.htm) | 通过指定的IP地址，端口号信息的接口服务器获取密钥，然后进行激活，适用局域网内一个主服务对其他电脑的激活操作。如果网络失败或是密钥不正确，则激活失败。  The interface server obtains the key through the specified IP address and port number information, and then activates it, which applies to the activation operation of one main service in the LAN to other computers. If the network fails or the key is incorrect, activation fails. |
| 公共方法静态成员 | [SetDllContact](bf5220e5-5a77-64a2-1533-e988b4fb0bd9.htm) | 设置本组件激活失败的情况时候的联系方式，可以设置成自定义的字符串，当且仅当企业用户激活成功后调用本方法有效，如果不设置默认为：联系QQ200962190，微信：13516702732，Email:hsl200909@163.com  Set the contact information when the activation of this component fails, can be set to a custom string, if and only if the enterprise user activation is successful to call this method is valid, if not set the default is: contact QQ200962190, WeChat: 13516702732,Email:hsl200909@163.com |
| 公共方法静态成员 | [SetHslCertificate](21895f0b-6ab8-a906-f5e7-28c784795ca2.htm) | 使用证书激活hslcommunication组件，证书请联系胡工科技获取，如果获取的是文件的话，那就是 SetHslCertificate( File.ReadAllBytes( "hsl.cert" ) ) 获取原始字节信息  Use the certificate to activate the hslcommunication component, please contact Hugong Technology to obtain the certificate, if the obtained is a file, that is, SetHslCertificate( File.ReadAllBytes( "hsl.cert") ) to obtain the raw byte information |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Authorization 构造函数 

[原文連結](http://api.hslcommunication.cn/html/e68d527d-1fdd-fffd-63a5-3fadaa948b88.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[Authorization 类](../html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm "Authorization 类")

[Authorization 构造函数](../html/e68d527d-1fdd-fffd-63a5-3fadaa948b88.htm "Authorization 构造函数 ")

[Authorization 方法](../html/d40a2af1-6e2f-d0cb-cf48-1096e1009001.htm "Authorization 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| Authorization 构造函数 |

初始化 [Authorization](e89d602f-541c-8b98-af2b-02ee7b19618b.htm) 类的一个新实例

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Authorization()
```

```
Public Sub New
```

```
public:
Authorization()
```

```
new : unit -> Authorization
```

![](../icons/SectionExpanded.png)参见

#### 引用

[Authorization 类](e89d602f-541c-8b98-af2b-02ee7b19618b.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Authorization 方法

[原文連結](http://api.hslcommunication.cn/html/d40a2af1-6e2f-d0cb-cf48-1096e1009001.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[Authorization 类](../html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm "Authorization 类")

[Authorization 方法](../html/d40a2af1-6e2f-d0cb-cf48-1096e1009001.htm "Authorization 方法")

[SetAuthorizationCode 方法](../html/0908319d-e4ef-b1c3-c302-8488e0a9fc93.htm "SetAuthorizationCode 方法 ")

[SetDllContact 方法](../html/bf5220e5-5a77-64a2-1533-e988b4fb0bd9.htm "SetDllContact 方法 ")

[SetHslCertificate 方法](../html/21895f0b-6ab8-a906-f5e7-28c784795ca2.htm "SetHslCertificate 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| Authorization 方法 |

[Authorization](e89d602f-541c-8b98-af2b-02ee7b19618b.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [SetAuthorizationCode(String)](a8ac31f8-be69-3288-f661-ced47c9710a2.htm) | 设置本组件系统的授权信息，如果激活失败，只能使用24小时，24小时后所有的网络通信不会成功  Set the authorization information of this component system. If the activation fails, it can only be used for 8 hours. All network communication will not succeed after 8 hours |
| 公共方法静态成员 | [SetAuthorizationCode(String, Int32, String)](ffa5f3ba-639e-cf82-c4e3-36916072b28c.htm) | 通过指定的IP地址，端口号信息的接口服务器获取密钥，然后进行激活，适用局域网内一个主服务对其他电脑的激活操作。如果网络失败或是密钥不正确，则激活失败。  The interface server obtains the key through the specified IP address and port number information, and then activates it, which applies to the activation operation of one main service in the LAN to other computers. If the network fails or the key is incorrect, activation fails. |
| 公共方法静态成员 | [SetDllContact](bf5220e5-5a77-64a2-1533-e988b4fb0bd9.htm) | 设置本组件激活失败的情况时候的联系方式，可以设置成自定义的字符串，当且仅当企业用户激活成功后调用本方法有效，如果不设置默认为：联系QQ200962190，微信：13516702732，Email:hsl200909@163.com  Set the contact information when the activation of this component fails, can be set to a custom string, if and only if the enterprise user activation is successful to call this method is valid, if not set the default is: contact QQ200962190, WeChat: 13516702732,Email:hsl200909@163.com |
| 公共方法静态成员 | [SetHslCertificate](21895f0b-6ab8-a906-f5e7-28c784795ca2.htm) | 使用证书激活hslcommunication组件，证书请联系胡工科技获取，如果获取的是文件的话，那就是 SetHslCertificate( File.ReadAllBytes( "hsl.cert" ) ) 获取原始字节信息  Use the certificate to activate the hslcommunication component, please contact Hugong Technology to obtain the certificate, if the obtained is a file, that is, SetHslCertificate( File.ReadAllBytes( "hsl.cert") ) to obtain the raw byte information |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[Authorization 类](e89d602f-541c-8b98-af2b-02ee7b19618b.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SetAuthorizationCode 方法 

[原文連結](http://api.hslcommunication.cn/html/0908319d-e4ef-b1c3-c302-8488e0a9fc93.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[Authorization 类](../html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm "Authorization 类")

[Authorization 方法](../html/d40a2af1-6e2f-d0cb-cf48-1096e1009001.htm "Authorization 方法")

[SetAuthorizationCode 方法](../html/0908319d-e4ef-b1c3-c302-8488e0a9fc93.htm "SetAuthorizationCode 方法 ")

[SetAuthorizationCode 方法 (String)](../html/a8ac31f8-be69-3288-f661-ced47c9710a2.htm "SetAuthorizationCode 方法 (String)")

[SetAuthorizationCode 方法 (String, Int32, String)](../html/ffa5f3ba-639e-cf82-c4e3-36916072b28c.htm "SetAuthorizationCode 方法 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AuthorizationSetAuthorizationCode 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [SetAuthorizationCode(String)](a8ac31f8-be69-3288-f661-ced47c9710a2.htm) | 设置本组件系统的授权信息，如果激活失败，只能使用24小时，24小时后所有的网络通信不会成功  Set the authorization information of this component system. If the activation fails, it can only be used for 8 hours. All network communication will not succeed after 8 hours |
| 公共方法静态成员 | [SetAuthorizationCode(String, Int32, String)](ffa5f3ba-639e-cf82-c4e3-36916072b28c.htm) | 通过指定的IP地址，端口号信息的接口服务器获取密钥，然后进行激活，适用局域网内一个主服务对其他电脑的激活操作。如果网络失败或是密钥不正确，则激活失败。  The interface server obtains the key through the specified IP address and port number information, and then activates it, which applies to the activation operation of one main service in the LAN to other computers. If the network fails or the key is incorrect, activation fails. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[Authorization 类](e89d602f-541c-8b98-af2b-02ee7b19618b.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SetAuthorizationCode 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/a8ac31f8-be69-3288-f661-ced47c9710a2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[Authorization 类](../html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm "Authorization 类")

[Authorization 方法](../html/d40a2af1-6e2f-d0cb-cf48-1096e1009001.htm "Authorization 方法")

[SetAuthorizationCode 方法](../html/0908319d-e4ef-b1c3-c302-8488e0a9fc93.htm "SetAuthorizationCode 方法 ")

[SetAuthorizationCode 方法 (String)](../html/a8ac31f8-be69-3288-f661-ced47c9710a2.htm "SetAuthorizationCode 方法 (String)")

[SetAuthorizationCode 方法 (String, Int32, String)](../html/ffa5f3ba-639e-cf82-c4e3-36916072b28c.htm "SetAuthorizationCode 方法 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AuthorizationSetAuthorizationCode 方法 (String) |

设置本组件系统的授权信息，如果激活失败，只能使用24小时，24小时后所有的网络通信不会成功  
Set the authorization information of this component system. If the activation fails, it can only be used for 8 hours. All network communication will not succeed after 8 hours

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool SetAuthorizationCode(
	string code
)
```

```
Public Shared Function SetAuthorizationCode ( 
	code As String
) As Boolean
```

```
public:
static bool SetAuthorizationCode(
	String^ code
)
```

```
static member SetAuthorizationCode : 
        code : string -> bool 
```

#### 参数

code
:   类型：SystemString  
    授权码，如果授权码为空，则使用在线激活的方式

#### 返回值

类型：Boolean  
是否激活成功

![](../icons/SectionExpanded.png)参见

#### 引用

[Authorization 类](e89d602f-541c-8b98-af2b-02ee7b19618b.htm)

[SetAuthorizationCode 重载](0908319d-e4ef-b1c3-c302-8488e0a9fc93.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SetAuthorizationCode 方法 (String, Int32, String)

[原文連結](http://api.hslcommunication.cn/html/ffa5f3ba-639e-cf82-c4e3-36916072b28c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[Authorization 类](../html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm "Authorization 类")

[Authorization 方法](../html/d40a2af1-6e2f-d0cb-cf48-1096e1009001.htm "Authorization 方法")

[SetAuthorizationCode 方法](../html/0908319d-e4ef-b1c3-c302-8488e0a9fc93.htm "SetAuthorizationCode 方法 ")

[SetAuthorizationCode 方法 (String)](../html/a8ac31f8-be69-3288-f661-ced47c9710a2.htm "SetAuthorizationCode 方法 (String)")

[SetAuthorizationCode 方法 (String, Int32, String)](../html/ffa5f3ba-639e-cf82-c4e3-36916072b28c.htm "SetAuthorizationCode 方法 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AuthorizationSetAuthorizationCode 方法 (String, Int32, String) |

通过指定的IP地址，端口号信息的接口服务器获取密钥，然后进行激活，适用局域网内一个主服务对其他电脑的激活操作。如果网络失败或是密钥不正确，则激活失败。  
The interface server obtains the key through the specified IP address and port number information, and then activates it,
which applies to the activation operation of one main service in the LAN to other computers. If the network fails or the key is incorrect, activation fails.

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool SetAuthorizationCode(
	string ip,
	int port,
	string token = null
)
```

```
Public Shared Function SetAuthorizationCode ( 
	ip As String,
	port As Integer,
	Optional token As String = Nothing
) As Boolean
```

```
public:
static bool SetAuthorizationCode(
	String^ ip, 
	int port, 
	String^ token = nullptr
)
```

```
static member SetAuthorizationCode : 
        ip : string * 
        port : int * 
        ?token : string 
(* Defaults:
        let _token = defaultArg token null
*)
-> bool 
```

#### 参数

ip
:   类型：SystemString  
    远程的IP地址信息

port
:   类型：SystemInt32  
    远程的端口号信息

token (Optional)
:   类型：SystemString  
    令牌信息，不是必须的，具体取决于服务器作用

#### 返回值

类型：Boolean  
是否激活成功

![](../icons/SectionExpanded.png)参见

#### 引用

[Authorization 类](e89d602f-541c-8b98-af2b-02ee7b19618b.htm)

[SetAuthorizationCode 重载](0908319d-e4ef-b1c3-c302-8488e0a9fc93.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SetDllContact 方法 

[原文連結](http://api.hslcommunication.cn/html/bf5220e5-5a77-64a2-1533-e988b4fb0bd9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[Authorization 类](../html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm "Authorization 类")

[Authorization 方法](../html/d40a2af1-6e2f-d0cb-cf48-1096e1009001.htm "Authorization 方法")

[SetAuthorizationCode 方法](../html/0908319d-e4ef-b1c3-c302-8488e0a9fc93.htm "SetAuthorizationCode 方法 ")

[SetDllContact 方法](../html/bf5220e5-5a77-64a2-1533-e988b4fb0bd9.htm "SetDllContact 方法 ")

[SetHslCertificate 方法](../html/21895f0b-6ab8-a906-f5e7-28c784795ca2.htm "SetHslCertificate 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AuthorizationSetDllContact 方法 |

设置本组件激活失败的情况时候的联系方式，可以设置成自定义的字符串，当且仅当企业用户激活成功后调用本方法有效，如果不设置默认为：联系QQ200962190，微信：13516702732，Email:hsl200909@163.com  
Set the contact information when the activation of this component fails, can be set to a custom string,
if and only if the enterprise user activation is successful to call this method is valid, if not set the default is: contact QQ200962190, WeChat: 13516702732,Email:hsl200909@163.com

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static void SetDllContact(
	string contact
)
```

```
Public Shared Sub SetDllContact ( 
	contact As String
)
```

```
public:
static void SetDllContact(
	String^ contact
)
```

```
static member SetDllContact : 
        contact : string -> unit 
```

#### 参数

contact
:   类型：SystemString  
    新的联系方式

![](../icons/SectionExpanded.png)参见

#### 引用

[Authorization 类](e89d602f-541c-8b98-af2b-02ee7b19618b.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SetHslCertificate 方法 

[原文連結](http://api.hslcommunication.cn/html/21895f0b-6ab8-a906-f5e7-28c784795ca2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[Authorization 类](../html/e89d602f-541c-8b98-af2b-02ee7b19618b.htm "Authorization 类")

[Authorization 方法](../html/d40a2af1-6e2f-d0cb-cf48-1096e1009001.htm "Authorization 方法")

[SetAuthorizationCode 方法](../html/0908319d-e4ef-b1c3-c302-8488e0a9fc93.htm "SetAuthorizationCode 方法 ")

[SetDllContact 方法](../html/bf5220e5-5a77-64a2-1533-e988b4fb0bd9.htm "SetDllContact 方法 ")

[SetHslCertificate 方法](../html/21895f0b-6ab8-a906-f5e7-28c784795ca2.htm "SetHslCertificate 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AuthorizationSetHslCertificate 方法 |

使用证书激活hslcommunication组件，证书请联系胡工科技获取，如果获取的是文件的话，那就是 SetHslCertificate( File.ReadAllBytes( "hsl.cert" ) ) 获取原始字节信息  
Use the certificate to activate the hslcommunication component, please contact Hugong Technology to obtain the certificate, if the obtained is a file,
that is, SetHslCertificate( File.ReadAllBytes( "hsl.cert") ) to obtain the raw byte information

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult SetHslCertificate(
	byte[] cert
)
```

```
Public Shared Function SetHslCertificate ( 
	cert As Byte()
) As OperateResult
```

```
public:
static OperateResult^ SetHslCertificate(
	array<unsigned char>^ cert
)
```

```
static member SetHslCertificate : 
        cert : byte[] -> OperateResult 
```

#### 参数

cert
:   类型：SystemByte  
    证书信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否激活成功

![](../icons/SectionExpanded.png)参见

#### 引用

[Authorization 类](e89d602f-541c-8b98-af2b-02ee7b19618b.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslExtension 类

[原文連結](http://api.hslcommunication.cn/html/05bbd254-20d0-2658-04fb-05c400c448dc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtension 类 |

扩展的辅助类方法

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunicationHslExtension

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static class HslExtension
```

```
<ExtensionAttribute>
Public NotInheritable Class HslExtension
```

```
[ExtensionAttribute]
public ref class HslExtension abstract sealed
```

```
[<AbstractClassAttribute>]
[<SealedAttribute>]
[<ExtensionAttribute>]
type HslExtension =  class end
```

HslExtension 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BeginReceiveResult(Socket, AsyncCallback)](eef13fbf-e165-2b82-bd82-2cba238e604d.htm) | 启动接收数据，需要传入回调方法，传递对象默认为socket本身  To start receiving data, you need to pass in a callback method. The default object is the socket itself. |
| 公共方法静态成员 | [BeginReceiveResult(Socket, AsyncCallback, Object)](af261659-d371-4e58-defb-f35323c8c716.htm) | 启动接收数据，需要传入回调方法，传递对象  To start receiving data, you need to pass in a callback method and pass an object |
| 公共方法静态成员 | [Contains](449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm) |  |
| 公共方法静态成员 | [CopyArrayT](8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm) | 拷贝当前的实例数组，是基于引用层的浅拷贝，如果类型为值类型，那就是深度拷贝，如果类型为引用类型，就是浅拷贝 |
| 公共方法静态成员 | [DecryptLargeData](bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm) | 对超过117字节限制的加密数据进行解密，解密出原始的字节数据，因为RSA本身限制了117字节，所以此处进行数据切割解密。  Decrypt the encrypted data that exceeds the 117-byte limit and decrypt the original byte data, because RSA itself limits 117 bytes, so the data is cut and decrypted here. |
| 公共方法静态成员 | [EncryptLargeData](120131e2-a737-5340-5b08-18e4632d6055.htm) | 对原始字节的数据进行加密，不限制长度，因为RSA本身限制了117字节，所以此处进行数据切割加密。  Encrypt the original byte data without limiting the length, because RSA itself limits 117 bytes, so the data is cut and encrypted here. |
| 公共方法静态成员 | [EndReceiveResult](190a250b-5fc7-2618-7255-a65bc6a328c8.htm) | 结束挂起的异步读取，返回读取的字节数，如果成功的情况。  Ends the pending asynchronous read and returns the number of bytes read, if successful. |
| 公共方法静态成员 | [EndsWith](0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm) |  |
| 公共方法静态成员 | [EveryByteAdd](000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm) | 将指定的数据添加到数组的每个元素上去，会改变每个元素的值 |
| 公共方法静态成员 | [GetBoolByIndex(Byte, Int32)](02376ace-8a06-10a4-600d-90db9f773e92.htm) | 获取Byte的第 boolIndex 偏移的bool值，比如3，就是第4位   Get the bool value of Byte's boolIndex offset, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(Byte, Int32)](1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm) | 获取Byte数组的第 boolIndex 偏移的bool值，这个偏移值可以为 10，就是第 1 个字节的 第3位   Get the bool value of the boolIndex offset of the Byte array. The offset value can be 10, which is the third bit of the first byte |
| 公共方法静态成员 | [GetBoolByIndex(Int16, Int32)](013b9689-81f1-1871-bd53-2a6ca2ae153a.htm) | 获取short类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the short type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(Int32, Int32)](65142f2b-8313-8bf6-db4c-5567801c19f3.htm) | 获取int类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the int type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(Int64, Int32)](fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm) | 获取long类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the long type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(UInt16, Int32)](4419a620-899f-6d5a-201a-2cb66c38969c.htm) | 获取ushort类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the ushort type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(UInt32, Int32)](19cac792-a632-936c-5bc3-052023e4e225.htm) | 获取uint类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the uint type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(UInt64, Int32)](5f91b230-1822-e4be-02c4-640045232711.htm) | 获取ulong类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the ulong type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolValue](dd581622-6e71-36a6-d734-dbfa30311d78.htm) | 获取Byte数组的第 bytIndex 个位置的，boolIndex偏移的bool值  Get the bool value of the bytIndex position of the Byte array and the boolIndex offset |
| 公共方法静态成员 | [GetBytes](adaf490f-7569-21af-cb68-70e8edb67449.htm) | 根据指定的字节长度信息，获取到随机的字节信息  Obtain random byte information according to the specified byte length information |
| 公共方法静态成员 | [GetPEMPrivateKey](db6c7187-5d63-6d79-dc01-136bc3158f22.htm) | 从RSA的算法对象里，获取到PEM格式的原始私钥数据，如果需要存储，或是显示，只需要 Convert.ToBase64String 方法  Obtain the original private key data in PEM format from the RSA algorithm object. If you need to store or display it, you only need the Convert.ToBase64String method |
| 公共方法静态成员 | [GetPEMPublicKey](567513da-22ad-4bd0-3d8a-b664b3f46109.htm) | 从RSA的算法对象里，获取到PEM格式的原始公钥数据，如果需要存储，或是显示，只需要 Convert.ToBase64String 方法  Obtain the original public key data in PEM format from the RSA algorithm object. If you need to store or display it, you only need the Convert.ToBase64String method |
| 公共方法静态成员 | [GetStringOrEndChar](aebf584d-117e-7e1d-d668-af3e6855a357.htm) | 从字节数组里提取字符串数据，如果碰到0x00字节，就直接结束 |
| 公共方法静态成员代码示例 | [GetValueOrDefaultT](3ebfcbba-df43-e997-10c0-316d68ee8534.htm) | 一个泛型方法，提供json对象的数据读取  A generic method that provides data read for a JSON object |
| 公共方法静态成员 | [IncreaseByT](bc6933f2-a930-5b8a-5226-b01cc8827f15.htm) | 将指定的数据添加到数组的每个元素上去，使用表达式树的形式实现，将会修改原数组。不适用byte类型 |
| 公共方法静态成员 | [IniSerialByFormatString](3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm) | 使用格式化的串口参数信息来初始化串口的参数，举例：9600-8-N-1，分别表示波特率，数据位，奇偶校验，停止位，当然也可以携带串口名称，例如：COM3-9600-8-N-1，linux环境也是支持的。  Use the formatted serial port parameter information to initialize the serial port parameters, for example: 9600-8-N-1, which means baud rate, data bit, parity, stop bit, of course, can also carry the serial port name, for example: COM3- 9600-8-N-1, linux environment is also supported. |
| 公共方法静态成员代码示例 | [RemoveBeginT](23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm) | 将一个数组的前面指定位数移除，返回新的一个数组  Removes the preceding specified number of bits in a array, returning a new array |
| 公共方法静态成员代码示例 | [RemoveDoubleT](1dba1036-8409-b99b-692d-af13e1f370ee.htm) | 将一个数组的前后移除指定位数，返回新的一个数组  Removes a array before and after the specified number of bits, returning a new array |
| 公共方法静态成员 | [RemoveLast(String, Int32)](fd57d397-fc74-c8fa-4318-578c874e30e3.htm) | 移除指定字符串数据的最后 length 个字符。如果字符串本身的长度不足 length，则返回为空字符串。  Remove the last "length" characters of the specified string data. If the length of the string itself is less than length, an empty string is returned. |
| 公共方法静态成员代码示例 | [RemoveLastT(T, Int32)](e6a15fc0-92b7-8813-5765-810c21b67dd2.htm) | 将一个数组的后面指定位数移除，返回新的一个数组  Removes the specified number of digits after a array, returning a new array |
| 公共方法静态成员代码示例 | [ReverseByWord](e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm) | 将byte数组按照双字节进行反转，如果为单数的情况，则自动补齐  Reverses the byte array by double byte, or if the singular is the case, automatically |
| 公共方法静态成员 | [ReverseNewT](20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm) | 获取当前数组的倒序数组，这是一个新的实例，不改变原来的数组值  Get the reversed array of the current byte array, this is a new instance, does not change the original array value |
| 公共方法静态成员代码示例 | [SelectBeginT](d92e877f-363f-a5d0-615d-851d06e867a5.htm) | 选择一个数组的前面的几个数据信息  Select the begin few items of data information of a array |
| 公共方法静态成员代码示例 | [SelectLastT](e7b564ce-2e61-37f4-a41a-77a82340736a.htm) | 选择一个数组的后面的几个数据信息  Select the last few items of data information of a array |
| 公共方法静态成员代码示例 | [SelectMiddleT](89b90303-188b-4924-08e0-a214f1f51eea.htm) | 获取到数组里面的中间指定长度的数组  Get an array of the specified length in the array |
| 公共方法静态成员 | [SetBoolByIndex(Byte, Int32, Boolean)](a099e2a4-b398-77c2-97f5-aaffb6f13749.htm) | 设置Byte的第 boolIndex 位的bool值，可以强制为 true 或是 false, 不影响其他的位  Set the bool value of the boolIndex bit of Byte, which can be forced to true or false, without affecting other bits |
| 公共方法静态成员 | [SetBoolByIndex(Byte, Int32, Boolean)](1395f586-f24c-594e-34c8-e284e2e03eae.htm) | 设置Byte[]的第 boolIndex 位的bool值，可以强制为 true 或是 false, 不影响其他的位，如果是第 10 位，则表示第 1 个字节的第 2 位（都是从 0 地址开始算的）  Set the bool value of the boolIndex bit of Byte[], which can be forced to true or false, without affecting other bits. If it is the 10th bit, it means the second bit of the first byte (both starting from the 0 address Calculated) |
| 公共方法静态成员 | [SetBoolByIndex(Int16, Int32, Boolean)](62c62b40-6de1-8515-821f-21c8bbac7ad3.htm) | 修改short数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~15，之外的值会引发异常  Modify a bit of short data and return the modified value without affecting the original value. Bit index is 0~15, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(Int32, Int32, Boolean)](20f0d9e3-698b-6338-b89e-be9ddb66b082.htm) | 修改int数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~31，之外的值会引发异常  Modify a bit of int data and return the modified value without affecting the original value. Bit index is 0~31, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(Int64, Int32, Boolean)](5b0cf80c-d2d9-e5eb-3261-0270081610c0.htm) | 修改long数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~63，之外的值会引发异常  Modify a bit of long data and return the modified value without affecting the original value. Bit index is 0~63, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(UInt16, Int32, Boolean)](b2a0cf23-b850-6389-792d-b9c314ee8bdb.htm) | 修改ushort数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~15，之外的值会引发异常  Modify a bit of ushort data and return the modified value without affecting the original value. Bit index is 0~15, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(UInt32, Int32, Boolean)](8fa64052-d88d-74b4-9494-b167c7c8802e.htm) | 修改uint数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~31，之外的值会引发异常  Modify a bit of uint data and return the modified value without affecting the original value. Bit index is 0~31, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(UInt64, Int32, Boolean)](f006b687-4094-7db7-71b0-c7773b59526d.htm) | 修改ulong数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~63，之外的值会引发异常  Modify a bit of ulong data and return the modified value without affecting the original value. Bit index is 0~63, values outside will raise an exception |
| 公共方法静态成员 | [SetKeepAlive](8f2949d7-33fd-6175-3d11-50e33acc9d15.htm) | 设置套接字的活动时间和活动间歇时间，此值会设置到socket低级别的控制中，传入值如果为负数，则表示不使用 KeepAlive 功能。  Set the active time and active intermittent time of the socket. This value will be set to the low-level control of the socket. If the incoming value is a negative number, it means that the KeepAlive function is not used. |
| 公共方法静态成员代码示例 | [SpliceArrayT](c34cc020-1531-9910-6143-8ea611240268.htm) | 拼接任意个泛型数组为一个总的泛型数组对象，采用深度拷贝实现。  Splicing any number of generic arrays into a total generic array object is implemented using deep copy. |
| 公共方法静态成员 | [SplitDot](b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm) | 根据英文小数点进行切割字符串，去除空白的字符  Cut the string according to the English decimal point and remove the blank characters |
| 公共方法静态成员 | [StartsWith](b6f380b3-fc65-1e8f-aa86-213d03447970.htm) |  |
| 公共方法静态成员 | [StartsWithAndNumber(String, String)](6dbe85fb-c0a7-caf1-b76f-43ba54900f8b.htm) | 判断一个地址是否是指定的字符串开头，并且后面跟着数字，如果是，返回 true，反之，返回 false  Determines whether an address starts with a specified string followed by a number, returns true if so, and false if not |
| 公共方法静态成员 | [StartsWithAndNumber(String, String)](f0ed3b14-3949-6115-a31c-c6962ef891a6.htm) | 判断一个地址是否是指定的字符串开头，并且后面跟着数字，如果是，返回 true，反之，返回 false  Determines whether an address starts with a specified string followed by a number, returns true if so, and false if not |
| 公共方法静态成员 | [ToArrayStringT(T)](847f5536-d9c8-81e3-4d1c-31338062da03.htm) | 将数组格式化为显示的字符串的信息，支持所有的类型对象  Formats the array into the displayed string information, supporting all types of objects |
| 公共方法静态成员 | [ToArrayStringT(T, String)](c4522974-3c09-4d30-c637-50d8381d5efe.htm) | 将数组格式化为显示的字符串的信息，支持所有的类型对象  Formats the array into the displayed string information, supporting all types of objects |
| 公共方法静态成员代码示例 | [ToBoolArray(Byte)](98315715-80c3-9e73-1117-5236dd1607a4.htm) | 从Byte数组中提取所有的位数组  Extracts a bit array from a byte array, length represents the number of digits |
| 公共方法静态成员代码示例 | [ToBoolArray(Byte, Int32)](f20b2a0c-806b-5a6a-6a62-6736685b6849.htm) | 从Byte数组中提取位数组，length代表位数，例如数组 03 A1 长度10转为 [1100 0000 10]  Extracts a bit array from a byte array, length represents the number of digits |
| 公共方法静态成员代码示例 | [ToByteArray](db7298ee-87e8-df3a-59fa-5892d82016d5.htm) | 将bool数组转换到byte数组  Converting a bool array to a byte array |
| 公共方法静态成员 | [ToFormatString](fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm) | 将一个串口对象的基本配置参数进行格式化字符串，例如 COM3-9600-8-N-1  Format the basic configuration parameters of a serial port object as a string, for example, COM3-9600-8-N-1 |
| 公共方法静态成员代码示例 | [ToHexBytes](ff506e59-b0c9-918d-06f6-1e175bbda013.htm) | 将16进制的字符串转化成Byte数据，将检测每2个字符转化，也就是说，中间可以是任意字符  Converts a 16-character string into byte data, which will detect every 2 characters converted, that is, the middle can be any character |
| 公共方法静态成员代码示例 | [ToHexString(Byte)](cd4d5893-7d5a-be81-daef-409a431c12c7.htm) | 字节数据转化成16进制表示的字符串  Byte data into a string of 16 binary representations |
| 公共方法静态成员代码示例 | [ToHexString(Byte, Char)](b62a7672-668a-899e-0e35-ce775437c014.htm) | 字节数据转化成16进制表示的字符串  Byte data into a string of 16 binary representations |
| 公共方法静态成员代码示例 | [ToHexString(Byte, Char, Int32, String)](4ba455df-be1f-54a7-a1ec-ad0954791b94.htm) | 字节数据转化成16进制表示的字符串  Byte data into a string of 16 binary representations |
| 公共方法静态成员 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. |
| 公共方法静态成员 | [ToStringArrayT(String)](1eef62af-6ee5-ab44-8f03-595fd933628d.htm) | 将字符串数组转换为实际的数据数组。支持byte,sbyte,bool,short,ushort,int,uint,long,ulong,float,double，使用默认的十进制，例如字符串格式[1,2,3,4,5]，可以转成实际的数组对象  Converts a string array into an actual data array. Support byte, sbyte, bool, short, ushort, int, uint, long, ulong, float, double, use the default decimal, such as the string format [1,2,3,4,5], which can be converted into an actual array Object |
| 公共方法静态成员 | [ToStringArrayT(String, FuncString, T)](3cc8bc2a-d0f8-0e41-488b-600f03fcc47f.htm) | 将字符串数组转换为实际的数据数组。例如字符串格式[1,2,3,4,5]，可以转成实际的数组对象  Converts a string array into an actual data array. For example, the string format [1,2,3,4,5] can be converted into an actual array object |
| 公共方法静态成员 | [Write](ea8cba73-29a6-0c12-de93-2960758fc8d7.htm) |  |
| 公共方法静态成员 | [WriteReverse](e51a3e83-1552-09c9-f01d-015fb5d6c226.htm) | 将UInt16数据写入到字节流，字节顺序为相反  Write UInt16 data to the byte stream, the byte order is reversed |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslExtension 方法

[原文連結](http://api.hslcommunication.cn/html/32366708-60b9-6e27-c174-0edbd522cbd4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[Contains 方法](../html/449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm "Contains 方法 ")

[CopyArray(T) 方法](../html/8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm "CopyArray(T) 方法 ")

[DecryptLargeData 方法](../html/bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm "DecryptLargeData 方法 ")

[EncryptLargeData 方法](../html/120131e2-a737-5340-5b08-18e4632d6055.htm "EncryptLargeData 方法 ")

[EndReceiveResult 方法](../html/190a250b-5fc7-2618-7255-a65bc6a328c8.htm "EndReceiveResult 方法 ")

[EndsWith 方法](../html/0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm "EndsWith 方法 ")

[EveryByteAdd 方法](../html/000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm "EveryByteAdd 方法 ")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolValue 方法](../html/dd581622-6e71-36a6-d734-dbfa30311d78.htm "GetBoolValue 方法 ")

[GetBytes 方法](../html/adaf490f-7569-21af-cb68-70e8edb67449.htm "GetBytes 方法 ")

[GetPEMPrivateKey 方法](../html/db6c7187-5d63-6d79-dc01-136bc3158f22.htm "GetPEMPrivateKey 方法 ")

[GetPEMPublicKey 方法](../html/567513da-22ad-4bd0-3d8a-b664b3f46109.htm "GetPEMPublicKey 方法 ")

[GetStringOrEndChar 方法](../html/aebf584d-117e-7e1d-d668-af3e6855a357.htm "GetStringOrEndChar 方法 ")

[GetValueOrDefault(T) 方法](../html/3ebfcbba-df43-e997-10c0-316d68ee8534.htm "GetValueOrDefault(T) 方法 ")

[IncreaseBy(T) 方法](../html/bc6933f2-a930-5b8a-5226-b01cc8827f15.htm "IncreaseBy(T) 方法 ")

[IniSerialByFormatString 方法](../html/3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm "IniSerialByFormatString 方法 ")

[RemoveBegin(T) 方法](../html/23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm "RemoveBegin(T) 方法 ")

[RemoveDouble(T) 方法](../html/1dba1036-8409-b99b-692d-af13e1f370ee.htm "RemoveDouble(T) 方法 ")

[RemoveLast 方法](../html/15db4c58-5c86-79a5-2098-146c39aa3436.htm "RemoveLast 方法 ")

[ReverseByWord 方法](../html/e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm "ReverseByWord 方法 ")

[ReverseNew(T) 方法](../html/20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm "ReverseNew(T) 方法 ")

[SelectBegin(T) 方法](../html/d92e877f-363f-a5d0-615d-851d06e867a5.htm "SelectBegin(T) 方法 ")

[SelectLast(T) 方法](../html/e7b564ce-2e61-37f4-a41a-77a82340736a.htm "SelectLast(T) 方法 ")

[SelectMiddle(T) 方法](../html/89b90303-188b-4924-08e0-a214f1f51eea.htm "SelectMiddle(T) 方法 ")

[SetBoolByIndex 方法](../html/92627e54-07a1-ca29-1a63-c1cd2b28b786.htm "SetBoolByIndex 方法 ")

[SetKeepAlive 方法](../html/8f2949d7-33fd-6175-3d11-50e33acc9d15.htm "SetKeepAlive 方法 ")

[SpliceArray(T) 方法](../html/c34cc020-1531-9910-6143-8ea611240268.htm "SpliceArray(T) 方法 ")

[SplitDot 方法](../html/b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm "SplitDot 方法 ")

[StartsWith 方法](../html/b6f380b3-fc65-1e8f-aa86-213d03447970.htm "StartsWith 方法 ")

[StartsWithAndNumber 方法](../html/f71b5bea-6b5e-89c7-e922-c7ad06cd4cb9.htm "StartsWithAndNumber 方法 ")

[ToArrayString 方法](../html/e7ca844f-e63d-1204-9832-a7f8f43ece36.htm "ToArrayString 方法 ")

[ToBoolArray 方法](../html/1bcba304-f35a-23cf-d17f-f539b806e433.htm "ToBoolArray 方法 ")

[ToByteArray 方法](../html/db7298ee-87e8-df3a-59fa-5892d82016d5.htm "ToByteArray 方法 ")

[ToFormatString 方法](../html/fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm "ToFormatString 方法 ")

[ToHexBytes 方法](../html/ff506e59-b0c9-918d-06f6-1e175bbda013.htm "ToHexBytes 方法 ")

[ToHexString 方法](../html/563d3090-af73-0aaf-aedc-047bd1dedb00.htm "ToHexString 方法 ")

[ToJsonString 方法](../html/828b0b70-5381-936e-e606-e6c0a905f688.htm "ToJsonString 方法 ")

[ToStringArray 方法](../html/58230997-75a3-0a6b-bd9a-695a51820891.htm "ToStringArray 方法 ")

[Write 方法](../html/ea8cba73-29a6-0c12-de93-2960758fc8d7.htm "Write 方法 ")

[WriteReverse 方法](../html/e51a3e83-1552-09c9-f01d-015fb5d6c226.htm "WriteReverse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtension 方法 |

[HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BeginReceiveResult(Socket, AsyncCallback)](eef13fbf-e165-2b82-bd82-2cba238e604d.htm) | 启动接收数据，需要传入回调方法，传递对象默认为socket本身  To start receiving data, you need to pass in a callback method. The default object is the socket itself. |
| 公共方法静态成员 | [BeginReceiveResult(Socket, AsyncCallback, Object)](af261659-d371-4e58-defb-f35323c8c716.htm) | 启动接收数据，需要传入回调方法，传递对象  To start receiving data, you need to pass in a callback method and pass an object |
| 公共方法静态成员 | [Contains](449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm) |  |
| 公共方法静态成员 | [CopyArrayT](8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm) | 拷贝当前的实例数组，是基于引用层的浅拷贝，如果类型为值类型，那就是深度拷贝，如果类型为引用类型，就是浅拷贝 |
| 公共方法静态成员 | [DecryptLargeData](bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm) | 对超过117字节限制的加密数据进行解密，解密出原始的字节数据，因为RSA本身限制了117字节，所以此处进行数据切割解密。  Decrypt the encrypted data that exceeds the 117-byte limit and decrypt the original byte data, because RSA itself limits 117 bytes, so the data is cut and decrypted here. |
| 公共方法静态成员 | [EncryptLargeData](120131e2-a737-5340-5b08-18e4632d6055.htm) | 对原始字节的数据进行加密，不限制长度，因为RSA本身限制了117字节，所以此处进行数据切割加密。  Encrypt the original byte data without limiting the length, because RSA itself limits 117 bytes, so the data is cut and encrypted here. |
| 公共方法静态成员 | [EndReceiveResult](190a250b-5fc7-2618-7255-a65bc6a328c8.htm) | 结束挂起的异步读取，返回读取的字节数，如果成功的情况。  Ends the pending asynchronous read and returns the number of bytes read, if successful. |
| 公共方法静态成员 | [EndsWith](0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm) |  |
| 公共方法静态成员 | [EveryByteAdd](000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm) | 将指定的数据添加到数组的每个元素上去，会改变每个元素的值 |
| 公共方法静态成员 | [GetBoolByIndex(Byte, Int32)](02376ace-8a06-10a4-600d-90db9f773e92.htm) | 获取Byte的第 boolIndex 偏移的bool值，比如3，就是第4位   Get the bool value of Byte's boolIndex offset, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(Byte, Int32)](1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm) | 获取Byte数组的第 boolIndex 偏移的bool值，这个偏移值可以为 10，就是第 1 个字节的 第3位   Get the bool value of the boolIndex offset of the Byte array. The offset value can be 10, which is the third bit of the first byte |
| 公共方法静态成员 | [GetBoolByIndex(Int16, Int32)](013b9689-81f1-1871-bd53-2a6ca2ae153a.htm) | 获取short类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the short type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(Int32, Int32)](65142f2b-8313-8bf6-db4c-5567801c19f3.htm) | 获取int类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the int type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(Int64, Int32)](fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm) | 获取long类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the long type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(UInt16, Int32)](4419a620-899f-6d5a-201a-2cb66c38969c.htm) | 获取ushort类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the ushort type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(UInt32, Int32)](19cac792-a632-936c-5bc3-052023e4e225.htm) | 获取uint类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the uint type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(UInt64, Int32)](5f91b230-1822-e4be-02c4-640045232711.htm) | 获取ulong类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the ulong type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolValue](dd581622-6e71-36a6-d734-dbfa30311d78.htm) | 获取Byte数组的第 bytIndex 个位置的，boolIndex偏移的bool值  Get the bool value of the bytIndex position of the Byte array and the boolIndex offset |
| 公共方法静态成员 | [GetBytes](adaf490f-7569-21af-cb68-70e8edb67449.htm) | 根据指定的字节长度信息，获取到随机的字节信息  Obtain random byte information according to the specified byte length information |
| 公共方法静态成员 | [GetPEMPrivateKey](db6c7187-5d63-6d79-dc01-136bc3158f22.htm) | 从RSA的算法对象里，获取到PEM格式的原始私钥数据，如果需要存储，或是显示，只需要 Convert.ToBase64String 方法  Obtain the original private key data in PEM format from the RSA algorithm object. If you need to store or display it, you only need the Convert.ToBase64String method |
| 公共方法静态成员 | [GetPEMPublicKey](567513da-22ad-4bd0-3d8a-b664b3f46109.htm) | 从RSA的算法对象里，获取到PEM格式的原始公钥数据，如果需要存储，或是显示，只需要 Convert.ToBase64String 方法  Obtain the original public key data in PEM format from the RSA algorithm object. If you need to store or display it, you only need the Convert.ToBase64String method |
| 公共方法静态成员 | [GetStringOrEndChar](aebf584d-117e-7e1d-d668-af3e6855a357.htm) | 从字节数组里提取字符串数据，如果碰到0x00字节，就直接结束 |
| 公共方法静态成员代码示例 | [GetValueOrDefaultT](3ebfcbba-df43-e997-10c0-316d68ee8534.htm) | 一个泛型方法，提供json对象的数据读取  A generic method that provides data read for a JSON object |
| 公共方法静态成员 | [IncreaseByT](bc6933f2-a930-5b8a-5226-b01cc8827f15.htm) | 将指定的数据添加到数组的每个元素上去，使用表达式树的形式实现，将会修改原数组。不适用byte类型 |
| 公共方法静态成员 | [IniSerialByFormatString](3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm) | 使用格式化的串口参数信息来初始化串口的参数，举例：9600-8-N-1，分别表示波特率，数据位，奇偶校验，停止位，当然也可以携带串口名称，例如：COM3-9600-8-N-1，linux环境也是支持的。  Use the formatted serial port parameter information to initialize the serial port parameters, for example: 9600-8-N-1, which means baud rate, data bit, parity, stop bit, of course, can also carry the serial port name, for example: COM3- 9600-8-N-1, linux environment is also supported. |
| 公共方法静态成员代码示例 | [RemoveBeginT](23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm) | 将一个数组的前面指定位数移除，返回新的一个数组  Removes the preceding specified number of bits in a array, returning a new array |
| 公共方法静态成员代码示例 | [RemoveDoubleT](1dba1036-8409-b99b-692d-af13e1f370ee.htm) | 将一个数组的前后移除指定位数，返回新的一个数组  Removes a array before and after the specified number of bits, returning a new array |
| 公共方法静态成员 | [RemoveLast(String, Int32)](fd57d397-fc74-c8fa-4318-578c874e30e3.htm) | 移除指定字符串数据的最后 length 个字符。如果字符串本身的长度不足 length，则返回为空字符串。  Remove the last "length" characters of the specified string data. If the length of the string itself is less than length, an empty string is returned. |
| 公共方法静态成员代码示例 | [RemoveLastT(T, Int32)](e6a15fc0-92b7-8813-5765-810c21b67dd2.htm) | 将一个数组的后面指定位数移除，返回新的一个数组  Removes the specified number of digits after a array, returning a new array |
| 公共方法静态成员代码示例 | [ReverseByWord](e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm) | 将byte数组按照双字节进行反转，如果为单数的情况，则自动补齐  Reverses the byte array by double byte, or if the singular is the case, automatically |
| 公共方法静态成员 | [ReverseNewT](20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm) | 获取当前数组的倒序数组，这是一个新的实例，不改变原来的数组值  Get the reversed array of the current byte array, this is a new instance, does not change the original array value |
| 公共方法静态成员代码示例 | [SelectBeginT](d92e877f-363f-a5d0-615d-851d06e867a5.htm) | 选择一个数组的前面的几个数据信息  Select the begin few items of data information of a array |
| 公共方法静态成员代码示例 | [SelectLastT](e7b564ce-2e61-37f4-a41a-77a82340736a.htm) | 选择一个数组的后面的几个数据信息  Select the last few items of data information of a array |
| 公共方法静态成员代码示例 | [SelectMiddleT](89b90303-188b-4924-08e0-a214f1f51eea.htm) | 获取到数组里面的中间指定长度的数组  Get an array of the specified length in the array |
| 公共方法静态成员 | [SetBoolByIndex(Byte, Int32, Boolean)](a099e2a4-b398-77c2-97f5-aaffb6f13749.htm) | 设置Byte的第 boolIndex 位的bool值，可以强制为 true 或是 false, 不影响其他的位  Set the bool value of the boolIndex bit of Byte, which can be forced to true or false, without affecting other bits |
| 公共方法静态成员 | [SetBoolByIndex(Byte, Int32, Boolean)](1395f586-f24c-594e-34c8-e284e2e03eae.htm) | 设置Byte[]的第 boolIndex 位的bool值，可以强制为 true 或是 false, 不影响其他的位，如果是第 10 位，则表示第 1 个字节的第 2 位（都是从 0 地址开始算的）  Set the bool value of the boolIndex bit of Byte[], which can be forced to true or false, without affecting other bits. If it is the 10th bit, it means the second bit of the first byte (both starting from the 0 address Calculated) |
| 公共方法静态成员 | [SetBoolByIndex(Int16, Int32, Boolean)](62c62b40-6de1-8515-821f-21c8bbac7ad3.htm) | 修改short数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~15，之外的值会引发异常  Modify a bit of short data and return the modified value without affecting the original value. Bit index is 0~15, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(Int32, Int32, Boolean)](20f0d9e3-698b-6338-b89e-be9ddb66b082.htm) | 修改int数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~31，之外的值会引发异常  Modify a bit of int data and return the modified value without affecting the original value. Bit index is 0~31, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(Int64, Int32, Boolean)](5b0cf80c-d2d9-e5eb-3261-0270081610c0.htm) | 修改long数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~63，之外的值会引发异常  Modify a bit of long data and return the modified value without affecting the original value. Bit index is 0~63, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(UInt16, Int32, Boolean)](b2a0cf23-b850-6389-792d-b9c314ee8bdb.htm) | 修改ushort数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~15，之外的值会引发异常  Modify a bit of ushort data and return the modified value without affecting the original value. Bit index is 0~15, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(UInt32, Int32, Boolean)](8fa64052-d88d-74b4-9494-b167c7c8802e.htm) | 修改uint数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~31，之外的值会引发异常  Modify a bit of uint data and return the modified value without affecting the original value. Bit index is 0~31, values outside will raise an exception |
| 公共方法静态成员 | [SetBoolByIndex(UInt64, Int32, Boolean)](f006b687-4094-7db7-71b0-c7773b59526d.htm) | 修改ulong数据的某个位，并且返回修改后的值，不影响原来的值。位索引为 0~63，之外的值会引发异常  Modify a bit of ulong data and return the modified value without affecting the original value. Bit index is 0~63, values outside will raise an exception |
| 公共方法静态成员 | [SetKeepAlive](8f2949d7-33fd-6175-3d11-50e33acc9d15.htm) | 设置套接字的活动时间和活动间歇时间，此值会设置到socket低级别的控制中，传入值如果为负数，则表示不使用 KeepAlive 功能。  Set the active time and active intermittent time of the socket. This value will be set to the low-level control of the socket. If the incoming value is a negative number, it means that the KeepAlive function is not used. |
| 公共方法静态成员代码示例 | [SpliceArrayT](c34cc020-1531-9910-6143-8ea611240268.htm) | 拼接任意个泛型数组为一个总的泛型数组对象，采用深度拷贝实现。  Splicing any number of generic arrays into a total generic array object is implemented using deep copy. |
| 公共方法静态成员 | [SplitDot](b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm) | 根据英文小数点进行切割字符串，去除空白的字符  Cut the string according to the English decimal point and remove the blank characters |
| 公共方法静态成员 | [StartsWith](b6f380b3-fc65-1e8f-aa86-213d03447970.htm) |  |
| 公共方法静态成员 | [StartsWithAndNumber(String, String)](6dbe85fb-c0a7-caf1-b76f-43ba54900f8b.htm) | 判断一个地址是否是指定的字符串开头，并且后面跟着数字，如果是，返回 true，反之，返回 false  Determines whether an address starts with a specified string followed by a number, returns true if so, and false if not |
| 公共方法静态成员 | [StartsWithAndNumber(String, String)](f0ed3b14-3949-6115-a31c-c6962ef891a6.htm) | 判断一个地址是否是指定的字符串开头，并且后面跟着数字，如果是，返回 true，反之，返回 false  Determines whether an address starts with a specified string followed by a number, returns true if so, and false if not |
| 公共方法静态成员 | [ToArrayStringT(T)](847f5536-d9c8-81e3-4d1c-31338062da03.htm) | 将数组格式化为显示的字符串的信息，支持所有的类型对象  Formats the array into the displayed string information, supporting all types of objects |
| 公共方法静态成员 | [ToArrayStringT(T, String)](c4522974-3c09-4d30-c637-50d8381d5efe.htm) | 将数组格式化为显示的字符串的信息，支持所有的类型对象  Formats the array into the displayed string information, supporting all types of objects |
| 公共方法静态成员代码示例 | [ToBoolArray(Byte)](98315715-80c3-9e73-1117-5236dd1607a4.htm) | 从Byte数组中提取所有的位数组  Extracts a bit array from a byte array, length represents the number of digits |
| 公共方法静态成员代码示例 | [ToBoolArray(Byte, Int32)](f20b2a0c-806b-5a6a-6a62-6736685b6849.htm) | 从Byte数组中提取位数组，length代表位数，例如数组 03 A1 长度10转为 [1100 0000 10]  Extracts a bit array from a byte array, length represents the number of digits |
| 公共方法静态成员代码示例 | [ToByteArray](db7298ee-87e8-df3a-59fa-5892d82016d5.htm) | 将bool数组转换到byte数组  Converting a bool array to a byte array |
| 公共方法静态成员 | [ToFormatString](fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm) | 将一个串口对象的基本配置参数进行格式化字符串，例如 COM3-9600-8-N-1  Format the basic configuration parameters of a serial port object as a string, for example, COM3-9600-8-N-1 |
| 公共方法静态成员代码示例 | [ToHexBytes](ff506e59-b0c9-918d-06f6-1e175bbda013.htm) | 将16进制的字符串转化成Byte数据，将检测每2个字符转化，也就是说，中间可以是任意字符  Converts a 16-character string into byte data, which will detect every 2 characters converted, that is, the middle can be any character |
| 公共方法静态成员代码示例 | [ToHexString(Byte)](cd4d5893-7d5a-be81-daef-409a431c12c7.htm) | 字节数据转化成16进制表示的字符串  Byte data into a string of 16 binary representations |
| 公共方法静态成员代码示例 | [ToHexString(Byte, Char)](b62a7672-668a-899e-0e35-ce775437c014.htm) | 字节数据转化成16进制表示的字符串  Byte data into a string of 16 binary representations |
| 公共方法静态成员代码示例 | [ToHexString(Byte, Char, Int32, String)](4ba455df-be1f-54a7-a1ec-ad0954791b94.htm) | 字节数据转化成16进制表示的字符串  Byte data into a string of 16 binary representations |
| 公共方法静态成员 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. |
| 公共方法静态成员 | [ToStringArrayT(String)](1eef62af-6ee5-ab44-8f03-595fd933628d.htm) | 将字符串数组转换为实际的数据数组。支持byte,sbyte,bool,short,ushort,int,uint,long,ulong,float,double，使用默认的十进制，例如字符串格式[1,2,3,4,5]，可以转成实际的数组对象  Converts a string array into an actual data array. Support byte, sbyte, bool, short, ushort, int, uint, long, ulong, float, double, use the default decimal, such as the string format [1,2,3,4,5], which can be converted into an actual array Object |
| 公共方法静态成员 | [ToStringArrayT(String, FuncString, T)](3cc8bc2a-d0f8-0e41-488b-600f03fcc47f.htm) | 将字符串数组转换为实际的数据数组。例如字符串格式[1,2,3,4,5]，可以转成实际的数组对象  Converts a string array into an actual data array. For example, the string format [1,2,3,4,5] can be converted into an actual array object |
| 公共方法静态成员 | [Write](ea8cba73-29a6-0c12-de93-2960758fc8d7.htm) |  |
| 公共方法静态成员 | [WriteReverse](e51a3e83-1552-09c9-f01d-015fb5d6c226.htm) | 将UInt16数据写入到字节流，字节顺序为相反  Write UInt16 data to the byte stream, the byte order is reversed |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeginReceiveResult 方法 

[原文連結](http://api.hslcommunication.cn/html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[BeginReceiveResult 方法 (Socket, AsyncCallback)](../html/eef13fbf-e165-2b82-bd82-2cba238e604d.htm "BeginReceiveResult 方法 (Socket, AsyncCallback)")

[BeginReceiveResult 方法 (Socket, AsyncCallback, Object)](../html/af261659-d371-4e58-defb-f35323c8c716.htm "BeginReceiveResult 方法 (Socket, AsyncCallback, Object)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionBeginReceiveResult 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BeginReceiveResult(Socket, AsyncCallback)](eef13fbf-e165-2b82-bd82-2cba238e604d.htm) | 启动接收数据，需要传入回调方法，传递对象默认为socket本身  To start receiving data, you need to pass in a callback method. The default object is the socket itself. |
| 公共方法静态成员 | [BeginReceiveResult(Socket, AsyncCallback, Object)](af261659-d371-4e58-defb-f35323c8c716.htm) | 启动接收数据，需要传入回调方法，传递对象  To start receiving data, you need to pass in a callback method and pass an object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeginReceiveResult 方法 (Socket, AsyncCallback)

[原文連結](http://api.hslcommunication.cn/html/eef13fbf-e165-2b82-bd82-2cba238e604d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[BeginReceiveResult 方法 (Socket, AsyncCallback)](../html/eef13fbf-e165-2b82-bd82-2cba238e604d.htm "BeginReceiveResult 方法 (Socket, AsyncCallback)")

[BeginReceiveResult 方法 (Socket, AsyncCallback, Object)](../html/af261659-d371-4e58-defb-f35323c8c716.htm "BeginReceiveResult 方法 (Socket, AsyncCallback, Object)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionBeginReceiveResult 方法 (Socket, AsyncCallback) |

启动接收数据，需要传入回调方法，传递对象默认为socket本身  
To start receiving data, you need to pass in a callback method. The default object is the socket itself.

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult BeginReceiveResult(
	this Socket socket,
	AsyncCallback callback
)
```

```
<ExtensionAttribute>
Public Shared Function BeginReceiveResult ( 
	socket As Socket,
	callback As AsyncCallback
) As OperateResult
```

```
public:
[ExtensionAttribute]
static OperateResult^ BeginReceiveResult(
	Socket^ socket, 
	AsyncCallback^ callback
)
```

```
[<ExtensionAttribute>]
static member BeginReceiveResult : 
        socket : Socket * 
        callback : AsyncCallback -> OperateResult 
```

#### 参数

socket
:   类型：System.Net.SocketsSocket  
    socket对象

callback
:   类型：SystemAsyncCallback  
    回调方法

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否启动成功

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 Socket 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[BeginReceiveResult 重载](cdcfdc10-5f34-00c6-142b-95265069e58a.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeginReceiveResult 方法 (Socket, AsyncCallback, Object)

[原文連結](http://api.hslcommunication.cn/html/af261659-d371-4e58-defb-f35323c8c716.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[BeginReceiveResult 方法 (Socket, AsyncCallback)](../html/eef13fbf-e165-2b82-bd82-2cba238e604d.htm "BeginReceiveResult 方法 (Socket, AsyncCallback)")

[BeginReceiveResult 方法 (Socket, AsyncCallback, Object)](../html/af261659-d371-4e58-defb-f35323c8c716.htm "BeginReceiveResult 方法 (Socket, AsyncCallback, Object)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionBeginReceiveResult 方法 (Socket, AsyncCallback, Object) |

启动接收数据，需要传入回调方法，传递对象  
To start receiving data, you need to pass in a callback method and pass an object

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult BeginReceiveResult(
	this Socket socket,
	AsyncCallback callback,
	Object obj
)
```

```
<ExtensionAttribute>
Public Shared Function BeginReceiveResult ( 
	socket As Socket,
	callback As AsyncCallback,
	obj As Object
) As OperateResult
```

```
public:
[ExtensionAttribute]
static OperateResult^ BeginReceiveResult(
	Socket^ socket, 
	AsyncCallback^ callback, 
	Object^ obj
)
```

```
[<ExtensionAttribute>]
static member BeginReceiveResult : 
        socket : Socket * 
        callback : AsyncCallback * 
        obj : Object -> OperateResult 
```

#### 参数

socket
:   类型：System.Net.SocketsSocket  
    socket对象

callback
:   类型：SystemAsyncCallback  
    回调方法

obj
:   类型：SystemObject  
    数据对象

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否启动成功

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 Socket 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[BeginReceiveResult 重载](cdcfdc10-5f34-00c6-142b-95265069e58a.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Contains 方法 

[原文連結](http://api.hslcommunication.cn/html/449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[Contains 方法](../html/449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm "Contains 方法 ")

[CopyArray(T) 方法](../html/8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm "CopyArray(T) 方法 ")

[DecryptLargeData 方法](../html/bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm "DecryptLargeData 方法 ")

[EncryptLargeData 方法](../html/120131e2-a737-5340-5b08-18e4632d6055.htm "EncryptLargeData 方法 ")

[EndReceiveResult 方法](../html/190a250b-5fc7-2618-7255-a65bc6a328c8.htm "EndReceiveResult 方法 ")

[EndsWith 方法](../html/0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm "EndsWith 方法 ")

[EveryByteAdd 方法](../html/000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm "EveryByteAdd 方法 ")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolValue 方法](../html/dd581622-6e71-36a6-d734-dbfa30311d78.htm "GetBoolValue 方法 ")

[GetBytes 方法](../html/adaf490f-7569-21af-cb68-70e8edb67449.htm "GetBytes 方法 ")

[GetPEMPrivateKey 方法](../html/db6c7187-5d63-6d79-dc01-136bc3158f22.htm "GetPEMPrivateKey 方法 ")

[GetPEMPublicKey 方法](../html/567513da-22ad-4bd0-3d8a-b664b3f46109.htm "GetPEMPublicKey 方法 ")

[GetStringOrEndChar 方法](../html/aebf584d-117e-7e1d-d668-af3e6855a357.htm "GetStringOrEndChar 方法 ")

[GetValueOrDefault(T) 方法](../html/3ebfcbba-df43-e997-10c0-316d68ee8534.htm "GetValueOrDefault(T) 方法 ")

[IncreaseBy(T) 方法](../html/bc6933f2-a930-5b8a-5226-b01cc8827f15.htm "IncreaseBy(T) 方法 ")

[IniSerialByFormatString 方法](../html/3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm "IniSerialByFormatString 方法 ")

[RemoveBegin(T) 方法](../html/23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm "RemoveBegin(T) 方法 ")

[RemoveDouble(T) 方法](../html/1dba1036-8409-b99b-692d-af13e1f370ee.htm "RemoveDouble(T) 方法 ")

[RemoveLast 方法](../html/15db4c58-5c86-79a5-2098-146c39aa3436.htm "RemoveLast 方法 ")

[ReverseByWord 方法](../html/e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm "ReverseByWord 方法 ")

[ReverseNew(T) 方法](../html/20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm "ReverseNew(T) 方法 ")

[SelectBegin(T) 方法](../html/d92e877f-363f-a5d0-615d-851d06e867a5.htm "SelectBegin(T) 方法 ")

[SelectLast(T) 方法](../html/e7b564ce-2e61-37f4-a41a-77a82340736a.htm "SelectLast(T) 方法 ")

[SelectMiddle(T) 方法](../html/89b90303-188b-4924-08e0-a214f1f51eea.htm "SelectMiddle(T) 方法 ")

[SetBoolByIndex 方法](../html/92627e54-07a1-ca29-1a63-c1cd2b28b786.htm "SetBoolByIndex 方法 ")

[SetKeepAlive 方法](../html/8f2949d7-33fd-6175-3d11-50e33acc9d15.htm "SetKeepAlive 方法 ")

[SpliceArray(T) 方法](../html/c34cc020-1531-9910-6143-8ea611240268.htm "SpliceArray(T) 方法 ")

[SplitDot 方法](../html/b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm "SplitDot 方法 ")

[StartsWith 方法](../html/b6f380b3-fc65-1e8f-aa86-213d03447970.htm "StartsWith 方法 ")

[StartsWithAndNumber 方法](../html/f71b5bea-6b5e-89c7-e922-c7ad06cd4cb9.htm "StartsWithAndNumber 方法 ")

[ToArrayString 方法](../html/e7ca844f-e63d-1204-9832-a7f8f43ece36.htm "ToArrayString 方法 ")

[ToBoolArray 方法](../html/1bcba304-f35a-23cf-d17f-f539b806e433.htm "ToBoolArray 方法 ")

[ToByteArray 方法](../html/db7298ee-87e8-df3a-59fa-5892d82016d5.htm "ToByteArray 方法 ")

[ToFormatString 方法](../html/fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm "ToFormatString 方法 ")

[ToHexBytes 方法](../html/ff506e59-b0c9-918d-06f6-1e175bbda013.htm "ToHexBytes 方法 ")

[ToHexString 方法](../html/563d3090-af73-0aaf-aedc-047bd1dedb00.htm "ToHexString 方法 ")

[ToJsonString 方法](../html/828b0b70-5381-936e-e606-e6c0a905f688.htm "ToJsonString 方法 ")

[ToStringArray 方法](../html/58230997-75a3-0a6b-bd9a-695a51820891.htm "ToStringArray 方法 ")

[Write 方法](../html/ea8cba73-29a6-0c12-de93-2960758fc8d7.htm "Write 方法 ")

[WriteReverse 方法](../html/e51a3e83-1552-09c9-f01d-015fb5d6c226.htm "WriteReverse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionContains 方法 |

[缺少 "M:HslCommunication.HslExtension.Contains(System.String,System.String[])" 的 <summary> 文档]

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool Contains(
	this string str,
	string[] value
)
```

```
<ExtensionAttribute>
Public Shared Function Contains ( 
	str As String,
	value As String()
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool Contains(
	String^ str, 
	array<String^>^ value
)
```

```
[<ExtensionAttribute>]
static member Contains : 
        str : string * 
        value : string[] -> bool 
```

#### 参数

str
:   类型：SystemString  

    [缺少 "M:HslCommunication.HslExtension.Contains(System.String,System.String[])" 的 <param name="str"/> 文档]

value
:   类型：SystemString  

    [缺少 "M:HslCommunication.HslExtension.Contains(System.String,System.String[])" 的 <param name="value"/> 文档]

#### 返回值

类型：Boolean  

[缺少 "M:HslCommunication.HslExtension.Contains(System.String,System.String[])" 的 <returns> 文档]

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 String 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CopyArray(T) 方法 

[原文連結](http://api.hslcommunication.cn/html/8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[Contains 方法](../html/449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm "Contains 方法 ")

[CopyArray(T) 方法](../html/8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm "CopyArray(T) 方法 ")

[DecryptLargeData 方法](../html/bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm "DecryptLargeData 方法 ")

[EncryptLargeData 方法](../html/120131e2-a737-5340-5b08-18e4632d6055.htm "EncryptLargeData 方法 ")

[EndReceiveResult 方法](../html/190a250b-5fc7-2618-7255-a65bc6a328c8.htm "EndReceiveResult 方法 ")

[EndsWith 方法](../html/0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm "EndsWith 方法 ")

[EveryByteAdd 方法](../html/000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm "EveryByteAdd 方法 ")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolValue 方法](../html/dd581622-6e71-36a6-d734-dbfa30311d78.htm "GetBoolValue 方法 ")

[GetBytes 方法](../html/adaf490f-7569-21af-cb68-70e8edb67449.htm "GetBytes 方法 ")

[GetPEMPrivateKey 方法](../html/db6c7187-5d63-6d79-dc01-136bc3158f22.htm "GetPEMPrivateKey 方法 ")

[GetPEMPublicKey 方法](../html/567513da-22ad-4bd0-3d8a-b664b3f46109.htm "GetPEMPublicKey 方法 ")

[GetStringOrEndChar 方法](../html/aebf584d-117e-7e1d-d668-af3e6855a357.htm "GetStringOrEndChar 方法 ")

[GetValueOrDefault(T) 方法](../html/3ebfcbba-df43-e997-10c0-316d68ee8534.htm "GetValueOrDefault(T) 方法 ")

[IncreaseBy(T) 方法](../html/bc6933f2-a930-5b8a-5226-b01cc8827f15.htm "IncreaseBy(T) 方法 ")

[IniSerialByFormatString 方法](../html/3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm "IniSerialByFormatString 方法 ")

[RemoveBegin(T) 方法](../html/23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm "RemoveBegin(T) 方法 ")

[RemoveDouble(T) 方法](../html/1dba1036-8409-b99b-692d-af13e1f370ee.htm "RemoveDouble(T) 方法 ")

[RemoveLast 方法](../html/15db4c58-5c86-79a5-2098-146c39aa3436.htm "RemoveLast 方法 ")

[ReverseByWord 方法](../html/e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm "ReverseByWord 方法 ")

[ReverseNew(T) 方法](../html/20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm "ReverseNew(T) 方法 ")

[SelectBegin(T) 方法](../html/d92e877f-363f-a5d0-615d-851d06e867a5.htm "SelectBegin(T) 方法 ")

[SelectLast(T) 方法](../html/e7b564ce-2e61-37f4-a41a-77a82340736a.htm "SelectLast(T) 方法 ")

[SelectMiddle(T) 方法](../html/89b90303-188b-4924-08e0-a214f1f51eea.htm "SelectMiddle(T) 方法 ")

[SetBoolByIndex 方法](../html/92627e54-07a1-ca29-1a63-c1cd2b28b786.htm "SetBoolByIndex 方法 ")

[SetKeepAlive 方法](../html/8f2949d7-33fd-6175-3d11-50e33acc9d15.htm "SetKeepAlive 方法 ")

[SpliceArray(T) 方法](../html/c34cc020-1531-9910-6143-8ea611240268.htm "SpliceArray(T) 方法 ")

[SplitDot 方法](../html/b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm "SplitDot 方法 ")

[StartsWith 方法](../html/b6f380b3-fc65-1e8f-aa86-213d03447970.htm "StartsWith 方法 ")

[StartsWithAndNumber 方法](../html/f71b5bea-6b5e-89c7-e922-c7ad06cd4cb9.htm "StartsWithAndNumber 方法 ")

[ToArrayString 方法](../html/e7ca844f-e63d-1204-9832-a7f8f43ece36.htm "ToArrayString 方法 ")

[ToBoolArray 方法](../html/1bcba304-f35a-23cf-d17f-f539b806e433.htm "ToBoolArray 方法 ")

[ToByteArray 方法](../html/db7298ee-87e8-df3a-59fa-5892d82016d5.htm "ToByteArray 方法 ")

[ToFormatString 方法](../html/fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm "ToFormatString 方法 ")

[ToHexBytes 方法](../html/ff506e59-b0c9-918d-06f6-1e175bbda013.htm "ToHexBytes 方法 ")

[ToHexString 方法](../html/563d3090-af73-0aaf-aedc-047bd1dedb00.htm "ToHexString 方法 ")

[ToJsonString 方法](../html/828b0b70-5381-936e-e606-e6c0a905f688.htm "ToJsonString 方法 ")

[ToStringArray 方法](../html/58230997-75a3-0a6b-bd9a-695a51820891.htm "ToStringArray 方法 ")

[Write 方法](../html/ea8cba73-29a6-0c12-de93-2960758fc8d7.htm "Write 方法 ")

[WriteReverse 方法](../html/e51a3e83-1552-09c9-f01d-015fb5d6c226.htm "WriteReverse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionCopyArrayT 方法 |

拷贝当前的实例数组，是基于引用层的浅拷贝，如果类型为值类型，那就是深度拷贝，如果类型为引用类型，就是浅拷贝

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static T[] CopyArray<T>(
	this T[] value
)
```

```
<ExtensionAttribute>
Public Shared Function CopyArray(Of T) ( 
	value As T()
) As T()
```

```
public:
[ExtensionAttribute]
generic<typename T>
static array<T>^ CopyArray(
	array<T>^ value
)
```

```
[<ExtensionAttribute>]
static member CopyArray : 
        value : 'T[] -> 'T[] 
```

#### 参数

value
:   类型：T  
    数组对象

#### 类型参数

T
:   类型对象

#### 返回值

类型：T  
拷贝的结果内容

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecryptLargeData 方法 

[原文連結](http://api.hslcommunication.cn/html/bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[Contains 方法](../html/449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm "Contains 方法 ")

[CopyArray(T) 方法](../html/8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm "CopyArray(T) 方法 ")

[DecryptLargeData 方法](../html/bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm "DecryptLargeData 方法 ")

[EncryptLargeData 方法](../html/120131e2-a737-5340-5b08-18e4632d6055.htm "EncryptLargeData 方法 ")

[EndReceiveResult 方法](../html/190a250b-5fc7-2618-7255-a65bc6a328c8.htm "EndReceiveResult 方法 ")

[EndsWith 方法](../html/0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm "EndsWith 方法 ")

[EveryByteAdd 方法](../html/000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm "EveryByteAdd 方法 ")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolValue 方法](../html/dd581622-6e71-36a6-d734-dbfa30311d78.htm "GetBoolValue 方法 ")

[GetBytes 方法](../html/adaf490f-7569-21af-cb68-70e8edb67449.htm "GetBytes 方法 ")

[GetPEMPrivateKey 方法](../html/db6c7187-5d63-6d79-dc01-136bc3158f22.htm "GetPEMPrivateKey 方法 ")

[GetPEMPublicKey 方法](../html/567513da-22ad-4bd0-3d8a-b664b3f46109.htm "GetPEMPublicKey 方法 ")

[GetStringOrEndChar 方法](../html/aebf584d-117e-7e1d-d668-af3e6855a357.htm "GetStringOrEndChar 方法 ")

[GetValueOrDefault(T) 方法](../html/3ebfcbba-df43-e997-10c0-316d68ee8534.htm "GetValueOrDefault(T) 方法 ")

[IncreaseBy(T) 方法](../html/bc6933f2-a930-5b8a-5226-b01cc8827f15.htm "IncreaseBy(T) 方法 ")

[IniSerialByFormatString 方法](../html/3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm "IniSerialByFormatString 方法 ")

[RemoveBegin(T) 方法](../html/23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm "RemoveBegin(T) 方法 ")

[RemoveDouble(T) 方法](../html/1dba1036-8409-b99b-692d-af13e1f370ee.htm "RemoveDouble(T) 方法 ")

[RemoveLast 方法](../html/15db4c58-5c86-79a5-2098-146c39aa3436.htm "RemoveLast 方法 ")

[ReverseByWord 方法](../html/e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm "ReverseByWord 方法 ")

[ReverseNew(T) 方法](../html/20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm "ReverseNew(T) 方法 ")

[SelectBegin(T) 方法](../html/d92e877f-363f-a5d0-615d-851d06e867a5.htm "SelectBegin(T) 方法 ")

[SelectLast(T) 方法](../html/e7b564ce-2e61-37f4-a41a-77a82340736a.htm "SelectLast(T) 方法 ")

[SelectMiddle(T) 方法](../html/89b90303-188b-4924-08e0-a214f1f51eea.htm "SelectMiddle(T) 方法 ")

[SetBoolByIndex 方法](../html/92627e54-07a1-ca29-1a63-c1cd2b28b786.htm "SetBoolByIndex 方法 ")

[SetKeepAlive 方法](../html/8f2949d7-33fd-6175-3d11-50e33acc9d15.htm "SetKeepAlive 方法 ")

[SpliceArray(T) 方法](../html/c34cc020-1531-9910-6143-8ea611240268.htm "SpliceArray(T) 方法 ")

[SplitDot 方法](../html/b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm "SplitDot 方法 ")

[StartsWith 方法](../html/b6f380b3-fc65-1e8f-aa86-213d03447970.htm "StartsWith 方法 ")

[StartsWithAndNumber 方法](../html/f71b5bea-6b5e-89c7-e922-c7ad06cd4cb9.htm "StartsWithAndNumber 方法 ")

[ToArrayString 方法](../html/e7ca844f-e63d-1204-9832-a7f8f43ece36.htm "ToArrayString 方法 ")

[ToBoolArray 方法](../html/1bcba304-f35a-23cf-d17f-f539b806e433.htm "ToBoolArray 方法 ")

[ToByteArray 方法](../html/db7298ee-87e8-df3a-59fa-5892d82016d5.htm "ToByteArray 方法 ")

[ToFormatString 方法](../html/fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm "ToFormatString 方法 ")

[ToHexBytes 方法](../html/ff506e59-b0c9-918d-06f6-1e175bbda013.htm "ToHexBytes 方法 ")

[ToHexString 方法](../html/563d3090-af73-0aaf-aedc-047bd1dedb00.htm "ToHexString 方法 ")

[ToJsonString 方法](../html/828b0b70-5381-936e-e606-e6c0a905f688.htm "ToJsonString 方法 ")

[ToStringArray 方法](../html/58230997-75a3-0a6b-bd9a-695a51820891.htm "ToStringArray 方法 ")

[Write 方法](../html/ea8cba73-29a6-0c12-de93-2960758fc8d7.htm "Write 方法 ")

[WriteReverse 方法](../html/e51a3e83-1552-09c9-f01d-015fb5d6c226.htm "WriteReverse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionDecryptLargeData 方法 |

对超过117字节限制的加密数据进行解密，解密出原始的字节数据，因为RSA本身限制了117字节，所以此处进行数据切割解密。  
Decrypt the encrypted data that exceeds the 117-byte limit and decrypt the original byte data, because RSA itself limits 117 bytes, so the data is cut and decrypted here.

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] DecryptLargeData(
	this RSACryptoServiceProvider rsa,
	byte[] data
)
```

```
<ExtensionAttribute>
Public Shared Function DecryptLargeData ( 
	rsa As RSACryptoServiceProvider,
	data As Byte()
) As Byte()
```

```
public:
[ExtensionAttribute]
static array<unsigned char>^ DecryptLargeData(
	RSACryptoServiceProvider^ rsa, 
	array<unsigned char>^ data
)
```

```
[<ExtensionAttribute>]
static member DecryptLargeData : 
        rsa : RSACryptoServiceProvider * 
        data : byte[] -> byte[] 
```

#### 参数

rsa
:   类型：System.Security.CryptographyRSACryptoServiceProvider  

    [缺少 "M:HslCommunication.HslExtension.DecryptLargeData(System.Security.Cryptography.RSACryptoServiceProvider,System.Byte[])" 的 <param name="rsa"/> 文档]

data
:   类型：SystemByte  
    等待解密的数据

#### 返回值

类型：Byte  
解密之后的结果数据

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 RSACryptoServiceProvider 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EncryptLargeData 方法 

[原文連結](http://api.hslcommunication.cn/html/120131e2-a737-5340-5b08-18e4632d6055.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[Contains 方法](../html/449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm "Contains 方法 ")

[CopyArray(T) 方法](../html/8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm "CopyArray(T) 方法 ")

[DecryptLargeData 方法](../html/bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm "DecryptLargeData 方法 ")

[EncryptLargeData 方法](../html/120131e2-a737-5340-5b08-18e4632d6055.htm "EncryptLargeData 方法 ")

[EndReceiveResult 方法](../html/190a250b-5fc7-2618-7255-a65bc6a328c8.htm "EndReceiveResult 方法 ")

[EndsWith 方法](../html/0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm "EndsWith 方法 ")

[EveryByteAdd 方法](../html/000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm "EveryByteAdd 方法 ")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolValue 方法](../html/dd581622-6e71-36a6-d734-dbfa30311d78.htm "GetBoolValue 方法 ")

[GetBytes 方法](../html/adaf490f-7569-21af-cb68-70e8edb67449.htm "GetBytes 方法 ")

[GetPEMPrivateKey 方法](../html/db6c7187-5d63-6d79-dc01-136bc3158f22.htm "GetPEMPrivateKey 方法 ")

[GetPEMPublicKey 方法](../html/567513da-22ad-4bd0-3d8a-b664b3f46109.htm "GetPEMPublicKey 方法 ")

[GetStringOrEndChar 方法](../html/aebf584d-117e-7e1d-d668-af3e6855a357.htm "GetStringOrEndChar 方法 ")

[GetValueOrDefault(T) 方法](../html/3ebfcbba-df43-e997-10c0-316d68ee8534.htm "GetValueOrDefault(T) 方法 ")

[IncreaseBy(T) 方法](../html/bc6933f2-a930-5b8a-5226-b01cc8827f15.htm "IncreaseBy(T) 方法 ")

[IniSerialByFormatString 方法](../html/3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm "IniSerialByFormatString 方法 ")

[RemoveBegin(T) 方法](../html/23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm "RemoveBegin(T) 方法 ")

[RemoveDouble(T) 方法](../html/1dba1036-8409-b99b-692d-af13e1f370ee.htm "RemoveDouble(T) 方法 ")

[RemoveLast 方法](../html/15db4c58-5c86-79a5-2098-146c39aa3436.htm "RemoveLast 方法 ")

[ReverseByWord 方法](../html/e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm "ReverseByWord 方法 ")

[ReverseNew(T) 方法](../html/20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm "ReverseNew(T) 方法 ")

[SelectBegin(T) 方法](../html/d92e877f-363f-a5d0-615d-851d06e867a5.htm "SelectBegin(T) 方法 ")

[SelectLast(T) 方法](../html/e7b564ce-2e61-37f4-a41a-77a82340736a.htm "SelectLast(T) 方法 ")

[SelectMiddle(T) 方法](../html/89b90303-188b-4924-08e0-a214f1f51eea.htm "SelectMiddle(T) 方法 ")

[SetBoolByIndex 方法](../html/92627e54-07a1-ca29-1a63-c1cd2b28b786.htm "SetBoolByIndex 方法 ")

[SetKeepAlive 方法](../html/8f2949d7-33fd-6175-3d11-50e33acc9d15.htm "SetKeepAlive 方法 ")

[SpliceArray(T) 方法](../html/c34cc020-1531-9910-6143-8ea611240268.htm "SpliceArray(T) 方法 ")

[SplitDot 方法](../html/b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm "SplitDot 方法 ")

[StartsWith 方法](../html/b6f380b3-fc65-1e8f-aa86-213d03447970.htm "StartsWith 方法 ")

[StartsWithAndNumber 方法](../html/f71b5bea-6b5e-89c7-e922-c7ad06cd4cb9.htm "StartsWithAndNumber 方法 ")

[ToArrayString 方法](../html/e7ca844f-e63d-1204-9832-a7f8f43ece36.htm "ToArrayString 方法 ")

[ToBoolArray 方法](../html/1bcba304-f35a-23cf-d17f-f539b806e433.htm "ToBoolArray 方法 ")

[ToByteArray 方法](../html/db7298ee-87e8-df3a-59fa-5892d82016d5.htm "ToByteArray 方法 ")

[ToFormatString 方法](../html/fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm "ToFormatString 方法 ")

[ToHexBytes 方法](../html/ff506e59-b0c9-918d-06f6-1e175bbda013.htm "ToHexBytes 方法 ")

[ToHexString 方法](../html/563d3090-af73-0aaf-aedc-047bd1dedb00.htm "ToHexString 方法 ")

[ToJsonString 方法](../html/828b0b70-5381-936e-e606-e6c0a905f688.htm "ToJsonString 方法 ")

[ToStringArray 方法](../html/58230997-75a3-0a6b-bd9a-695a51820891.htm "ToStringArray 方法 ")

[Write 方法](../html/ea8cba73-29a6-0c12-de93-2960758fc8d7.htm "Write 方法 ")

[WriteReverse 方法](../html/e51a3e83-1552-09c9-f01d-015fb5d6c226.htm "WriteReverse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionEncryptLargeData 方法 |

对原始字节的数据进行加密，不限制长度，因为RSA本身限制了117字节，所以此处进行数据切割加密。  
Encrypt the original byte data without limiting the length, because RSA itself limits 117 bytes, so the data is cut and encrypted here.

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] EncryptLargeData(
	this RSACryptoServiceProvider rsa,
	byte[] data
)
```

```
<ExtensionAttribute>
Public Shared Function EncryptLargeData ( 
	rsa As RSACryptoServiceProvider,
	data As Byte()
) As Byte()
```

```
public:
[ExtensionAttribute]
static array<unsigned char>^ EncryptLargeData(
	RSACryptoServiceProvider^ rsa, 
	array<unsigned char>^ data
)
```

```
[<ExtensionAttribute>]
static member EncryptLargeData : 
        rsa : RSACryptoServiceProvider * 
        data : byte[] -> byte[] 
```

#### 参数

rsa
:   类型：System.Security.CryptographyRSACryptoServiceProvider  

    [缺少 "M:HslCommunication.HslExtension.EncryptLargeData(System.Security.Cryptography.RSACryptoServiceProvider,System.Byte[])" 的 <param name="rsa"/> 文档]

data
:   类型：SystemByte  
    等待加密的原始数据

#### 返回值

类型：Byte  
加密之后的结果信息

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 RSACryptoServiceProvider 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EndReceiveResult 方法 

[原文連結](http://api.hslcommunication.cn/html/190a250b-5fc7-2618-7255-a65bc6a328c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[Contains 方法](../html/449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm "Contains 方法 ")

[CopyArray(T) 方法](../html/8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm "CopyArray(T) 方法 ")

[DecryptLargeData 方法](../html/bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm "DecryptLargeData 方法 ")

[EncryptLargeData 方法](../html/120131e2-a737-5340-5b08-18e4632d6055.htm "EncryptLargeData 方法 ")

[EndReceiveResult 方法](../html/190a250b-5fc7-2618-7255-a65bc6a328c8.htm "EndReceiveResult 方法 ")

[EndsWith 方法](../html/0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm "EndsWith 方法 ")

[EveryByteAdd 方法](../html/000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm "EveryByteAdd 方法 ")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolValue 方法](../html/dd581622-6e71-36a6-d734-dbfa30311d78.htm "GetBoolValue 方法 ")

[GetBytes 方法](../html/adaf490f-7569-21af-cb68-70e8edb67449.htm "GetBytes 方法 ")

[GetPEMPrivateKey 方法](../html/db6c7187-5d63-6d79-dc01-136bc3158f22.htm "GetPEMPrivateKey 方法 ")

[GetPEMPublicKey 方法](../html/567513da-22ad-4bd0-3d8a-b664b3f46109.htm "GetPEMPublicKey 方法 ")

[GetStringOrEndChar 方法](../html/aebf584d-117e-7e1d-d668-af3e6855a357.htm "GetStringOrEndChar 方法 ")

[GetValueOrDefault(T) 方法](../html/3ebfcbba-df43-e997-10c0-316d68ee8534.htm "GetValueOrDefault(T) 方法 ")

[IncreaseBy(T) 方法](../html/bc6933f2-a930-5b8a-5226-b01cc8827f15.htm "IncreaseBy(T) 方法 ")

[IniSerialByFormatString 方法](../html/3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm "IniSerialByFormatString 方法 ")

[RemoveBegin(T) 方法](../html/23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm "RemoveBegin(T) 方法 ")

[RemoveDouble(T) 方法](../html/1dba1036-8409-b99b-692d-af13e1f370ee.htm "RemoveDouble(T) 方法 ")

[RemoveLast 方法](../html/15db4c58-5c86-79a5-2098-146c39aa3436.htm "RemoveLast 方法 ")

[ReverseByWord 方法](../html/e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm "ReverseByWord 方法 ")

[ReverseNew(T) 方法](../html/20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm "ReverseNew(T) 方法 ")

[SelectBegin(T) 方法](../html/d92e877f-363f-a5d0-615d-851d06e867a5.htm "SelectBegin(T) 方法 ")

[SelectLast(T) 方法](../html/e7b564ce-2e61-37f4-a41a-77a82340736a.htm "SelectLast(T) 方法 ")

[SelectMiddle(T) 方法](../html/89b90303-188b-4924-08e0-a214f1f51eea.htm "SelectMiddle(T) 方法 ")

[SetBoolByIndex 方法](../html/92627e54-07a1-ca29-1a63-c1cd2b28b786.htm "SetBoolByIndex 方法 ")

[SetKeepAlive 方法](../html/8f2949d7-33fd-6175-3d11-50e33acc9d15.htm "SetKeepAlive 方法 ")

[SpliceArray(T) 方法](../html/c34cc020-1531-9910-6143-8ea611240268.htm "SpliceArray(T) 方法 ")

[SplitDot 方法](../html/b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm "SplitDot 方法 ")

[StartsWith 方法](../html/b6f380b3-fc65-1e8f-aa86-213d03447970.htm "StartsWith 方法 ")

[StartsWithAndNumber 方法](../html/f71b5bea-6b5e-89c7-e922-c7ad06cd4cb9.htm "StartsWithAndNumber 方法 ")

[ToArrayString 方法](../html/e7ca844f-e63d-1204-9832-a7f8f43ece36.htm "ToArrayString 方法 ")

[ToBoolArray 方法](../html/1bcba304-f35a-23cf-d17f-f539b806e433.htm "ToBoolArray 方法 ")

[ToByteArray 方法](../html/db7298ee-87e8-df3a-59fa-5892d82016d5.htm "ToByteArray 方法 ")

[ToFormatString 方法](../html/fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm "ToFormatString 方法 ")

[ToHexBytes 方法](../html/ff506e59-b0c9-918d-06f6-1e175bbda013.htm "ToHexBytes 方法 ")

[ToHexString 方法](../html/563d3090-af73-0aaf-aedc-047bd1dedb00.htm "ToHexString 方法 ")

[ToJsonString 方法](../html/828b0b70-5381-936e-e606-e6c0a905f688.htm "ToJsonString 方法 ")

[ToStringArray 方法](../html/58230997-75a3-0a6b-bd9a-695a51820891.htm "ToStringArray 方法 ")

[Write 方法](../html/ea8cba73-29a6-0c12-de93-2960758fc8d7.htm "Write 方法 ")

[WriteReverse 方法](../html/e51a3e83-1552-09c9-f01d-015fb5d6c226.htm "WriteReverse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionEndReceiveResult 方法 |

结束挂起的异步读取，返回读取的字节数，如果成功的情况。  
Ends the pending asynchronous read and returns the number of bytes read, if successful.

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<int> EndReceiveResult(
	this Socket socket,
	IAsyncResult ar
)
```

```
<ExtensionAttribute>
Public Shared Function EndReceiveResult ( 
	socket As Socket,
	ar As IAsyncResult
) As OperateResult(Of Integer)
```

```
public:
[ExtensionAttribute]
static OperateResult<int>^ EndReceiveResult(
	Socket^ socket, 
	IAsyncResult^ ar
)
```

```
[<ExtensionAttribute>]
static member EndReceiveResult : 
        socket : Socket * 
        ar : IAsyncResult -> OperateResult<int> 
```

#### 参数

socket
:   类型：System.Net.SocketsSocket  
    socket对象

ar
:   类型：SystemIAsyncResult  
    回调方法

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int32  
是否启动成功

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 Socket 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EndsWith 方法 

[原文連結](http://api.hslcommunication.cn/html/0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[Contains 方法](../html/449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm "Contains 方法 ")

[CopyArray(T) 方法](../html/8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm "CopyArray(T) 方法 ")

[DecryptLargeData 方法](../html/bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm "DecryptLargeData 方法 ")

[EncryptLargeData 方法](../html/120131e2-a737-5340-5b08-18e4632d6055.htm "EncryptLargeData 方法 ")

[EndReceiveResult 方法](../html/190a250b-5fc7-2618-7255-a65bc6a328c8.htm "EndReceiveResult 方法 ")

[EndsWith 方法](../html/0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm "EndsWith 方法 ")

[EveryByteAdd 方法](../html/000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm "EveryByteAdd 方法 ")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolValue 方法](../html/dd581622-6e71-36a6-d734-dbfa30311d78.htm "GetBoolValue 方法 ")

[GetBytes 方法](../html/adaf490f-7569-21af-cb68-70e8edb67449.htm "GetBytes 方法 ")

[GetPEMPrivateKey 方法](../html/db6c7187-5d63-6d79-dc01-136bc3158f22.htm "GetPEMPrivateKey 方法 ")

[GetPEMPublicKey 方法](../html/567513da-22ad-4bd0-3d8a-b664b3f46109.htm "GetPEMPublicKey 方法 ")

[GetStringOrEndChar 方法](../html/aebf584d-117e-7e1d-d668-af3e6855a357.htm "GetStringOrEndChar 方法 ")

[GetValueOrDefault(T) 方法](../html/3ebfcbba-df43-e997-10c0-316d68ee8534.htm "GetValueOrDefault(T) 方法 ")

[IncreaseBy(T) 方法](../html/bc6933f2-a930-5b8a-5226-b01cc8827f15.htm "IncreaseBy(T) 方法 ")

[IniSerialByFormatString 方法](../html/3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm "IniSerialByFormatString 方法 ")

[RemoveBegin(T) 方法](../html/23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm "RemoveBegin(T) 方法 ")

[RemoveDouble(T) 方法](../html/1dba1036-8409-b99b-692d-af13e1f370ee.htm "RemoveDouble(T) 方法 ")

[RemoveLast 方法](../html/15db4c58-5c86-79a5-2098-146c39aa3436.htm "RemoveLast 方法 ")

[ReverseByWord 方法](../html/e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm "ReverseByWord 方法 ")

[ReverseNew(T) 方法](../html/20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm "ReverseNew(T) 方法 ")

[SelectBegin(T) 方法](../html/d92e877f-363f-a5d0-615d-851d06e867a5.htm "SelectBegin(T) 方法 ")

[SelectLast(T) 方法](../html/e7b564ce-2e61-37f4-a41a-77a82340736a.htm "SelectLast(T) 方法 ")

[SelectMiddle(T) 方法](../html/89b90303-188b-4924-08e0-a214f1f51eea.htm "SelectMiddle(T) 方法 ")

[SetBoolByIndex 方法](../html/92627e54-07a1-ca29-1a63-c1cd2b28b786.htm "SetBoolByIndex 方法 ")

[SetKeepAlive 方法](../html/8f2949d7-33fd-6175-3d11-50e33acc9d15.htm "SetKeepAlive 方法 ")

[SpliceArray(T) 方法](../html/c34cc020-1531-9910-6143-8ea611240268.htm "SpliceArray(T) 方法 ")

[SplitDot 方法](../html/b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm "SplitDot 方法 ")

[StartsWith 方法](../html/b6f380b3-fc65-1e8f-aa86-213d03447970.htm "StartsWith 方法 ")

[StartsWithAndNumber 方法](../html/f71b5bea-6b5e-89c7-e922-c7ad06cd4cb9.htm "StartsWithAndNumber 方法 ")

[ToArrayString 方法](../html/e7ca844f-e63d-1204-9832-a7f8f43ece36.htm "ToArrayString 方法 ")

[ToBoolArray 方法](../html/1bcba304-f35a-23cf-d17f-f539b806e433.htm "ToBoolArray 方法 ")

[ToByteArray 方法](../html/db7298ee-87e8-df3a-59fa-5892d82016d5.htm "ToByteArray 方法 ")

[ToFormatString 方法](../html/fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm "ToFormatString 方法 ")

[ToHexBytes 方法](../html/ff506e59-b0c9-918d-06f6-1e175bbda013.htm "ToHexBytes 方法 ")

[ToHexString 方法](../html/563d3090-af73-0aaf-aedc-047bd1dedb00.htm "ToHexString 方法 ")

[ToJsonString 方法](../html/828b0b70-5381-936e-e606-e6c0a905f688.htm "ToJsonString 方法 ")

[ToStringArray 方法](../html/58230997-75a3-0a6b-bd9a-695a51820891.htm "ToStringArray 方法 ")

[Write 方法](../html/ea8cba73-29a6-0c12-de93-2960758fc8d7.htm "Write 方法 ")

[WriteReverse 方法](../html/e51a3e83-1552-09c9-f01d-015fb5d6c226.htm "WriteReverse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionEndsWith 方法 |

[缺少 "M:HslCommunication.HslExtension.EndsWith(System.String,System.String[])" 的 <summary> 文档]

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool EndsWith(
	this string str,
	string[] value
)
```

```
<ExtensionAttribute>
Public Shared Function EndsWith ( 
	str As String,
	value As String()
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool EndsWith(
	String^ str, 
	array<String^>^ value
)
```

```
[<ExtensionAttribute>]
static member EndsWith : 
        str : string * 
        value : string[] -> bool 
```

#### 参数

str
:   类型：SystemString  

    [缺少 "M:HslCommunication.HslExtension.EndsWith(System.String,System.String[])" 的 <param name="str"/> 文档]

value
:   类型：SystemString  

    [缺少 "M:HslCommunication.HslExtension.EndsWith(System.String,System.String[])" 的 <param name="value"/> 文档]

#### 返回值

类型：Boolean  

[缺少 "M:HslCommunication.HslExtension.EndsWith(System.String,System.String[])" 的 <returns> 文档]

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 String 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EveryByteAdd 方法 

[原文連結](http://api.hslcommunication.cn/html/000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[BeginReceiveResult 方法](../html/cdcfdc10-5f34-00c6-142b-95265069e58a.htm "BeginReceiveResult 方法 ")

[Contains 方法](../html/449a0ed7-bd1e-0f20-2299-fa3cf82a27d9.htm "Contains 方法 ")

[CopyArray(T) 方法](../html/8c80bd38-e70d-df82-5e1d-b7cc2b06de80.htm "CopyArray(T) 方法 ")

[DecryptLargeData 方法](../html/bbdb6fa9-9b12-5bad-1b06-64c4131e37a1.htm "DecryptLargeData 方法 ")

[EncryptLargeData 方法](../html/120131e2-a737-5340-5b08-18e4632d6055.htm "EncryptLargeData 方法 ")

[EndReceiveResult 方法](../html/190a250b-5fc7-2618-7255-a65bc6a328c8.htm "EndReceiveResult 方法 ")

[EndsWith 方法](../html/0d983099-ace5-e47a-6b7f-f0625e16b7f2.htm "EndsWith 方法 ")

[EveryByteAdd 方法](../html/000dfa8f-4c1f-54f3-f0a6-bbe31adf2007.htm "EveryByteAdd 方法 ")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolValue 方法](../html/dd581622-6e71-36a6-d734-dbfa30311d78.htm "GetBoolValue 方法 ")

[GetBytes 方法](../html/adaf490f-7569-21af-cb68-70e8edb67449.htm "GetBytes 方法 ")

[GetPEMPrivateKey 方法](../html/db6c7187-5d63-6d79-dc01-136bc3158f22.htm "GetPEMPrivateKey 方法 ")

[GetPEMPublicKey 方法](../html/567513da-22ad-4bd0-3d8a-b664b3f46109.htm "GetPEMPublicKey 方法 ")

[GetStringOrEndChar 方法](../html/aebf584d-117e-7e1d-d668-af3e6855a357.htm "GetStringOrEndChar 方法 ")

[GetValueOrDefault(T) 方法](../html/3ebfcbba-df43-e997-10c0-316d68ee8534.htm "GetValueOrDefault(T) 方法 ")

[IncreaseBy(T) 方法](../html/bc6933f2-a930-5b8a-5226-b01cc8827f15.htm "IncreaseBy(T) 方法 ")

[IniSerialByFormatString 方法](../html/3ff7d3f8-7b87-dfff-bef0-33de6c3ffcf7.htm "IniSerialByFormatString 方法 ")

[RemoveBegin(T) 方法](../html/23c11ce3-1d2f-7c5e-f255-7e805d6fdef6.htm "RemoveBegin(T) 方法 ")

[RemoveDouble(T) 方法](../html/1dba1036-8409-b99b-692d-af13e1f370ee.htm "RemoveDouble(T) 方法 ")

[RemoveLast 方法](../html/15db4c58-5c86-79a5-2098-146c39aa3436.htm "RemoveLast 方法 ")

[ReverseByWord 方法](../html/e4cdaf02-f504-637e-4b17-ac30a6e1f82e.htm "ReverseByWord 方法 ")

[ReverseNew(T) 方法](../html/20bc1760-44d6-04d6-2c1e-e00a3bab3c9d.htm "ReverseNew(T) 方法 ")

[SelectBegin(T) 方法](../html/d92e877f-363f-a5d0-615d-851d06e867a5.htm "SelectBegin(T) 方法 ")

[SelectLast(T) 方法](../html/e7b564ce-2e61-37f4-a41a-77a82340736a.htm "SelectLast(T) 方法 ")

[SelectMiddle(T) 方法](../html/89b90303-188b-4924-08e0-a214f1f51eea.htm "SelectMiddle(T) 方法 ")

[SetBoolByIndex 方法](../html/92627e54-07a1-ca29-1a63-c1cd2b28b786.htm "SetBoolByIndex 方法 ")

[SetKeepAlive 方法](../html/8f2949d7-33fd-6175-3d11-50e33acc9d15.htm "SetKeepAlive 方法 ")

[SpliceArray(T) 方法](../html/c34cc020-1531-9910-6143-8ea611240268.htm "SpliceArray(T) 方法 ")

[SplitDot 方法](../html/b45ffd2b-defb-15b3-3c91-34ff61aa0c52.htm "SplitDot 方法 ")

[StartsWith 方法](../html/b6f380b3-fc65-1e8f-aa86-213d03447970.htm "StartsWith 方法 ")

[StartsWithAndNumber 方法](../html/f71b5bea-6b5e-89c7-e922-c7ad06cd4cb9.htm "StartsWithAndNumber 方法 ")

[ToArrayString 方法](../html/e7ca844f-e63d-1204-9832-a7f8f43ece36.htm "ToArrayString 方法 ")

[ToBoolArray 方法](../html/1bcba304-f35a-23cf-d17f-f539b806e433.htm "ToBoolArray 方法 ")

[ToByteArray 方法](../html/db7298ee-87e8-df3a-59fa-5892d82016d5.htm "ToByteArray 方法 ")

[ToFormatString 方法](../html/fdb4a3d1-e810-fdf9-a13a-854e977d34c8.htm "ToFormatString 方法 ")

[ToHexBytes 方法](../html/ff506e59-b0c9-918d-06f6-1e175bbda013.htm "ToHexBytes 方法 ")

[ToHexString 方法](../html/563d3090-af73-0aaf-aedc-047bd1dedb00.htm "ToHexString 方法 ")

[ToJsonString 方法](../html/828b0b70-5381-936e-e606-e6c0a905f688.htm "ToJsonString 方法 ")

[ToStringArray 方法](../html/58230997-75a3-0a6b-bd9a-695a51820891.htm "ToStringArray 方法 ")

[Write 方法](../html/ea8cba73-29a6-0c12-de93-2960758fc8d7.htm "Write 方法 ")

[WriteReverse 方法](../html/e51a3e83-1552-09c9-f01d-015fb5d6c226.htm "WriteReverse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionEveryByteAdd 方法 |

将指定的数据添加到数组的每个元素上去，会改变每个元素的值

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] EveryByteAdd(
	this byte[] array,
	int value
)
```

```
<ExtensionAttribute>
Public Shared Function EveryByteAdd ( 
	array As Byte(),
	value As Integer
) As Byte()
```

```
public:
[ExtensionAttribute]
static array<unsigned char>^ EveryByteAdd(
	array<unsigned char>^ array, 
	int value
)
```

```
[<ExtensionAttribute>]
static member EveryByteAdd : 
        array : byte[] * 
        value : int -> byte[] 
```

#### 参数

array
:   类型：SystemByte  
    原始数组

value
:   类型：SystemInt32  
    值

#### 返回值

类型：Byte  
修改后的数组信息

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolByIndex 方法 

[原文連結](http://api.hslcommunication.cn/html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolByIndex 方法 (Byte, Int32)](../html/02376ace-8a06-10a4-600d-90db9f773e92.htm "GetBoolByIndex 方法 (Byte, Int32)")

[GetBoolByIndex 方法 (Byte[], Int32)](../html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm "GetBoolByIndex 方法 (Byte[], Int32)")

[GetBoolByIndex 方法 (Int16, Int32)](../html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm "GetBoolByIndex 方法 (Int16, Int32)")

[GetBoolByIndex 方法 (Int32, Int32)](../html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm "GetBoolByIndex 方法 (Int32, Int32)")

[GetBoolByIndex 方法 (Int64, Int32)](../html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm "GetBoolByIndex 方法 (Int64, Int32)")

[GetBoolByIndex 方法 (UInt16, Int32)](../html/4419a620-899f-6d5a-201a-2cb66c38969c.htm "GetBoolByIndex 方法 (UInt16, Int32)")

[GetBoolByIndex 方法 (UInt32, Int32)](../html/19cac792-a632-936c-5bc3-052023e4e225.htm "GetBoolByIndex 方法 (UInt32, Int32)")

[GetBoolByIndex 方法 (UInt64, Int32)](../html/5f91b230-1822-e4be-02c4-640045232711.htm "GetBoolByIndex 方法 (UInt64, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionGetBoolByIndex 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [GetBoolByIndex(Byte, Int32)](02376ace-8a06-10a4-600d-90db9f773e92.htm) | 获取Byte的第 boolIndex 偏移的bool值，比如3，就是第4位   Get the bool value of Byte's boolIndex offset, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(Byte, Int32)](1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm) | 获取Byte数组的第 boolIndex 偏移的bool值，这个偏移值可以为 10，就是第 1 个字节的 第3位   Get the bool value of the boolIndex offset of the Byte array. The offset value can be 10, which is the third bit of the first byte |
| 公共方法静态成员 | [GetBoolByIndex(Int16, Int32)](013b9689-81f1-1871-bd53-2a6ca2ae153a.htm) | 获取short类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the short type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(Int32, Int32)](65142f2b-8313-8bf6-db4c-5567801c19f3.htm) | 获取int类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the int type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(Int64, Int32)](fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm) | 获取long类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the long type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(UInt16, Int32)](4419a620-899f-6d5a-201a-2cb66c38969c.htm) | 获取ushort类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the ushort type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(UInt32, Int32)](19cac792-a632-936c-5bc3-052023e4e225.htm) | 获取uint类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the uint type data, such as 3, which is the 4th bit |
| 公共方法静态成员 | [GetBoolByIndex(UInt64, Int32)](5f91b230-1822-e4be-02c4-640045232711.htm) | 获取ulong类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   Get the bool value of the boolIndex (starting from 0) offset of the ulong type data, such as 3, which is the 4th bit |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolByIndex 方法 (Byte, Int32)

[原文連結](http://api.hslcommunication.cn/html/02376ace-8a06-10a4-600d-90db9f773e92.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolByIndex 方法 (Byte, Int32)](../html/02376ace-8a06-10a4-600d-90db9f773e92.htm "GetBoolByIndex 方法 (Byte, Int32)")

[GetBoolByIndex 方法 (Byte[], Int32)](../html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm "GetBoolByIndex 方法 (Byte[], Int32)")

[GetBoolByIndex 方法 (Int16, Int32)](../html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm "GetBoolByIndex 方法 (Int16, Int32)")

[GetBoolByIndex 方法 (Int32, Int32)](../html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm "GetBoolByIndex 方法 (Int32, Int32)")

[GetBoolByIndex 方法 (Int64, Int32)](../html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm "GetBoolByIndex 方法 (Int64, Int32)")

[GetBoolByIndex 方法 (UInt16, Int32)](../html/4419a620-899f-6d5a-201a-2cb66c38969c.htm "GetBoolByIndex 方法 (UInt16, Int32)")

[GetBoolByIndex 方法 (UInt32, Int32)](../html/19cac792-a632-936c-5bc3-052023e4e225.htm "GetBoolByIndex 方法 (UInt32, Int32)")

[GetBoolByIndex 方法 (UInt64, Int32)](../html/5f91b230-1822-e4be-02c4-640045232711.htm "GetBoolByIndex 方法 (UInt64, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionGetBoolByIndex 方法 (Byte, Int32) |

获取Byte的第 boolIndex 偏移的bool值，比如3，就是第4位   
Get the bool value of Byte's boolIndex offset, such as 3, which is the 4th bit

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool GetBoolByIndex(
	this byte byt,
	int boolIndex
)
```

```
<ExtensionAttribute>
Public Shared Function GetBoolByIndex ( 
	byt As Byte,
	boolIndex As Integer
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool GetBoolByIndex(
	unsigned char byt, 
	int boolIndex
)
```

```
[<ExtensionAttribute>]
static member GetBoolByIndex : 
        byt : byte * 
        boolIndex : int -> bool 
```

#### 参数

byt
:   类型：SystemByte  
    字节信息

boolIndex
:   类型：SystemInt32  
    指定字节的位偏移

#### 返回值

类型：Boolean  
bool值

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 Byte 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[GetBoolByIndex 重载](d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolByIndex 方法 (Byte[], Int32)

[原文連結](http://api.hslcommunication.cn/html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolByIndex 方法 (Byte, Int32)](../html/02376ace-8a06-10a4-600d-90db9f773e92.htm "GetBoolByIndex 方法 (Byte, Int32)")

[GetBoolByIndex 方法 (Byte[], Int32)](../html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm "GetBoolByIndex 方法 (Byte[], Int32)")

[GetBoolByIndex 方法 (Int16, Int32)](../html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm "GetBoolByIndex 方法 (Int16, Int32)")

[GetBoolByIndex 方法 (Int32, Int32)](../html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm "GetBoolByIndex 方法 (Int32, Int32)")

[GetBoolByIndex 方法 (Int64, Int32)](../html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm "GetBoolByIndex 方法 (Int64, Int32)")

[GetBoolByIndex 方法 (UInt16, Int32)](../html/4419a620-899f-6d5a-201a-2cb66c38969c.htm "GetBoolByIndex 方法 (UInt16, Int32)")

[GetBoolByIndex 方法 (UInt32, Int32)](../html/19cac792-a632-936c-5bc3-052023e4e225.htm "GetBoolByIndex 方法 (UInt32, Int32)")

[GetBoolByIndex 方法 (UInt64, Int32)](../html/5f91b230-1822-e4be-02c4-640045232711.htm "GetBoolByIndex 方法 (UInt64, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionGetBoolByIndex 方法 (Byte, Int32) |

获取Byte数组的第 boolIndex 偏移的bool值，这个偏移值可以为 10，就是第 1 个字节的 第3位   
Get the bool value of the boolIndex offset of the Byte array. The offset value can be 10, which is the third bit of the first byte

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool GetBoolByIndex(
	this byte[] bytes,
	int boolIndex
)
```

```
<ExtensionAttribute>
Public Shared Function GetBoolByIndex ( 
	bytes As Byte(),
	boolIndex As Integer
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool GetBoolByIndex(
	array<unsigned char>^ bytes, 
	int boolIndex
)
```

```
[<ExtensionAttribute>]
static member GetBoolByIndex : 
        bytes : byte[] * 
        boolIndex : int -> bool 
```

#### 参数

bytes
:   类型：SystemByte  
    字节数组信息

boolIndex
:   类型：SystemInt32  
    指定字节的位偏移

#### 返回值

类型：Boolean  
bool值

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[GetBoolByIndex 重载](d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolByIndex 方法 (Int16, Int32)

[原文連結](http://api.hslcommunication.cn/html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolByIndex 方法 (Byte, Int32)](../html/02376ace-8a06-10a4-600d-90db9f773e92.htm "GetBoolByIndex 方法 (Byte, Int32)")

[GetBoolByIndex 方法 (Byte[], Int32)](../html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm "GetBoolByIndex 方法 (Byte[], Int32)")

[GetBoolByIndex 方法 (Int16, Int32)](../html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm "GetBoolByIndex 方法 (Int16, Int32)")

[GetBoolByIndex 方法 (Int32, Int32)](../html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm "GetBoolByIndex 方法 (Int32, Int32)")

[GetBoolByIndex 方法 (Int64, Int32)](../html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm "GetBoolByIndex 方法 (Int64, Int32)")

[GetBoolByIndex 方法 (UInt16, Int32)](../html/4419a620-899f-6d5a-201a-2cb66c38969c.htm "GetBoolByIndex 方法 (UInt16, Int32)")

[GetBoolByIndex 方法 (UInt32, Int32)](../html/19cac792-a632-936c-5bc3-052023e4e225.htm "GetBoolByIndex 方法 (UInt32, Int32)")

[GetBoolByIndex 方法 (UInt64, Int32)](../html/5f91b230-1822-e4be-02c4-640045232711.htm "GetBoolByIndex 方法 (UInt64, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionGetBoolByIndex 方法 (Int16, Int32) |

获取short类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   
Get the bool value of the boolIndex (starting from 0) offset of the short type data, such as 3, which is the 4th bit

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool GetBoolByIndex(
	this short value,
	int boolIndex
)
```

```
<ExtensionAttribute>
Public Shared Function GetBoolByIndex ( 
	value As Short,
	boolIndex As Integer
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool GetBoolByIndex(
	short value, 
	int boolIndex
)
```

```
[<ExtensionAttribute>]
static member GetBoolByIndex : 
        value : int16 * 
        boolIndex : int -> bool 
```

#### 参数

value
:   类型：SystemInt16  
    short数据值

boolIndex
:   类型：SystemInt32  
    位偏移索引，从0开始

#### 返回值

类型：Boolean  
bool值

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 Int16 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[GetBoolByIndex 重载](d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolByIndex 方法 (Int32, Int32)

[原文連結](http://api.hslcommunication.cn/html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolByIndex 方法 (Byte, Int32)](../html/02376ace-8a06-10a4-600d-90db9f773e92.htm "GetBoolByIndex 方法 (Byte, Int32)")

[GetBoolByIndex 方法 (Byte[], Int32)](../html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm "GetBoolByIndex 方法 (Byte[], Int32)")

[GetBoolByIndex 方法 (Int16, Int32)](../html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm "GetBoolByIndex 方法 (Int16, Int32)")

[GetBoolByIndex 方法 (Int32, Int32)](../html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm "GetBoolByIndex 方法 (Int32, Int32)")

[GetBoolByIndex 方法 (Int64, Int32)](../html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm "GetBoolByIndex 方法 (Int64, Int32)")

[GetBoolByIndex 方法 (UInt16, Int32)](../html/4419a620-899f-6d5a-201a-2cb66c38969c.htm "GetBoolByIndex 方法 (UInt16, Int32)")

[GetBoolByIndex 方法 (UInt32, Int32)](../html/19cac792-a632-936c-5bc3-052023e4e225.htm "GetBoolByIndex 方法 (UInt32, Int32)")

[GetBoolByIndex 方法 (UInt64, Int32)](../html/5f91b230-1822-e4be-02c4-640045232711.htm "GetBoolByIndex 方法 (UInt64, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionGetBoolByIndex 方法 (Int32, Int32) |

获取int类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   
Get the bool value of the boolIndex (starting from 0) offset of the int type data, such as 3, which is the 4th bit

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool GetBoolByIndex(
	this int value,
	int boolIndex
)
```

```
<ExtensionAttribute>
Public Shared Function GetBoolByIndex ( 
	value As Integer,
	boolIndex As Integer
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool GetBoolByIndex(
	int value, 
	int boolIndex
)
```

```
[<ExtensionAttribute>]
static member GetBoolByIndex : 
        value : int * 
        boolIndex : int -> bool 
```

#### 参数

value
:   类型：SystemInt32  
    int数据值

boolIndex
:   类型：SystemInt32  
    位偏移索引，从0开始

#### 返回值

类型：Boolean  
bool值

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 Int32 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[GetBoolByIndex 重载](d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolByIndex 方法 (Int64, Int32)

[原文連結](http://api.hslcommunication.cn/html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolByIndex 方法 (Byte, Int32)](../html/02376ace-8a06-10a4-600d-90db9f773e92.htm "GetBoolByIndex 方法 (Byte, Int32)")

[GetBoolByIndex 方法 (Byte[], Int32)](../html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm "GetBoolByIndex 方法 (Byte[], Int32)")

[GetBoolByIndex 方法 (Int16, Int32)](../html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm "GetBoolByIndex 方法 (Int16, Int32)")

[GetBoolByIndex 方法 (Int32, Int32)](../html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm "GetBoolByIndex 方法 (Int32, Int32)")

[GetBoolByIndex 方法 (Int64, Int32)](../html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm "GetBoolByIndex 方法 (Int64, Int32)")

[GetBoolByIndex 方法 (UInt16, Int32)](../html/4419a620-899f-6d5a-201a-2cb66c38969c.htm "GetBoolByIndex 方法 (UInt16, Int32)")

[GetBoolByIndex 方法 (UInt32, Int32)](../html/19cac792-a632-936c-5bc3-052023e4e225.htm "GetBoolByIndex 方法 (UInt32, Int32)")

[GetBoolByIndex 方法 (UInt64, Int32)](../html/5f91b230-1822-e4be-02c4-640045232711.htm "GetBoolByIndex 方法 (UInt64, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionGetBoolByIndex 方法 (Int64, Int32) |

获取long类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   
Get the bool value of the boolIndex (starting from 0) offset of the long type data, such as 3, which is the 4th bit

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool GetBoolByIndex(
	this long value,
	int boolIndex
)
```

```
<ExtensionAttribute>
Public Shared Function GetBoolByIndex ( 
	value As Long,
	boolIndex As Integer
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool GetBoolByIndex(
	long long value, 
	int boolIndex
)
```

```
[<ExtensionAttribute>]
static member GetBoolByIndex : 
        value : int64 * 
        boolIndex : int -> bool 
```

#### 参数

value
:   类型：SystemInt64  
    long数据值

boolIndex
:   类型：SystemInt32  
    位偏移索引，从0开始

#### 返回值

类型：Boolean  
bool值

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 Int64 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[GetBoolByIndex 重载](d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolByIndex 方法 (UInt16, Int32)

[原文連結](http://api.hslcommunication.cn/html/4419a620-899f-6d5a-201a-2cb66c38969c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolByIndex 方法 (Byte, Int32)](../html/02376ace-8a06-10a4-600d-90db9f773e92.htm "GetBoolByIndex 方法 (Byte, Int32)")

[GetBoolByIndex 方法 (Byte[], Int32)](../html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm "GetBoolByIndex 方法 (Byte[], Int32)")

[GetBoolByIndex 方法 (Int16, Int32)](../html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm "GetBoolByIndex 方法 (Int16, Int32)")

[GetBoolByIndex 方法 (Int32, Int32)](../html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm "GetBoolByIndex 方法 (Int32, Int32)")

[GetBoolByIndex 方法 (Int64, Int32)](../html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm "GetBoolByIndex 方法 (Int64, Int32)")

[GetBoolByIndex 方法 (UInt16, Int32)](../html/4419a620-899f-6d5a-201a-2cb66c38969c.htm "GetBoolByIndex 方法 (UInt16, Int32)")

[GetBoolByIndex 方法 (UInt32, Int32)](../html/19cac792-a632-936c-5bc3-052023e4e225.htm "GetBoolByIndex 方法 (UInt32, Int32)")

[GetBoolByIndex 方法 (UInt64, Int32)](../html/5f91b230-1822-e4be-02c4-640045232711.htm "GetBoolByIndex 方法 (UInt64, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionGetBoolByIndex 方法 (UInt16, Int32) |

获取ushort类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   
Get the bool value of the boolIndex (starting from 0) offset of the ushort type data, such as 3, which is the 4th bit

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool GetBoolByIndex(
	this ushort value,
	int boolIndex
)
```

```
<ExtensionAttribute>
Public Shared Function GetBoolByIndex ( 
	value As UShort,
	boolIndex As Integer
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool GetBoolByIndex(
	unsigned short value, 
	int boolIndex
)
```

```
[<ExtensionAttribute>]
static member GetBoolByIndex : 
        value : uint16 * 
        boolIndex : int -> bool 
```

#### 参数

value
:   类型：SystemUInt16  
    ushort数据值

boolIndex
:   类型：SystemInt32  
    位偏移索引，从0开始

#### 返回值

类型：Boolean  
bool值

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 UInt16 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[GetBoolByIndex 重载](d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolByIndex 方法 (UInt32, Int32)

[原文連結](http://api.hslcommunication.cn/html/19cac792-a632-936c-5bc3-052023e4e225.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolByIndex 方法 (Byte, Int32)](../html/02376ace-8a06-10a4-600d-90db9f773e92.htm "GetBoolByIndex 方法 (Byte, Int32)")

[GetBoolByIndex 方法 (Byte[], Int32)](../html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm "GetBoolByIndex 方法 (Byte[], Int32)")

[GetBoolByIndex 方法 (Int16, Int32)](../html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm "GetBoolByIndex 方法 (Int16, Int32)")

[GetBoolByIndex 方法 (Int32, Int32)](../html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm "GetBoolByIndex 方法 (Int32, Int32)")

[GetBoolByIndex 方法 (Int64, Int32)](../html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm "GetBoolByIndex 方法 (Int64, Int32)")

[GetBoolByIndex 方法 (UInt16, Int32)](../html/4419a620-899f-6d5a-201a-2cb66c38969c.htm "GetBoolByIndex 方法 (UInt16, Int32)")

[GetBoolByIndex 方法 (UInt32, Int32)](../html/19cac792-a632-936c-5bc3-052023e4e225.htm "GetBoolByIndex 方法 (UInt32, Int32)")

[GetBoolByIndex 方法 (UInt64, Int32)](../html/5f91b230-1822-e4be-02c4-640045232711.htm "GetBoolByIndex 方法 (UInt64, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionGetBoolByIndex 方法 (UInt32, Int32) |

获取uint类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   
Get the bool value of the boolIndex (starting from 0) offset of the uint type data, such as 3, which is the 4th bit

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool GetBoolByIndex(
	this uint value,
	int boolIndex
)
```

```
<ExtensionAttribute>
Public Shared Function GetBoolByIndex ( 
	value As UInteger,
	boolIndex As Integer
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool GetBoolByIndex(
	unsigned int value, 
	int boolIndex
)
```

```
[<ExtensionAttribute>]
static member GetBoolByIndex : 
        value : uint32 * 
        boolIndex : int -> bool 
```

#### 参数

value
:   类型：SystemUInt32  
    uint数据值

boolIndex
:   类型：SystemInt32  
    位偏移索引，从0开始

#### 返回值

类型：Boolean  
bool值

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 UInt32 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[GetBoolByIndex 重载](d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetBoolByIndex 方法 (UInt64, Int32)

[原文連結](http://api.hslcommunication.cn/html/5f91b230-1822-e4be-02c4-640045232711.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication](../html/c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication")

[HslExtension 类](../html/05bbd254-20d0-2658-04fb-05c400c448dc.htm "HslExtension 类")

[HslExtension 方法](../html/32366708-60b9-6e27-c174-0edbd522cbd4.htm "HslExtension 方法")

[GetBoolByIndex 方法](../html/d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm "GetBoolByIndex 方法 ")

[GetBoolByIndex 方法 (Byte, Int32)](../html/02376ace-8a06-10a4-600d-90db9f773e92.htm "GetBoolByIndex 方法 (Byte, Int32)")

[GetBoolByIndex 方法 (Byte[], Int32)](../html/1b15d4d6-a193-0bff-3c9d-a8c47066ca2a.htm "GetBoolByIndex 方法 (Byte[], Int32)")

[GetBoolByIndex 方法 (Int16, Int32)](../html/013b9689-81f1-1871-bd53-2a6ca2ae153a.htm "GetBoolByIndex 方法 (Int16, Int32)")

[GetBoolByIndex 方法 (Int32, Int32)](../html/65142f2b-8313-8bf6-db4c-5567801c19f3.htm "GetBoolByIndex 方法 (Int32, Int32)")

[GetBoolByIndex 方法 (Int64, Int32)](../html/fd4a7676-b7b5-2880-72c6-76e9ed4e0467.htm "GetBoolByIndex 方法 (Int64, Int32)")

[GetBoolByIndex 方法 (UInt16, Int32)](../html/4419a620-899f-6d5a-201a-2cb66c38969c.htm "GetBoolByIndex 方法 (UInt16, Int32)")

[GetBoolByIndex 方法 (UInt32, Int32)](../html/19cac792-a632-936c-5bc3-052023e4e225.htm "GetBoolByIndex 方法 (UInt32, Int32)")

[GetBoolByIndex 方法 (UInt64, Int32)](../html/5f91b230-1822-e4be-02c4-640045232711.htm "GetBoolByIndex 方法 (UInt64, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslExtensionGetBoolByIndex 方法 (UInt64, Int32) |

获取ulong类型数据的第 boolIndex (从0起始)偏移的bool值，比如3，就是第4位   
Get the bool value of the boolIndex (starting from 0) offset of the ulong type data, such as 3, which is the 4th bit

**命名空间：**
 [HslCommunication](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool GetBoolByIndex(
	this ulong value,
	int boolIndex
)
```

```
<ExtensionAttribute>
Public Shared Function GetBoolByIndex ( 
	value As ULong,
	boolIndex As Integer
) As Boolean
```

```
public:
[ExtensionAttribute]
static bool GetBoolByIndex(
	unsigned long long value, 
	int boolIndex
)
```

```
[<ExtensionAttribute>]
static member GetBoolByIndex : 
        value : uint64 * 
        boolIndex : int -> bool 
```

#### 参数

value
:   类型：SystemUInt64  
    ulong数据值

boolIndex
:   类型：SystemInt32  
    位偏移索引，从0开始

#### 返回值

类型：Boolean  
bool值

#### 备注

在 Visual Basic 和 C# 中，这个方法可以当成为类型 UInt64 的实例方法来调用。在采用实例方法语法调用这个方法时，请省略第一个参数。请参考 [扩展方法 (Visual Basic)](http://msdn.microsoft.com/zh-CN/library/bb384936.aspx) 或 [扩展方法 (C# 编程指南)](http://msdn.microsoft.com/zh-CN/library/bb383977.aspx) 获取更多信息。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslExtension 类](05bbd254-20d0-2658-04fb-05c400c448dc.htm)

[GetBoolByIndex 重载](d56f77a2-cf9a-5e04-ff46-b7eb8202b1f7.htm)

[HslCommunication 命名空间](c136d3de-eab7-9b0f-4bdf-d891297c8018.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)