# HslCommunication - HslCommunication.LogNet

> 分類頁數: 30



---
## HslCommunication.LogNet

[原文連結](http://api.hslcommunication.cn/html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[GenerateMode 枚举](../html/daf3c59a-f81d-0cdb-2da8-51b40dbb99e4.htm "GenerateMode 枚举")

[HslEventArgs 类](../html/9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm "HslEventArgs 类")

[HslMessageDegree 枚举](../html/3ce37fee-7678-e278-6a0a-7656126003b6.htm "HslMessageDegree 枚举")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[ILogNet 接口](../html/d95d4704-db0d-b379-41e6-03879927a543.htm "ILogNet 接口")

[LogNetAnalysisControl 类](../html/cfda4292-54d6-ce11-0ada-b50021ca4ab7.htm "LogNetAnalysisControl 类")

[LogNetBase 类](../html/591f2a6c-f9cd-f53a-7abc-9d3c7adc563b.htm "LogNetBase 类")

[LogNetDateTime 类](../html/90b8130d-b017-4254-40f3-591d12226ddf.htm "LogNetDateTime 类")

[LogNetException 类](../html/30c57f15-6b01-d816-b1cc-29e7fed02717.htm "LogNetException 类")

[LogNetFileSize 类](../html/b36c7f80-f741-7a9a-e71a-dc2abb9714fe.htm "LogNetFileSize 类")

[LogNetManagment 类](../html/ad49f067-176c-6857-ff05-fb7a6f8fab36.htm "LogNetManagment 类")

[LogNetSingle 类](../html/2b98518d-c8fa-fc59-9900-3ccf00acfdba.htm "LogNetSingle 类")

[LogPathBase 类](../html/24a9b63f-53d8-a3aa-efc2-88ed24655ce0.htm "LogPathBase 类")

[LogSaveMode 枚举](../html/9ea4013b-4452-92f4-dd9a-3f2e12285927.htm "LogSaveMode 枚举")

[LogStatistics 类](../html/00e19edf-0fc1-28d3-17c5-43bf92ea094f.htm "LogStatistics 类")

[LogStatisticsBase(T) 类](../html/4b6cfa8a-7a7b-8151-eba6-e8bd57a62be5.htm "LogStatisticsBase(T) 类")

[LogStatisticsDict 类](../html/28d1bdaf-29f2-c51a-3e8b-0b41d1b98d15.htm "LogStatisticsDict 类")

[LogValueLimit 类](../html/4edfa9e7-8910-aaa0-566b-1fcbfa9c2b14.htm "LogValueLimit 类")

[LogValueLimitDict 类](../html/4e33e792-63f0-90f0-c872-9641fbb18154.htm "LogValueLimitDict 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.LogNet 命名空间 |

[缺少 "N:HslCommunication.LogNet" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [FormLogNetView](48091323-2df2-b52d-1420-224b3bf75e1b.htm) | 日志查看器的窗口类，用于分析统计日志数据，实例化的时候可以直接日志文件，然后直接显示出文件内容出来，然后可以根据日志的等级，或是关键字进行搜索信息  The window class of the log viewer is used to analyze the statistical log data, and when instantiating, you can directly log the file, and then directly display the file content, and then you can search for information according to the level of the log or keyword |
| 公共类 | [HslEventArgs](9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm) | 带有日志消息的事件 |
| 公共类 | [HslMessageItem](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm) | 单条日志的记录信息，包含了消息等级，线程号，关键字，文本信息  Record information of a single log, including message level, thread number, keywords, text information |
| 公共类 | [LogNetAnalysisControl](cfda4292-54d6-ce11-0ada-b50021ca4ab7.htm) | 一个用于日志分析的控件 |
| 公共类 | [LogNetBase](591f2a6c-f9cd-f53a-7abc-9d3c7adc563b.htm) | 日志存储类的基类，提供一些基础的服务 |
| 公共类代码示例 | [LogNetDateTime](90b8130d-b017-4254-40f3-591d12226ddf.htm) | 一个日志组件，可以根据时间来区分不同的文件存储  A log component that can distinguish different file storages based on time |
| 公共类 | [LogNetException](30c57f15-6b01-d816-b1cc-29e7fed02717.htm) | 日志存储回调的异常信息 |
| 公共类代码示例 | [LogNetFileSize](b36c7f80-f741-7a9a-e71a-dc2abb9714fe.htm) | 根据文件的大小来存储日志信息，当前的文件大小增长超过设定值，就会创建新的文件来存储，新的文件命名为当前时间。  Log information is stored according to the size of the file. If the current file size exceeds the set value, a new file is created for storage, and the new file is named the current time. |
| 公共类 | [LogNetManagment](ad49f067-176c-6857-ff05-fb7a6f8fab36.htm) | 日志类的管理器，提供了基本的功能代码。  The manager of the log class provides the basic function code. |
| 公共类代码示例 | [LogNetSingle](2b98518d-c8fa-fc59-9900-3ccf00acfdba.htm) | 单日志文件对象，所有的日志信息的记录都会写入一个文件里面去。文件名指定为空的时候，自动不存储文件。  Single log file object, all log information records will be written to a file. When the file name is specified as empty, the file is not stored automatically. |
| 公共类 | [LogPathBase](24a9b63f-53d8-a3aa-efc2-88ed24655ce0.htm) | 基于路径实现的日志类的基类，提供几个基础的方法信息。  The base class of the log class implemented based on the path provides several basic method information. |
| 公共类代码示例 | [LogStatistics](00e19edf-0fc1-28d3-17c5-43bf92ea094f.htm) | 一个统计次数的辅助类，可用于实现分析一些次数统计信息，比如统计某个API最近每天的访问次数， 统计日志组件最近每天访问的次数，调用者只需要关心统计方式和数据个数，详细参照API文档。  An auxiliary class for counting the number of times, which can be used to realize the analysis of some number of times statistical information, such as counting the number of daily visits of an API, and counting the number of daily visits of the log component. The caller only needs to care about the statistical method and the number of data. Refer to details API documentation. |
| 公共类 | [LogStatisticsBaseT](4b6cfa8a-7a7b-8151-eba6-e8bd57a62be5.htm) | 一个按照实际进行数据分割的辅助基类，可以用于实现对某个的API按照每天进行调用次数统计，也可以实现对某个设备数据按照天进行最大最小值均值分析，这些都是需要继承实现。  An auxiliary base class that divides the data according to the actual data can be used to implement statistics on the number of calls per day for a certain API, and it can also implement the maximum and minimum average analysis of a certain device data according to the day. These all need to be inherited. . |
| 公共类 | [LogStatisticsDict](28d1bdaf-29f2-c51a-3e8b-0b41d1b98d15.htm) | 的词典集合类，用于多个数据的统计信息，例如可以统计多个规格的产量信息，统计多个方法的调用次数信息  The dictionary collection class of is used for the statistical information of multiple data, for example, the output information of multiple specifications can be counted, and the number of calls of multiple methods can be counted. |
| 公共类代码示例 | [LogValueLimit](4edfa9e7-8910-aaa0-566b-1fcbfa9c2b14.htm) | 一个用于数值范围记录的类，可以按照时间进行分类统计，比如计算一个温度值的每天的开始值，结束值，最大值，最小值，平均值信息。详细见API文档信息。  A class used to record the value range, which can be classified according to time, such as calculating the start value, end value, maximum value, minimum value, and average value of a temperature value. See the API documentation for details. |
| 公共类 | [LogValueLimitDict](4e33e792-63f0-90f0-c872-9641fbb18154.htm) | 的词典集合类，用于多个数据的统计信息，例如可以统计多个温度变量的最大值，最小值，平均值  The dictionary collection class, used for statistical information of multiple data, for example, it can count the maximum, minimum, and average values of multiple temperature variables |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口代码示例 | [ILogNet](d95d4704-db0d-b379-41e6-03879927a543.htm) | 一个通用的日志接口，支持5个等级的日志消息写入，支持设置当前的消息等级，定义一个消息存储前的触发事件。  A general-purpose log interface, supports the writing of 5 levels of log messages, supports setting the current message level, and defining a trigger event before a message is stored. |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [GenerateMode](daf3c59a-f81d-0cdb-2da8-51b40dbb99e4.htm) | 日志文件输出模式 |
| 公共枚举 | [HslMessageDegree](3ce37fee-7678-e278-6a0a-7656126003b6.htm) | 记录消息的等级 |
| 公共枚举 | [LogSaveMode](9ea4013b-4452-92f4-dd9a-3f2e12285927.htm) | 日志文件的存储模式 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FormLogNetView 类

[原文連結](http://api.hslcommunication.cn/html/48091323-2df2-b52d-1420-224b3bf75e1b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 构造函数](../html/abc13d39-7d83-9e57-aadf-14131c4954a9.htm "FormLogNetView 构造函数 ")

[FormLogNetView 属性](../html/e39cf0a7-5e4c-0d82-39b6-11af4dcce5b6.htm "FormLogNetView 属性")

[FormLogNetView 方法](../html/b42f6f6e-241d-b7ef-f727-e75a68d8407e.htm "FormLogNetView 方法")

[FormLogNetView 事件](../html/9a0e39b9-2117-2a86-ef9f-754fb5ee728f.htm "FormLogNetView 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetView 类 |

日志查看器的窗口类，用于分析统计日志数据，实例化的时候可以直接日志文件，然后直接显示出文件内容出来，然后可以根据日志的等级，或是关键字进行搜索信息  
The window class of the log viewer is used to analyze the statistical log data, and when instantiating, you can directly log the file,
and then directly display the file content, and then you can search for information according to the level of the log or keyword

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  SystemMarshalByRefObject  
    System.ComponentModelComponent  
      System.Windows.FormsControl  
        System.Windows.FormsScrollableControl  
          System.Windows.FormsContainerControl  
            System.Windows.FormsForm  
              HslCommunication.LogNetFormLogNetView

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FormLogNetView : Form
```

```
Public Class FormLogNetView
	Inherits Form
```

```
public ref class FormLogNetView : public Form
```

```
type FormLogNetView =  
    class
        inherit Form
    end
```

FormLogNetView 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FormLogNetView](0d518468-9879-b5eb-b57d-f37359a21405.htm) | 实例化一个默认的日志查看器的窗口  Instantiates a default log viewer window |
| 公共方法 | [FormLogNetView(String)](111ffd4c-432c-af08-3f28-12f4151f672e.htm) | 指定一个日志路径实例化一个日志查看界面 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | AcceptButton | (继承自 Form。) |
| 公共属性 | AccessibilityObject | (继承自 Control。) |
| 公共属性 | AccessibleDefaultActionDescription | (继承自 Control。) |
| 公共属性 | AccessibleDescription | (继承自 Control。) |
| 公共属性 | AccessibleName | (继承自 Control。) |
| 公共属性 | AccessibleRole | (继承自 Control。) |
| 公共属性 | ActiveControl | (继承自 ContainerControl。) |
| 公共属性 | ActiveMdiChild | (继承自 Form。) |
| 公共属性 | AllowDrop | (继承自 Control。) |
| 公共属性 | AllowTransparency | (继承自 Form。) |
| 公共属性 | Anchor | (继承自 Control。) |
| 公共属性 | AutoScaleDimensions | (继承自 ContainerControl。) |
| 受保护的属性 | AutoScaleFactor | (继承自 ContainerControl。) |
| 公共属性 | AutoScaleMode | (继承自 ContainerControl。) |
| 公共属性 | AutoScroll | (继承自 Form。) |
| 公共属性 | AutoScrollMargin | (继承自 ScrollableControl。) |
| 公共属性 | AutoScrollMinSize | (继承自 ScrollableControl。) |
| 公共属性 | AutoScrollOffset | (继承自 Control。) |
| 公共属性 | AutoScrollPosition | (继承自 ScrollableControl。) |
| 公共属性 | AutoSize | (继承自 Form。) |
| 公共属性 | AutoSizeMode | (继承自 Form。) |
| 公共属性 | AutoValidate | (继承自 Form。) |
| 公共属性 | BackColor | (继承自 Form。) |
| 公共属性 | BackgroundImage | (继承自 Control。) |
| 公共属性 | BackgroundImageLayout | (继承自 Control。) |
| 公共属性 | BindingContext | (继承自 ContainerControl。) |
| 公共属性 | Bottom | (继承自 Control。) |
| 公共属性 | Bounds | (继承自 Control。) |
| 公共属性 | CancelButton | (继承自 Form。) |
| 受保护的属性 | CanEnableIme | (继承自 ContainerControl。) |
| 公共属性 | CanFocus | (继承自 Control。) |
| 受保护的属性 | CanRaiseEvents | (继承自 Control。) |
| 公共属性 | CanSelect | (继承自 Control。) |
| 公共属性 | Capture | (继承自 Control。) |
| 公共属性 | CausesValidation | (继承自 Control。) |
| 公共属性 | ClientRectangle | (继承自 Control。) |
| 公共属性 | ClientSize | (继承自 Form。) |
| 公共属性 | CompanyName | (继承自 Control。) |
| 公共属性 | Container | (继承自 Component。) |
| 公共属性 | ContainsFocus | (继承自 Control。) |
| 公共属性 | ContextMenu | (继承自 Control。) |
| 公共属性 | ContextMenuStrip | (继承自 Control。) |
| 公共属性 | ControlBox | (继承自 Form。) |
| 公共属性 | Controls | (继承自 Control。) |
| 公共属性 | Created | (继承自 Control。) |
| 受保护的属性 | CreateParams | (继承自 Form。) |
| 公共属性 | CurrentAutoScaleDimensions | (继承自 ContainerControl。) |
| 公共属性 | Cursor | (继承自 Control。) |
| 公共属性 | DataBindings | (继承自 Control。) |
| 受保护的属性 | DefaultCursor | (继承自 Control。) |
| 受保护的属性 | DefaultImeMode | (继承自 Form。) |
| 受保护的属性 | DefaultMargin | (继承自 Control。) |
| 受保护的属性 | DefaultMaximumSize | (继承自 Control。) |
| 受保护的属性 | DefaultMinimumSize | (继承自 Control。) |
| 受保护的属性 | DefaultPadding | (继承自 Control。) |
| 受保护的属性 | DefaultSize | (继承自 Form。) |
| 受保护的属性 | DesignMode | (继承自 Component。) |
| 公共属性 | DesktopBounds | (继承自 Form。) |
| 公共属性 | DesktopLocation | (继承自 Form。) |
| 公共属性 | DialogResult | (继承自 Form。) |
| 公共属性 | DisplayRectangle | (继承自 ScrollableControl。) |
| 公共属性 | Disposing | (继承自 Control。) |
| 公共属性 | Dock | (继承自 Control。) |
| 受保护的属性 | DoubleBuffered | (继承自 Control。) |
| 公共属性 | Enabled | (继承自 Control。) |
| 受保护的属性 | Events | (继承自 Component。) |
| 公共属性 | Focused | (继承自 Control。) |
| 公共属性 | Font | (继承自 Control。) |
| 受保护的属性 | FontHeight | (继承自 Control。) |
| 公共属性 | ForeColor | (继承自 Control。) |
| 公共属性 | FormBorderStyle | (继承自 Form。) |
| 公共属性 | Handle | (继承自 Control。) |
| 公共属性 | HasChildren | (继承自 Control。) |
| 公共属性 | Height | (继承自 Control。) |
| 公共属性 | HelpButton | (继承自 Form。) |
| 公共属性 | HorizontalScroll | (继承自 ScrollableControl。) |
| 受保护的属性 | HScroll | (继承自 ScrollableControl。) |
| 公共属性 | Icon | (继承自 Form。) |
| 公共属性 | ImeMode | (继承自 Control。) |
| 受保护的属性 | ImeModeBase | (继承自 Control。) |
| 公共属性 | InvokeRequired | (继承自 Control。) |
| 公共属性 | IsAccessible | (继承自 Control。) |
| 公共属性 | IsDisposed | (继承自 Control。) |
| 公共属性 | IsHandleCreated | (继承自 Control。) |
| 公共属性 | IsMdiChild | (继承自 Form。) |
| 公共属性 | IsMdiContainer | (继承自 Form。) |
| 公共属性 | IsMirrored | (继承自 Control。) |
| 公共属性 | IsRestrictedWindow | (继承自 Form。) |
| 公共属性 | KeyPreview | (继承自 Form。) |
| 公共属性 | LayoutEngine | (继承自 Control。) |
| 公共属性 | Left | (继承自 Control。) |
| 公共属性 | Location | (继承自 Form。) |
| 公共属性 | MainMenuStrip | (继承自 Form。) |
| 公共属性 | Margin | (继承自 Control。) |
| 公共属性 | MaximizeBox | (继承自 Form。) |
| 受保护的属性 | MaximizedBounds | (继承自 Form。) |
| 公共属性 | MaximumSize | (继承自 Form。) |
| 公共属性 | MdiChildren | (继承自 Form。) |
| 公共属性 | MdiParent | (继承自 Form。) |
| 公共属性 | Menu | (继承自 Form。) |
| 公共属性 | MergedMenu | (继承自 Form。) |
| 公共属性 | MinimizeBox | (继承自 Form。) |
| 公共属性 | MinimumSize | (继承自 Form。) |
| 公共属性 | Modal | (继承自 Form。) |
| 公共属性 | Name | (继承自 Control。) |
| 公共属性 | Opacity | (继承自 Form。) |
| 公共属性 | [OpenDialogDefaultPath](2d6498c4-77b3-0c3d-0ae9-62b7d21e3eaf.htm) | 获取或设置当前日志选择窗口的默认的路径信息  Get or set the default path information of the current log selection window |
| 公共属性 | OwnedForms | (继承自 Form。) |
| 公共属性 | Owner | (继承自 Form。) |
| 公共属性 | Padding | (继承自 Control。) |
| 公共属性 | Parent | (继承自 Control。) |
| 公共属性 | ParentForm | (继承自 ContainerControl。) |
| 公共属性 | PreferredSize | (继承自 Control。) |
| 公共属性 | ProductName | (继承自 Control。) |
| 公共属性 | ProductVersion | (继承自 Control。) |
| 公共属性 | RecreatingHandle | (继承自 Control。) |
| 公共属性 | Region | (继承自 Control。) |
| 受保护的属性 | RenderRightToLeft | **已过时。** (继承自 Control。) |
| 受保护的属性 | ResizeRedraw | (继承自 Control。) |
| 公共属性 | RestoreBounds | (继承自 Form。) |
| 公共属性 | Right | (继承自 Control。) |
| 公共属性 | RightToLeft | (继承自 Control。) |
| 公共属性 | RightToLeftLayout | (继承自 Form。) |
| 受保护的属性 | ScaleChildren | (继承自 Control。) |
| 受保护的属性 | ShowFocusCues | (继承自 Control。) |
| 公共属性 | ShowIcon | (继承自 Form。) |
| 公共属性 | ShowInTaskbar | (继承自 Form。) |
| 受保护的属性 | ShowKeyboardCues | (继承自 Control。) |
| 受保护的属性 | ShowWithoutActivation | (继承自 Form。) |
| 公共属性 | Site | (继承自 Control。) |
| 公共属性 | Size | (继承自 Form。) |
| 公共属性 | SizeGripStyle | (继承自 Form。) |
| 公共属性 | StartPosition | (继承自 Form。) |
| 公共属性 | TabIndex | (继承自 Control。) |
| 公共属性 | TabStop | (继承自 Control。) |
| 公共属性 | Tag | (继承自 Control。) |
| 公共属性 | Text | (继承自 Form。) |
| 公共属性 | Top | (继承自 Control。) |
| 公共属性 | TopLevel | (继承自 Form。) |
| 公共属性 | TopLevelControl | (继承自 Control。) |
| 公共属性 | TopMost | (继承自 Form。) |
| 公共属性 | TransparencyKey | (继承自 Form。) |
| 公共属性 | UseWaitCursor | (继承自 Control。) |
| 公共属性 | VerticalScroll | (继承自 ScrollableControl。) |
| 公共属性 | Visible | (继承自 Control。) |
| 受保护的属性 | VScroll | (继承自 ScrollableControl。) |
| 公共属性 | Width | (继承自 Control。) |
| 公共属性 | WindowState | (继承自 Form。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | AccessibilityNotifyClients(AccessibleEvents, Int32) | (继承自 Control。) |
| 受保护的方法 | AccessibilityNotifyClients(AccessibleEvents, Int32, Int32) | (继承自 Control。) |
| 公共方法 | Activate | (继承自 Form。) |
| 受保护的方法 | ActivateMdiChild | (继承自 Form。) |
| 公共方法 | AddOwnedForm | (继承自 Form。) |
| 受保护的方法 | AdjustFormScrollbars | (继承自 Form。) |
| 公共方法 | BeginInvoke(Delegate) | (继承自 Control。) |
| 公共方法 | BeginInvoke(Delegate, Object) | (继承自 Control。) |
| 公共方法 | BringToFront | (继承自 Control。) |
| 受保护的方法 | CenterToParent | (继承自 Form。) |
| 受保护的方法 | CenterToScreen | (继承自 Form。) |
| 公共方法 | Close | (继承自 Form。) |
| 公共方法 | Contains | (继承自 Control。) |
| 受保护的方法 | CreateAccessibilityInstance | (继承自 Control。) |
| 公共方法 | CreateControl | (继承自 Control。) |
| 受保护的方法 | CreateControlsInstance | (继承自 Form。) |
| 公共方法 | CreateGraphics | (继承自 Control。) |
| 受保护的方法 | CreateHandle | (继承自 Form。) |
| 公共方法 | CreateObjRef | (继承自 MarshalByRefObject。) |
| 受保护的方法 | DefWndProc | (继承自 Form。) |
| 受保护的方法 | DestroyHandle | (继承自 Control。) |
| 公共方法 | Dispose | (继承自 Component。) |
| 受保护的方法 | [Dispose(Boolean)](312ebbad-220a-fb04-c1f7-428e90018da4.htm) | Clean up any resources being used. (重写 FormDispose(Boolean).) |
| 公共方法 | DoDragDrop | (继承自 Control。) |
| 公共方法 | DrawToBitmap | (继承自 Control。) |
| 公共方法 | EndInvoke | (继承自 Control。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Component。) |
| 公共方法 | FindForm | (继承自 Control。) |
| 公共方法 | Focus | (继承自 Control。) |
| 受保护的方法 | GetAccessibilityObjectById | (继承自 Control。) |
| 受保护的方法 | GetAutoSizeMode | (继承自 Control。) |
| 公共方法 | GetChildAtPoint(Point) | (继承自 Control。) |
| 公共方法 | GetChildAtPoint(Point, GetChildAtPointSkip) | (继承自 Control。) |
| 公共方法 | GetContainerControl | (继承自 Control。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetLifetimeService | (继承自 MarshalByRefObject。) |
| 公共方法 | GetNextControl | (继承自 Control。) |
| 公共方法 | GetPreferredSize | (继承自 Control。) |
| 受保护的方法 | GetScaledBounds | (继承自 Form。) |
| 受保护的方法 | GetScrollState | (继承自 ScrollableControl。) |
| 受保护的方法 | GetService | (继承自 Component。) |
| 受保护的方法 | GetStyle | (继承自 Control。) |
| 受保护的方法 | GetTopLevel | (继承自 Control。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | Hide | (继承自 Control。) |
| 公共方法 | InitializeLifetimeService | (继承自 MarshalByRefObject。) |
| 受保护的方法 | InitLayout | (继承自 Control。) |
| 公共方法 | Invalidate | (继承自 Control。) |
| 公共方法 | Invalidate(Region) | (继承自 Control。) |
| 公共方法 | Invalidate(Boolean) | (继承自 Control。) |
| 公共方法 | Invalidate(Rectangle) | (继承自 Control。) |
| 公共方法 | Invalidate(Region, Boolean) | (继承自 Control。) |
| 公共方法 | Invalidate(Rectangle, Boolean) | (继承自 Control。) |
| 公共方法 | Invoke(Delegate) | (继承自 Control。) |
| 公共方法 | Invoke(Delegate, Object) | (继承自 Control。) |
| 受保护的方法 | InvokeGotFocus | (继承自 Control。) |
| 受保护的方法 | InvokeLostFocus | (继承自 Control。) |
| 受保护的方法 | InvokeOnClick | (继承自 Control。) |
| 受保护的方法 | InvokePaint | (继承自 Control。) |
| 受保护的方法 | InvokePaintBackground | (继承自 Control。) |
| 受保护的方法 | IsInputChar | (继承自 Control。) |
| 受保护的方法 | IsInputKey | (继承自 Control。) |
| 公共方法 | LayoutMdi | (继承自 Form。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone(Boolean) | (继承自 MarshalByRefObject。) |
| 受保护的方法 | NotifyInvalidate | (继承自 Control。) |
| 受保护的方法 | OnActivated | (继承自 Form。) |
| 受保护的方法 | OnAutoSizeChanged | (继承自 Control。) |
| 受保护的方法 | OnAutoValidateChanged | (继承自 ContainerControl。) |
| 受保护的方法 | OnBackColorChanged | (继承自 Control。) |
| 受保护的方法 | OnBackgroundImageChanged | (继承自 Form。) |
| 受保护的方法 | OnBackgroundImageLayoutChanged | (继承自 Form。) |
| 受保护的方法 | OnBindingContextChanged | (继承自 Control。) |
| 受保护的方法 | OnCausesValidationChanged | (继承自 Control。) |
| 受保护的方法 | OnChangeUICues | (继承自 Control。) |
| 受保护的方法 | OnClick | (继承自 Control。) |
| 受保护的方法 | OnClientSizeChanged | (继承自 Control。) |
| 受保护的方法 | OnClosed | (继承自 Form。) |
| 受保护的方法 | OnClosing | (继承自 Form。) |
| 受保护的方法 | OnContextMenuChanged | (继承自 Control。) |
| 受保护的方法 | OnContextMenuStripChanged | (继承自 Control。) |
| 受保护的方法 | OnControlAdded | (继承自 Control。) |
| 受保护的方法 | OnControlRemoved | (继承自 Control。) |
| 受保护的方法 | OnCreateControl | (继承自 Form。) |
| 受保护的方法 | OnCursorChanged | (继承自 Control。) |
| 受保护的方法 | OnDeactivate | (继承自 Form。) |
| 受保护的方法 | OnDockChanged | (继承自 Control。) |
| 受保护的方法 | OnDoubleClick | (继承自 Control。) |
| 受保护的方法 | OnDragDrop | (继承自 Control。) |
| 受保护的方法 | OnDragEnter | (继承自 Control。) |
| 受保护的方法 | OnDragLeave | (继承自 Control。) |
| 受保护的方法 | OnDragOver | (继承自 Control。) |
| 受保护的方法 | OnEnabledChanged | (继承自 Form。) |
| 受保护的方法 | OnEnter | (继承自 Form。) |
| 受保护的方法 | OnFontChanged | (继承自 Form。) |
| 受保护的方法 | OnForeColorChanged | (继承自 Control。) |
| 受保护的方法 | OnFormClosed | (继承自 Form。) |
| 受保护的方法 | OnFormClosing | (继承自 Form。) |
| 受保护的方法 | OnGiveFeedback | (继承自 Control。) |
| 受保护的方法 | OnGotFocus | (继承自 Control。) |
| 受保护的方法 | OnHandleCreated | (继承自 Form。) |
| 受保护的方法 | OnHandleDestroyed | (继承自 Form。) |
| 受保护的方法 | OnHelpButtonClicked | (继承自 Form。) |
| 受保护的方法 | OnHelpRequested | (继承自 Control。) |
| 受保护的方法 | OnImeModeChanged | (继承自 Control。) |
| 受保护的方法 | OnInputLanguageChanged | (继承自 Form。) |
| 受保护的方法 | OnInputLanguageChanging | (继承自 Form。) |
| 受保护的方法 | OnInvalidated | (继承自 Control。) |
| 受保护的方法 | OnKeyDown | (继承自 Control。) |
| 受保护的方法 | OnKeyPress | (继承自 Control。) |
| 受保护的方法 | OnKeyUp | (继承自 Control。) |
| 受保护的方法 | OnLayout | (继承自 Form。) |
| 受保护的方法 | OnLeave | (继承自 Control。) |
| 受保护的方法 | OnLoad | (继承自 Form。) |
| 受保护的方法 | OnLocationChanged | (继承自 Control。) |
| 受保护的方法 | OnLostFocus | (继承自 Control。) |
| 受保护的方法 | OnMarginChanged | (继承自 Control。) |
| 受保护的方法 | OnMaximizedBoundsChanged | (继承自 Form。) |
| 受保护的方法 | OnMaximumSizeChanged | (继承自 Form。) |
| 受保护的方法 | OnMdiChildActivate | (继承自 Form。) |
| 受保护的方法 | OnMenuComplete | (继承自 Form。) |
| 受保护的方法 | OnMenuStart | (继承自 Form。) |
| 受保护的方法 | OnMinimumSizeChanged | (继承自 Form。) |
| 受保护的方法 | OnMouseCaptureChanged | (继承自 Control。) |
| 受保护的方法 | OnMouseClick | (继承自 Control。) |
| 受保护的方法 | OnMouseDoubleClick | (继承自 Control。) |
| 受保护的方法 | OnMouseDown | (继承自 Control。) |
| 受保护的方法 | OnMouseEnter | (继承自 Control。) |
| 受保护的方法 | OnMouseHover | (继承自 Control。) |
| 受保护的方法 | OnMouseLeave | (继承自 Control。) |
| 受保护的方法 | OnMouseMove | (继承自 Control。) |
| 受保护的方法 | OnMouseUp | (继承自 Control。) |
| 受保护的方法 | OnMouseWheel | (继承自 ScrollableControl。) |
| 受保护的方法 | OnMove | (继承自 Control。) |
| 受保护的方法 | OnNotifyMessage | (继承自 Control。) |
| 受保护的方法 | OnPaddingChanged | (继承自 ScrollableControl。) |
| 受保护的方法 | OnPaint | (继承自 Form。) |
| 受保护的方法 | OnPaintBackground | (继承自 ScrollableControl。) |
| 受保护的方法 | OnParentBackColorChanged | (继承自 Control。) |
| 受保护的方法 | OnParentBackgroundImageChanged | (继承自 Control。) |
| 受保护的方法 | OnParentBindingContextChanged | (继承自 Control。) |
| 受保护的方法 | OnParentChanged | (继承自 ContainerControl。) |
| 受保护的方法 | OnParentCursorChanged | (继承自 Control。) |
| 受保护的方法 | OnParentEnabledChanged | (继承自 Control。) |
| 受保护的方法 | OnParentFontChanged | (继承自 Control。) |
| 受保护的方法 | OnParentForeColorChanged | (继承自 Control。) |
| 受保护的方法 | OnParentRightToLeftChanged | (继承自 Control。) |
| 受保护的方法 | OnParentVisibleChanged | (继承自 Control。) |
| 受保护的方法 | OnPreviewKeyDown | (继承自 Control。) |
| 受保护的方法 | OnPrint | (继承自 Control。) |
| 受保护的方法 | OnQueryContinueDrag | (继承自 Control。) |
| 受保护的方法 | OnRegionChanged | (继承自 Control。) |
| 受保护的方法 | OnResize | (继承自 Form。) |
| 受保护的方法 | OnResizeBegin | (继承自 Form。) |
| 受保护的方法 | OnResizeEnd | (继承自 Form。) |
| 受保护的方法 | OnRightToLeftChanged | (继承自 ScrollableControl。) |
| 受保护的方法 | OnRightToLeftLayoutChanged | (继承自 Form。) |
| 受保护的方法 | OnScroll | (继承自 ScrollableControl。) |
| 受保护的方法 | [OnShown](211bbe9f-1f05-37d9-2df1-5f0ac93474ed.htm) | (重写 FormOnShown(EventArgs).) |
| 受保护的方法 | OnSizeChanged | (继承自 Control。) |
| 受保护的方法 | OnStyleChanged | (继承自 Form。) |
| 受保护的方法 | OnSystemColorsChanged | (继承自 Control。) |
| 受保护的方法 | OnTabIndexChanged | (继承自 Control。) |
| 受保护的方法 | OnTabStopChanged | (继承自 Control。) |
| 受保护的方法 | OnTextChanged | (继承自 Form。) |
| 受保护的方法 | OnValidated | (继承自 Control。) |
| 受保护的方法 | OnValidating | (继承自 Control。) |
| 受保护的方法 | OnVisibleChanged | (继承自 Form。) |
| 公共方法 | PerformAutoScale | (继承自 ContainerControl。) |
| 公共方法 | PerformLayout | (继承自 Control。) |
| 公共方法 | PerformLayout(Control, String) | (继承自 Control。) |
| 公共方法 | PointToClient | (继承自 Control。) |
| 公共方法 | PointToScreen | (继承自 Control。) |
| 公共方法 | PreProcessControlMessage | (继承自 Control。) |
| 公共方法 | PreProcessMessage | (继承自 Control。) |
| 受保护的方法 | ProcessCmdKey | (继承自 Form。) |
| 受保护的方法 | ProcessDialogChar | (继承自 Form。) |
| 受保护的方法 | ProcessDialogKey | (继承自 Form。) |
| 受保护的方法 | ProcessKeyEventArgs | (继承自 Control。) |
| 受保护的方法 | ProcessKeyMessage | (继承自 Control。) |
| 受保护的方法 | ProcessKeyPreview | (继承自 Form。) |
| 受保护的方法 | ProcessMnemonic | (继承自 Form。) |
| 受保护的方法 | ProcessTabKey | (继承自 Form。) |
| 受保护的方法 | RaiseDragEvent | (继承自 Control。) |
| 受保护的方法 | RaiseKeyEvent | (继承自 Control。) |
| 受保护的方法 | RaiseMouseEvent | (继承自 Control。) |
| 受保护的方法 | RaisePaintEvent | (继承自 Control。) |
| 受保护的方法 | RecreateHandle | (继承自 Control。) |
| 公共方法 | RectangleToClient | (继承自 Control。) |
| 公共方法 | RectangleToScreen | (继承自 Control。) |
| 公共方法 | Refresh | (继承自 Control。) |
| 公共方法 | RemoveOwnedForm | (继承自 Form。) |
| 受保护的方法 | ResetMouseEventArgs | (继承自 Control。) |
| 公共方法 | ResetText | (继承自 Control。) |
| 公共方法 | ResumeLayout | (继承自 Control。) |
| 公共方法 | ResumeLayout(Boolean) | (继承自 Control。) |
| 受保护的方法 | RtlTranslateAlignment(HorizontalAlignment) | (继承自 Control。) |
| 受保护的方法 | RtlTranslateAlignment(LeftRightAlignment) | (继承自 Control。) |
| 受保护的方法 | RtlTranslateAlignment(ContentAlignment) | (继承自 Control。) |
| 受保护的方法 | RtlTranslateContent | (继承自 Control。) |
| 受保护的方法 | RtlTranslateHorizontal | (继承自 Control。) |
| 受保护的方法 | RtlTranslateLeftRight | (继承自 Control。) |
| 公共方法 | Scale | (继承自 Control。) |
| 受保护的方法 | ScaleControl | (继承自 Form。) |
| 公共方法 | ScrollControlIntoView | (继承自 ScrollableControl。) |
| 受保护的方法 | ScrollToControl | (继承自 ScrollableControl。) |
| 公共方法 | Select | (继承自 Control。) |
| 受保护的方法 | Select(Boolean, Boolean) | (继承自 Form。) |
| 公共方法 | SelectNextControl | (继承自 Control。) |
| 公共方法 | SendToBack | (继承自 Control。) |
| 公共方法 | SetAutoScrollMargin | (继承自 ScrollableControl。) |
| 受保护的方法 | SetAutoSizeMode | (继承自 Control。) |
| 公共方法 | SetBounds(Int32, Int32, Int32, Int32) | (继承自 Control。) |
| 公共方法 | SetBounds(Int32, Int32, Int32, Int32, BoundsSpecified) | (继承自 Control。) |
| 受保护的方法 | SetBoundsCore | (继承自 Form。) |
| 受保护的方法 | SetClientSizeCore | (继承自 Form。) |
| 公共方法 | SetDesktopBounds | (继承自 Form。) |
| 公共方法 | SetDesktopLocation | (继承自 Form。) |
| 受保护的方法 | SetDisplayRectLocation | (继承自 ScrollableControl。) |
| 受保护的方法 | SetScrollState | (继承自 ScrollableControl。) |
| 受保护的方法 | SetStyle | (继承自 Control。) |
| 受保护的方法 | SetTopLevel | (继承自 Control。) |
| 受保护的方法 | SetVisibleCore | (继承自 Form。) |
| 公共方法 | Show | (继承自 Control。) |
| 公共方法 | Show(IWin32Window) | (继承自 Form。) |
| 公共方法 | ShowDialog | (继承自 Form。) |
| 公共方法 | ShowDialog(IWin32Window) | (继承自 Form。) |
| 受保护的方法 | SizeFromClientSize | (继承自 Control。) |
| 公共方法 | SuspendLayout | (继承自 Control。) |
| 公共方法 | ToString | (继承自 Form。) |
| 公共方法 | Update | (继承自 Control。) |
| 受保护的方法 | UpdateBounds | (继承自 Control。) |
| 受保护的方法 | UpdateBounds(Int32, Int32, Int32, Int32) | (继承自 Control。) |
| 受保护的方法 | UpdateBounds(Int32, Int32, Int32, Int32, Int32, Int32) | (继承自 Control。) |
| 受保护的方法 | UpdateDefaultButton | (继承自 Form。) |
| 受保护的方法 | UpdateStyles | (继承自 Control。) |
| 受保护的方法 | UpdateZOrder | (继承自 Control。) |
| 公共方法 | Validate | (继承自 ContainerControl。) |
| 公共方法 | Validate(Boolean) | (继承自 ContainerControl。) |
| 公共方法 | ValidateChildren | (继承自 Form。) |
| 公共方法 | ValidateChildren(ValidationConstraints) | (继承自 Form。) |
| 受保护的方法 | WndProc | (继承自 Form。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | Activated | (继承自 Form。) |
| 公共事件 | AutoSizeChanged | (继承自 Form。) |
| 公共事件 | AutoValidateChanged | (继承自 Form。) |
| 公共事件 | BackColorChanged | (继承自 Control。) |
| 公共事件 | BackgroundImageChanged | (继承自 Control。) |
| 公共事件 | BackgroundImageLayoutChanged | (继承自 Control。) |
| 公共事件 | BindingContextChanged | (继承自 Control。) |
| 公共事件 | CausesValidationChanged | (继承自 Control。) |
| 公共事件 | ChangeUICues | (继承自 Control。) |
| 公共事件 | Click | (继承自 Control。) |
| 公共事件 | ClientSizeChanged | (继承自 Control。) |
| 公共事件 | ContextMenuChanged | (继承自 Control。) |
| 公共事件 | ContextMenuStripChanged | (继承自 Control。) |
| 公共事件 | ControlAdded | (继承自 Control。) |
| 公共事件 | ControlRemoved | (继承自 Control。) |
| 公共事件 | CursorChanged | (继承自 Control。) |
| 公共事件 | Deactivate | (继承自 Form。) |
| 公共事件 | Disposed | (继承自 Component。) |
| 公共事件 | DockChanged | (继承自 Control。) |
| 公共事件 | DoubleClick | (继承自 Control。) |
| 公共事件 | DragDrop | (继承自 Control。) |
| 公共事件 | DragEnter | (继承自 Control。) |
| 公共事件 | DragLeave | (继承自 Control。) |
| 公共事件 | DragOver | (继承自 Control。) |
| 公共事件 | EnabledChanged | (继承自 Control。) |
| 公共事件 | Enter | (继承自 Control。) |
| 公共事件 | FontChanged | (继承自 Control。) |
| 公共事件 | ForeColorChanged | (继承自 Control。) |
| 公共事件 | FormClosed | (继承自 Form。) |
| 公共事件 | FormClosing | (继承自 Form。) |
| 公共事件 | GiveFeedback | (继承自 Control。) |
| 公共事件 | GotFocus | (继承自 Control。) |
| 公共事件 | HandleCreated | (继承自 Control。) |
| 公共事件 | HandleDestroyed | (继承自 Control。) |
| 公共事件 | HelpButtonClicked | (继承自 Form。) |
| 公共事件 | HelpRequested | (继承自 Control。) |
| 公共事件 | ImeModeChanged | (继承自 Control。) |
| 公共事件 | InputLanguageChanged | (继承自 Form。) |
| 公共事件 | InputLanguageChanging | (继承自 Form。) |
| 公共事件 | Invalidated | (继承自 Control。) |
| 公共事件 | KeyDown | (继承自 Control。) |
| 公共事件 | KeyPress | (继承自 Control。) |
| 公共事件 | KeyUp | (继承自 Control。) |
| 公共事件 | Layout | (继承自 Control。) |
| 公共事件 | Leave | (继承自 Control。) |
| 公共事件 | Load | (继承自 Form。) |
| 公共事件 | LocationChanged | (继承自 Control。) |
| 公共事件 | LostFocus | (继承自 Control。) |
| 公共事件 | MarginChanged | (继承自 Control。) |
| 公共事件 | MaximizedBoundsChanged | (继承自 Form。) |
| 公共事件 | MaximumSizeChanged | (继承自 Form。) |
| 公共事件 | MdiChildActivate | (继承自 Form。) |
| 公共事件 | MenuComplete | (继承自 Form。) |
| 公共事件 | MenuStart | (继承自 Form。) |
| 公共事件 | MinimumSizeChanged | (继承自 Form。) |
| 公共事件 | MouseCaptureChanged | (继承自 Control。) |
| 公共事件 | MouseClick | (继承自 Control。) |
| 公共事件 | MouseDoubleClick | (继承自 Control。) |
| 公共事件 | MouseDown | (继承自 Control。) |
| 公共事件 | MouseEnter | (继承自 Control。) |
| 公共事件 | MouseHover | (继承自 Control。) |
| 公共事件 | MouseLeave | (继承自 Control。) |
| 公共事件 | MouseMove | (继承自 Control。) |
| 公共事件 | MouseUp | (继承自 Control。) |
| 公共事件 | MouseWheel | (继承自 Control。) |
| 公共事件 | Move | (继承自 Control。) |
| 公共事件 | PaddingChanged | (继承自 Control。) |
| 公共事件 | Paint | (继承自 Control。) |
| 公共事件 | ParentChanged | (继承自 Control。) |
| 公共事件 | PreviewKeyDown | (继承自 Control。) |
| 公共事件 | QueryAccessibilityHelp | (继承自 Control。) |
| 公共事件 | QueryContinueDrag | (继承自 Control。) |
| 公共事件 | RegionChanged | (继承自 Control。) |
| 公共事件 | Resize | (继承自 Control。) |
| 公共事件 | ResizeBegin | (继承自 Form。) |
| 公共事件 | ResizeEnd | (继承自 Form。) |
| 公共事件 | RightToLeftChanged | (继承自 Control。) |
| 公共事件 | RightToLeftLayoutChanged | (继承自 Form。) |
| 公共事件 | Scroll | (继承自 ScrollableControl。) |
| 公共事件 | Shown | (继承自 Form。) |
| 公共事件 | SizeChanged | (继承自 Control。) |
| 公共事件 | StyleChanged | (继承自 Control。) |
| 公共事件 | SystemColorsChanged | (继承自 Control。) |
| 公共事件 | TabIndexChanged | (继承自 Control。) |
| 公共事件 | TabStopChanged | (继承自 Control。) |
| 公共事件 | TextChanged | (继承自 Control。) |
| 公共事件 | Validated | (继承自 Control。) |
| 公共事件 | Validating | (继承自 Control。) |
| 公共事件 | VisibleChanged | (继承自 Control。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FormLogNetView 构造函数 

[原文連結](http://api.hslcommunication.cn/html/abc13d39-7d83-9e57-aadf-14131c4954a9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 构造函数](../html/abc13d39-7d83-9e57-aadf-14131c4954a9.htm "FormLogNetView 构造函数 ")

[FormLogNetView 构造函数](../html/0d518468-9879-b5eb-b57d-f37359a21405.htm "FormLogNetView 构造函数 ")

[FormLogNetView 构造函数 (String)](../html/111ffd4c-432c-af08-3f28-12f4151f672e.htm "FormLogNetView 构造函数 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetView 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FormLogNetView](0d518468-9879-b5eb-b57d-f37359a21405.htm) | 实例化一个默认的日志查看器的窗口  Instantiates a default log viewer window |
| 公共方法 | [FormLogNetView(String)](111ffd4c-432c-af08-3f28-12f4151f672e.htm) | 指定一个日志路径实例化一个日志查看界面 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FormLogNetView 构造函数 

[原文連結](http://api.hslcommunication.cn/html/0d518468-9879-b5eb-b57d-f37359a21405.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 构造函数](../html/abc13d39-7d83-9e57-aadf-14131c4954a9.htm "FormLogNetView 构造函数 ")

[FormLogNetView 构造函数](../html/0d518468-9879-b5eb-b57d-f37359a21405.htm "FormLogNetView 构造函数 ")

[FormLogNetView 构造函数 (String)](../html/111ffd4c-432c-af08-3f28-12f4151f672e.htm "FormLogNetView 构造函数 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetView 构造函数 |

实例化一个默认的日志查看器的窗口  
Instantiates a default log viewer window

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FormLogNetView()
```

```
Public Sub New
```

```
public:
FormLogNetView()
```

```
new : unit -> FormLogNetView
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[FormLogNetView 重载](abc13d39-7d83-9e57-aadf-14131c4954a9.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FormLogNetView 构造函数 (String)

[原文連結](http://api.hslcommunication.cn/html/111ffd4c-432c-af08-3f28-12f4151f672e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 构造函数](../html/abc13d39-7d83-9e57-aadf-14131c4954a9.htm "FormLogNetView 构造函数 ")

[FormLogNetView 构造函数](../html/0d518468-9879-b5eb-b57d-f37359a21405.htm "FormLogNetView 构造函数 ")

[FormLogNetView 构造函数 (String)](../html/111ffd4c-432c-af08-3f28-12f4151f672e.htm "FormLogNetView 构造函数 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetView 构造函数 (String) |

指定一个日志路径实例化一个日志查看界面

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FormLogNetView(
	string log
)
```

```
Public Sub New ( 
	log As String
)
```

```
public:
FormLogNetView(
	String^ log
)
```

```
new : 
        log : string -> FormLogNetView
```

#### 参数

log
:   类型：SystemString  
    日志的路径

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[FormLogNetView 重载](abc13d39-7d83-9e57-aadf-14131c4954a9.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FormLogNetView 属性

[原文連結](http://api.hslcommunication.cn/html/e39cf0a7-5e4c-0d82-39b6-11af4dcce5b6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 属性](../html/e39cf0a7-5e4c-0d82-39b6-11af4dcce5b6.htm "FormLogNetView 属性")

[OpenDialogDefaultPath 属性](../html/2d6498c4-77b3-0c3d-0ae9-62b7d21e3eaf.htm "OpenDialogDefaultPath 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetView 属性 |

[FormLogNetView](48091323-2df2-b52d-1420-224b3bf75e1b.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | AcceptButton | (继承自 Form。) |
| 公共属性 | AccessibilityObject | (继承自 Control。) |
| 公共属性 | AccessibleDefaultActionDescription | (继承自 Control。) |
| 公共属性 | AccessibleDescription | (继承自 Control。) |
| 公共属性 | AccessibleName | (继承自 Control。) |
| 公共属性 | AccessibleRole | (继承自 Control。) |
| 公共属性 | ActiveControl | (继承自 ContainerControl。) |
| 公共属性 | ActiveMdiChild | (继承自 Form。) |
| 公共属性 | AllowDrop | (继承自 Control。) |
| 公共属性 | AllowTransparency | (继承自 Form。) |
| 公共属性 | Anchor | (继承自 Control。) |
| 公共属性 | AutoScaleDimensions | (继承自 ContainerControl。) |
| 受保护的属性 | AutoScaleFactor | (继承自 ContainerControl。) |
| 公共属性 | AutoScaleMode | (继承自 ContainerControl。) |
| 公共属性 | AutoScroll | (继承自 Form。) |
| 公共属性 | AutoScrollMargin | (继承自 ScrollableControl。) |
| 公共属性 | AutoScrollMinSize | (继承自 ScrollableControl。) |
| 公共属性 | AutoScrollOffset | (继承自 Control。) |
| 公共属性 | AutoScrollPosition | (继承自 ScrollableControl。) |
| 公共属性 | AutoSize | (继承自 Form。) |
| 公共属性 | AutoSizeMode | (继承自 Form。) |
| 公共属性 | AutoValidate | (继承自 Form。) |
| 公共属性 | BackColor | (继承自 Form。) |
| 公共属性 | BackgroundImage | (继承自 Control。) |
| 公共属性 | BackgroundImageLayout | (继承自 Control。) |
| 公共属性 | BindingContext | (继承自 ContainerControl。) |
| 公共属性 | Bottom | (继承自 Control。) |
| 公共属性 | Bounds | (继承自 Control。) |
| 公共属性 | CancelButton | (继承自 Form。) |
| 受保护的属性 | CanEnableIme | (继承自 ContainerControl。) |
| 公共属性 | CanFocus | (继承自 Control。) |
| 受保护的属性 | CanRaiseEvents | (继承自 Control。) |
| 公共属性 | CanSelect | (继承自 Control。) |
| 公共属性 | Capture | (继承自 Control。) |
| 公共属性 | CausesValidation | (继承自 Control。) |
| 公共属性 | ClientRectangle | (继承自 Control。) |
| 公共属性 | ClientSize | (继承自 Form。) |
| 公共属性 | CompanyName | (继承自 Control。) |
| 公共属性 | Container | (继承自 Component。) |
| 公共属性 | ContainsFocus | (继承自 Control。) |
| 公共属性 | ContextMenu | (继承自 Control。) |
| 公共属性 | ContextMenuStrip | (继承自 Control。) |
| 公共属性 | ControlBox | (继承自 Form。) |
| 公共属性 | Controls | (继承自 Control。) |
| 公共属性 | Created | (继承自 Control。) |
| 受保护的属性 | CreateParams | (继承自 Form。) |
| 公共属性 | CurrentAutoScaleDimensions | (继承自 ContainerControl。) |
| 公共属性 | Cursor | (继承自 Control。) |
| 公共属性 | DataBindings | (继承自 Control。) |
| 受保护的属性 | DefaultCursor | (继承自 Control。) |
| 受保护的属性 | DefaultImeMode | (继承自 Form。) |
| 受保护的属性 | DefaultMargin | (继承自 Control。) |
| 受保护的属性 | DefaultMaximumSize | (继承自 Control。) |
| 受保护的属性 | DefaultMinimumSize | (继承自 Control。) |
| 受保护的属性 | DefaultPadding | (继承自 Control。) |
| 受保护的属性 | DefaultSize | (继承自 Form。) |
| 受保护的属性 | DesignMode | (继承自 Component。) |
| 公共属性 | DesktopBounds | (继承自 Form。) |
| 公共属性 | DesktopLocation | (继承自 Form。) |
| 公共属性 | DialogResult | (继承自 Form。) |
| 公共属性 | DisplayRectangle | (继承自 ScrollableControl。) |
| 公共属性 | Disposing | (继承自 Control。) |
| 公共属性 | Dock | (继承自 Control。) |
| 受保护的属性 | DoubleBuffered | (继承自 Control。) |
| 公共属性 | Enabled | (继承自 Control。) |
| 受保护的属性 | Events | (继承自 Component。) |
| 公共属性 | Focused | (继承自 Control。) |
| 公共属性 | Font | (继承自 Control。) |
| 受保护的属性 | FontHeight | (继承自 Control。) |
| 公共属性 | ForeColor | (继承自 Control。) |
| 公共属性 | FormBorderStyle | (继承自 Form。) |
| 公共属性 | Handle | (继承自 Control。) |
| 公共属性 | HasChildren | (继承自 Control。) |
| 公共属性 | Height | (继承自 Control。) |
| 公共属性 | HelpButton | (继承自 Form。) |
| 公共属性 | HorizontalScroll | (继承自 ScrollableControl。) |
| 受保护的属性 | HScroll | (继承自 ScrollableControl。) |
| 公共属性 | Icon | (继承自 Form。) |
| 公共属性 | ImeMode | (继承自 Control。) |
| 受保护的属性 | ImeModeBase | (继承自 Control。) |
| 公共属性 | InvokeRequired | (继承自 Control。) |
| 公共属性 | IsAccessible | (继承自 Control。) |
| 公共属性 | IsDisposed | (继承自 Control。) |
| 公共属性 | IsHandleCreated | (继承自 Control。) |
| 公共属性 | IsMdiChild | (继承自 Form。) |
| 公共属性 | IsMdiContainer | (继承自 Form。) |
| 公共属性 | IsMirrored | (继承自 Control。) |
| 公共属性 | IsRestrictedWindow | (继承自 Form。) |
| 公共属性 | KeyPreview | (继承自 Form。) |
| 公共属性 | LayoutEngine | (继承自 Control。) |
| 公共属性 | Left | (继承自 Control。) |
| 公共属性 | Location | (继承自 Form。) |
| 公共属性 | MainMenuStrip | (继承自 Form。) |
| 公共属性 | Margin | (继承自 Control。) |
| 公共属性 | MaximizeBox | (继承自 Form。) |
| 受保护的属性 | MaximizedBounds | (继承自 Form。) |
| 公共属性 | MaximumSize | (继承自 Form。) |
| 公共属性 | MdiChildren | (继承自 Form。) |
| 公共属性 | MdiParent | (继承自 Form。) |
| 公共属性 | Menu | (继承自 Form。) |
| 公共属性 | MergedMenu | (继承自 Form。) |
| 公共属性 | MinimizeBox | (继承自 Form。) |
| 公共属性 | MinimumSize | (继承自 Form。) |
| 公共属性 | Modal | (继承自 Form。) |
| 公共属性 | Name | (继承自 Control。) |
| 公共属性 | Opacity | (继承自 Form。) |
| 公共属性 | [OpenDialogDefaultPath](2d6498c4-77b3-0c3d-0ae9-62b7d21e3eaf.htm) | 获取或设置当前日志选择窗口的默认的路径信息  Get or set the default path information of the current log selection window |
| 公共属性 | OwnedForms | (继承自 Form。) |
| 公共属性 | Owner | (继承自 Form。) |
| 公共属性 | Padding | (继承自 Control。) |
| 公共属性 | Parent | (继承自 Control。) |
| 公共属性 | ParentForm | (继承自 ContainerControl。) |
| 公共属性 | PreferredSize | (继承自 Control。) |
| 公共属性 | ProductName | (继承自 Control。) |
| 公共属性 | ProductVersion | (继承自 Control。) |
| 公共属性 | RecreatingHandle | (继承自 Control。) |
| 公共属性 | Region | (继承自 Control。) |
| 受保护的属性 | RenderRightToLeft | **已过时。** (继承自 Control。) |
| 受保护的属性 | ResizeRedraw | (继承自 Control。) |
| 公共属性 | RestoreBounds | (继承自 Form。) |
| 公共属性 | Right | (继承自 Control。) |
| 公共属性 | RightToLeft | (继承自 Control。) |
| 公共属性 | RightToLeftLayout | (继承自 Form。) |
| 受保护的属性 | ScaleChildren | (继承自 Control。) |
| 受保护的属性 | ShowFocusCues | (继承自 Control。) |
| 公共属性 | ShowIcon | (继承自 Form。) |
| 公共属性 | ShowInTaskbar | (继承自 Form。) |
| 受保护的属性 | ShowKeyboardCues | (继承自 Control。) |
| 受保护的属性 | ShowWithoutActivation | (继承自 Form。) |
| 公共属性 | Site | (继承自 Control。) |
| 公共属性 | Size | (继承自 Form。) |
| 公共属性 | SizeGripStyle | (继承自 Form。) |
| 公共属性 | StartPosition | (继承自 Form。) |
| 公共属性 | TabIndex | (继承自 Control。) |
| 公共属性 | TabStop | (继承自 Control。) |
| 公共属性 | Tag | (继承自 Control。) |
| 公共属性 | Text | (继承自 Form。) |
| 公共属性 | Top | (继承自 Control。) |
| 公共属性 | TopLevel | (继承自 Form。) |
| 公共属性 | TopLevelControl | (继承自 Control。) |
| 公共属性 | TopMost | (继承自 Form。) |
| 公共属性 | TransparencyKey | (继承自 Form。) |
| 公共属性 | UseWaitCursor | (继承自 Control。) |
| 公共属性 | VerticalScroll | (继承自 ScrollableControl。) |
| 公共属性 | Visible | (继承自 Control。) |
| 受保护的属性 | VScroll | (继承自 ScrollableControl。) |
| 公共属性 | Width | (继承自 Control。) |
| 公共属性 | WindowState | (继承自 Form。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OpenDialogDefaultPath 属性 

[原文連結](http://api.hslcommunication.cn/html/2d6498c4-77b3-0c3d-0ae9-62b7d21e3eaf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 属性](../html/e39cf0a7-5e4c-0d82-39b6-11af4dcce5b6.htm "FormLogNetView 属性")

[OpenDialogDefaultPath 属性](../html/2d6498c4-77b3-0c3d-0ae9-62b7d21e3eaf.htm "OpenDialogDefaultPath 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetViewOpenDialogDefaultPath 属性 |

获取或设置当前日志选择窗口的默认的路径信息  
Get or set the default path information of the current log selection window

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string OpenDialogDefaultPath { get; set; }
```

```
Public Property OpenDialogDefaultPath As String
	Get
	Set
```

```
public:
property String^ OpenDialogDefaultPath {
	String^ get ();
	void set (String^ value);
}
```

```
member OpenDialogDefaultPath : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FormLogNetView 方法

[原文連結](http://api.hslcommunication.cn/html/b42f6f6e-241d-b7ef-f727-e75a68d8407e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 方法](../html/b42f6f6e-241d-b7ef-f727-e75a68d8407e.htm "FormLogNetView 方法")

[Dispose 方法](../html/f956a136-0f5a-dcd0-a9b8-7819ba4b63bd.htm "Dispose 方法 ")

[OnShown 方法](../html/211bbe9f-1f05-37d9-2df1-5f0ac93474ed.htm "OnShown 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetView 方法 |

[FormLogNetView](48091323-2df2-b52d-1420-224b3bf75e1b.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | AccessibilityNotifyClients(AccessibleEvents, Int32) | (继承自 Control。) |
| 受保护的方法 | AccessibilityNotifyClients(AccessibleEvents, Int32, Int32) | (继承自 Control。) |
| 公共方法 | Activate | (继承自 Form。) |
| 受保护的方法 | ActivateMdiChild | (继承自 Form。) |
| 公共方法 | AddOwnedForm | (继承自 Form。) |
| 受保护的方法 | AdjustFormScrollbars | (继承自 Form。) |
| 公共方法 | BeginInvoke(Delegate) | (继承自 Control。) |
| 公共方法 | BeginInvoke(Delegate, Object) | (继承自 Control。) |
| 公共方法 | BringToFront | (继承自 Control。) |
| 受保护的方法 | CenterToParent | (继承自 Form。) |
| 受保护的方法 | CenterToScreen | (继承自 Form。) |
| 公共方法 | Close | (继承自 Form。) |
| 公共方法 | Contains | (继承自 Control。) |
| 受保护的方法 | CreateAccessibilityInstance | (继承自 Control。) |
| 公共方法 | CreateControl | (继承自 Control。) |
| 受保护的方法 | CreateControlsInstance | (继承自 Form。) |
| 公共方法 | CreateGraphics | (继承自 Control。) |
| 受保护的方法 | CreateHandle | (继承自 Form。) |
| 公共方法 | CreateObjRef | (继承自 MarshalByRefObject。) |
| 受保护的方法 | DefWndProc | (继承自 Form。) |
| 受保护的方法 | DestroyHandle | (继承自 Control。) |
| 公共方法 | Dispose | (继承自 Component。) |
| 受保护的方法 | [Dispose(Boolean)](312ebbad-220a-fb04-c1f7-428e90018da4.htm) | Clean up any resources being used. (重写 FormDispose(Boolean).) |
| 公共方法 | DoDragDrop | (继承自 Control。) |
| 公共方法 | DrawToBitmap | (继承自 Control。) |
| 公共方法 | EndInvoke | (继承自 Control。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Component。) |
| 公共方法 | FindForm | (继承自 Control。) |
| 公共方法 | Focus | (继承自 Control。) |
| 受保护的方法 | GetAccessibilityObjectById | (继承自 Control。) |
| 受保护的方法 | GetAutoSizeMode | (继承自 Control。) |
| 公共方法 | GetChildAtPoint(Point) | (继承自 Control。) |
| 公共方法 | GetChildAtPoint(Point, GetChildAtPointSkip) | (继承自 Control。) |
| 公共方法 | GetContainerControl | (继承自 Control。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetLifetimeService | (继承自 MarshalByRefObject。) |
| 公共方法 | GetNextControl | (继承自 Control。) |
| 公共方法 | GetPreferredSize | (继承自 Control。) |
| 受保护的方法 | GetScaledBounds | (继承自 Form。) |
| 受保护的方法 | GetScrollState | (继承自 ScrollableControl。) |
| 受保护的方法 | GetService | (继承自 Component。) |
| 受保护的方法 | GetStyle | (继承自 Control。) |
| 受保护的方法 | GetTopLevel | (继承自 Control。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | Hide | (继承自 Control。) |
| 公共方法 | InitializeLifetimeService | (继承自 MarshalByRefObject。) |
| 受保护的方法 | InitLayout | (继承自 Control。) |
| 公共方法 | Invalidate | (继承自 Control。) |
| 公共方法 | Invalidate(Region) | (继承自 Control。) |
| 公共方法 | Invalidate(Boolean) | (继承自 Control。) |
| 公共方法 | Invalidate(Rectangle) | (继承自 Control。) |
| 公共方法 | Invalidate(Region, Boolean) | (继承自 Control。) |
| 公共方法 | Invalidate(Rectangle, Boolean) | (继承自 Control。) |
| 公共方法 | Invoke(Delegate) | (继承自 Control。) |
| 公共方法 | Invoke(Delegate, Object) | (继承自 Control。) |
| 受保护的方法 | InvokeGotFocus | (继承自 Control。) |
| 受保护的方法 | InvokeLostFocus | (继承自 Control。) |
| 受保护的方法 | InvokeOnClick | (继承自 Control。) |
| 受保护的方法 | InvokePaint | (继承自 Control。) |
| 受保护的方法 | InvokePaintBackground | (继承自 Control。) |
| 受保护的方法 | IsInputChar | (继承自 Control。) |
| 受保护的方法 | IsInputKey | (继承自 Control。) |
| 公共方法 | LayoutMdi | (继承自 Form。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone(Boolean) | (继承自 MarshalByRefObject。) |
| 受保护的方法 | NotifyInvalidate | (继承自 Control。) |
| 受保护的方法 | OnActivated | (继承自 Form。) |
| 受保护的方法 | OnAutoSizeChanged | (继承自 Control。) |
| 受保护的方法 | OnAutoValidateChanged | (继承自 ContainerControl。) |
| 受保护的方法 | OnBackColorChanged | (继承自 Control。) |
| 受保护的方法 | OnBackgroundImageChanged | (继承自 Form。) |
| 受保护的方法 | OnBackgroundImageLayoutChanged | (继承自 Form。) |
| 受保护的方法 | OnBindingContextChanged | (继承自 Control。) |
| 受保护的方法 | OnCausesValidationChanged | (继承自 Control。) |
| 受保护的方法 | OnChangeUICues | (继承自 Control。) |
| 受保护的方法 | OnClick | (继承自 Control。) |
| 受保护的方法 | OnClientSizeChanged | (继承自 Control。) |
| 受保护的方法 | OnClosed | (继承自 Form。) |
| 受保护的方法 | OnClosing | (继承自 Form。) |
| 受保护的方法 | OnContextMenuChanged | (继承自 Control。) |
| 受保护的方法 | OnContextMenuStripChanged | (继承自 Control。) |
| 受保护的方法 | OnControlAdded | (继承自 Control。) |
| 受保护的方法 | OnControlRemoved | (继承自 Control。) |
| 受保护的方法 | OnCreateControl | (继承自 Form。) |
| 受保护的方法 | OnCursorChanged | (继承自 Control。) |
| 受保护的方法 | OnDeactivate | (继承自 Form。) |
| 受保护的方法 | OnDockChanged | (继承自 Control。) |
| 受保护的方法 | OnDoubleClick | (继承自 Control。) |
| 受保护的方法 | OnDragDrop | (继承自 Control。) |
| 受保护的方法 | OnDragEnter | (继承自 Control。) |
| 受保护的方法 | OnDragLeave | (继承自 Control。) |
| 受保护的方法 | OnDragOver | (继承自 Control。) |
| 受保护的方法 | OnEnabledChanged | (继承自 Form。) |
| 受保护的方法 | OnEnter | (继承自 Form。) |
| 受保护的方法 | OnFontChanged | (继承自 Form。) |
| 受保护的方法 | OnForeColorChanged | (继承自 Control。) |
| 受保护的方法 | OnFormClosed | (继承自 Form。) |
| 受保护的方法 | OnFormClosing | (继承自 Form。) |
| 受保护的方法 | OnGiveFeedback | (继承自 Control。) |
| 受保护的方法 | OnGotFocus | (继承自 Control。) |
| 受保护的方法 | OnHandleCreated | (继承自 Form。) |
| 受保护的方法 | OnHandleDestroyed | (继承自 Form。) |
| 受保护的方法 | OnHelpButtonClicked | (继承自 Form。) |
| 受保护的方法 | OnHelpRequested | (继承自 Control。) |
| 受保护的方法 | OnImeModeChanged | (继承自 Control。) |
| 受保护的方法 | OnInputLanguageChanged | (继承自 Form。) |
| 受保护的方法 | OnInputLanguageChanging | (继承自 Form。) |
| 受保护的方法 | OnInvalidated | (继承自 Control。) |
| 受保护的方法 | OnKeyDown | (继承自 Control。) |
| 受保护的方法 | OnKeyPress | (继承自 Control。) |
| 受保护的方法 | OnKeyUp | (继承自 Control。) |
| 受保护的方法 | OnLayout | (继承自 Form。) |
| 受保护的方法 | OnLeave | (继承自 Control。) |
| 受保护的方法 | OnLoad | (继承自 Form。) |
| 受保护的方法 | OnLocationChanged | (继承自 Control。) |
| 受保护的方法 | OnLostFocus | (继承自 Control。) |
| 受保护的方法 | OnMarginChanged | (继承自 Control。) |
| 受保护的方法 | OnMaximizedBoundsChanged | (继承自 Form。) |
| 受保护的方法 | OnMaximumSizeChanged | (继承自 Form。) |
| 受保护的方法 | OnMdiChildActivate | (继承自 Form。) |
| 受保护的方法 | OnMenuComplete | (继承自 Form。) |
| 受保护的方法 | OnMenuStart | (继承自 Form。) |
| 受保护的方法 | OnMinimumSizeChanged | (继承自 Form。) |
| 受保护的方法 | OnMouseCaptureChanged | (继承自 Control。) |
| 受保护的方法 | OnMouseClick | (继承自 Control。) |
| 受保护的方法 | OnMouseDoubleClick | (继承自 Control。) |
| 受保护的方法 | OnMouseDown | (继承自 Control。) |
| 受保护的方法 | OnMouseEnter | (继承自 Control。) |
| 受保护的方法 | OnMouseHover | (继承自 Control。) |
| 受保护的方法 | OnMouseLeave | (继承自 Control。) |
| 受保护的方法 | OnMouseMove | (继承自 Control。) |
| 受保护的方法 | OnMouseUp | (继承自 Control。) |
| 受保护的方法 | OnMouseWheel | (继承自 ScrollableControl。) |
| 受保护的方法 | OnMove | (继承自 Control。) |
| 受保护的方法 | OnNotifyMessage | (继承自 Control。) |
| 受保护的方法 | OnPaddingChanged | (继承自 ScrollableControl。) |
| 受保护的方法 | OnPaint | (继承自 Form。) |
| 受保护的方法 | OnPaintBackground | (继承自 ScrollableControl。) |
| 受保护的方法 | OnParentBackColorChanged | (继承自 Control。) |
| 受保护的方法 | OnParentBackgroundImageChanged | (继承自 Control。) |
| 受保护的方法 | OnParentBindingContextChanged | (继承自 Control。) |
| 受保护的方法 | OnParentChanged | (继承自 ContainerControl。) |
| 受保护的方法 | OnParentCursorChanged | (继承自 Control。) |
| 受保护的方法 | OnParentEnabledChanged | (继承自 Control。) |
| 受保护的方法 | OnParentFontChanged | (继承自 Control。) |
| 受保护的方法 | OnParentForeColorChanged | (继承自 Control。) |
| 受保护的方法 | OnParentRightToLeftChanged | (继承自 Control。) |
| 受保护的方法 | OnParentVisibleChanged | (继承自 Control。) |
| 受保护的方法 | OnPreviewKeyDown | (继承自 Control。) |
| 受保护的方法 | OnPrint | (继承自 Control。) |
| 受保护的方法 | OnQueryContinueDrag | (继承自 Control。) |
| 受保护的方法 | OnRegionChanged | (继承自 Control。) |
| 受保护的方法 | OnResize | (继承自 Form。) |
| 受保护的方法 | OnResizeBegin | (继承自 Form。) |
| 受保护的方法 | OnResizeEnd | (继承自 Form。) |
| 受保护的方法 | OnRightToLeftChanged | (继承自 ScrollableControl。) |
| 受保护的方法 | OnRightToLeftLayoutChanged | (继承自 Form。) |
| 受保护的方法 | OnScroll | (继承自 ScrollableControl。) |
| 受保护的方法 | [OnShown](211bbe9f-1f05-37d9-2df1-5f0ac93474ed.htm) | (重写 FormOnShown(EventArgs).) |
| 受保护的方法 | OnSizeChanged | (继承自 Control。) |
| 受保护的方法 | OnStyleChanged | (继承自 Form。) |
| 受保护的方法 | OnSystemColorsChanged | (继承自 Control。) |
| 受保护的方法 | OnTabIndexChanged | (继承自 Control。) |
| 受保护的方法 | OnTabStopChanged | (继承自 Control。) |
| 受保护的方法 | OnTextChanged | (继承自 Form。) |
| 受保护的方法 | OnValidated | (继承自 Control。) |
| 受保护的方法 | OnValidating | (继承自 Control。) |
| 受保护的方法 | OnVisibleChanged | (继承自 Form。) |
| 公共方法 | PerformAutoScale | (继承自 ContainerControl。) |
| 公共方法 | PerformLayout | (继承自 Control。) |
| 公共方法 | PerformLayout(Control, String) | (继承自 Control。) |
| 公共方法 | PointToClient | (继承自 Control。) |
| 公共方法 | PointToScreen | (继承自 Control。) |
| 公共方法 | PreProcessControlMessage | (继承自 Control。) |
| 公共方法 | PreProcessMessage | (继承自 Control。) |
| 受保护的方法 | ProcessCmdKey | (继承自 Form。) |
| 受保护的方法 | ProcessDialogChar | (继承自 Form。) |
| 受保护的方法 | ProcessDialogKey | (继承自 Form。) |
| 受保护的方法 | ProcessKeyEventArgs | (继承自 Control。) |
| 受保护的方法 | ProcessKeyMessage | (继承自 Control。) |
| 受保护的方法 | ProcessKeyPreview | (继承自 Form。) |
| 受保护的方法 | ProcessMnemonic | (继承自 Form。) |
| 受保护的方法 | ProcessTabKey | (继承自 Form。) |
| 受保护的方法 | RaiseDragEvent | (继承自 Control。) |
| 受保护的方法 | RaiseKeyEvent | (继承自 Control。) |
| 受保护的方法 | RaiseMouseEvent | (继承自 Control。) |
| 受保护的方法 | RaisePaintEvent | (继承自 Control。) |
| 受保护的方法 | RecreateHandle | (继承自 Control。) |
| 公共方法 | RectangleToClient | (继承自 Control。) |
| 公共方法 | RectangleToScreen | (继承自 Control。) |
| 公共方法 | Refresh | (继承自 Control。) |
| 公共方法 | RemoveOwnedForm | (继承自 Form。) |
| 受保护的方法 | ResetMouseEventArgs | (继承自 Control。) |
| 公共方法 | ResetText | (继承自 Control。) |
| 公共方法 | ResumeLayout | (继承自 Control。) |
| 公共方法 | ResumeLayout(Boolean) | (继承自 Control。) |
| 受保护的方法 | RtlTranslateAlignment(HorizontalAlignment) | (继承自 Control。) |
| 受保护的方法 | RtlTranslateAlignment(LeftRightAlignment) | (继承自 Control。) |
| 受保护的方法 | RtlTranslateAlignment(ContentAlignment) | (继承自 Control。) |
| 受保护的方法 | RtlTranslateContent | (继承自 Control。) |
| 受保护的方法 | RtlTranslateHorizontal | (继承自 Control。) |
| 受保护的方法 | RtlTranslateLeftRight | (继承自 Control。) |
| 公共方法 | Scale | (继承自 Control。) |
| 受保护的方法 | ScaleControl | (继承自 Form。) |
| 公共方法 | ScrollControlIntoView | (继承自 ScrollableControl。) |
| 受保护的方法 | ScrollToControl | (继承自 ScrollableControl。) |
| 公共方法 | Select | (继承自 Control。) |
| 受保护的方法 | Select(Boolean, Boolean) | (继承自 Form。) |
| 公共方法 | SelectNextControl | (继承自 Control。) |
| 公共方法 | SendToBack | (继承自 Control。) |
| 公共方法 | SetAutoScrollMargin | (继承自 ScrollableControl。) |
| 受保护的方法 | SetAutoSizeMode | (继承自 Control。) |
| 公共方法 | SetBounds(Int32, Int32, Int32, Int32) | (继承自 Control。) |
| 公共方法 | SetBounds(Int32, Int32, Int32, Int32, BoundsSpecified) | (继承自 Control。) |
| 受保护的方法 | SetBoundsCore | (继承自 Form。) |
| 受保护的方法 | SetClientSizeCore | (继承自 Form。) |
| 公共方法 | SetDesktopBounds | (继承自 Form。) |
| 公共方法 | SetDesktopLocation | (继承自 Form。) |
| 受保护的方法 | SetDisplayRectLocation | (继承自 ScrollableControl。) |
| 受保护的方法 | SetScrollState | (继承自 ScrollableControl。) |
| 受保护的方法 | SetStyle | (继承自 Control。) |
| 受保护的方法 | SetTopLevel | (继承自 Control。) |
| 受保护的方法 | SetVisibleCore | (继承自 Form。) |
| 公共方法 | Show | (继承自 Control。) |
| 公共方法 | Show(IWin32Window) | (继承自 Form。) |
| 公共方法 | ShowDialog | (继承自 Form。) |
| 公共方法 | ShowDialog(IWin32Window) | (继承自 Form。) |
| 受保护的方法 | SizeFromClientSize | (继承自 Control。) |
| 公共方法 | SuspendLayout | (继承自 Control。) |
| 公共方法 | ToString | (继承自 Form。) |
| 公共方法 | Update | (继承自 Control。) |
| 受保护的方法 | UpdateBounds | (继承自 Control。) |
| 受保护的方法 | UpdateBounds(Int32, Int32, Int32, Int32) | (继承自 Control。) |
| 受保护的方法 | UpdateBounds(Int32, Int32, Int32, Int32, Int32, Int32) | (继承自 Control。) |
| 受保护的方法 | UpdateDefaultButton | (继承自 Form。) |
| 受保护的方法 | UpdateStyles | (继承自 Control。) |
| 受保护的方法 | UpdateZOrder | (继承自 Control。) |
| 公共方法 | Validate | (继承自 ContainerControl。) |
| 公共方法 | Validate(Boolean) | (继承自 ContainerControl。) |
| 公共方法 | ValidateChildren | (继承自 Form。) |
| 公共方法 | ValidateChildren(ValidationConstraints) | (继承自 Form。) |
| 受保护的方法 | WndProc | (继承自 Form。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/f956a136-0f5a-dcd0-a9b8-7819ba4b63bd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 方法](../html/b42f6f6e-241d-b7ef-f727-e75a68d8407e.htm "FormLogNetView 方法")

[Dispose 方法](../html/f956a136-0f5a-dcd0-a9b8-7819ba4b63bd.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/312ebbad-220a-fb04-c1f7-428e90018da4.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetViewDispose 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Dispose | (继承自 Component。) |
| 受保护的方法 | [Dispose(Boolean)](312ebbad-220a-fb04-c1f7-428e90018da4.htm) | Clean up any resources being used. (重写 FormDispose(Boolean).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 (Boolean)

[原文連結](http://api.hslcommunication.cn/html/312ebbad-220a-fb04-c1f7-428e90018da4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 方法](../html/b42f6f6e-241d-b7ef-f727-e75a68d8407e.htm "FormLogNetView 方法")

[Dispose 方法](../html/f956a136-0f5a-dcd0-a9b8-7819ba4b63bd.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/312ebbad-220a-fb04-c1f7-428e90018da4.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetViewDispose 方法 (Boolean) |

Clean up any resources being used.

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override void Dispose(
	bool disposing
)
```

```
Protected Overrides Sub Dispose ( 
	disposing As Boolean
)
```

```
protected:
virtual void Dispose(
	bool disposing
) override
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
    true if managed resources should be disposed; otherwise, false.

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[Dispose 重载](f956a136-0f5a-dcd0-a9b8-7819ba4b63bd.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnShown 方法 

[原文連結](http://api.hslcommunication.cn/html/211bbe9f-1f05-37d9-2df1-5f0ac93474ed.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 方法](../html/b42f6f6e-241d-b7ef-f727-e75a68d8407e.htm "FormLogNetView 方法")

[Dispose 方法](../html/f956a136-0f5a-dcd0-a9b8-7819ba4b63bd.htm "Dispose 方法 ")

[OnShown 方法](../html/211bbe9f-1f05-37d9-2df1-5f0ac93474ed.htm "OnShown 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetViewOnShown 方法 |

[缺少 "M:HslCommunication.LogNet.FormLogNetView.OnShown(System.EventArgs)" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override void OnShown(
	EventArgs e
)
```

```
Protected Overrides Sub OnShown ( 
	e As EventArgs
)
```

```
protected:
virtual void OnShown(
	EventArgs^ e
) override
```

```
abstract OnShown : 
        e : EventArgs -> unit 
override OnShown : 
        e : EventArgs -> unit
```

#### 参数

e
:   类型：SystemEventArgs  

    [缺少 "M:HslCommunication.LogNet.FormLogNetView.OnShown(System.EventArgs)" 的 <param name="e"/> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FormLogNetView 事件

[原文連結](http://api.hslcommunication.cn/html/9a0e39b9-2117-2a86-ef9f-754fb5ee728f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[FormLogNetView 构造函数](../html/abc13d39-7d83-9e57-aadf-14131c4954a9.htm "FormLogNetView 构造函数 ")

[FormLogNetView 属性](../html/e39cf0a7-5e4c-0d82-39b6-11af4dcce5b6.htm "FormLogNetView 属性")

[FormLogNetView 方法](../html/b42f6f6e-241d-b7ef-f727-e75a68d8407e.htm "FormLogNetView 方法")

[FormLogNetView 事件](../html/9a0e39b9-2117-2a86-ef9f-754fb5ee728f.htm "FormLogNetView 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FormLogNetView 事件 |

[FormLogNetView](48091323-2df2-b52d-1420-224b3bf75e1b.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | Activated | (继承自 Form。) |
| 公共事件 | AutoSizeChanged | (继承自 Form。) |
| 公共事件 | AutoValidateChanged | (继承自 Form。) |
| 公共事件 | BackColorChanged | (继承自 Control。) |
| 公共事件 | BackgroundImageChanged | (继承自 Control。) |
| 公共事件 | BackgroundImageLayoutChanged | (继承自 Control。) |
| 公共事件 | BindingContextChanged | (继承自 Control。) |
| 公共事件 | CausesValidationChanged | (继承自 Control。) |
| 公共事件 | ChangeUICues | (继承自 Control。) |
| 公共事件 | Click | (继承自 Control。) |
| 公共事件 | ClientSizeChanged | (继承自 Control。) |
| 公共事件 | ContextMenuChanged | (继承自 Control。) |
| 公共事件 | ContextMenuStripChanged | (继承自 Control。) |
| 公共事件 | ControlAdded | (继承自 Control。) |
| 公共事件 | ControlRemoved | (继承自 Control。) |
| 公共事件 | CursorChanged | (继承自 Control。) |
| 公共事件 | Deactivate | (继承自 Form。) |
| 公共事件 | Disposed | (继承自 Component。) |
| 公共事件 | DockChanged | (继承自 Control。) |
| 公共事件 | DoubleClick | (继承自 Control。) |
| 公共事件 | DragDrop | (继承自 Control。) |
| 公共事件 | DragEnter | (继承自 Control。) |
| 公共事件 | DragLeave | (继承自 Control。) |
| 公共事件 | DragOver | (继承自 Control。) |
| 公共事件 | EnabledChanged | (继承自 Control。) |
| 公共事件 | Enter | (继承自 Control。) |
| 公共事件 | FontChanged | (继承自 Control。) |
| 公共事件 | ForeColorChanged | (继承自 Control。) |
| 公共事件 | FormClosed | (继承自 Form。) |
| 公共事件 | FormClosing | (继承自 Form。) |
| 公共事件 | GiveFeedback | (继承自 Control。) |
| 公共事件 | GotFocus | (继承自 Control。) |
| 公共事件 | HandleCreated | (继承自 Control。) |
| 公共事件 | HandleDestroyed | (继承自 Control。) |
| 公共事件 | HelpButtonClicked | (继承自 Form。) |
| 公共事件 | HelpRequested | (继承自 Control。) |
| 公共事件 | ImeModeChanged | (继承自 Control。) |
| 公共事件 | InputLanguageChanged | (继承自 Form。) |
| 公共事件 | InputLanguageChanging | (继承自 Form。) |
| 公共事件 | Invalidated | (继承自 Control。) |
| 公共事件 | KeyDown | (继承自 Control。) |
| 公共事件 | KeyPress | (继承自 Control。) |
| 公共事件 | KeyUp | (继承自 Control。) |
| 公共事件 | Layout | (继承自 Control。) |
| 公共事件 | Leave | (继承自 Control。) |
| 公共事件 | Load | (继承自 Form。) |
| 公共事件 | LocationChanged | (继承自 Control。) |
| 公共事件 | LostFocus | (继承自 Control。) |
| 公共事件 | MarginChanged | (继承自 Control。) |
| 公共事件 | MaximizedBoundsChanged | (继承自 Form。) |
| 公共事件 | MaximumSizeChanged | (继承自 Form。) |
| 公共事件 | MdiChildActivate | (继承自 Form。) |
| 公共事件 | MenuComplete | (继承自 Form。) |
| 公共事件 | MenuStart | (继承自 Form。) |
| 公共事件 | MinimumSizeChanged | (继承自 Form。) |
| 公共事件 | MouseCaptureChanged | (继承自 Control。) |
| 公共事件 | MouseClick | (继承自 Control。) |
| 公共事件 | MouseDoubleClick | (继承自 Control。) |
| 公共事件 | MouseDown | (继承自 Control。) |
| 公共事件 | MouseEnter | (继承自 Control。) |
| 公共事件 | MouseHover | (继承自 Control。) |
| 公共事件 | MouseLeave | (继承自 Control。) |
| 公共事件 | MouseMove | (继承自 Control。) |
| 公共事件 | MouseUp | (继承自 Control。) |
| 公共事件 | MouseWheel | (继承自 Control。) |
| 公共事件 | Move | (继承自 Control。) |
| 公共事件 | PaddingChanged | (继承自 Control。) |
| 公共事件 | Paint | (继承自 Control。) |
| 公共事件 | ParentChanged | (继承自 Control。) |
| 公共事件 | PreviewKeyDown | (继承自 Control。) |
| 公共事件 | QueryAccessibilityHelp | (继承自 Control。) |
| 公共事件 | QueryContinueDrag | (继承自 Control。) |
| 公共事件 | RegionChanged | (继承自 Control。) |
| 公共事件 | Resize | (继承自 Control。) |
| 公共事件 | ResizeBegin | (继承自 Form。) |
| 公共事件 | ResizeEnd | (继承自 Form。) |
| 公共事件 | RightToLeftChanged | (继承自 Control。) |
| 公共事件 | RightToLeftLayoutChanged | (继承自 Form。) |
| 公共事件 | Scroll | (继承自 ScrollableControl。) |
| 公共事件 | Shown | (继承自 Form。) |
| 公共事件 | SizeChanged | (继承自 Control。) |
| 公共事件 | StyleChanged | (继承自 Control。) |
| 公共事件 | SystemColorsChanged | (继承自 Control。) |
| 公共事件 | TabIndexChanged | (继承自 Control。) |
| 公共事件 | TabStopChanged | (继承自 Control。) |
| 公共事件 | TextChanged | (继承自 Control。) |
| 公共事件 | Validated | (继承自 Control。) |
| 公共事件 | Validating | (继承自 Control。) |
| 公共事件 | VisibleChanged | (继承自 Control。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FormLogNetView 类](48091323-2df2-b52d-1420-224b3bf75e1b.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GenerateMode 枚举

[原文連結](http://api.hslcommunication.cn/html/daf3c59a-f81d-0cdb-2da8-51b40dbb99e4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[GenerateMode 枚举](../html/daf3c59a-f81d-0cdb-2da8-51b40dbb99e4.htm "GenerateMode 枚举")

[HslEventArgs 类](../html/9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm "HslEventArgs 类")

[HslMessageDegree 枚举](../html/3ce37fee-7678-e278-6a0a-7656126003b6.htm "HslMessageDegree 枚举")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[ILogNet 接口](../html/d95d4704-db0d-b379-41e6-03879927a543.htm "ILogNet 接口")

[LogNetAnalysisControl 类](../html/cfda4292-54d6-ce11-0ada-b50021ca4ab7.htm "LogNetAnalysisControl 类")

[LogNetBase 类](../html/591f2a6c-f9cd-f53a-7abc-9d3c7adc563b.htm "LogNetBase 类")

[LogNetDateTime 类](../html/90b8130d-b017-4254-40f3-591d12226ddf.htm "LogNetDateTime 类")

[LogNetException 类](../html/30c57f15-6b01-d816-b1cc-29e7fed02717.htm "LogNetException 类")

[LogNetFileSize 类](../html/b36c7f80-f741-7a9a-e71a-dc2abb9714fe.htm "LogNetFileSize 类")

[LogNetManagment 类](../html/ad49f067-176c-6857-ff05-fb7a6f8fab36.htm "LogNetManagment 类")

[LogNetSingle 类](../html/2b98518d-c8fa-fc59-9900-3ccf00acfdba.htm "LogNetSingle 类")

[LogPathBase 类](../html/24a9b63f-53d8-a3aa-efc2-88ed24655ce0.htm "LogPathBase 类")

[LogSaveMode 枚举](../html/9ea4013b-4452-92f4-dd9a-3f2e12285927.htm "LogSaveMode 枚举")

[LogStatistics 类](../html/00e19edf-0fc1-28d3-17c5-43bf92ea094f.htm "LogStatistics 类")

[LogStatisticsBase(T) 类](../html/4b6cfa8a-7a7b-8151-eba6-e8bd57a62be5.htm "LogStatisticsBase(T) 类")

[LogStatisticsDict 类](../html/28d1bdaf-29f2-c51a-3e8b-0b41d1b98d15.htm "LogStatisticsDict 类")

[LogValueLimit 类](../html/4edfa9e7-8910-aaa0-566b-1fcbfa9c2b14.htm "LogValueLimit 类")

[LogValueLimitDict 类](../html/4e33e792-63f0-90f0-c872-9641fbb18154.htm "LogValueLimitDict 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GenerateMode 枚举 |

日志文件输出模式

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public enum GenerateMode
```

```
Public Enumeration GenerateMode
```

```
public enum class GenerateMode
```

```
type GenerateMode
```

![](../icons/SectionExpanded.png)成员

|  | 成员名称 | 值 | 说明 |
| --- | --- | --- | --- |
|  | ByEveryMinute | 1 | 按每分钟生成日志文件 |
|  | ByEveryHour | 2 | 按每个小时生成日志文件 |
|  | ByEveryDay | 3 | 按每天生成日志文件 |
|  | ByEveryWeek | 4 | 按每个周生成日志文件 |
|  | ByEveryMonth | 5 | 按每个月生成日志文件 |
|  | ByEverySeason | 6 | 按每季度生成日志文件 |
|  | ByEveryYear | 7 | 按每年生成日志文件 |

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslEventArgs 类

[原文連結](http://api.hslcommunication.cn/html/9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslEventArgs 类](../html/9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm "HslEventArgs 类")

[HslEventArgs 构造函数](../html/8f70428a-7f1a-e868-5186-f2d34115c5da.htm "HslEventArgs 构造函数 ")

[HslEventArgs 属性](../html/f89275c2-88a4-4307-5cea-5510212d7343.htm "HslEventArgs 属性")

[HslEventArgs 方法](../html/8f5f698b-2d2a-a533-e3c5-f5c4b57e6c3d.htm "HslEventArgs 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslEventArgs 类 |

带有日志消息的事件

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  SystemEventArgs  
    HslCommunication.LogNetHslEventArgs

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class HslEventArgs : EventArgs
```

```
Public Class HslEventArgs
	Inherits EventArgs
```

```
public ref class HslEventArgs : public EventArgs
```

```
type HslEventArgs =  
    class
        inherit EventArgs
    end
```

HslEventArgs 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HslEventArgs](8f70428a-7f1a-e868-5186-f2d34115c5da.htm) | 初始化 HslEventArgs 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [HslMessage](7de544e0-d319-bfe9-4730-427fc7d7dce0.htm) | 消息信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
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

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslEventArgs 构造函数 

[原文連結](http://api.hslcommunication.cn/html/8f70428a-7f1a-e868-5186-f2d34115c5da.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslEventArgs 类](../html/9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm "HslEventArgs 类")

[HslEventArgs 构造函数](../html/8f70428a-7f1a-e868-5186-f2d34115c5da.htm "HslEventArgs 构造函数 ")

[HslEventArgs 属性](../html/f89275c2-88a4-4307-5cea-5510212d7343.htm "HslEventArgs 属性")

[HslEventArgs 方法](../html/8f5f698b-2d2a-a533-e3c5-f5c4b57e6c3d.htm "HslEventArgs 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslEventArgs 构造函数 |

初始化 [HslEventArgs](9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslEventArgs()
```

```
Public Sub New
```

```
public:
HslEventArgs()
```

```
new : unit -> HslEventArgs
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslEventArgs 类](9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslEventArgs 属性

[原文連結](http://api.hslcommunication.cn/html/f89275c2-88a4-4307-5cea-5510212d7343.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslEventArgs 类](../html/9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm "HslEventArgs 类")

[HslEventArgs 属性](../html/f89275c2-88a4-4307-5cea-5510212d7343.htm "HslEventArgs 属性")

[HslMessage 属性](../html/7de544e0-d319-bfe9-4730-427fc7d7dce0.htm "HslMessage 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslEventArgs 属性 |

[HslEventArgs](9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [HslMessage](7de544e0-d319-bfe9-4730-427fc7d7dce0.htm) | 消息信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslEventArgs 类](9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslMessage 属性 

[原文連結](http://api.hslcommunication.cn/html/7de544e0-d319-bfe9-4730-427fc7d7dce0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslEventArgs 类](../html/9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm "HslEventArgs 类")

[HslEventArgs 属性](../html/f89275c2-88a4-4307-5cea-5510212d7343.htm "HslEventArgs 属性")

[HslMessage 属性](../html/7de544e0-d319-bfe9-4730-427fc7d7dce0.htm "HslMessage 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslEventArgsHslMessage 属性 |

消息信息

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslMessageItem HslMessage { get; set; }
```

```
Public Property HslMessage As HslMessageItem
	Get
	Set
```

```
public:
property HslMessageItem^ HslMessage {
	HslMessageItem^ get ();
	void set (HslMessageItem^ value);
}
```

```
member HslMessage : HslMessageItem with get, set
```

#### 属性值

类型：[HslMessageItem](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslEventArgs 类](9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslEventArgs 方法

[原文連結](http://api.hslcommunication.cn/html/8f5f698b-2d2a-a533-e3c5-f5c4b57e6c3d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslEventArgs 类](../html/9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm "HslEventArgs 类")

[HslEventArgs 构造函数](../html/8f70428a-7f1a-e868-5186-f2d34115c5da.htm "HslEventArgs 构造函数 ")

[HslEventArgs 属性](../html/f89275c2-88a4-4307-5cea-5510212d7343.htm "HslEventArgs 属性")

[HslEventArgs 方法](../html/8f5f698b-2d2a-a533-e3c5-f5c4b57e6c3d.htm "HslEventArgs 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslEventArgs 方法 |

[HslEventArgs](9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
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

[HslEventArgs 类](9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslMessageDegree 枚举

[原文連結](http://api.hslcommunication.cn/html/3ce37fee-7678-e278-6a0a-7656126003b6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[FormLogNetView 类](../html/48091323-2df2-b52d-1420-224b3bf75e1b.htm "FormLogNetView 类")

[GenerateMode 枚举](../html/daf3c59a-f81d-0cdb-2da8-51b40dbb99e4.htm "GenerateMode 枚举")

[HslEventArgs 类](../html/9c19de27-52cd-8a37-3a55-2fb1e4c63854.htm "HslEventArgs 类")

[HslMessageDegree 枚举](../html/3ce37fee-7678-e278-6a0a-7656126003b6.htm "HslMessageDegree 枚举")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[ILogNet 接口](../html/d95d4704-db0d-b379-41e6-03879927a543.htm "ILogNet 接口")

[LogNetAnalysisControl 类](../html/cfda4292-54d6-ce11-0ada-b50021ca4ab7.htm "LogNetAnalysisControl 类")

[LogNetBase 类](../html/591f2a6c-f9cd-f53a-7abc-9d3c7adc563b.htm "LogNetBase 类")

[LogNetDateTime 类](../html/90b8130d-b017-4254-40f3-591d12226ddf.htm "LogNetDateTime 类")

[LogNetException 类](../html/30c57f15-6b01-d816-b1cc-29e7fed02717.htm "LogNetException 类")

[LogNetFileSize 类](../html/b36c7f80-f741-7a9a-e71a-dc2abb9714fe.htm "LogNetFileSize 类")

[LogNetManagment 类](../html/ad49f067-176c-6857-ff05-fb7a6f8fab36.htm "LogNetManagment 类")

[LogNetSingle 类](../html/2b98518d-c8fa-fc59-9900-3ccf00acfdba.htm "LogNetSingle 类")

[LogPathBase 类](../html/24a9b63f-53d8-a3aa-efc2-88ed24655ce0.htm "LogPathBase 类")

[LogSaveMode 枚举](../html/9ea4013b-4452-92f4-dd9a-3f2e12285927.htm "LogSaveMode 枚举")

[LogStatistics 类](../html/00e19edf-0fc1-28d3-17c5-43bf92ea094f.htm "LogStatistics 类")

[LogStatisticsBase(T) 类](../html/4b6cfa8a-7a7b-8151-eba6-e8bd57a62be5.htm "LogStatisticsBase(T) 类")

[LogStatisticsDict 类](../html/28d1bdaf-29f2-c51a-3e8b-0b41d1b98d15.htm "LogStatisticsDict 类")

[LogValueLimit 类](../html/4edfa9e7-8910-aaa0-566b-1fcbfa9c2b14.htm "LogValueLimit 类")

[LogValueLimitDict 类](../html/4e33e792-63f0-90f0-c872-9641fbb18154.htm "LogValueLimitDict 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageDegree 枚举 |

记录消息的等级

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public enum HslMessageDegree
```

```
Public Enumeration HslMessageDegree
```

```
public enum class HslMessageDegree
```

```
type HslMessageDegree
```

![](../icons/SectionExpanded.png)成员

|  | 成员名称 | 值 | 说明 |
| --- | --- | --- | --- |
|  | None | 1 | 一条消息都不记录 |
|  | FATAL | 2 | 记录致命等级及以上日志的消息 |
|  | ERROR | 3 | 记录异常等级及以上日志的消息 |
|  | WARN | 4 | 记录警告等级及以上日志的消息 |
|  | INFO | 5 | 记录信息等级及以上日志的消息 |
|  | DEBUG | 6 | 记录调试等级及以上日志的信息 |

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslMessageItem 类

[原文連結](http://api.hslcommunication.cn/html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 构造函数](../html/45da0987-0e84-e2f6-15e2-1cd64691ff07.htm "HslMessageItem 构造函数 ")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[HslMessageItem 方法](../html/0f4067df-4673-373e-197d-ae07692704b9.htm "HslMessageItem 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItem 类 |

单条日志的记录信息，包含了消息等级，线程号，关键字，文本信息  
Record information of a single log, including message level, thread number, keywords, text information

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.LogNetHslMessageItem

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class HslMessageItem
```

```
Public Class HslMessageItem
```

```
public ref class HslMessageItem
```

```
type HslMessageItem =  class end
```

HslMessageItem 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HslMessageItem](45da0987-0e84-e2f6-15e2-1cd64691ff07.htm) | 实例化一个默认的对象  Instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Cancel](cf183617-0736-3a2d-e117-6b47bfe392c8.htm) | 是否取消写入到文件中去，在事件 [BeforeSaveToFile](66e453d5-37db-e790-a2b2-a15e4d3484c9.htm) 触发的时候捕获即可设置。  Whether to cancel writing to the file, can be set when the event [BeforeSaveToFile](66e453d5-37db-e790-a2b2-a15e4d3484c9.htm) is triggered. |
| 公共属性 | [Degree](d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm) | 消息的等级，包括DEBUG，INFO，WARN，ERROR，FATAL，NONE共计六个等级  Message levels, including DEBUG, INFO, WARN, ERROR, FATAL, NONE total six levels |
| 公共属性 | [Id](3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm) | 单个记录信息的标识ID，程序重新运行时清空，代表程序从运行以来的日志计数，不管存储的或是未存储的  The ID of a single record of information. It is cleared when the program is re-run. It represents the log count of the program since it was run, whether stored or unstored. |
| 公共属性 | [KeyWord](d033d1d2-5115-2210-494b-6a05f9881ea2.htm) | 消息的关键字  Keyword of the message |
| 公共属性 | [Text](d143ad83-2190-7f80-17a0-faf9aa5d230e.htm) | 消息文本，记录日志的时候给定  Message text, given when logging |
| 公共属性 | [ThreadId](15ffacd5-2269-bae9-e19b-33a02c03ce69.htm) | 线程ID，发生异常时的线程号  Thread ID, the thread number when the exception occurred |
| 公共属性 | [Time](b3460e41-6401-86cc-e04b-a76b856e8d3c.htm) | 记录日志的时间，而非存储日志的时间  The time the log was recorded, not the time it was stored |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](fcb2ad14-3cef-6539-eff9-5244d5236255.htm) | (重写 ObjectToString.) |
| 公共方法 | [ToStringWithoutKeyword](4df35b59-7a90-1dcb-fcdc-ad321f2bd3b1.htm) | 返回表示当前对象的字符串，剔除了关键字  Returns a string representing the current object, excluding keywords |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslMessageItem 构造函数 

[原文連結](http://api.hslcommunication.cn/html/45da0987-0e84-e2f6-15e2-1cd64691ff07.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 构造函数](../html/45da0987-0e84-e2f6-15e2-1cd64691ff07.htm "HslMessageItem 构造函数 ")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[HslMessageItem 方法](../html/0f4067df-4673-373e-197d-ae07692704b9.htm "HslMessageItem 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItem 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslMessageItem()
```

```
Public Sub New
```

```
public:
HslMessageItem()
```

```
new : unit -> HslMessageItem
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslMessageItem 属性

[原文連結](http://api.hslcommunication.cn/html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[Cancel 属性](../html/cf183617-0736-3a2d-e117-6b47bfe392c8.htm "Cancel 属性 ")

[Degree 属性](../html/d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm "Degree 属性 ")

[Id 属性](../html/3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm "Id 属性 ")

[KeyWord 属性](../html/d033d1d2-5115-2210-494b-6a05f9881ea2.htm "KeyWord 属性 ")

[Text 属性](../html/d143ad83-2190-7f80-17a0-faf9aa5d230e.htm "Text 属性 ")

[ThreadId 属性](../html/15ffacd5-2269-bae9-e19b-33a02c03ce69.htm "ThreadId 属性 ")

[Time 属性](../html/b3460e41-6401-86cc-e04b-a76b856e8d3c.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItem 属性 |

[HslMessageItem](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Cancel](cf183617-0736-3a2d-e117-6b47bfe392c8.htm) | 是否取消写入到文件中去，在事件 [BeforeSaveToFile](66e453d5-37db-e790-a2b2-a15e4d3484c9.htm) 触发的时候捕获即可设置。  Whether to cancel writing to the file, can be set when the event [BeforeSaveToFile](66e453d5-37db-e790-a2b2-a15e4d3484c9.htm) is triggered. |
| 公共属性 | [Degree](d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm) | 消息的等级，包括DEBUG，INFO，WARN，ERROR，FATAL，NONE共计六个等级  Message levels, including DEBUG, INFO, WARN, ERROR, FATAL, NONE total six levels |
| 公共属性 | [Id](3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm) | 单个记录信息的标识ID，程序重新运行时清空，代表程序从运行以来的日志计数，不管存储的或是未存储的  The ID of a single record of information. It is cleared when the program is re-run. It represents the log count of the program since it was run, whether stored or unstored. |
| 公共属性 | [KeyWord](d033d1d2-5115-2210-494b-6a05f9881ea2.htm) | 消息的关键字  Keyword of the message |
| 公共属性 | [Text](d143ad83-2190-7f80-17a0-faf9aa5d230e.htm) | 消息文本，记录日志的时候给定  Message text, given when logging |
| 公共属性 | [ThreadId](15ffacd5-2269-bae9-e19b-33a02c03ce69.htm) | 线程ID，发生异常时的线程号  Thread ID, the thread number when the exception occurred |
| 公共属性 | [Time](b3460e41-6401-86cc-e04b-a76b856e8d3c.htm) | 记录日志的时间，而非存储日志的时间  The time the log was recorded, not the time it was stored |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Cancel 属性 

[原文連結](http://api.hslcommunication.cn/html/cf183617-0736-3a2d-e117-6b47bfe392c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[Cancel 属性](../html/cf183617-0736-3a2d-e117-6b47bfe392c8.htm "Cancel 属性 ")

[Degree 属性](../html/d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm "Degree 属性 ")

[Id 属性](../html/3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm "Id 属性 ")

[KeyWord 属性](../html/d033d1d2-5115-2210-494b-6a05f9881ea2.htm "KeyWord 属性 ")

[Text 属性](../html/d143ad83-2190-7f80-17a0-faf9aa5d230e.htm "Text 属性 ")

[ThreadId 属性](../html/15ffacd5-2269-bae9-e19b-33a02c03ce69.htm "ThreadId 属性 ")

[Time 属性](../html/b3460e41-6401-86cc-e04b-a76b856e8d3c.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItemCancel 属性 |

是否取消写入到文件中去，在事件 [BeforeSaveToFile](66e453d5-37db-e790-a2b2-a15e4d3484c9.htm) 触发的时候捕获即可设置。  
Whether to cancel writing to the file, can be set when the event [BeforeSaveToFile](66e453d5-37db-e790-a2b2-a15e4d3484c9.htm) is triggered.

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool Cancel { get; set; }
```

```
Public Property Cancel As Boolean
	Get
	Set
```

```
public:
property bool Cancel {
	bool get ();
	void set (bool value);
}
```

```
member Cancel : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Degree 属性 

[原文連結](http://api.hslcommunication.cn/html/d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[Cancel 属性](../html/cf183617-0736-3a2d-e117-6b47bfe392c8.htm "Cancel 属性 ")

[Degree 属性](../html/d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm "Degree 属性 ")

[Id 属性](../html/3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm "Id 属性 ")

[KeyWord 属性](../html/d033d1d2-5115-2210-494b-6a05f9881ea2.htm "KeyWord 属性 ")

[Text 属性](../html/d143ad83-2190-7f80-17a0-faf9aa5d230e.htm "Text 属性 ")

[ThreadId 属性](../html/15ffacd5-2269-bae9-e19b-33a02c03ce69.htm "ThreadId 属性 ")

[Time 属性](../html/b3460e41-6401-86cc-e04b-a76b856e8d3c.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItemDegree 属性 |

消息的等级，包括DEBUG，INFO，WARN，ERROR，FATAL，NONE共计六个等级  
Message levels, including DEBUG, INFO, WARN, ERROR, FATAL, NONE total six levels

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslMessageDegree Degree { get; set; }
```

```
Public Property Degree As HslMessageDegree
	Get
	Set
```

```
public:
property HslMessageDegree Degree {
	HslMessageDegree get ();
	void set (HslMessageDegree value);
}
```

```
member Degree : HslMessageDegree with get, set
```

#### 属性值

类型：[HslMessageDegree](3ce37fee-7678-e278-6a0a-7656126003b6.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Id 属性 

[原文連結](http://api.hslcommunication.cn/html/3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[Cancel 属性](../html/cf183617-0736-3a2d-e117-6b47bfe392c8.htm "Cancel 属性 ")

[Degree 属性](../html/d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm "Degree 属性 ")

[Id 属性](../html/3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm "Id 属性 ")

[KeyWord 属性](../html/d033d1d2-5115-2210-494b-6a05f9881ea2.htm "KeyWord 属性 ")

[Text 属性](../html/d143ad83-2190-7f80-17a0-faf9aa5d230e.htm "Text 属性 ")

[ThreadId 属性](../html/15ffacd5-2269-bae9-e19b-33a02c03ce69.htm "ThreadId 属性 ")

[Time 属性](../html/b3460e41-6401-86cc-e04b-a76b856e8d3c.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItemId 属性 |

单个记录信息的标识ID，程序重新运行时清空，代表程序从运行以来的日志计数，不管存储的或是未存储的  
The ID of a single record of information. It is cleared when the program is re-run.
It represents the log count of the program since it was run, whether stored or unstored.

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public long Id { get; }
```

```
Public ReadOnly Property Id As Long
	Get
```

```
public:
property long long Id {
	long long get ();
}
```

```
member Id : int64 with get
```

#### 属性值

类型：Int64

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyWord 属性 

[原文連結](http://api.hslcommunication.cn/html/d033d1d2-5115-2210-494b-6a05f9881ea2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[Cancel 属性](../html/cf183617-0736-3a2d-e117-6b47bfe392c8.htm "Cancel 属性 ")

[Degree 属性](../html/d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm "Degree 属性 ")

[Id 属性](../html/3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm "Id 属性 ")

[KeyWord 属性](../html/d033d1d2-5115-2210-494b-6a05f9881ea2.htm "KeyWord 属性 ")

[Text 属性](../html/d143ad83-2190-7f80-17a0-faf9aa5d230e.htm "Text 属性 ")

[ThreadId 属性](../html/15ffacd5-2269-bae9-e19b-33a02c03ce69.htm "ThreadId 属性 ")

[Time 属性](../html/b3460e41-6401-86cc-e04b-a76b856e8d3c.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItemKeyWord 属性 |

消息的关键字  
Keyword of the message

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string KeyWord { get; set; }
```

```
Public Property KeyWord As String
	Get
	Set
```

```
public:
property String^ KeyWord {
	String^ get ();
	void set (String^ value);
}
```

```
member KeyWord : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Text 属性 

[原文連結](http://api.hslcommunication.cn/html/d143ad83-2190-7f80-17a0-faf9aa5d230e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[Cancel 属性](../html/cf183617-0736-3a2d-e117-6b47bfe392c8.htm "Cancel 属性 ")

[Degree 属性](../html/d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm "Degree 属性 ")

[Id 属性](../html/3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm "Id 属性 ")

[KeyWord 属性](../html/d033d1d2-5115-2210-494b-6a05f9881ea2.htm "KeyWord 属性 ")

[Text 属性](../html/d143ad83-2190-7f80-17a0-faf9aa5d230e.htm "Text 属性 ")

[ThreadId 属性](../html/15ffacd5-2269-bae9-e19b-33a02c03ce69.htm "ThreadId 属性 ")

[Time 属性](../html/b3460e41-6401-86cc-e04b-a76b856e8d3c.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItemText 属性 |

消息文本，记录日志的时候给定  
Message text, given when logging

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Text { get; set; }
```

```
Public Property Text As String
	Get
	Set
```

```
public:
property String^ Text {
	String^ get ();
	void set (String^ value);
}
```

```
member Text : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ThreadId 属性 

[原文連結](http://api.hslcommunication.cn/html/15ffacd5-2269-bae9-e19b-33a02c03ce69.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[Cancel 属性](../html/cf183617-0736-3a2d-e117-6b47bfe392c8.htm "Cancel 属性 ")

[Degree 属性](../html/d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm "Degree 属性 ")

[Id 属性](../html/3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm "Id 属性 ")

[KeyWord 属性](../html/d033d1d2-5115-2210-494b-6a05f9881ea2.htm "KeyWord 属性 ")

[Text 属性](../html/d143ad83-2190-7f80-17a0-faf9aa5d230e.htm "Text 属性 ")

[ThreadId 属性](../html/15ffacd5-2269-bae9-e19b-33a02c03ce69.htm "ThreadId 属性 ")

[Time 属性](../html/b3460e41-6401-86cc-e04b-a76b856e8d3c.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItemThreadId 属性 |

线程ID，发生异常时的线程号  
Thread ID, the thread number when the exception occurred

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ThreadId { get; set; }
```

```
Public Property ThreadId As Integer
	Get
	Set
```

```
public:
property int ThreadId {
	int get ();
	void set (int value);
}
```

```
member ThreadId : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Time 属性 

[原文連結](http://api.hslcommunication.cn/html/b3460e41-6401-86cc-e04b-a76b856e8d3c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 属性](../html/c941ae70-3cb5-a940-34b0-54dd8633f1d7.htm "HslMessageItem 属性")

[Cancel 属性](../html/cf183617-0736-3a2d-e117-6b47bfe392c8.htm "Cancel 属性 ")

[Degree 属性](../html/d6c7b854-6daf-2bfd-e423-b82cb2898c00.htm "Degree 属性 ")

[Id 属性](../html/3a0bfa8c-ed3d-a0d9-d3a8-48f561155723.htm "Id 属性 ")

[KeyWord 属性](../html/d033d1d2-5115-2210-494b-6a05f9881ea2.htm "KeyWord 属性 ")

[Text 属性](../html/d143ad83-2190-7f80-17a0-faf9aa5d230e.htm "Text 属性 ")

[ThreadId 属性](../html/15ffacd5-2269-bae9-e19b-33a02c03ce69.htm "ThreadId 属性 ")

[Time 属性](../html/b3460e41-6401-86cc-e04b-a76b856e8d3c.htm "Time 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItemTime 属性 |

记录日志的时间，而非存储日志的时间  
The time the log was recorded, not the time it was stored

**命名空间：**
 [HslCommunication.LogNet](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)  
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

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslMessageItem 方法

[原文連結](http://api.hslcommunication.cn/html/0f4067df-4673-373e-197d-ae07692704b9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.LogNet](../html/ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm "HslCommunication.LogNet")

[HslMessageItem 类](../html/7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm "HslMessageItem 类")

[HslMessageItem 方法](../html/0f4067df-4673-373e-197d-ae07692704b9.htm "HslMessageItem 方法")

[ToString 方法](../html/fcb2ad14-3cef-6539-eff9-5244d5236255.htm "ToString 方法 ")

[ToStringWithoutKeyword 方法](../html/4df35b59-7a90-1dcb-fcdc-ad321f2bd3b1.htm "ToStringWithoutKeyword 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslMessageItem 方法 |

[HslMessageItem](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](fcb2ad14-3cef-6539-eff9-5244d5236255.htm) | (重写 ObjectToString.) |
| 公共方法 | [ToStringWithoutKeyword](4df35b59-7a90-1dcb-fcdc-ad321f2bd3b1.htm) | 返回表示当前对象的字符串，剔除了关键字  Returns a string representing the current object, excluding keywords |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslMessageItem 类](7bcefc29-8bc7-5727-3757-4a3530a3ed53.htm)

[HslCommunication.LogNet 命名空间](ee5b896e-8dba-c0c2-1020-0569b6d56f8e.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)