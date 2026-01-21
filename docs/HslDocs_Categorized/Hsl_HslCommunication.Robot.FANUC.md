# HslCommunication - HslCommunication.Robot.FANUC

> 分類頁數: 30



---
## HslCommunication.Robot.FANUC

[原文連結](http://api.hslcommunication.cn/html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucHelper 类](../html/df938f03-f1d3-a426-8206-df01bc0cea01.htm "FanucHelper 类")

[FanucInterfaceNet 类](../html/9a252da5-4341-0437-0fb7-27da2b49d3c4.htm "FanucInterfaceNet 类")

[FanucPose 类](../html/cdcdf7d8-60e1-ed91-7b3c-b14717cb31e6.htm "FanucPose 类")

[FanucRobotServer 类](../html/a625abd9-9a3f-168a-271b-2654d6b0ee6b.htm "FanucRobotServer 类")

[FanucTask 类](../html/29d5f632-9411-0665-b8e6-a9c9194dc2d5.htm "FanucTask 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Robot.FANUC 命名空间 |

[缺少 "N:HslCommunication.Robot.FANUC" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [FanucAlarm](f669181a-2617-5fda-b19c-e416aa573f3b.htm) | Fanuc机器人的报警对象 |
| 公共类 | [FanucData](fd9d131d-3177-76a1-ab73-e521883434f5.htm) | Fanuc机器人的所有的数据信息 |
| 公共类 | [FanucHelper](df938f03-f1d3-a426-8206-df01bc0cea01.htm) | Fanuc的辅助方法信息 |
| 公共类代码示例 | [FanucInterfaceNet](9a252da5-4341-0437-0fb7-27da2b49d3c4.htm) | Fanuc机器人的PC Interface实现，在R-30iB mate plus型号上测试通过，支持读写任意的数据，写入操作务必谨慎调用，写入数据不当造成生命财产损失，作者概不负责。读写任意的地址见api文档信息  The Fanuc robot's PC Interface implementation has been tested on R-30iB mate plus models. It supports reading and writing arbitrary data. The writing operation must be called carefully. Improper writing of data will cause loss of life and property. The author is not responsible. Read and write arbitrary addresses see api documentation information |
| 公共类 | [FanucPose](cdcdf7d8-60e1-ed91-7b3c-b14717cb31e6.htm) | 机器人的姿态数据 |
| 公共类 | [FanucRobotServer](a625abd9-9a3f-168a-271b-2654d6b0ee6b.htm) | 虚拟的FANUC机器人的服务器对象，支持I,Q,M,D,AI,AQ数据区的数据读写，其中D区是机器人数据存放的区域，相关的数据需要去机器人区读取。详细参见api文档信息。  The server object of the virtual FANUC robot supports data reading and writing in I, Q, M, D, AI, and AQ data areas, where D area is the area where the robot data is stored, and related data needs to be read in the robot area. See the api documentation for details. |
| 公共类 | [FanucTask](29d5f632-9411-0665-b8e6-a9c9194dc2d5.htm) | Fanuc机器人的任务类 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucAlarm 类

[原文連結](http://api.hslcommunication.cn/html/f669181a-2617-5fda-b19c-e416aa573f3b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 构造函数](../html/101a483d-8cc1-b00f-4acc-6c614cd9c325.htm "FanucAlarm 构造函数 ")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[FanucAlarm 方法](../html/41a98067-ebb7-ef37-a75d-8476821d4905.htm "FanucAlarm 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarm 类 |

Fanuc机器人的报警对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Robot.FANUCFanucAlarm

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FanucAlarm
```

```
Public Class FanucAlarm
```

```
public ref class FanucAlarm
```

```
type FanucAlarm =  class end
```

FanucAlarm 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FanucAlarm](101a483d-8cc1-b00f-4acc-6c614cd9c325.htm) | 初始化 FanucAlarm 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AlarmID](f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm) | AlarmID |
| 公共属性 | [AlarmMessage](9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm) | AlarmMessage |
| 公共属性 | [AlarmNumber](96c50845-b1f5-f059-be51-ebc16c193055.htm) | AlarmNumber |
| 公共属性 | [CauseAlarmID](9cb339cc-66cd-21e3-4655-8e52a6044324.htm) | CauseAlarmID |
| 公共属性 | [CauseAlarmMessage](9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm) | CauseAlarmMessage |
| 公共属性 | [CauseAlarmNumber](16559f6e-c155-56a9-db87-833d0c1fc2d0.htm) | CauseAlarmNumber |
| 公共属性 | [Severity](81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm) | Severity |
| 公共属性 | [SeverityMessage](be885209-d3c6-5288-f35b-33acb0ac9ec9.htm) | SeverityMessage |
| 公共属性 | [Time](611a9f82-fd7a-f030-b348-f7e4086ded63.htm) | Time |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [LoadByContent](ac61b4f5-d53a-a20b-5134-a1d178841760.htm) | 从字节数据加载真实的信息 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PraseFrom](52295cde-1f4d-f655-f86b-bb7b5a7102c9.htm) | 从数据内容创建报警信息 |
| 公共方法 | [ToString](23a27ae0-0f73-4bcd-b649-a02833ba9e95.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucAlarm 构造函数 

[原文連結](http://api.hslcommunication.cn/html/101a483d-8cc1-b00f-4acc-6c614cd9c325.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 构造函数](../html/101a483d-8cc1-b00f-4acc-6c614cd9c325.htm "FanucAlarm 构造函数 ")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[FanucAlarm 方法](../html/41a98067-ebb7-ef37-a75d-8476821d4905.htm "FanucAlarm 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarm 构造函数 |

初始化 [FanucAlarm](f669181a-2617-5fda-b19c-e416aa573f3b.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucAlarm()
```

```
Public Sub New
```

```
public:
FanucAlarm()
```

```
new : unit -> FanucAlarm
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucAlarm 属性

[原文連結](http://api.hslcommunication.cn/html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarm 属性 |

[FanucAlarm](f669181a-2617-5fda-b19c-e416aa573f3b.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AlarmID](f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm) | AlarmID |
| 公共属性 | [AlarmMessage](9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm) | AlarmMessage |
| 公共属性 | [AlarmNumber](96c50845-b1f5-f059-be51-ebc16c193055.htm) | AlarmNumber |
| 公共属性 | [CauseAlarmID](9cb339cc-66cd-21e3-4655-8e52a6044324.htm) | CauseAlarmID |
| 公共属性 | [CauseAlarmMessage](9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm) | CauseAlarmMessage |
| 公共属性 | [CauseAlarmNumber](16559f6e-c155-56a9-db87-833d0c1fc2d0.htm) | CauseAlarmNumber |
| 公共属性 | [Severity](81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm) | Severity |
| 公共属性 | [SeverityMessage](be885209-d3c6-5288-f35b-33acb0ac9ec9.htm) | SeverityMessage |
| 公共属性 | [Time](611a9f82-fd7a-f030-b348-f7e4086ded63.htm) | Time |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AlarmID 属性 

[原文連結](http://api.hslcommunication.cn/html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmAlarmID 属性 |

AlarmID

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public short AlarmID { get; set; }
```

```
Public Property AlarmID As Short
	Get
	Set
```

```
public:
property short AlarmID {
	short get ();
	void set (short value);
}
```

```
member AlarmID : int16 with get, set
```

#### 属性值

类型：Int16

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AlarmMessage 属性 

[原文連結](http://api.hslcommunication.cn/html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmAlarmMessage 属性 |

AlarmMessage

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string AlarmMessage { get; set; }
```

```
Public Property AlarmMessage As String
	Get
	Set
```

```
public:
property String^ AlarmMessage {
	String^ get ();
	void set (String^ value);
}
```

```
member AlarmMessage : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AlarmNumber 属性 

[原文連結](http://api.hslcommunication.cn/html/96c50845-b1f5-f059-be51-ebc16c193055.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmAlarmNumber 属性 |

AlarmNumber

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public short AlarmNumber { get; set; }
```

```
Public Property AlarmNumber As Short
	Get
	Set
```

```
public:
property short AlarmNumber {
	short get ();
	void set (short value);
}
```

```
member AlarmNumber : int16 with get, set
```

#### 属性值

类型：Int16

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CauseAlarmID 属性 

[原文連結](http://api.hslcommunication.cn/html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmCauseAlarmID 属性 |

CauseAlarmID

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public short CauseAlarmID { get; set; }
```

```
Public Property CauseAlarmID As Short
	Get
	Set
```

```
public:
property short CauseAlarmID {
	short get ();
	void set (short value);
}
```

```
member CauseAlarmID : int16 with get, set
```

#### 属性值

类型：Int16

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CauseAlarmMessage 属性 

[原文連結](http://api.hslcommunication.cn/html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmCauseAlarmMessage 属性 |

CauseAlarmMessage

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string CauseAlarmMessage { get; set; }
```

```
Public Property CauseAlarmMessage As String
	Get
	Set
```

```
public:
property String^ CauseAlarmMessage {
	String^ get ();
	void set (String^ value);
}
```

```
member CauseAlarmMessage : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CauseAlarmNumber 属性 

[原文連結](http://api.hslcommunication.cn/html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmCauseAlarmNumber 属性 |

CauseAlarmNumber

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public short CauseAlarmNumber { get; set; }
```

```
Public Property CauseAlarmNumber As Short
	Get
	Set
```

```
public:
property short CauseAlarmNumber {
	short get ();
	void set (short value);
}
```

```
member CauseAlarmNumber : int16 with get, set
```

#### 属性值

类型：Int16

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Severity 属性 

[原文連結](http://api.hslcommunication.cn/html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmSeverity 属性 |

Severity

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public short Severity { get; set; }
```

```
Public Property Severity As Short
	Get
	Set
```

```
public:
property short Severity {
	short get ();
	void set (short value);
}
```

```
member Severity : int16 with get, set
```

#### 属性值

类型：Int16

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SeverityMessage 属性 

[原文連結](http://api.hslcommunication.cn/html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmSeverityMessage 属性 |

SeverityMessage

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string SeverityMessage { get; set; }
```

```
Public Property SeverityMessage As String
	Get
	Set
```

```
public:
property String^ SeverityMessage {
	String^ get ();
	void set (String^ value);
}
```

```
member SeverityMessage : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Time 属性 

[原文連結](http://api.hslcommunication.cn/html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 属性](../html/c68f5d73-83ae-8e19-9d55-3ea181067935.htm "FanucAlarm 属性")

[AlarmID 属性](../html/f30c8722-7ef3-d8b3-4a63-939663b3b7b3.htm "AlarmID 属性 ")

[AlarmMessage 属性](../html/9af78d58-0341-e9aa-cdbd-49c27ff4ac72.htm "AlarmMessage 属性 ")

[AlarmNumber 属性](../html/96c50845-b1f5-f059-be51-ebc16c193055.htm "AlarmNumber 属性 ")

[CauseAlarmID 属性](../html/9cb339cc-66cd-21e3-4655-8e52a6044324.htm "CauseAlarmID 属性 ")

[CauseAlarmMessage 属性](../html/9b51dd62-0b1b-d36a-9715-4a6a29d092dd.htm "CauseAlarmMessage 属性 ")

[CauseAlarmNumber 属性](../html/16559f6e-c155-56a9-db87-833d0c1fc2d0.htm "CauseAlarmNumber 属性 ")

[Severity 属性](../html/81de3547-08bb-cbb1-d5a3-76f8e7c0ebfb.htm "Severity 属性 ")

[SeverityMessage 属性](../html/be885209-d3c6-5288-f35b-33acb0ac9ec9.htm "SeverityMessage 属性 ")

[Time 属性](../html/611a9f82-fd7a-f030-b348-f7e4086ded63.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmTime 属性 |

Time

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DateTime Time { get; set; }
```

```
Public Property Time As DateTime
	Get
	Set
```

```
public:
property DateTime Time {
	DateTime get ();
	void set (DateTime value);
}
```

```
member Time : DateTime with get, set
```

#### 属性值

类型：DateTime

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucAlarm 方法

[原文連結](http://api.hslcommunication.cn/html/41a98067-ebb7-ef37-a75d-8476821d4905.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 方法](../html/41a98067-ebb7-ef37-a75d-8476821d4905.htm "FanucAlarm 方法")

[LoadByContent 方法](../html/ac61b4f5-d53a-a20b-5134-a1d178841760.htm "LoadByContent 方法 ")

[PraseFrom 方法](../html/52295cde-1f4d-f655-f86b-bb7b5a7102c9.htm "PraseFrom 方法 ")

[ToString 方法](../html/23a27ae0-0f73-4bcd-b649-a02833ba9e95.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarm 方法 |

[FanucAlarm](f669181a-2617-5fda-b19c-e416aa573f3b.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [LoadByContent](ac61b4f5-d53a-a20b-5134-a1d178841760.htm) | 从字节数据加载真实的信息 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PraseFrom](52295cde-1f4d-f655-f86b-bb7b5a7102c9.htm) | 从数据内容创建报警信息 |
| 公共方法 | [ToString](23a27ae0-0f73-4bcd-b649-a02833ba9e95.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LoadByContent 方法 

[原文連結](http://api.hslcommunication.cn/html/ac61b4f5-d53a-a20b-5134-a1d178841760.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 方法](../html/41a98067-ebb7-ef37-a75d-8476821d4905.htm "FanucAlarm 方法")

[LoadByContent 方法](../html/ac61b4f5-d53a-a20b-5134-a1d178841760.htm "LoadByContent 方法 ")

[PraseFrom 方法](../html/52295cde-1f4d-f655-f86b-bb7b5a7102c9.htm "PraseFrom 方法 ")

[ToString 方法](../html/23a27ae0-0f73-4bcd-b649-a02833ba9e95.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmLoadByContent 方法 |

从字节数据加载真实的信息

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void LoadByContent(
	IByteTransform byteTransform,
	byte[] content,
	int index,
	Encoding encoding
)
```

```
Public Sub LoadByContent ( 
	byteTransform As IByteTransform,
	content As Byte(),
	index As Integer,
	encoding As Encoding
)
```

```
public:
void LoadByContent(
	IByteTransform^ byteTransform, 
	array<unsigned char>^ content, 
	int index, 
	Encoding^ encoding
)
```

```
member LoadByContent : 
        byteTransform : IByteTransform * 
        content : byte[] * 
        index : int * 
        encoding : Encoding -> unit 
```

#### 参数

byteTransform
:   类型：[HslCommunication.CoreIByteTransform](56c55574-bb2a-fe66-e7d5-1332a1bc18f0.htm)  
    字节变换

content
:   类型：SystemByte  
    原始的字节内容

index
:   类型：SystemInt32  
    索引

encoding
:   类型：System.TextEncoding  
    编码

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PraseFrom 方法 

[原文連結](http://api.hslcommunication.cn/html/52295cde-1f4d-f655-f86b-bb7b5a7102c9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 方法](../html/41a98067-ebb7-ef37-a75d-8476821d4905.htm "FanucAlarm 方法")

[LoadByContent 方法](../html/ac61b4f5-d53a-a20b-5134-a1d178841760.htm "LoadByContent 方法 ")

[PraseFrom 方法](../html/52295cde-1f4d-f655-f86b-bb7b5a7102c9.htm "PraseFrom 方法 ")

[ToString 方法](../html/23a27ae0-0f73-4bcd-b649-a02833ba9e95.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmPraseFrom 方法 |

从数据内容创建报警信息

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static FanucAlarm PraseFrom(
	IByteTransform byteTransform,
	byte[] content,
	int index,
	Encoding encoding
)
```

```
Public Shared Function PraseFrom ( 
	byteTransform As IByteTransform,
	content As Byte(),
	index As Integer,
	encoding As Encoding
) As FanucAlarm
```

```
public:
static FanucAlarm^ PraseFrom(
	IByteTransform^ byteTransform, 
	array<unsigned char>^ content, 
	int index, 
	Encoding^ encoding
)
```

```
static member PraseFrom : 
        byteTransform : IByteTransform * 
        content : byte[] * 
        index : int * 
        encoding : Encoding -> FanucAlarm 
```

#### 参数

byteTransform
:   类型：[HslCommunication.CoreIByteTransform](56c55574-bb2a-fe66-e7d5-1332a1bc18f0.htm)  
    字节变换

content
:   类型：SystemByte  
    原始的字节内容

index
:   类型：SystemInt32  
    索引

encoding
:   类型：System.TextEncoding  
    编码

#### 返回值

类型：[FanucAlarm](f669181a-2617-5fda-b19c-e416aa573f3b.htm)  
报警信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/23a27ae0-0f73-4bcd-b649-a02833ba9e95.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucAlarm 类](../html/f669181a-2617-5fda-b19c-e416aa573f3b.htm "FanucAlarm 类")

[FanucAlarm 方法](../html/41a98067-ebb7-ef37-a75d-8476821d4905.htm "FanucAlarm 方法")

[LoadByContent 方法](../html/ac61b4f5-d53a-a20b-5134-a1d178841760.htm "LoadByContent 方法 ")

[PraseFrom 方法](../html/52295cde-1f4d-f655-f86b-bb7b5a7102c9.htm "PraseFrom 方法 ")

[ToString 方法](../html/23a27ae0-0f73-4bcd-b649-a02833ba9e95.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucAlarmToString 方法 |

[缺少 "M:HslCommunication.Robot.FANUC.FanucAlarm.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
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

[缺少 "M:HslCommunication.Robot.FANUC.FanucAlarm.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucAlarm 类](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucData 类

[原文連結](http://api.hslcommunication.cn/html/fd9d131d-3177-76a1-ab73-e521883434f5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 构造函数](../html/c7413752-c310-7fb0-b4a1-2558d28592f3.htm "FanucData 构造函数 ")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[FanucData 方法](../html/8044317d-04e5-1527-cc48-ef1304549afc.htm "FanucData 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucData 类 |

Fanuc机器人的所有的数据信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Robot.FANUCFanucData

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FanucData
```

```
Public Class FanucData
```

```
public ref class FanucData
```

```
type FanucData =  class end
```

FanucData 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FanucData](c7413752-c310-7fb0-b4a1-2558d28592f3.htm) | 初始化 FanucData 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AIComment](caa1ce46-6807-7b6f-0295-822722285f18.htm) |  |
| 公共属性 | [AlarmCurrent](845a692e-e041-8cd2-61d8-739a07890e8a.htm) |  |
| 公共属性 | [AlarmList](4c36388c-6992-1133-51e2-2c197a98663a.htm) |  |
| 公共属性 | [AlarmPassword](7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm) |  |
| 公共属性 | [AOComment](fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm) |  |
| 公共属性 | [CurrentPose](7c03669c-5e00-3769-b016-bb8bb68ba896.htm) |  |
| 公共属性 | [CurrentPose2](0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm) |  |
| 公共属性 | [CurrentPose3](0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm) |  |
| 公共属性 | [CurrentPose4](8d49d24a-422e-8ed5-1cc8-a109bf757671.htm) |  |
| 公共属性 | [CurrentPose5](6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm) |  |
| 公共属性 | [CurrentPoseUF](716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm) |  |
| 公共属性 | [DataPosRegMG](43b750c7-1763-cc14-e6a8-36dcad959b90.htm) |  |
| 公共属性 | [DIComment](c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm) |  |
| 公共属性 | [DOComment](f67f3bd6-7d81-347e-d8bd-92fce6656329.htm) |  |
| 公共属性 | [DUTY\_TEMP](db3365e3-04f5-6a50-79d2-c180364338b9.htm) |  |
| 公共属性 | [FAST\_CLOCK](d85a83dd-ae39-25c3-1d82-39df52aacd72.htm) |  |
| 公共属性 | [GIComment](f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm) |  |
| 公共属性 | [GOComment](cb612119-5f3e-1119-2a34-401c00b48698.htm) |  |
| 公共属性 | [HTTPKCL\_CMDS](46901d3b-052b-1d58-5ef7-b6adae0008c9.htm) |  |
| 公共属性 | [MNUTOOL1\_1](a788bc36-efbe-efc4-21ee-c32af667e9aa.htm) |  |
| 公共属性 | [MOR\_GRP\_CURRENT\_ANG](796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm) |  |
| 公共属性 | [NumReg1](e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm) |  |
| 公共属性 | [NumReg2](17a87be2-0d12-3387-e20c-93f849ccd2db.htm) |  |
| 公共属性 | [PosRegGP1](fe3d3b49-ab48-18db-b360-7f307188a61a.htm) |  |
| 公共属性 | [PosRegGP2](4b0d2c25-1013-0a06-b131-a3146b27a847.htm) |  |
| 公共属性 | [PosRegGP3](ab0cd01e-908b-d86f-c228-da78ee470517.htm) |  |
| 公共属性 | [PosRegGP4](7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm) |  |
| 公共属性 | [PosRegGP5](179e2f5d-5715-57de-2330-cbaadcf200e5.htm) |  |
| 公共属性 | [RIComment](0140fd93-e229-3764-febd-3dd4fa151911.htm) |  |
| 公共属性 | [ROComment](33a5de78-9dba-c2ff-9d31-8c2943a34512.htm) |  |
| 公共属性 | [SIComment](9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm) |  |
| 公共属性 | [SOComment](3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm) |  |
| 公共属性 | [STRREG\_COMMENT\_Comment](5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm) |  |
| 公共属性 | [STRREGComment](46e7f07e-7631-a95c-4875-9576a6684559.htm) |  |
| 公共属性 | [Task](a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm) |  |
| 公共属性 | [TaskIgnoreKarel](5052b3ff-a907-c384-8f24-94fd3974443c.htm) |  |
| 公共属性 | [TaskIgnoreMacro](8502d4b6-23f5-08ab-2575-6374c51c54b8.htm) |  |
| 公共属性 | [TaskIgnoreMacroKarel](96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm) |  |
| 公共属性 | [TIMER10\_COMMENT](653ccde8-6961-7266-4164-b908aa4f5499.htm) |  |
| 公共属性 | [Timer10\_TIMER\_VAL](5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm) |  |
| 公共属性 | [TIMER2\_COMMENT](34698164-5021-1efc-6aaa-14dc3a5f8c58.htm) |  |
| 公共属性 | [UIComment](e9931afb-18e5-50a6-e891-d5c0372a831b.htm) |  |
| 公共属性 | [UOComment](72f57607-083d-ab26-aaa4-e1ec36278ca6.htm) |  |
| 公共属性 | [WIComment](b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm) |  |
| 公共属性 | [WOComment](f70478df-40cb-97aa-3245-41e98674d171.htm) |  |
| 公共属性 | [WSIComment](5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm) |  |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [LoadByContent](3cfd607d-4c22-e502-ed73-89b99c0675ec.htm) | 从原始的数据内容加载数据 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PraseFrom](69a3664c-9da8-8eaa-66c5-dedb244728ed.htm) | 从字节数组解析出fanuc的数据信息 |
| 公共方法 | [ToString](5db25a38-3066-4e02-a92b-e18dd2a3642e.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c7413752-c310-7fb0-b4a1-2558d28592f3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 构造函数](../html/c7413752-c310-7fb0-b4a1-2558d28592f3.htm "FanucData 构造函数 ")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[FanucData 方法](../html/8044317d-04e5-1527-cc48-ef1304549afc.htm "FanucData 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucData 构造函数 |

初始化 [FanucData](fd9d131d-3177-76a1-ab73-e521883434f5.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucData()
```

```
Public Sub New
```

```
public:
FanucData()
```

```
new : unit -> FanucData
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucData 属性

[原文連結](http://api.hslcommunication.cn/html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucData 属性 |

[FanucData](fd9d131d-3177-76a1-ab73-e521883434f5.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AIComment](caa1ce46-6807-7b6f-0295-822722285f18.htm) |  |
| 公共属性 | [AlarmCurrent](845a692e-e041-8cd2-61d8-739a07890e8a.htm) |  |
| 公共属性 | [AlarmList](4c36388c-6992-1133-51e2-2c197a98663a.htm) |  |
| 公共属性 | [AlarmPassword](7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm) |  |
| 公共属性 | [AOComment](fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm) |  |
| 公共属性 | [CurrentPose](7c03669c-5e00-3769-b016-bb8bb68ba896.htm) |  |
| 公共属性 | [CurrentPose2](0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm) |  |
| 公共属性 | [CurrentPose3](0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm) |  |
| 公共属性 | [CurrentPose4](8d49d24a-422e-8ed5-1cc8-a109bf757671.htm) |  |
| 公共属性 | [CurrentPose5](6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm) |  |
| 公共属性 | [CurrentPoseUF](716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm) |  |
| 公共属性 | [DataPosRegMG](43b750c7-1763-cc14-e6a8-36dcad959b90.htm) |  |
| 公共属性 | [DIComment](c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm) |  |
| 公共属性 | [DOComment](f67f3bd6-7d81-347e-d8bd-92fce6656329.htm) |  |
| 公共属性 | [DUTY\_TEMP](db3365e3-04f5-6a50-79d2-c180364338b9.htm) |  |
| 公共属性 | [FAST\_CLOCK](d85a83dd-ae39-25c3-1d82-39df52aacd72.htm) |  |
| 公共属性 | [GIComment](f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm) |  |
| 公共属性 | [GOComment](cb612119-5f3e-1119-2a34-401c00b48698.htm) |  |
| 公共属性 | [HTTPKCL\_CMDS](46901d3b-052b-1d58-5ef7-b6adae0008c9.htm) |  |
| 公共属性 | [MNUTOOL1\_1](a788bc36-efbe-efc4-21ee-c32af667e9aa.htm) |  |
| 公共属性 | [MOR\_GRP\_CURRENT\_ANG](796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm) |  |
| 公共属性 | [NumReg1](e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm) |  |
| 公共属性 | [NumReg2](17a87be2-0d12-3387-e20c-93f849ccd2db.htm) |  |
| 公共属性 | [PosRegGP1](fe3d3b49-ab48-18db-b360-7f307188a61a.htm) |  |
| 公共属性 | [PosRegGP2](4b0d2c25-1013-0a06-b131-a3146b27a847.htm) |  |
| 公共属性 | [PosRegGP3](ab0cd01e-908b-d86f-c228-da78ee470517.htm) |  |
| 公共属性 | [PosRegGP4](7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm) |  |
| 公共属性 | [PosRegGP5](179e2f5d-5715-57de-2330-cbaadcf200e5.htm) |  |
| 公共属性 | [RIComment](0140fd93-e229-3764-febd-3dd4fa151911.htm) |  |
| 公共属性 | [ROComment](33a5de78-9dba-c2ff-9d31-8c2943a34512.htm) |  |
| 公共属性 | [SIComment](9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm) |  |
| 公共属性 | [SOComment](3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm) |  |
| 公共属性 | [STRREG\_COMMENT\_Comment](5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm) |  |
| 公共属性 | [STRREGComment](46e7f07e-7631-a95c-4875-9576a6684559.htm) |  |
| 公共属性 | [Task](a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm) |  |
| 公共属性 | [TaskIgnoreKarel](5052b3ff-a907-c384-8f24-94fd3974443c.htm) |  |
| 公共属性 | [TaskIgnoreMacro](8502d4b6-23f5-08ab-2575-6374c51c54b8.htm) |  |
| 公共属性 | [TaskIgnoreMacroKarel](96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm) |  |
| 公共属性 | [TIMER10\_COMMENT](653ccde8-6961-7266-4164-b908aa4f5499.htm) |  |
| 公共属性 | [Timer10\_TIMER\_VAL](5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm) |  |
| 公共属性 | [TIMER2\_COMMENT](34698164-5021-1efc-6aaa-14dc3a5f8c58.htm) |  |
| 公共属性 | [UIComment](e9931afb-18e5-50a6-e891-d5c0372a831b.htm) |  |
| 公共属性 | [UOComment](72f57607-083d-ab26-aaa4-e1ec36278ca6.htm) |  |
| 公共属性 | [WIComment](b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm) |  |
| 公共属性 | [WOComment](f70478df-40cb-97aa-3245-41e98674d171.htm) |  |
| 公共属性 | [WSIComment](5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm) |  |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AIComment 属性 

[原文連結](http://api.hslcommunication.cn/html/caa1ce46-6807-7b6f-0295-822722285f18.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataAIComment 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.AIComment" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string[] AIComment { get; set; }
```

```
Public Property AIComment As String()
	Get
	Set
```

```
public:
property array<String^>^ AIComment {
	array<String^>^ get ();
	void set (array<String^>^ value);
}
```

```
member AIComment : string[] with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AlarmCurrent 属性 

[原文連結](http://api.hslcommunication.cn/html/845a692e-e041-8cd2-61d8-739a07890e8a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataAlarmCurrent 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.AlarmCurrent" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucAlarm AlarmCurrent { get; set; }
```

```
Public Property AlarmCurrent As FanucAlarm
	Get
	Set
```

```
public:
property FanucAlarm^ AlarmCurrent {
	FanucAlarm^ get ();
	void set (FanucAlarm^ value);
}
```

```
member AlarmCurrent : FanucAlarm with get, set
```

#### 属性值

类型：[FanucAlarm](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AlarmList 属性 

[原文連結](http://api.hslcommunication.cn/html/4c36388c-6992-1133-51e2-2c197a98663a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataAlarmList 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.AlarmList" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucAlarm[] AlarmList { get; set; }
```

```
Public Property AlarmList As FanucAlarm()
	Get
	Set
```

```
public:
property array<FanucAlarm^>^ AlarmList {
	array<FanucAlarm^>^ get ();
	void set (array<FanucAlarm^>^ value);
}
```

```
member AlarmList : FanucAlarm[] with get, set
```

#### 属性值

类型：[FanucAlarm](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AlarmPassword 属性 

[原文連結](http://api.hslcommunication.cn/html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataAlarmPassword 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.AlarmPassword" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucAlarm AlarmPassword { get; set; }
```

```
Public Property AlarmPassword As FanucAlarm
	Get
	Set
```

```
public:
property FanucAlarm^ AlarmPassword {
	FanucAlarm^ get ();
	void set (FanucAlarm^ value);
}
```

```
member AlarmPassword : FanucAlarm with get, set
```

#### 属性值

类型：[FanucAlarm](f669181a-2617-5fda-b19c-e416aa573f3b.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AOComment 属性 

[原文連結](http://api.hslcommunication.cn/html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataAOComment 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.AOComment" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string[] AOComment { get; set; }
```

```
Public Property AOComment As String()
	Get
	Set
```

```
public:
property array<String^>^ AOComment {
	array<String^>^ get ();
	void set (array<String^>^ value);
}
```

```
member AOComment : string[] with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CurrentPose 属性 

[原文連結](http://api.hslcommunication.cn/html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataCurrentPose 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.CurrentPose" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucPose CurrentPose { get; set; }
```

```
Public Property CurrentPose As FanucPose
	Get
	Set
```

```
public:
property FanucPose^ CurrentPose {
	FanucPose^ get ();
	void set (FanucPose^ value);
}
```

```
member CurrentPose : FanucPose with get, set
```

#### 属性值

类型：[FanucPose](cdcdf7d8-60e1-ed91-7b3c-b14717cb31e6.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CurrentPose2 属性 

[原文連結](http://api.hslcommunication.cn/html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataCurrentPose2 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.CurrentPose2" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucPose CurrentPose2 { get; set; }
```

```
Public Property CurrentPose2 As FanucPose
	Get
	Set
```

```
public:
property FanucPose^ CurrentPose2 {
	FanucPose^ get ();
	void set (FanucPose^ value);
}
```

```
member CurrentPose2 : FanucPose with get, set
```

#### 属性值

类型：[FanucPose](cdcdf7d8-60e1-ed91-7b3c-b14717cb31e6.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CurrentPose3 属性 

[原文連結](http://api.hslcommunication.cn/html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataCurrentPose3 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.CurrentPose3" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucPose CurrentPose3 { get; set; }
```

```
Public Property CurrentPose3 As FanucPose
	Get
	Set
```

```
public:
property FanucPose^ CurrentPose3 {
	FanucPose^ get ();
	void set (FanucPose^ value);
}
```

```
member CurrentPose3 : FanucPose with get, set
```

#### 属性值

类型：[FanucPose](cdcdf7d8-60e1-ed91-7b3c-b14717cb31e6.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CurrentPose4 属性 

[原文連結](http://api.hslcommunication.cn/html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataCurrentPose4 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.CurrentPose4" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucPose CurrentPose4 { get; set; }
```

```
Public Property CurrentPose4 As FanucPose
	Get
	Set
```

```
public:
property FanucPose^ CurrentPose4 {
	FanucPose^ get ();
	void set (FanucPose^ value);
}
```

```
member CurrentPose4 : FanucPose with get, set
```

#### 属性值

类型：[FanucPose](cdcdf7d8-60e1-ed91-7b3c-b14717cb31e6.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CurrentPose5 属性 

[原文連結](http://api.hslcommunication.cn/html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.FANUC](../html/6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm "HslCommunication.Robot.FANUC")

[FanucData 类](../html/fd9d131d-3177-76a1-ab73-e521883434f5.htm "FanucData 类")

[FanucData 属性](../html/e3b9204b-ae26-491b-b1d6-319d32e1549f.htm "FanucData 属性")

[AIComment 属性](../html/caa1ce46-6807-7b6f-0295-822722285f18.htm "AIComment 属性 ")

[AlarmCurrent 属性](../html/845a692e-e041-8cd2-61d8-739a07890e8a.htm "AlarmCurrent 属性 ")

[AlarmList 属性](../html/4c36388c-6992-1133-51e2-2c197a98663a.htm "AlarmList 属性 ")

[AlarmPassword 属性](../html/7a928dda-0b66-70f1-7f2b-b48da768d1fa.htm "AlarmPassword 属性 ")

[AOComment 属性](../html/fc7b9e67-d55d-5e48-cd57-888145dfb0d7.htm "AOComment 属性 ")

[CurrentPose 属性](../html/7c03669c-5e00-3769-b016-bb8bb68ba896.htm "CurrentPose 属性 ")

[CurrentPose2 属性](../html/0693e664-aeed-a36e-a91a-d0ddd84b0bf5.htm "CurrentPose2 属性 ")

[CurrentPose3 属性](../html/0449e93c-31e3-50cc-cb58-c35fd1b7daaf.htm "CurrentPose3 属性 ")

[CurrentPose4 属性](../html/8d49d24a-422e-8ed5-1cc8-a109bf757671.htm "CurrentPose4 属性 ")

[CurrentPose5 属性](../html/6c9735cc-5de9-ddfa-af58-d51dfc67dee4.htm "CurrentPose5 属性 ")

[CurrentPoseUF 属性](../html/716c4a2f-8d40-32c5-b5a5-cf9240fa5fd8.htm "CurrentPoseUF 属性 ")

[DataPosRegMG 属性](../html/43b750c7-1763-cc14-e6a8-36dcad959b90.htm "DataPosRegMG 属性 ")

[DIComment 属性](../html/c68c0a50-5a87-95ea-4b2c-8861a0c21804.htm "DIComment 属性 ")

[DOComment 属性](../html/f67f3bd6-7d81-347e-d8bd-92fce6656329.htm "DOComment 属性 ")

[DUTY\_TEMP 属性](../html/db3365e3-04f5-6a50-79d2-c180364338b9.htm "DUTY_TEMP 属性 ")

[FAST\_CLOCK 属性](../html/d85a83dd-ae39-25c3-1d82-39df52aacd72.htm "FAST_CLOCK 属性 ")

[GIComment 属性](../html/f0fb3bb4-777a-7077-60bd-73b0bf2ef274.htm "GIComment 属性 ")

[GOComment 属性](../html/cb612119-5f3e-1119-2a34-401c00b48698.htm "GOComment 属性 ")

[HTTPKCL\_CMDS 属性](../html/46901d3b-052b-1d58-5ef7-b6adae0008c9.htm "HTTPKCL_CMDS 属性 ")

[MNUTOOL1\_1 属性](../html/a788bc36-efbe-efc4-21ee-c32af667e9aa.htm "MNUTOOL1_1 属性 ")

[MOR\_GRP\_CURRENT\_ANG 属性](../html/796dd9bd-1d53-62fe-b51d-07720ae92cbd.htm "MOR_GRP_CURRENT_ANG 属性 ")

[NumReg1 属性](../html/e7e50cb2-7263-bf61-228f-5e8cb94bff17.htm "NumReg1 属性 ")

[NumReg2 属性](../html/17a87be2-0d12-3387-e20c-93f849ccd2db.htm "NumReg2 属性 ")

[PosRegGP1 属性](../html/fe3d3b49-ab48-18db-b360-7f307188a61a.htm "PosRegGP1 属性 ")

[PosRegGP2 属性](../html/4b0d2c25-1013-0a06-b131-a3146b27a847.htm "PosRegGP2 属性 ")

[PosRegGP3 属性](../html/ab0cd01e-908b-d86f-c228-da78ee470517.htm "PosRegGP3 属性 ")

[PosRegGP4 属性](../html/7ff443c6-9db9-40e5-1b72-59217b6a8ed1.htm "PosRegGP4 属性 ")

[PosRegGP5 属性](../html/179e2f5d-5715-57de-2330-cbaadcf200e5.htm "PosRegGP5 属性 ")

[RIComment 属性](../html/0140fd93-e229-3764-febd-3dd4fa151911.htm "RIComment 属性 ")

[ROComment 属性](../html/33a5de78-9dba-c2ff-9d31-8c2943a34512.htm "ROComment 属性 ")

[SIComment 属性](../html/9fff2dc9-ba26-cb6d-aee4-2f6d19b26756.htm "SIComment 属性 ")

[SOComment 属性](../html/3ea4c011-6bd8-1589-b8bc-5a87afa6628f.htm "SOComment 属性 ")

[STRREG\_COMMENT\_Comment 属性](../html/5f25bdb1-0f49-f812-a381-87fdf0c1de37.htm "STRREG_COMMENT_Comment 属性 ")

[STRREGComment 属性](../html/46e7f07e-7631-a95c-4875-9576a6684559.htm "STRREGComment 属性 ")

[Task 属性](../html/a91d35e1-5db0-0625-5e8f-1b32690e7eb2.htm "Task 属性 ")

[TaskIgnoreKarel 属性](../html/5052b3ff-a907-c384-8f24-94fd3974443c.htm "TaskIgnoreKarel 属性 ")

[TaskIgnoreMacro 属性](../html/8502d4b6-23f5-08ab-2575-6374c51c54b8.htm "TaskIgnoreMacro 属性 ")

[TaskIgnoreMacroKarel 属性](../html/96d036ac-e7e0-c289-4f6a-3f0846cb832b.htm "TaskIgnoreMacroKarel 属性 ")

[TIMER10\_COMMENT 属性](../html/653ccde8-6961-7266-4164-b908aa4f5499.htm "TIMER10_COMMENT 属性 ")

[Timer10\_TIMER\_VAL 属性](../html/5e04f3f2-4721-4911-5e81-8677d0acfc1e.htm "Timer10_TIMER_VAL 属性 ")

[TIMER2\_COMMENT 属性](../html/34698164-5021-1efc-6aaa-14dc3a5f8c58.htm "TIMER2_COMMENT 属性 ")

[UIComment 属性](../html/e9931afb-18e5-50a6-e891-d5c0372a831b.htm "UIComment 属性 ")

[UOComment 属性](../html/72f57607-083d-ab26-aaa4-e1ec36278ca6.htm "UOComment 属性 ")

[WIComment 属性](../html/b1b2b3aa-2889-bd38-445a-4b5e79e6eaa4.htm "WIComment 属性 ")

[WOComment 属性](../html/f70478df-40cb-97aa-3245-41e98674d171.htm "WOComment 属性 ")

[WSIComment 属性](../html/5ecb0404-32ad-2fcc-a1e4-b8a50e1d2366.htm "WSIComment 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucDataCurrentPose5 属性 |

[缺少 "P:HslCommunication.Robot.FANUC.FanucData.CurrentPose5" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.FANUC](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucPose CurrentPose5 { get; set; }
```

```
Public Property CurrentPose5 As FanucPose
	Get
	Set
```

```
public:
property FanucPose^ CurrentPose5 {
	FanucPose^ get ();
	void set (FanucPose^ value);
}
```

```
member CurrentPose5 : FanucPose with get, set
```

#### 属性值

类型：[FanucPose](cdcdf7d8-60e1-ed91-7b3c-b14717cb31e6.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucData 类](fd9d131d-3177-76a1-ab73-e521883434f5.htm)

[HslCommunication.Robot.FANUC 命名空间](6bac49e5-7b3f-b2ab-f0c1-86310e053d5a.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)