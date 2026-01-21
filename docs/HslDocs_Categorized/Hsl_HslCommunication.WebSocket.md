# HslCommunication - HslCommunication.WebSocket

> 分類頁數: 30



---
## HslCommunication.WebSocket

[原文連結](http://api.hslcommunication.cn/html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient.OnClientApplicationMessageReceiveDelegate 委托](../html/afc01b80-cc10-5563-6629-87bfeb3ed6e0.htm "WebSocketClient.OnClientApplicationMessageReceiveDelegate 委托")

[WebSocketClient.OnClientConnectedDelegate 委托](../html/96ca6387-967b-f29f-4883-862da1f75d60.htm "WebSocketClient.OnClientConnectedDelegate 委托")

[WebSocketHelper 类](../html/82f6c790-57d7-18c0-f7bd-b7a9c9e2661b.htm "WebSocketHelper 类")

[WebSocketMessage 类](../html/17ae7d0d-9840-3bf6-9780-62d69422744a.htm "WebSocketMessage 类")

[WebSocketQANet 类](../html/e6cdc29f-e609-c678-4a73-32c909fcfda8.htm "WebSocketQANet 类")

[WebSocketServer 类](../html/a46e31a5-1011-9c93-ffcf-fa8404df0a65.htm "WebSocketServer 类")

[WebSocketServer.OnClientApplicationMessageReceiveDelegate 委托](../html/a15be011-9e53-ed50-7a1f-647571adf47a.htm "WebSocketServer.OnClientApplicationMessageReceiveDelegate 委托")

[WebSocketServer.OnClientConnectedDelegate 委托](../html/de91bef1-8558-8517-7f99-3b1df66a4bdf.htm "WebSocketServer.OnClientConnectedDelegate 委托")

[WebSocketSession 类](../html/26a36542-0771-80a4-ed4d-e6d93c341107.htm "WebSocketSession 类")

[WSOpCode 枚举](../html/f5b46717-4684-558b-d616-9d9c51f5c8d7.htm "WSOpCode 枚举")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.WebSocket 命名空间 |

[缺少 "N:HslCommunication.WebSocket" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类代码示例 | [WebSocketClient](28c4182b-6004-535e-ffa7-d7352570e9ea.htm) | websocket协议的客户端实现，支持从服务器订阅，发布数据内容信息，详细参考api文档信息  Client implementation of the websocket protocol. It supports subscribing from the server and publishing data content information. |
| 公共类 | [WebSocketHelper](82f6c790-57d7-18c0-f7bd-b7a9c9e2661b.htm) | websocket的相关辅助的方法 |
| 公共类 | [WebSocketMessage](17ae7d0d-9840-3bf6-9780-62d69422744a.htm) | websocket 协议下的单个消息的数据对象  Data object for a single message under the websocket protocol |
| 公共类 | [WebSocketQANet](e6cdc29f-e609-c678-4a73-32c909fcfda8.htm) | WebSocket的问答机制的客户端，本客户端将会在请求头上追加 RequestAndAnswer: true，本客户端将会请求服务器的信息，然后等待服务器的返回  Client of WebSocket Q & A mechanism, this client will append RequestAndAnswer: true to the request header, this client will request the server information, and then wait for the server to return |
| 公共类代码示例 | [WebSocketServer](a46e31a5-1011-9c93-ffcf-fa8404df0a65.htm) | WebSocket协议的实现，支持创建自定义的websocket服务器，直接给其他的网页端，客户端，手机端发送数据信息，详细看api文档说明  The implementation of the WebSocket protocol supports the creation of custom websocket servers and sends data information directly to other web pages, clients, and mobile phones. See the API documentation for details. |
| 公共类 | [WebSocketSession](26a36542-0771-80a4-ed4d-e6d93c341107.htm) | websocket 的会话客户端 |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [WebSocketClientOnClientApplicationMessageReceiveDelegate](afc01b80-cc10-5563-6629-87bfeb3ed6e0.htm) | websocket的消息收到委托  websocket message received delegate |
| 公共委托 | [WebSocketClientOnClientConnectedDelegate](96ca6387-967b-f29f-4883-862da1f75d60.htm) | 连接服务器成功的委托  Connection server successfully delegated |
| 公共委托 | [WebSocketServerOnClientApplicationMessageReceiveDelegate](a15be011-9e53-ed50-7a1f-647571adf47a.htm) | websocket的消息收到委托  websocket message received delegate |
| 公共委托 | [WebSocketServerOnClientConnectedDelegate](de91bef1-8558-8517-7f99-3b1df66a4bdf.htm) | 当前websocket连接上服务器的事件委托  Event delegation of the server on the current websocket connection |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [WSOpCode](f5b46717-4684-558b-d616-9d9c51f5c8d7.htm) | websocket 协议的 op的枚举信息 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WebSocketClient 类

[原文連結](http://api.hslcommunication.cn/html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 构造函数](../html/e87bc526-08c8-14eb-eccb-236265a4467b.htm "WebSocketClient 构造函数 ")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[WebSocketClient 事件](../html/a55182f6-5e52-67b9-9c41-876b1f8f3290.htm "WebSocketClient 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClient 类 |

websocket协议的客户端实现，支持从服务器订阅，发布数据内容信息，详细参考api文档信息  
Client implementation of the websocket protocol. It supports subscribing from the server and publishing data content information.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.WebSocketWebSocketClient

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class WebSocketClient : IDisposable
```

```
Public Class WebSocketClient
	Implements IDisposable
```

```
public ref class WebSocketClient : IDisposable
```

```
type WebSocketClient =  
    class
        interface IDisposable
    end
```

WebSocketClient 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [WebSocketClient(String)](c32d9d82-99ee-9a70-199a-20c55f3b9887.htm) | 使用指定的url来实例化一个默认的对象，例如 ws://127.0.0.1:1883/A/B?C=123456 或是 ws://www.hslcommunication.cn:1883  Use the specified url to instantiate a default object, such as ws://127.0.0.1:1883/A/B?C=123456 or ws://www.hslcommunication.cn:1883s |
| 公共方法 | [WebSocketClient(String, Int32)](53c2077a-f5fd-4244-ae4a-0903b3cc31ff.htm) | 使用指定的ip，端口来实例化一个默认的对象  Use the specified ip and port to instantiate a default objects |
| 公共方法 | [WebSocketClient(String, Int32, String)](d1ddf0d8-88b8-d1ac-095e-e7b22a7303c4.htm) | 使用指定的ip，端口，额外的url信息来实例化一个默认的对象  Use the specified ip, port, and additional url information to instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AutoDecompress](0c7971d1-9dca-6855-5751-46a8935781a8.htm) | 获取或设置数据自动解压缩功能，默认开启自动解压缩功能，如果需要自己手动解压缩数据的话，请设置为 false。  Get or set the automatic data decompression function. By default, the automatic decompression function is enabled. If you need to manually decompress the data, please set it to false. |
| 公共属性 | [ConnectTimeOut](6d811cfb-79cc-a484-9c76-1db7131b629f.htm) | 获取或设置当前客户端的连接超时时间，默认10,000毫秒，单位ms  Gets or sets the connection timeout of the current client. The default is 10,000 milliseconds. The unit is ms. |
| 公共属性 | [GetCarryHostAndPort](28ae13f2-6772-45bf-195c-5fcee79fb431.htm) | 获取或设置连接时候的 GET 命令后面是否协议Host及端口号信息，默认为 False，不携带。  Obtain or set the protocol host and port number information after the GET command when connecting, which is False by default and does not carry it. |
| 公共属性 | [IpAddress](9aba5840-5dda-e642-b3c0-4546344607fb.htm) | Mqtt服务器的ip地址  IP address of Mqtt server |
| 公共属性 | [IsClosed](d84dc9d8-c23b-1e71-adaa-026921fd6844.htm) | 获取当前的客户端状态是否关闭了连接，当自己手动处理网络异常事件的时候，在重连之前就需要判断是否关闭了连接。  Obtain whether the current client status has closed the connection. When manually handling network abnormal events, you need to determine whether the connection is closed before reconnecting. |
| 公共属性代码示例 | [LogNet](04e8ead6-533c-4b51-c58c-33f0572b80eb.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) |
| 公共属性 | [Port](880da6be-6ca9-c507-f4c8-59f3f53d119f.htm) | 端口号。默认1883  The port number. Default 1883 |
| 公共属性 | [SupportDeflate](57017a28-6348-bd7b-80dd-502558db693a.htm) | 获取或设置是否支持数据压缩功能，默认支持压缩功能，如果服务器不支持压缩功能，请设置为 false。  Get or set whether data compression function is supported. By default, compression function is supported. If the server does not support compression function, please set it to false. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ConnectClose](77d60360-6674-11a7-181f-3790417f1ff2.htm) | 关闭与Websocket服务器的连接。  Close the connection to the Mqtt server. |
| 公共方法 | [ConnectCloseAsync](f8dcad2a-d6f2-a62d-52ac-d9e60c43b039.htm) | 关闭与Websocket服务器的连接。  Close the connection to the Mqtt server. |
| 公共方法 | [ConnectServer](d50494a4-c56e-e78d-4564-aa8fc4d014bb.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServer(String)](1239b330-9b8f-046e-1e48-f36226b94c3d.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServerAsync](e80ad7fe-314c-7b3f-ee1c-27ebdac8a30c.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServerAsync(String)](ad11585d-4f8b-c6c1-bc56-a20519b925e6.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [Dispose](1266d90f-2ca3-e8ca-8de5-6e1a1a06764d.htm) | 释放被 WebSocketClient 使用的所有资源 |
| 受保护的方法 | [Dispose(Boolean)](07f93192-03c2-0a68-c4e3-e9f948ec46c4.htm) | 释放当前的对象 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [SendServer(String)](7348ddc5-9af8-8fa3-98ff-08e626436673.htm) | 发送数据到WebSocket的服务器  Send data to WebSocket server |
| 公共方法 | [SendServer(Boolean, String)](37d7ec3b-3ad5-3891-b672-0d35be806d0f.htm) | 发送数据到WebSocket的服务器，可以指定是否进行掩码操作  Send data to the WebSocket server, you can specify whether to perform a mask operation |
| 公共方法 | [SendServer(Int32, Boolean, Byte)](f2434409-540a-3096-f13b-c4ccb6f82872.htm) | 发送自定义的命令到WebSocket服务器，可以指定操作码，是否掩码操作，原始字节数据  Send custom commands to the WebSocket server, you can specify the operation code, whether to mask operation, raw byte data |
| 公共方法 | [ToString](de57374b-55ef-d4e0-7447-ea2ff7e3dd7a.htm) | (重写 ObjectToString.) |
| 公共方法 | [UseSSL(X509Certificate)](484fce1a-01da-3867-80a3-88f2c48e2513.htm) | 使用一个证书来初始化 SSL/TLS 通信  Use a certificate to initiate SSL/TLS communication |
| 公共方法 | [UseSSL(String)](805c7a5c-179b-0e6f-1139-3cf539c11bb8.htm) | 使用一个证书路径来初始化 SSL/TLS 通信  Use a certificate path to initialize SSL/TLS communication |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnClientApplicationMessageReceive](68bacb0e-54d1-93dd-fae1-35b78c34b076.htm) | websocket的消息收到时触发  Triggered when a websocket message is received |
| 公共事件 | [OnClientConnected](e85d73c0-bbc3-9649-1601-20e987d20d8f.htm) | 当客户端连接成功触发事件，就算是重新连接服务器后，也是会触发的  The event is triggered when the client is connected successfully, even after reconnecting to the server. |
| 公共事件 | [OnNetworkError](26145134-6945-5f35-075e-d911c8db0841.htm) | 当网络发生异常的时候触发的事件，用户应该在事件里进行重连服务器 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)示例

本客户端使用起来非常的方便，基本就是实例化，绑定一个数据接收的事件即可，如下所示

简单的实例化

[复制](# "复制")

```
private WebSocketClient webSocketClient;

public void Start( )
{
    webSocketClient = new WebSocketClient( "127.0.0.1", 1883 );
    webSocketClient.OnClientApplicationMessageReceive += WebSocketClient_OnClientApplicationMessageReceive;
    OperateResult connect = webSocketClient.ConnectServer( );
    if (connect.IsSuccess)
    {
        Console.WriteLine( "connect successful" );
    }
    else
    {
        // 连接失败。需要等会重新连接，需要注意的是，对同一个实例来说，OnClientApplicationMessageReceive不能重复绑定事件
        // Connection failed. Need to wait for reconnection, it should be noted that, for the same instance, OnClientApplicationMessageReceive can not repeatedly bind events
        Console.WriteLine( "connect failed" );
    }

    // 当连接成功后，网络发生了异常。客户端会自动重新连接的。
    // When the connection was successful, an exception occurred on the network. The client will automatically reconnect.
}

private void WebSocketClient_OnClientApplicationMessageReceive( WebSocketMessage message )
{
    Console.WriteLine( message.ToString( ) );
}
```

假设我们需要发数据给服务端，那么可以参考如下的方式

发送数据

[复制](# "复制")

```
public void button_Click(object sender, EventArgs e )
{
    webSocketClient.SendServer( "This is a message from hslcommunication" );
}
```

如果我们需要搭配服务器来做订阅推送的功能的话，写法上会稍微有点区别，按照下面的代码来写。

订阅操作

[复制](# "复制")

```
private WebSocketClient webSocketClient;

public void Start( )
{
    webSocketClient = new WebSocketClient( "127.0.0.1", 1883 );
    // 连接上服务器的时候触发，在断线重连的时候也会触发，如果使用发送服务器实现订阅的方式，在下面的事件里订阅是合理的
    // Triggered when connected to the server, it will also trigger when disconnected and reconnected. 
    // If the sending server is used to implement the subscription, the subscription is reasonable in the following events
    webSocketClient.OnClientConnected += ( ) => 
    {
        // 订阅的内容添加在这里
        // Subscribed content added here
        webSocketClient.SendServer( "A" );
    };
    webSocketClient.OnClientApplicationMessageReceive += WebSocketClient_OnClientApplicationMessageReceive;
    OperateResult connect = webSocketClient.ConnectServer( );
    if (connect.IsSuccess)
    {
        Console.WriteLine( "connect successful" );
    }
    else
    {
        // 连接失败。需要等会重新连接，需要注意的是，对同一个实例来说，OnClientApplicationMessageReceive不能重复绑定事件
        // Connection failed. Need to wait for reconnection, it should be noted that, for the same instance, OnClientApplicationMessageReceive can not repeatedly bind events
        Console.WriteLine( "connect failed" );
    }

    // 当连接成功后，网络发生了异常。客户端会自动重新连接的。
    // When the connection was successful, an exception occurred on the network. The client will automatically reconnect.

}

private void WebSocketClient_OnClientApplicationMessageReceive( WebSocketMessage message )
{
    // 一般来说，一个客户端订阅一个topic，如果要订阅多个的话，message就要区分主题，需要采用json格式的数据
    // Generally, a client subscribes to a topic. If you want to subscribe to multiple topics, 
    // the message must distinguish between topics, and payload in json format is required.
    Console.WriteLine( message.ToString( ) );
}
```

当网络发生异常的时候，我们需要这么来进行重新连接。

异常重连

[复制](# "复制")

```
private WebSocketClient webSocketClient;

public void Start( )
{
    webSocketClient = new WebSocketClient( "127.0.0.1", 1883 );
    // 连接上服务器的时候触发，在断线重连的时候也会触发，如果使用发送服务器实现订阅的方式，在下面的事件里订阅是合理的
    // Triggered when connected to the server, it will also trigger when disconnected and reconnected. 
    // If the sending server is used to implement the subscription, the subscription is reasonable in the following events
    webSocketClient.OnClientConnected += ( ) =>
    {
        // 订阅的内容添加在这里
        // Subscribed content added here
        webSocketClient.SendServer( "A" );
    };
    webSocketClient.OnNetworkError += ( object sender, EventArgs e ) =>
    {
        // 当网络异常的时候触发，可以在此处重连服务器
        if (sender is WebSocketClient client)
        {
            // 开始重连服务器，直到连接成功为止
            client.LogNet?.WriteInfo( "网络异常，准备10秒后重新连接。" );
            while (true)
            {
                // 每隔10秒重连
                System.Threading.Thread.Sleep( 10_000 );
                client.LogNet?.WriteInfo( "准备重新连接服务器..." );
                OperateResult connectResult = client.ConnectServer( );
                if (connectResult.IsSuccess)
                {
                    client.LogNet?.WriteInfo( "连接服务器成功！" );
                    break;
                }
                client.LogNet?.WriteInfo( "连接失败，准备10秒后重新连接。" );
            }
        }
    };
    webSocketClient.OnClientApplicationMessageReceive += WebSocketClient_OnClientApplicationMessageReceive;
    OperateResult connect = webSocketClient.ConnectServer( );
    if (connect.IsSuccess)
    {
        Console.WriteLine( "connect successful" );
    }
    else
    {
        // 连接失败。需要等会重新连接，需要注意的是，对同一个实例来说，OnClientApplicationMessageReceive不能重复绑定事件
        // Connection failed. Need to wait for reconnection, it should be noted that, for the same instance, OnClientApplicationMessageReceive can not repeatedly bind events
        Console.WriteLine( "connect failed" );
    }

    // 当连接成功后，网络发生了异常。客户端会自动重新连接的。
    // When the connection was successful, an exception occurred on the network. The client will automatically reconnect.

}

private void WebSocketClient_OnClientApplicationMessageReceive( WebSocketMessage message )
{
    // 一般来说，一个客户端订阅一个topic，如果要订阅多个的话，message就要区分主题，需要采用json格式的数据
    // Generally, a client subscribes to a topic. If you want to subscribe to multiple topics, 
    // the message must distinguish between topics, and payload in json format is required.
    Console.WriteLine( message.ToString( ) );
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WebSocketClient 构造函数 

[原文連結](http://api.hslcommunication.cn/html/e87bc526-08c8-14eb-eccb-236265a4467b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 构造函数](../html/e87bc526-08c8-14eb-eccb-236265a4467b.htm "WebSocketClient 构造函数 ")

[WebSocketClient 构造函数 (String)](../html/c32d9d82-99ee-9a70-199a-20c55f3b9887.htm "WebSocketClient 构造函数 (String)")

[WebSocketClient 构造函数 (String, Int32)](../html/53c2077a-f5fd-4244-ae4a-0903b3cc31ff.htm "WebSocketClient 构造函数 (String, Int32)")

[WebSocketClient 构造函数 (String, Int32, String)](../html/d1ddf0d8-88b8-d1ac-095e-e7b22a7303c4.htm "WebSocketClient 构造函数 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClient 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [WebSocketClient(String)](c32d9d82-99ee-9a70-199a-20c55f3b9887.htm) | 使用指定的url来实例化一个默认的对象，例如 ws://127.0.0.1:1883/A/B?C=123456 或是 ws://www.hslcommunication.cn:1883  Use the specified url to instantiate a default object, such as ws://127.0.0.1:1883/A/B?C=123456 or ws://www.hslcommunication.cn:1883s |
| 公共方法 | [WebSocketClient(String, Int32)](53c2077a-f5fd-4244-ae4a-0903b3cc31ff.htm) | 使用指定的ip，端口来实例化一个默认的对象  Use the specified ip and port to instantiate a default objects |
| 公共方法 | [WebSocketClient(String, Int32, String)](d1ddf0d8-88b8-d1ac-095e-e7b22a7303c4.htm) | 使用指定的ip，端口，额外的url信息来实例化一个默认的对象  Use the specified ip, port, and additional url information to instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WebSocketClient 构造函数 (String)

[原文連結](http://api.hslcommunication.cn/html/c32d9d82-99ee-9a70-199a-20c55f3b9887.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 构造函数](../html/e87bc526-08c8-14eb-eccb-236265a4467b.htm "WebSocketClient 构造函数 ")

[WebSocketClient 构造函数 (String)](../html/c32d9d82-99ee-9a70-199a-20c55f3b9887.htm "WebSocketClient 构造函数 (String)")

[WebSocketClient 构造函数 (String, Int32)](../html/53c2077a-f5fd-4244-ae4a-0903b3cc31ff.htm "WebSocketClient 构造函数 (String, Int32)")

[WebSocketClient 构造函数 (String, Int32, String)](../html/d1ddf0d8-88b8-d1ac-095e-e7b22a7303c4.htm "WebSocketClient 构造函数 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClient 构造函数 (String) |

使用指定的url来实例化一个默认的对象，例如 ws://127.0.0.1:1883/A/B?C=123456 或是 ws://www.hslcommunication.cn:1883  
Use the specified url to instantiate a default object, such as ws://127.0.0.1:1883/A/B?C=123456 or ws://www.hslcommunication.cn:1883s

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public WebSocketClient(
	string url
)
```

```
Public Sub New ( 
	url As String
)
```

```
public:
WebSocketClient(
	String^ url
)
```

```
new : 
        url : string -> WebSocketClient
```

#### 参数

url
:   类型：SystemString  
    完整的ws地址

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[WebSocketClient 重载](e87bc526-08c8-14eb-eccb-236265a4467b.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WebSocketClient 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/53c2077a-f5fd-4244-ae4a-0903b3cc31ff.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 构造函数](../html/e87bc526-08c8-14eb-eccb-236265a4467b.htm "WebSocketClient 构造函数 ")

[WebSocketClient 构造函数 (String)](../html/c32d9d82-99ee-9a70-199a-20c55f3b9887.htm "WebSocketClient 构造函数 (String)")

[WebSocketClient 构造函数 (String, Int32)](../html/53c2077a-f5fd-4244-ae4a-0903b3cc31ff.htm "WebSocketClient 构造函数 (String, Int32)")

[WebSocketClient 构造函数 (String, Int32, String)](../html/d1ddf0d8-88b8-d1ac-095e-e7b22a7303c4.htm "WebSocketClient 构造函数 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClient 构造函数 (String, Int32) |

使用指定的ip，端口来实例化一个默认的对象  
Use the specified ip and port to instantiate a default objects

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public WebSocketClient(
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
WebSocketClient(
	String^ ipAddress, 
	int port
)
```

```
new : 
        ipAddress : string * 
        port : int -> WebSocketClient
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

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[WebSocketClient 重载](e87bc526-08c8-14eb-eccb-236265a4467b.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WebSocketClient 构造函数 (String, Int32, String)

[原文連結](http://api.hslcommunication.cn/html/d1ddf0d8-88b8-d1ac-095e-e7b22a7303c4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 构造函数](../html/e87bc526-08c8-14eb-eccb-236265a4467b.htm "WebSocketClient 构造函数 ")

[WebSocketClient 构造函数 (String)](../html/c32d9d82-99ee-9a70-199a-20c55f3b9887.htm "WebSocketClient 构造函数 (String)")

[WebSocketClient 构造函数 (String, Int32)](../html/53c2077a-f5fd-4244-ae4a-0903b3cc31ff.htm "WebSocketClient 构造函数 (String, Int32)")

[WebSocketClient 构造函数 (String, Int32, String)](../html/d1ddf0d8-88b8-d1ac-095e-e7b22a7303c4.htm "WebSocketClient 构造函数 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClient 构造函数 (String, Int32, String) |

使用指定的ip，端口，额外的url信息来实例化一个默认的对象  
Use the specified ip, port, and additional url information to instantiate a default object

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public WebSocketClient(
	string ipAddress,
	int port,
	string url
)
```

```
Public Sub New ( 
	ipAddress As String,
	port As Integer,
	url As String
)
```

```
public:
WebSocketClient(
	String^ ipAddress, 
	int port, 
	String^ url
)
```

```
new : 
        ipAddress : string * 
        port : int * 
        url : string -> WebSocketClient
```

#### 参数

ipAddress
:   类型：SystemString  
    Ip地址信息

port
:   类型：SystemInt32  
    端口号信息

url
:   类型：SystemString  
    额外的信息，比如 /A/B?C=123456

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[WebSocketClient 重载](e87bc526-08c8-14eb-eccb-236265a4467b.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WebSocketClient 属性

[原文連結](http://api.hslcommunication.cn/html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[AutoDecompress 属性](../html/0c7971d1-9dca-6855-5751-46a8935781a8.htm "AutoDecompress 属性 ")

[ConnectTimeOut 属性](../html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm "ConnectTimeOut 属性 ")

[GetCarryHostAndPort 属性](../html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm "GetCarryHostAndPort 属性 ")

[IpAddress 属性](../html/9aba5840-5dda-e642-b3c0-4546344607fb.htm "IpAddress 属性 ")

[IsClosed 属性](../html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm "IsClosed 属性 ")

[LogNet 属性](../html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm "LogNet 属性 ")

[Port 属性](../html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm "Port 属性 ")

[SupportDeflate 属性](../html/57017a28-6348-bd7b-80dd-502558db693a.htm "SupportDeflate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClient 属性 |

[WebSocketClient](28c4182b-6004-535e-ffa7-d7352570e9ea.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AutoDecompress](0c7971d1-9dca-6855-5751-46a8935781a8.htm) | 获取或设置数据自动解压缩功能，默认开启自动解压缩功能，如果需要自己手动解压缩数据的话，请设置为 false。  Get or set the automatic data decompression function. By default, the automatic decompression function is enabled. If you need to manually decompress the data, please set it to false. |
| 公共属性 | [ConnectTimeOut](6d811cfb-79cc-a484-9c76-1db7131b629f.htm) | 获取或设置当前客户端的连接超时时间，默认10,000毫秒，单位ms  Gets or sets the connection timeout of the current client. The default is 10,000 milliseconds. The unit is ms. |
| 公共属性 | [GetCarryHostAndPort](28ae13f2-6772-45bf-195c-5fcee79fb431.htm) | 获取或设置连接时候的 GET 命令后面是否协议Host及端口号信息，默认为 False，不携带。  Obtain or set the protocol host and port number information after the GET command when connecting, which is False by default and does not carry it. |
| 公共属性 | [IpAddress](9aba5840-5dda-e642-b3c0-4546344607fb.htm) | Mqtt服务器的ip地址  IP address of Mqtt server |
| 公共属性 | [IsClosed](d84dc9d8-c23b-1e71-adaa-026921fd6844.htm) | 获取当前的客户端状态是否关闭了连接，当自己手动处理网络异常事件的时候，在重连之前就需要判断是否关闭了连接。  Obtain whether the current client status has closed the connection. When manually handling network abnormal events, you need to determine whether the connection is closed before reconnecting. |
| 公共属性代码示例 | [LogNet](04e8ead6-533c-4b51-c58c-33f0572b80eb.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) |
| 公共属性 | [Port](880da6be-6ca9-c507-f4c8-59f3f53d119f.htm) | 端口号。默认1883  The port number. Default 1883 |
| 公共属性 | [SupportDeflate](57017a28-6348-bd7b-80dd-502558db693a.htm) | 获取或设置是否支持数据压缩功能，默认支持压缩功能，如果服务器不支持压缩功能，请设置为 false。  Get or set whether data compression function is supported. By default, compression function is supported. If the server does not support compression function, please set it to false. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AutoDecompress 属性 

[原文連結](http://api.hslcommunication.cn/html/0c7971d1-9dca-6855-5751-46a8935781a8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[AutoDecompress 属性](../html/0c7971d1-9dca-6855-5751-46a8935781a8.htm "AutoDecompress 属性 ")

[ConnectTimeOut 属性](../html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm "ConnectTimeOut 属性 ")

[GetCarryHostAndPort 属性](../html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm "GetCarryHostAndPort 属性 ")

[IpAddress 属性](../html/9aba5840-5dda-e642-b3c0-4546344607fb.htm "IpAddress 属性 ")

[IsClosed 属性](../html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm "IsClosed 属性 ")

[LogNet 属性](../html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm "LogNet 属性 ")

[Port 属性](../html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm "Port 属性 ")

[SupportDeflate 属性](../html/57017a28-6348-bd7b-80dd-502558db693a.htm "SupportDeflate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientAutoDecompress 属性 |

获取或设置数据自动解压缩功能，默认开启自动解压缩功能，如果需要自己手动解压缩数据的话，请设置为 false。  
Get or set the automatic data decompression function. By default, the automatic decompression function is enabled. If you need to manually decompress the data, please set it to false.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool AutoDecompress { get; set; }
```

```
Public Property AutoDecompress As Boolean
	Get
	Set
```

```
public:
property bool AutoDecompress {
	bool get ();
	void set (bool value);
}
```

```
member AutoDecompress : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectTimeOut 属性 

[原文連結](http://api.hslcommunication.cn/html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[AutoDecompress 属性](../html/0c7971d1-9dca-6855-5751-46a8935781a8.htm "AutoDecompress 属性 ")

[ConnectTimeOut 属性](../html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm "ConnectTimeOut 属性 ")

[GetCarryHostAndPort 属性](../html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm "GetCarryHostAndPort 属性 ")

[IpAddress 属性](../html/9aba5840-5dda-e642-b3c0-4546344607fb.htm "IpAddress 属性 ")

[IsClosed 属性](../html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm "IsClosed 属性 ")

[LogNet 属性](../html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm "LogNet 属性 ")

[Port 属性](../html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm "Port 属性 ")

[SupportDeflate 属性](../html/57017a28-6348-bd7b-80dd-502558db693a.htm "SupportDeflate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientConnectTimeOut 属性 |

获取或设置当前客户端的连接超时时间，默认10,000毫秒，单位ms  
Gets or sets the connection timeout of the current client. The default is 10,000 milliseconds. The unit is ms.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ConnectTimeOut { get; set; }
```

```
Public Property ConnectTimeOut As Integer
	Get
	Set
```

```
public:
property int ConnectTimeOut {
	int get ();
	void set (int value);
}
```

```
member ConnectTimeOut : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetCarryHostAndPort 属性 

[原文連結](http://api.hslcommunication.cn/html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[AutoDecompress 属性](../html/0c7971d1-9dca-6855-5751-46a8935781a8.htm "AutoDecompress 属性 ")

[ConnectTimeOut 属性](../html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm "ConnectTimeOut 属性 ")

[GetCarryHostAndPort 属性](../html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm "GetCarryHostAndPort 属性 ")

[IpAddress 属性](../html/9aba5840-5dda-e642-b3c0-4546344607fb.htm "IpAddress 属性 ")

[IsClosed 属性](../html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm "IsClosed 属性 ")

[LogNet 属性](../html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm "LogNet 属性 ")

[Port 属性](../html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm "Port 属性 ")

[SupportDeflate 属性](../html/57017a28-6348-bd7b-80dd-502558db693a.htm "SupportDeflate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientGetCarryHostAndPort 属性 |

获取或设置连接时候的 GET 命令后面是否协议Host及端口号信息，默认为 False，不携带。  
Obtain or set the protocol host and port number information after the GET command when connecting, which is False by default and does not carry it.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool GetCarryHostAndPort { get; set; }
```

```
Public Property GetCarryHostAndPort As Boolean
	Get
	Set
```

```
public:
property bool GetCarryHostAndPort {
	bool get ();
	void set (bool value);
}
```

```
member GetCarryHostAndPort : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IpAddress 属性 

[原文連結](http://api.hslcommunication.cn/html/9aba5840-5dda-e642-b3c0-4546344607fb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[AutoDecompress 属性](../html/0c7971d1-9dca-6855-5751-46a8935781a8.htm "AutoDecompress 属性 ")

[ConnectTimeOut 属性](../html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm "ConnectTimeOut 属性 ")

[GetCarryHostAndPort 属性](../html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm "GetCarryHostAndPort 属性 ")

[IpAddress 属性](../html/9aba5840-5dda-e642-b3c0-4546344607fb.htm "IpAddress 属性 ")

[IsClosed 属性](../html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm "IsClosed 属性 ")

[LogNet 属性](../html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm "LogNet 属性 ")

[Port 属性](../html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm "Port 属性 ")

[SupportDeflate 属性](../html/57017a28-6348-bd7b-80dd-502558db693a.htm "SupportDeflate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientIpAddress 属性 |

Mqtt服务器的ip地址  
IP address of Mqtt server

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string IpAddress { get; set; }
```

```
Public Property IpAddress As String
	Get
	Set
```

```
public:
property String^ IpAddress {
	String^ get ();
	void set (String^ value);
}
```

```
member IpAddress : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsClosed 属性 

[原文連結](http://api.hslcommunication.cn/html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[AutoDecompress 属性](../html/0c7971d1-9dca-6855-5751-46a8935781a8.htm "AutoDecompress 属性 ")

[ConnectTimeOut 属性](../html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm "ConnectTimeOut 属性 ")

[GetCarryHostAndPort 属性](../html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm "GetCarryHostAndPort 属性 ")

[IpAddress 属性](../html/9aba5840-5dda-e642-b3c0-4546344607fb.htm "IpAddress 属性 ")

[IsClosed 属性](../html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm "IsClosed 属性 ")

[LogNet 属性](../html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm "LogNet 属性 ")

[Port 属性](../html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm "Port 属性 ")

[SupportDeflate 属性](../html/57017a28-6348-bd7b-80dd-502558db693a.htm "SupportDeflate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientIsClosed 属性 |

获取当前的客户端状态是否关闭了连接，当自己手动处理网络异常事件的时候，在重连之前就需要判断是否关闭了连接。  
Obtain whether the current client status has closed the connection. When manually handling network abnormal events, you need to determine whether the connection is closed before reconnecting.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsClosed { get; }
```

```
Public ReadOnly Property IsClosed As Boolean
	Get
```

```
public:
property bool IsClosed {
	bool get ();
}
```

```
member IsClosed : bool with get
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LogNet 属性 

[原文連結](http://api.hslcommunication.cn/html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[AutoDecompress 属性](../html/0c7971d1-9dca-6855-5751-46a8935781a8.htm "AutoDecompress 属性 ")

[ConnectTimeOut 属性](../html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm "ConnectTimeOut 属性 ")

[GetCarryHostAndPort 属性](../html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm "GetCarryHostAndPort 属性 ")

[IpAddress 属性](../html/9aba5840-5dda-e642-b3c0-4546344607fb.htm "IpAddress 属性 ")

[IsClosed 属性](../html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm "IsClosed 属性 ")

[LogNet 属性](../html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm "LogNet 属性 ")

[Port 属性](../html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm "Port 属性 ")

[SupportDeflate 属性](../html/57017a28-6348-bd7b-80dd-502558db693a.htm "SupportDeflate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientLogNet 属性 |

组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  
The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ILogNet LogNet { get; set; }
```

```
Public Property LogNet As ILogNet
	Get
	Set
```

```
public:
property ILogNet^ LogNet {
	ILogNet^ get ();
	void set (ILogNet^ value);
}
```

```
member LogNet : ILogNet with get, set
```

#### 属性值

类型：[ILogNet](d95d4704-db0d-b379-41e6-03879927a543.htm)

![](../icons/SectionExpanded.png)备注

只要实例化即可以记录日志，实例化的对象需要实现接口 [ILogNet](d95d4704-db0d-b379-41e6-03879927a543.htm) ，本组件提供了三个日志记录类，你可以实现基于 [ILogNet](d95d4704-db0d-b379-41e6-03879927a543.htm) 的对象。

![](../icons/SectionExpanded.png)示例

如下的实例化适用于所有的Network及其派生类，以下举两个例子，三菱的设备类及服务器类

LogNet示例

[复制](# "复制")

```
// 设备连接对象的日志
MelsecMcNet melsec = new MelsecMcNet( "192.168.0.100", 6000 );

// 举例实现日志文件为单日志文件
melsec.LogNet = new HslCommunication.LogNet.LogNetSingle( "D://123.txt" );
```

LogNet示例

[复制](# "复制")

```
// 一般服务器对象的
NetSimplifyServer simplifyServer = new NetSimplifyServer( );
simplifyServer.LogNet = new HslCommunication.LogNet.LogNetSingle( "D://log.txt" );
simplifyServer.ReceiveStringEvent += ( HslCommunication.Core.Net.AppSession session, HslCommunication.NetHandle handle, string data ) =>
{
    simplifyServer.SendMessage( session, handle, "Back:" + data );
};
simplifyServer.ServerStart( 45678 );
```

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Port 属性 

[原文連結](http://api.hslcommunication.cn/html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[AutoDecompress 属性](../html/0c7971d1-9dca-6855-5751-46a8935781a8.htm "AutoDecompress 属性 ")

[ConnectTimeOut 属性](../html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm "ConnectTimeOut 属性 ")

[GetCarryHostAndPort 属性](../html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm "GetCarryHostAndPort 属性 ")

[IpAddress 属性](../html/9aba5840-5dda-e642-b3c0-4546344607fb.htm "IpAddress 属性 ")

[IsClosed 属性](../html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm "IsClosed 属性 ")

[LogNet 属性](../html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm "LogNet 属性 ")

[Port 属性](../html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm "Port 属性 ")

[SupportDeflate 属性](../html/57017a28-6348-bd7b-80dd-502558db693a.htm "SupportDeflate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientPort 属性 |

端口号。默认1883  
The port number. Default 1883

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int Port { get; set; }
```

```
Public Property Port As Integer
	Get
	Set
```

```
public:
property int Port {
	int get ();
	void set (int value);
}
```

```
member Port : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SupportDeflate 属性 

[原文連結](http://api.hslcommunication.cn/html/57017a28-6348-bd7b-80dd-502558db693a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 属性](../html/2dd721ee-7a66-5d35-26c2-ce669b614a1b.htm "WebSocketClient 属性")

[AutoDecompress 属性](../html/0c7971d1-9dca-6855-5751-46a8935781a8.htm "AutoDecompress 属性 ")

[ConnectTimeOut 属性](../html/6d811cfb-79cc-a484-9c76-1db7131b629f.htm "ConnectTimeOut 属性 ")

[GetCarryHostAndPort 属性](../html/28ae13f2-6772-45bf-195c-5fcee79fb431.htm "GetCarryHostAndPort 属性 ")

[IpAddress 属性](../html/9aba5840-5dda-e642-b3c0-4546344607fb.htm "IpAddress 属性 ")

[IsClosed 属性](../html/d84dc9d8-c23b-1e71-adaa-026921fd6844.htm "IsClosed 属性 ")

[LogNet 属性](../html/04e8ead6-533c-4b51-c58c-33f0572b80eb.htm "LogNet 属性 ")

[Port 属性](../html/880da6be-6ca9-c507-f4c8-59f3f53d119f.htm "Port 属性 ")

[SupportDeflate 属性](../html/57017a28-6348-bd7b-80dd-502558db693a.htm "SupportDeflate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientSupportDeflate 属性 |

获取或设置是否支持数据压缩功能，默认支持压缩功能，如果服务器不支持压缩功能，请设置为 false。  
Get or set whether data compression function is supported. By default, compression function is supported. If the server does not support compression function, please set it to false.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool SupportDeflate { get; set; }
```

```
Public Property SupportDeflate As Boolean
	Get
	Set
```

```
public:
property bool SupportDeflate {
	bool get ();
	void set (bool value);
}
```

```
member SupportDeflate : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WebSocketClient 方法

[原文連結](http://api.hslcommunication.cn/html/422964ac-706f-1646-0f73-ef356b611dec.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[ConnectClose 方法](../html/77d60360-6674-11a7-181f-3790417f1ff2.htm "ConnectClose 方法 ")

[ConnectCloseAsync 方法](../html/f8dcad2a-d6f2-a62d-52ac-d9e60c43b039.htm "ConnectCloseAsync 方法 ")

[ConnectServer 方法](../html/97b6f80d-507d-5a37-5d50-d3622a9d90dc.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/414d7aa7-7730-c84e-bf74-44a17c90c92e.htm "ConnectServerAsync 方法 ")

[Dispose 方法](../html/7851c3a0-38ce-70a2-7233-d5c5b330c729.htm "Dispose 方法 ")

[SendServer 方法](../html/d39c4bbe-baa9-15b5-f011-08937629cdd4.htm "SendServer 方法 ")

[ToString 方法](../html/de57374b-55ef-d4e0-7447-ea2ff7e3dd7a.htm "ToString 方法 ")

[UseSSL 方法](../html/87632aa9-9f3e-2f85-d194-6e02895af793.htm "UseSSL 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClient 方法 |

[WebSocketClient](28c4182b-6004-535e-ffa7-d7352570e9ea.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ConnectClose](77d60360-6674-11a7-181f-3790417f1ff2.htm) | 关闭与Websocket服务器的连接。  Close the connection to the Mqtt server. |
| 公共方法 | [ConnectCloseAsync](f8dcad2a-d6f2-a62d-52ac-d9e60c43b039.htm) | 关闭与Websocket服务器的连接。  Close the connection to the Mqtt server. |
| 公共方法 | [ConnectServer](d50494a4-c56e-e78d-4564-aa8fc4d014bb.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServer(String)](1239b330-9b8f-046e-1e48-f36226b94c3d.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServerAsync](e80ad7fe-314c-7b3f-ee1c-27ebdac8a30c.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServerAsync(String)](ad11585d-4f8b-c6c1-bc56-a20519b925e6.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [Dispose](1266d90f-2ca3-e8ca-8de5-6e1a1a06764d.htm) |  |
| 受保护的方法 | [Dispose(Boolean)](07f93192-03c2-0a68-c4e3-e9f948ec46c4.htm) | 释放当前的对象 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [SendServer(String)](7348ddc5-9af8-8fa3-98ff-08e626436673.htm) | 发送数据到WebSocket的服务器  Send data to WebSocket server |
| 公共方法 | [SendServer(Boolean, String)](37d7ec3b-3ad5-3891-b672-0d35be806d0f.htm) | 发送数据到WebSocket的服务器，可以指定是否进行掩码操作  Send data to the WebSocket server, you can specify whether to perform a mask operation |
| 公共方法 | [SendServer(Int32, Boolean, Byte)](f2434409-540a-3096-f13b-c4ccb6f82872.htm) | 发送自定义的命令到WebSocket服务器，可以指定操作码，是否掩码操作，原始字节数据  Send custom commands to the WebSocket server, you can specify the operation code, whether to mask operation, raw byte data |
| 公共方法 | [ToString](de57374b-55ef-d4e0-7447-ea2ff7e3dd7a.htm) | (重写 ObjectToString.) |
| 公共方法 | [UseSSL(X509Certificate)](484fce1a-01da-3867-80a3-88f2c48e2513.htm) | 使用一个证书来初始化 SSL/TLS 通信  Use a certificate to initiate SSL/TLS communication |
| 公共方法 | [UseSSL(String)](805c7a5c-179b-0e6f-1139-3cf539c11bb8.htm) | 使用一个证书路径来初始化 SSL/TLS 通信  Use a certificate path to initialize SSL/TLS communication |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectClose 方法 

[原文連結](http://api.hslcommunication.cn/html/77d60360-6674-11a7-181f-3790417f1ff2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[ConnectClose 方法](../html/77d60360-6674-11a7-181f-3790417f1ff2.htm "ConnectClose 方法 ")

[ConnectCloseAsync 方法](../html/f8dcad2a-d6f2-a62d-52ac-d9e60c43b039.htm "ConnectCloseAsync 方法 ")

[ConnectServer 方法](../html/97b6f80d-507d-5a37-5d50-d3622a9d90dc.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/414d7aa7-7730-c84e-bf74-44a17c90c92e.htm "ConnectServerAsync 方法 ")

[Dispose 方法](../html/7851c3a0-38ce-70a2-7233-d5c5b330c729.htm "Dispose 方法 ")

[SendServer 方法](../html/d39c4bbe-baa9-15b5-f011-08937629cdd4.htm "SendServer 方法 ")

[ToString 方法](../html/de57374b-55ef-d4e0-7447-ea2ff7e3dd7a.htm "ToString 方法 ")

[UseSSL 方法](../html/87632aa9-9f3e-2f85-d194-6e02895af793.htm "UseSSL 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientConnectClose 方法 |

关闭与Websocket服务器的连接。  
Close the connection to the Mqtt server.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void ConnectClose()
```

```
Public Sub ConnectClose
```

```
public:
void ConnectClose()
```

```
member ConnectClose : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectCloseAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/f8dcad2a-d6f2-a62d-52ac-d9e60c43b039.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[ConnectClose 方法](../html/77d60360-6674-11a7-181f-3790417f1ff2.htm "ConnectClose 方法 ")

[ConnectCloseAsync 方法](../html/f8dcad2a-d6f2-a62d-52ac-d9e60c43b039.htm "ConnectCloseAsync 方法 ")

[ConnectServer 方法](../html/97b6f80d-507d-5a37-5d50-d3622a9d90dc.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/414d7aa7-7730-c84e-bf74-44a17c90c92e.htm "ConnectServerAsync 方法 ")

[Dispose 方法](../html/7851c3a0-38ce-70a2-7233-d5c5b330c729.htm "Dispose 方法 ")

[SendServer 方法](../html/d39c4bbe-baa9-15b5-f011-08937629cdd4.htm "SendServer 方法 ")

[ToString 方法](../html/de57374b-55ef-d4e0-7447-ea2ff7e3dd7a.htm "ToString 方法 ")

[UseSSL 方法](../html/87632aa9-9f3e-2f85-d194-6e02895af793.htm "UseSSL 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientConnectCloseAsync 方法 |

关闭与Websocket服务器的连接。  
Close the connection to the Mqtt server.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task ConnectCloseAsync()
```

```
Public Function ConnectCloseAsync As Task
```

```
public:
Task^ ConnectCloseAsync()
```

```
member ConnectCloseAsync : unit -> Task 
```

#### 返回值

类型：Task  

[缺少 "M:HslCommunication.WebSocket.WebSocketClient.ConnectCloseAsync" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectServer 方法 

[原文連結](http://api.hslcommunication.cn/html/97b6f80d-507d-5a37-5d50-d3622a9d90dc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[ConnectServer 方法](../html/97b6f80d-507d-5a37-5d50-d3622a9d90dc.htm "ConnectServer 方法 ")

[ConnectServer 方法](../html/d50494a4-c56e-e78d-4564-aa8fc4d014bb.htm "ConnectServer 方法 ")

[ConnectServer 方法 (String[])](../html/1239b330-9b8f-046e-1e48-f36226b94c3d.htm "ConnectServer 方法 (String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientConnectServer 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ConnectServer](d50494a4-c56e-e78d-4564-aa8fc4d014bb.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServer(String)](1239b330-9b8f-046e-1e48-f36226b94c3d.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectServer 方法 

[原文連結](http://api.hslcommunication.cn/html/d50494a4-c56e-e78d-4564-aa8fc4d014bb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[ConnectServer 方法](../html/97b6f80d-507d-5a37-5d50-d3622a9d90dc.htm "ConnectServer 方法 ")

[ConnectServer 方法](../html/d50494a4-c56e-e78d-4564-aa8fc4d014bb.htm "ConnectServer 方法 ")

[ConnectServer 方法 (String[])](../html/1239b330-9b8f-046e-1e48-f36226b94c3d.htm "ConnectServer 方法 (String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientConnectServer 方法 |

连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  
After connecting to the server, the client must be called at least once after instantiating the client.
If the return fails, please call this method to connect again after a period of time.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult ConnectServer()
```

```
Public Function ConnectServer As OperateResult
```

```
public:
OperateResult^ ConnectServer()
```

```
member ConnectServer : unit -> OperateResult 
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
连接是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[ConnectServer 重载](97b6f80d-507d-5a37-5d50-d3622a9d90dc.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectServer 方法 (String[])

[原文連結](http://api.hslcommunication.cn/html/1239b330-9b8f-046e-1e48-f36226b94c3d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[ConnectServer 方法](../html/97b6f80d-507d-5a37-5d50-d3622a9d90dc.htm "ConnectServer 方法 ")

[ConnectServer 方法](../html/d50494a4-c56e-e78d-4564-aa8fc4d014bb.htm "ConnectServer 方法 ")

[ConnectServer 方法 (String[])](../html/1239b330-9b8f-046e-1e48-f36226b94c3d.htm "ConnectServer 方法 (String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientConnectServer 方法 (String) |

连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  
After connecting to the server, the client must be called at least once after instantiating the client.
If the return fails, please call this method to connect again after a period of time.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult ConnectServer(
	string[] subscribes
)
```

```
Public Function ConnectServer ( 
	subscribes As String()
) As OperateResult
```

```
public:
OperateResult^ ConnectServer(
	array<String^>^ subscribes
)
```

```
member ConnectServer : 
        subscribes : string[] -> OperateResult 
```

#### 参数

subscribes
:   类型：SystemString  
    订阅的消息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
连接是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[ConnectServer 重载](97b6f80d-507d-5a37-5d50-d3622a9d90dc.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectServerAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/414d7aa7-7730-c84e-bf74-44a17c90c92e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[ConnectServerAsync 方法](../html/414d7aa7-7730-c84e-bf74-44a17c90c92e.htm "ConnectServerAsync 方法 ")

[ConnectServerAsync 方法](../html/e80ad7fe-314c-7b3f-ee1c-27ebdac8a30c.htm "ConnectServerAsync 方法 ")

[ConnectServerAsync 方法 (String[])](../html/ad11585d-4f8b-c6c1-bc56-a20519b925e6.htm "ConnectServerAsync 方法 (String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientConnectServerAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ConnectServerAsync](e80ad7fe-314c-7b3f-ee1c-27ebdac8a30c.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServerAsync(String)](ad11585d-4f8b-c6c1-bc56-a20519b925e6.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectServerAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/e80ad7fe-314c-7b3f-ee1c-27ebdac8a30c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[ConnectServerAsync 方法](../html/414d7aa7-7730-c84e-bf74-44a17c90c92e.htm "ConnectServerAsync 方法 ")

[ConnectServerAsync 方法](../html/e80ad7fe-314c-7b3f-ee1c-27ebdac8a30c.htm "ConnectServerAsync 方法 ")

[ConnectServerAsync 方法 (String[])](../html/ad11585d-4f8b-c6c1-bc56-a20519b925e6.htm "ConnectServerAsync 方法 (String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientConnectServerAsync 方法 |

连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  
After connecting to the server, the client must be called at least once after instantiating the client.
If the return fails, please call this method to connect again after a period of time.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> ConnectServerAsync()
```

```
Public Function ConnectServerAsync As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ ConnectServerAsync()
```

```
member ConnectServerAsync : unit -> Task<OperateResult> 
```

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
连接是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[ConnectServerAsync 重载](414d7aa7-7730-c84e-bf74-44a17c90c92e.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectServerAsync 方法 (String[])

[原文連結](http://api.hslcommunication.cn/html/ad11585d-4f8b-c6c1-bc56-a20519b925e6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[ConnectServerAsync 方法](../html/414d7aa7-7730-c84e-bf74-44a17c90c92e.htm "ConnectServerAsync 方法 ")

[ConnectServerAsync 方法](../html/e80ad7fe-314c-7b3f-ee1c-27ebdac8a30c.htm "ConnectServerAsync 方法 ")

[ConnectServerAsync 方法 (String[])](../html/ad11585d-4f8b-c6c1-bc56-a20519b925e6.htm "ConnectServerAsync 方法 (String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientConnectServerAsync 方法 (String) |

连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  
After connecting to the server, the client must be called at least once after instantiating the client.
If the return fails, please call this method to connect again after a period of time.

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> ConnectServerAsync(
	string[] subscribes
)
```

```
Public Function ConnectServerAsync ( 
	subscribes As String()
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ ConnectServerAsync(
	array<String^>^ subscribes
)
```

```
member ConnectServerAsync : 
        subscribes : string[] -> Task<OperateResult> 
```

#### 参数

subscribes
:   类型：SystemString  
    订阅的消息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
连接是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[ConnectServerAsync 重载](414d7aa7-7730-c84e-bf74-44a17c90c92e.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/7851c3a0-38ce-70a2-7233-d5c5b330c729.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[Dispose 方法](../html/7851c3a0-38ce-70a2-7233-d5c5b330c729.htm "Dispose 方法 ")

[Dispose 方法](../html/1266d90f-2ca3-e8ca-8de5-6e1a1a06764d.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/07f93192-03c2-0a68-c4e3-e9f948ec46c4.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientDispose 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Dispose](1266d90f-2ca3-e8ca-8de5-6e1a1a06764d.htm) | 释放被 [WebSocketClient](28c4182b-6004-535e-ffa7-d7352570e9ea.htm) 使用的所有资源 |
| 受保护的方法 | [Dispose(Boolean)](07f93192-03c2-0a68-c4e3-e9f948ec46c4.htm) | 释放当前的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/1266d90f-2ca3-e8ca-8de5-6e1a1a06764d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[Dispose 方法](../html/7851c3a0-38ce-70a2-7233-d5c5b330c729.htm "Dispose 方法 ")

[Dispose 方法](../html/1266d90f-2ca3-e8ca-8de5-6e1a1a06764d.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/07f93192-03c2-0a68-c4e3-e9f948ec46c4.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientDispose 方法 |

释放被 [WebSocketClient](28c4182b-6004-535e-ffa7-d7352570e9ea.htm) 使用的所有资源

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
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

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[Dispose 重载](7851c3a0-38ce-70a2-7233-d5c5b330c729.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 (Boolean)

[原文連結](http://api.hslcommunication.cn/html/07f93192-03c2-0a68-c4e3-e9f948ec46c4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[Dispose 方法](../html/7851c3a0-38ce-70a2-7233-d5c5b330c729.htm "Dispose 方法 ")

[Dispose 方法](../html/1266d90f-2ca3-e8ca-8de5-6e1a1a06764d.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/07f93192-03c2-0a68-c4e3-e9f948ec46c4.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientDispose 方法 (Boolean) |

释放当前的对象

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
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

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[Dispose 重载](7851c3a0-38ce-70a2-7233-d5c5b330c729.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SendServer 方法 

[原文連結](http://api.hslcommunication.cn/html/d39c4bbe-baa9-15b5-f011-08937629cdd4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[SendServer 方法](../html/d39c4bbe-baa9-15b5-f011-08937629cdd4.htm "SendServer 方法 ")

[SendServer 方法 (String)](../html/7348ddc5-9af8-8fa3-98ff-08e626436673.htm "SendServer 方法 (String)")

[SendServer 方法 (Boolean, String)](../html/37d7ec3b-3ad5-3891-b672-0d35be806d0f.htm "SendServer 方法 (Boolean, String)")

[SendServer 方法 (Int32, Boolean, Byte[])](../html/f2434409-540a-3096-f13b-c4ccb6f82872.htm "SendServer 方法 (Int32, Boolean, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientSendServer 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SendServer(String)](7348ddc5-9af8-8fa3-98ff-08e626436673.htm) | 发送数据到WebSocket的服务器  Send data to WebSocket server |
| 公共方法 | [SendServer(Boolean, String)](37d7ec3b-3ad5-3891-b672-0d35be806d0f.htm) | 发送数据到WebSocket的服务器，可以指定是否进行掩码操作  Send data to the WebSocket server, you can specify whether to perform a mask operation |
| 公共方法 | [SendServer(Int32, Boolean, Byte)](f2434409-540a-3096-f13b-c4ccb6f82872.htm) | 发送自定义的命令到WebSocket服务器，可以指定操作码，是否掩码操作，原始字节数据  Send custom commands to the WebSocket server, you can specify the operation code, whether to mask operation, raw byte data |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SendServer 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/7348ddc5-9af8-8fa3-98ff-08e626436673.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[SendServer 方法](../html/d39c4bbe-baa9-15b5-f011-08937629cdd4.htm "SendServer 方法 ")

[SendServer 方法 (String)](../html/7348ddc5-9af8-8fa3-98ff-08e626436673.htm "SendServer 方法 (String)")

[SendServer 方法 (Boolean, String)](../html/37d7ec3b-3ad5-3891-b672-0d35be806d0f.htm "SendServer 方法 (Boolean, String)")

[SendServer 方法 (Int32, Boolean, Byte[])](../html/f2434409-540a-3096-f13b-c4ccb6f82872.htm "SendServer 方法 (Int32, Boolean, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientSendServer 方法 (String) |

发送数据到WebSocket的服务器  
Send data to WebSocket server

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult SendServer(
	string message
)
```

```
Public Function SendServer ( 
	message As String
) As OperateResult
```

```
public:
OperateResult^ SendServer(
	String^ message
)
```

```
member SendServer : 
        message : string -> OperateResult 
```

#### 参数

message
:   类型：SystemString  
    消息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[SendServer 重载](d39c4bbe-baa9-15b5-f011-08937629cdd4.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SendServer 方法 (Boolean, String)

[原文連結](http://api.hslcommunication.cn/html/37d7ec3b-3ad5-3891-b672-0d35be806d0f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.WebSocket](../html/3ed09122-1aec-556c-4a95-abaff5a02c7a.htm "HslCommunication.WebSocket")

[WebSocketClient 类](../html/28c4182b-6004-535e-ffa7-d7352570e9ea.htm "WebSocketClient 类")

[WebSocketClient 方法](../html/422964ac-706f-1646-0f73-ef356b611dec.htm "WebSocketClient 方法")

[SendServer 方法](../html/d39c4bbe-baa9-15b5-f011-08937629cdd4.htm "SendServer 方法 ")

[SendServer 方法 (String)](../html/7348ddc5-9af8-8fa3-98ff-08e626436673.htm "SendServer 方法 (String)")

[SendServer 方法 (Boolean, String)](../html/37d7ec3b-3ad5-3891-b672-0d35be806d0f.htm "SendServer 方法 (Boolean, String)")

[SendServer 方法 (Int32, Boolean, Byte[])](../html/f2434409-540a-3096-f13b-c4ccb6f82872.htm "SendServer 方法 (Int32, Boolean, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| WebSocketClientSendServer 方法 (Boolean, String) |

发送数据到WebSocket的服务器，可以指定是否进行掩码操作  
Send data to the WebSocket server, you can specify whether to perform a mask operation

**命名空间：**
 [HslCommunication.WebSocket](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult SendServer(
	bool mask,
	string message
)
```

```
Public Function SendServer ( 
	mask As Boolean,
	message As String
) As OperateResult
```

```
public:
OperateResult^ SendServer(
	bool mask, 
	String^ message
)
```

```
member SendServer : 
        mask : bool * 
        message : string -> OperateResult 
```

#### 参数

mask
:   类型：SystemBoolean  
    是否进行掩码操作

message
:   类型：SystemString  
    消息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

![](../icons/SectionExpanded.png)参见

#### 引用

[WebSocketClient 类](28c4182b-6004-535e-ffa7-d7352570e9ea.htm)

[SendServer 重载](d39c4bbe-baa9-15b5-f011-08937629cdd4.htm)

[HslCommunication.WebSocket 命名空间](3ed09122-1aec-556c-4a95-abaff5a02c7a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)