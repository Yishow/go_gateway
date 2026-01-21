# HslCommunication - HslCommunication.Robot.YASKAWA.Helper

> 分類頁數: 12



---
## HslCommunication.Robot.YASKAWA.Helper

[原文連結](http://api.hslcommunication.cn/html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHelper 类](../html/8523c1c6-0016-c48e-53c0-8fb85588da78.htm "YRCHelper 类")

[YRCHighEthernetHelper 类](../html/8237ed72-337e-1003-e8b8-d3113335c0b5.htm "YRCHighEthernetHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Robot.YASKAWA.Helper 命名空间 |

[缺少 "N:HslCommunication.Robot.YASKAWA.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [YRCHelper](8523c1c6-0016-c48e-53c0-8fb85588da78.htm) | 安川机器人的静态辅助方法 |
| 公共类 | [YRCHighEthernetHelper](8237ed72-337e-1003-e8b8-d3113335c0b5.htm) | 安川机器人的高速以太网的辅助方法 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRCHelper 类

[原文連結](http://api.hslcommunication.cn/html/8523c1c6-0016-c48e-53c0-8fb85588da78.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHelper 类](../html/8523c1c6-0016-c48e-53c0-8fb85588da78.htm "YRCHelper 类")

[YRCHelper 构造函数](../html/12cfaf20-f4fd-3b55-f9f3-748953409efb.htm "YRCHelper 构造函数 ")

[YRCHelper 方法](../html/e0a9d014-a64c-472c-bfab-56e05f53fb7c.htm "YRCHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHelper 类 |

安川机器人的静态辅助方法

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Robot.YASKAWA.HelperYRCHelper

**命名空间：**
 [HslCommunication.Robot.YASKAWA.Helper](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class YRCHelper
```

```
Public Class YRCHelper
```

```
public ref class YRCHelper
```

```
type YRCHelper =  class end
```

YRCHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [YRCHelper](12cfaf20-f4fd-3b55-f9f3-748953409efb.htm) | 初始化 YRCHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraErrorMessage](6990777a-6d4e-2f06-33dd-662c0fd96d29.htm) | 当机器人返回ERROR的错误指令后，检测消息里面是否有相关的错误码数据，如果存在，就解析出错误对应的文本  When the robot returns the error instruction of ERROR, it checks whether there is related error code data in the message, and if it exists, it parses out the text corresponding to the error |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorMessage](e5b5293f-4387-967d-9041-e61471032796.htm) | 根据错误信息获取安川机器人的错误信息文本  Obtain the error message text of the Yaskawa robot according to the error message |
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

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRCHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/12cfaf20-f4fd-3b55-f9f3-748953409efb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHelper 类](../html/8523c1c6-0016-c48e-53c0-8fb85588da78.htm "YRCHelper 类")

[YRCHelper 构造函数](../html/12cfaf20-f4fd-3b55-f9f3-748953409efb.htm "YRCHelper 构造函数 ")

[YRCHelper 方法](../html/e0a9d014-a64c-472c-bfab-56e05f53fb7c.htm "YRCHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHelper 构造函数 |

初始化 [YRCHelper](8523c1c6-0016-c48e-53c0-8fb85588da78.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Robot.YASKAWA.Helper](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public YRCHelper()
```

```
Public Sub New
```

```
public:
YRCHelper()
```

```
new : unit -> YRCHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[YRCHelper 类](8523c1c6-0016-c48e-53c0-8fb85588da78.htm)

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRCHelper 方法

[原文連結](http://api.hslcommunication.cn/html/e0a9d014-a64c-472c-bfab-56e05f53fb7c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHelper 类](../html/8523c1c6-0016-c48e-53c0-8fb85588da78.htm "YRCHelper 类")

[YRCHelper 方法](../html/e0a9d014-a64c-472c-bfab-56e05f53fb7c.htm "YRCHelper 方法")

[ExtraErrorMessage 方法](../html/6990777a-6d4e-2f06-33dd-662c0fd96d29.htm "ExtraErrorMessage 方法 ")

[GetErrorMessage 方法](../html/e5b5293f-4387-967d-9041-e61471032796.htm "GetErrorMessage 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHelper 方法 |

[YRCHelper](8523c1c6-0016-c48e-53c0-8fb85588da78.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraErrorMessage](6990777a-6d4e-2f06-33dd-662c0fd96d29.htm) | 当机器人返回ERROR的错误指令后，检测消息里面是否有相关的错误码数据，如果存在，就解析出错误对应的文本  When the robot returns the error instruction of ERROR, it checks whether there is related error code data in the message, and if it exists, it parses out the text corresponding to the error |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorMessage](e5b5293f-4387-967d-9041-e61471032796.htm) | 根据错误信息获取安川机器人的错误信息文本  Obtain the error message text of the Yaskawa robot according to the error message |
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

[YRCHelper 类](8523c1c6-0016-c48e-53c0-8fb85588da78.htm)

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraErrorMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/6990777a-6d4e-2f06-33dd-662c0fd96d29.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHelper 类](../html/8523c1c6-0016-c48e-53c0-8fb85588da78.htm "YRCHelper 类")

[YRCHelper 方法](../html/e0a9d014-a64c-472c-bfab-56e05f53fb7c.htm "YRCHelper 方法")

[ExtraErrorMessage 方法](../html/6990777a-6d4e-2f06-33dd-662c0fd96d29.htm "ExtraErrorMessage 方法 ")

[GetErrorMessage 方法](../html/e5b5293f-4387-967d-9041-e61471032796.htm "GetErrorMessage 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHelperExtraErrorMessage 方法 |

当机器人返回ERROR的错误指令后，检测消息里面是否有相关的错误码数据，如果存在，就解析出错误对应的文本  
When the robot returns the error instruction of ERROR, it checks whether there is related error code data in the message,
and if it exists, it parses out the text corresponding to the error

**命名空间：**
 [HslCommunication.Robot.YASKAWA.Helper](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> ExtraErrorMessage(
	string errText
)
```

```
Public Shared Function ExtraErrorMessage ( 
	errText As String
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ ExtraErrorMessage(
	String^ errText
)
```

```
static member ExtraErrorMessage : 
        errText : string -> OperateResult<string> 
```

#### 参数

errText
:   类型：SystemString  
    返回的完整的报文

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有错误文本的数据信息

![](../icons/SectionExpanded.png)参见

#### 引用

[YRCHelper 类](8523c1c6-0016-c48e-53c0-8fb85588da78.htm)

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/e5b5293f-4387-967d-9041-e61471032796.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHelper 类](../html/8523c1c6-0016-c48e-53c0-8fb85588da78.htm "YRCHelper 类")

[YRCHelper 方法](../html/e0a9d014-a64c-472c-bfab-56e05f53fb7c.htm "YRCHelper 方法")

[ExtraErrorMessage 方法](../html/6990777a-6d4e-2f06-33dd-662c0fd96d29.htm "ExtraErrorMessage 方法 ")

[GetErrorMessage 方法](../html/e5b5293f-4387-967d-9041-e61471032796.htm "GetErrorMessage 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHelperGetErrorMessage 方法 |

根据错误信息获取安川机器人的错误信息文本  
Obtain the error message text of the Yaskawa robot according to the error message

**命名空间：**
 [HslCommunication.Robot.YASKAWA.Helper](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetErrorMessage(
	int err
)
```

```
Public Shared Function GetErrorMessage ( 
	err As Integer
) As String
```

```
public:
static String^ GetErrorMessage(
	int err
)
```

```
static member GetErrorMessage : 
        err : int -> string 
```

#### 参数

err
:   类型：SystemInt32  
    错误号

#### 返回值

类型：String  
错误文本信息

![](../icons/SectionExpanded.png)参见

#### 引用

[YRCHelper 类](8523c1c6-0016-c48e-53c0-8fb85588da78.htm)

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRCHighEthernetHelper 类

[原文連結](http://api.hslcommunication.cn/html/8237ed72-337e-1003-e8b8-d3113335c0b5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHighEthernetHelper 类](../html/8237ed72-337e-1003-e8b8-d3113335c0b5.htm "YRCHighEthernetHelper 类")

[YRCHighEthernetHelper 构造函数](../html/e66fbc46-5362-7757-f489-c3eb7a2f7459.htm "YRCHighEthernetHelper 构造函数 ")

[YRCHighEthernetHelper 方法](../html/b4751333-6b49-c4ff-e8a5-25cfc16a5bfb.htm "YRCHighEthernetHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHighEthernetHelper 类 |

安川机器人的高速以太网的辅助方法

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Robot.YASKAWA.HelperYRCHighEthernetHelper

**命名空间：**
 [HslCommunication.Robot.YASKAWA.Helper](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class YRCHighEthernetHelper
```

```
Public Class YRCHighEthernetHelper
```

```
public ref class YRCHighEthernetHelper
```

```
type YRCHighEthernetHelper =  class end
```

YRCHighEthernetHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [YRCHighEthernetHelper](e66fbc46-5362-7757-f489-c3eb7a2f7459.htm) | 初始化 YRCHighEthernetHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildCommand](205dbb30-8adc-3d4c-075c-95a732d9e01b.htm) | 构建完整的读取指令 |
| 公共方法静态成员 | [CheckResponseContent](3a6047f9-70d7-ed78-055d-445614075ffc.htm) | 检查当前的机器人反馈的数据是否正确 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorText](1a78c033-1e42-1621-35fe-7b91f30e7f77.htm) | 根据状态信息及附加状态信息来获取错误的文本描述信息 |
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

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRCHighEthernetHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/e66fbc46-5362-7757-f489-c3eb7a2f7459.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHighEthernetHelper 类](../html/8237ed72-337e-1003-e8b8-d3113335c0b5.htm "YRCHighEthernetHelper 类")

[YRCHighEthernetHelper 构造函数](../html/e66fbc46-5362-7757-f489-c3eb7a2f7459.htm "YRCHighEthernetHelper 构造函数 ")

[YRCHighEthernetHelper 方法](../html/b4751333-6b49-c4ff-e8a5-25cfc16a5bfb.htm "YRCHighEthernetHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHighEthernetHelper 构造函数 |

初始化 [YRCHighEthernetHelper](8237ed72-337e-1003-e8b8-d3113335c0b5.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Robot.YASKAWA.Helper](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public YRCHighEthernetHelper()
```

```
Public Sub New
```

```
public:
YRCHighEthernetHelper()
```

```
new : unit -> YRCHighEthernetHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[YRCHighEthernetHelper 类](8237ed72-337e-1003-e8b8-d3113335c0b5.htm)

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRCHighEthernetHelper 方法

[原文連結](http://api.hslcommunication.cn/html/b4751333-6b49-c4ff-e8a5-25cfc16a5bfb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHighEthernetHelper 类](../html/8237ed72-337e-1003-e8b8-d3113335c0b5.htm "YRCHighEthernetHelper 类")

[YRCHighEthernetHelper 方法](../html/b4751333-6b49-c4ff-e8a5-25cfc16a5bfb.htm "YRCHighEthernetHelper 方法")

[BuildCommand 方法](../html/205dbb30-8adc-3d4c-075c-95a732d9e01b.htm "BuildCommand 方法 ")

[CheckResponseContent 方法](../html/3a6047f9-70d7-ed78-055d-445614075ffc.htm "CheckResponseContent 方法 ")

[GetErrorText 方法](../html/1a78c033-1e42-1621-35fe-7b91f30e7f77.htm "GetErrorText 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHighEthernetHelper 方法 |

[YRCHighEthernetHelper](8237ed72-337e-1003-e8b8-d3113335c0b5.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildCommand](205dbb30-8adc-3d4c-075c-95a732d9e01b.htm) | 构建完整的读取指令 |
| 公共方法静态成员 | [CheckResponseContent](3a6047f9-70d7-ed78-055d-445614075ffc.htm) | 检查当前的机器人反馈的数据是否正确 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorText](1a78c033-1e42-1621-35fe-7b91f30e7f77.htm) | 根据状态信息及附加状态信息来获取错误的文本描述信息 |
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

[YRCHighEthernetHelper 类](8237ed72-337e-1003-e8b8-d3113335c0b5.htm)

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/205dbb30-8adc-3d4c-075c-95a732d9e01b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHighEthernetHelper 类](../html/8237ed72-337e-1003-e8b8-d3113335c0b5.htm "YRCHighEthernetHelper 类")

[YRCHighEthernetHelper 方法](../html/b4751333-6b49-c4ff-e8a5-25cfc16a5bfb.htm "YRCHighEthernetHelper 方法")

[BuildCommand 方法](../html/205dbb30-8adc-3d4c-075c-95a732d9e01b.htm "BuildCommand 方法 ")

[CheckResponseContent 方法](../html/3a6047f9-70d7-ed78-055d-445614075ffc.htm "CheckResponseContent 方法 ")

[GetErrorText 方法](../html/1a78c033-1e42-1621-35fe-7b91f30e7f77.htm "GetErrorText 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHighEthernetHelperBuildCommand 方法 |

构建完整的读取指令

**命名空间：**
 [HslCommunication.Robot.YASKAWA.Helper](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] BuildCommand(
	byte handle,
	byte requestID,
	ushort command,
	ushort dataAddress,
	byte dataAttribute,
	byte dataHandle,
	byte[] dataPart
)
```

```
Public Shared Function BuildCommand ( 
	handle As Byte,
	requestID As Byte,
	command As UShort,
	dataAddress As UShort,
	dataAttribute As Byte,
	dataHandle As Byte,
	dataPart As Byte()
) As Byte()
```

```
public:
static array<unsigned char>^ BuildCommand(
	unsigned char handle, 
	unsigned char requestID, 
	unsigned short command, 
	unsigned short dataAddress, 
	unsigned char dataAttribute, 
	unsigned char dataHandle, 
	array<unsigned char>^ dataPart
)
```

```
static member BuildCommand : 
        handle : byte * 
        requestID : byte * 
        command : uint16 * 
        dataAddress : uint16 * 
        dataAttribute : byte * 
        dataHandle : byte * 
        dataPart : byte[] -> byte[] 
```

#### 参数

handle
:   类型：SystemByte  
    处理分区，1:机器人控制 2:文件控制

requestID
:   类型：SystemByte  
    请求ID， 客户端每次命令输出的时请增量

command
:   类型：SystemUInt16  
    命令编号，相当于CIP通信的CLASS

dataAddress
:   类型：SystemUInt16  
    数据队列编号，相当于CIP通信的Instance

dataAttribute
:   类型：SystemByte  
    单元编号，相当于CIP通信协议的Attribute

dataHandle
:   类型：SystemByte  
    处理请求，定义数据的请方法

dataPart
:   类型：SystemByte  
    数据部分的内容

#### 返回值

类型：Byte  
构建结果

![](../icons/SectionExpanded.png)参见

#### 引用

[YRCHighEthernetHelper 类](8237ed72-337e-1003-e8b8-d3113335c0b5.htm)

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckResponseContent 方法 

[原文連結](http://api.hslcommunication.cn/html/3a6047f9-70d7-ed78-055d-445614075ffc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHighEthernetHelper 类](../html/8237ed72-337e-1003-e8b8-d3113335c0b5.htm "YRCHighEthernetHelper 类")

[YRCHighEthernetHelper 方法](../html/b4751333-6b49-c4ff-e8a5-25cfc16a5bfb.htm "YRCHighEthernetHelper 方法")

[BuildCommand 方法](../html/205dbb30-8adc-3d4c-075c-95a732d9e01b.htm "BuildCommand 方法 ")

[CheckResponseContent 方法](../html/3a6047f9-70d7-ed78-055d-445614075ffc.htm "CheckResponseContent 方法 ")

[GetErrorText 方法](../html/1a78c033-1e42-1621-35fe-7b91f30e7f77.htm "GetErrorText 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHighEthernetHelperCheckResponseContent 方法 |

检查当前的机器人反馈的数据是否正确

**命名空间：**
 [HslCommunication.Robot.YASKAWA.Helper](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult CheckResponseContent(
	byte[] response
)
```

```
Public Shared Function CheckResponseContent ( 
	response As Byte()
) As OperateResult
```

```
public:
static OperateResult^ CheckResponseContent(
	array<unsigned char>^ response
)
```

```
static member CheckResponseContent : 
        response : byte[] -> OperateResult 
```

#### 参数

response
:   类型：SystemByte  
    从机器人反馈的数据

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否检查正确

![](../icons/SectionExpanded.png)参见

#### 引用

[YRCHighEthernetHelper 类](8237ed72-337e-1003-e8b8-d3113335c0b5.htm)

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorText 方法 

[原文連結](http://api.hslcommunication.cn/html/1a78c033-1e42-1621-35fe-7b91f30e7f77.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA.Helper](../html/008672cf-5411-8219-fa26-dfc0fe5e5c94.htm "HslCommunication.Robot.YASKAWA.Helper")

[YRCHighEthernetHelper 类](../html/8237ed72-337e-1003-e8b8-d3113335c0b5.htm "YRCHighEthernetHelper 类")

[YRCHighEthernetHelper 方法](../html/b4751333-6b49-c4ff-e8a5-25cfc16a5bfb.htm "YRCHighEthernetHelper 方法")

[BuildCommand 方法](../html/205dbb30-8adc-3d4c-075c-95a732d9e01b.htm "BuildCommand 方法 ")

[CheckResponseContent 方法](../html/3a6047f9-70d7-ed78-055d-445614075ffc.htm "CheckResponseContent 方法 ")

[GetErrorText 方法](../html/1a78c033-1e42-1621-35fe-7b91f30e7f77.htm "GetErrorText 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRCHighEthernetHelperGetErrorText 方法 |

根据状态信息及附加状态信息来获取错误的文本描述信息

**命名空间：**
 [HslCommunication.Robot.YASKAWA.Helper](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetErrorText(
	byte status,
	int affix
)
```

```
Public Shared Function GetErrorText ( 
	status As Byte,
	affix As Integer
) As String
```

```
public:
static String^ GetErrorText(
	unsigned char status, 
	int affix
)
```

```
static member GetErrorText : 
        status : byte * 
        affix : int -> string 
```

#### 参数

status
:   类型：SystemByte  
    状态信息

affix
:   类型：SystemInt32  
    附加状态信息

#### 返回值

类型：String  
错误的文本描述信息

![](../icons/SectionExpanded.png)参见

#### 引用

[YRCHighEthernetHelper 类](8237ed72-337e-1003-e8b8-d3113335c0b5.htm)

[HslCommunication.Robot.YASKAWA.Helper 命名空间](008672cf-5411-8219-fa26-dfc0fe5e5c94.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)