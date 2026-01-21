# HslCommunication - HslCommunication.Robot.EFORT

> 分類頁數: 30



---
## HslCommunication.Robot.EFORT

[原文連結](http://api.hslcommunication.cn/html/a25bf846-86de-f385-03a1-570078f0c95b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[ER7BC10 类](../html/22bc0593-b20d-8f6d-bcdb-ddfb6a6085cd.htm "ER7BC10 类")

[ER7BC10Previous 类](../html/cf4dc864-ff47-c7a9-1473-5f9c2ef80529.htm "ER7BC10Previous 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Robot.EFORT 命名空间 |

[缺少 "N:HslCommunication.Robot.EFORT" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [EfortData](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm) | 埃夫特机器人的数据结构  The data structure of the efort robot |
| 公共类 | [ER7BC10](22bc0593-b20d-8f6d-bcdb-ddfb6a6085cd.htm) | 埃夫特机器人对应型号为ER7B-C10，此协议为定制版，新版报文对齐  The corresponding model of efort robot is er7b-c10. This protocol is the customized version, and the new version is the message alignment |
| 公共类 | [ER7BC10Previous](cf4dc864-ff47-c7a9-1473-5f9c2ef80529.htm) | 埃夫特机器人对应型号为ER7B-C10，此协议为旧版的定制版，报文未对齐的版本  The corresponding model of the efort robot is er7b-c10. This protocol is a customized version of the old version, and the message is not aligned |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EfortData 类

[原文連結](http://api.hslcommunication.cn/html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 构造函数](../html/155d2b48-035c-59a4-7798-d7ad70205ed0.htm "EfortData 构造函数 ")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[EfortData 方法](../html/6bf8cb08-fdd5-ce76-b5af-cb366990f8f4.htm "EfortData 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortData 类 |

埃夫特机器人的数据结构  
The data structure of the efort robot

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Robot.EFORTEfortData

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class EfortData
```

```
Public Class EfortData
```

```
public ref class EfortData
```

```
type EfortData =  class end
```

EfortData 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [EfortData](155d2b48-035c-59a4-7798-d7ad70205ed0.htm) | 实例化一个默认的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AuthorityStatus](97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm) | 权限状态，1：有权限，0：无权限 |
| 公共属性 | [AxisMoveStatus](e6082238-491c-708a-09b2-b3e1994980f3.htm) | 轴运动状态，1：有运动，0：未运动 |
| 公共属性 | [DbAxisAcc](4fc653f6-1979-57c5-480b-fd184236e397.htm) | 一到七轴的加速度 |
| 公共属性 | [DbAxisAccAcc](88bac00f-0011-58c5-973d-c582220bc9ce.htm) | 一到七轴的加加速度 |
| 公共属性 | [DbAxisDirCnt](b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm) | 轴反向计数 |
| 公共属性 | [DbAxisPos](2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm) | 一到七轴的角度 |
| 公共属性 | [DbAxisSpeed](047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm) | 一到七轴的速度 |
| 公共属性 | [DbAxisTime](ed4d34c0-894d-d267-cf3e-fac0676797c4.htm) | 轴工作总时长 |
| 公共属性 | [DbAxisTorque](7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm) | 一到七轴的力矩 |
| 公共属性 | [DbCartPos](f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm) | X,Y,Z,A,B,C方向，也叫笛卡尔坐标系 |
| 公共属性 | [DbDeviceTime](a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm) | 设备开机总时长 |
| 公共属性 | [ErrorStatus](d28c3c3e-b4ea-8064-3d62-10554dd80555.htm) | 报警状态，1:有报警，0：无报警 |
| 公共属性 | [ErrorText](bfd84369-977c-f9b1-da1d-e5d2722227d4.htm) | 错误信息 |
| 公共属性 | [HstopStatus](b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm) | 急停状态，1：无急停，0：有急停 |
| 公共属性 | [IoDIn](056aa978-0740-28a1-7591-ae38d28234c8.htm) | IoDIn状态 |
| 公共属性 | [IoDOut](5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm) | IoDOut状态 |
| 公共属性 | [IoIIn](8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm) | IoIIn状态 |
| 公共属性 | [IoIOut](14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm) | IoIOut状态 |
| 公共属性 | [ModeStatus](2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm) | 模式状态，1:手动，2:自动，3:远程 |
| 公共属性 | [PacketEnd](f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm) | 报文结束标记 |
| 公共属性 | [PacketHeartbeat](0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm) | 数据心跳 |
| 公共属性 | [PacketOrders](f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm) | 数据命令 |
| 公共属性 | [PacketStart](fa9c65dc-ee91-8616-6b04-bfa89285d082.htm) | 报文开始的字符串 |
| 公共属性 | [ProgHoldStatus](f0e77644-2606-422c-a2da-06aaadaf9d3e.htm) | 程序暂停状态，1：有暂停，0：无暂停 |
| 公共属性 | [ProgLoadStatus](795d513d-df00-ede3-898d-c8ee109bd3d1.htm) | 程序加载状态，1：有加载，0：无加载 |
| 公共属性 | [ProgMoveStatus](e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm) | 程序运行状态，1：有运行，0：未运行 |
| 公共属性 | [ProgramName](ae5a4e14-1c87-0710-60ef-3f8108efe575.htm) | 加载程序名 |
| 公共属性 | [ProjectName](3aadbce0-04f2-4ec8-e804-e628791978e8.htm) | 加载工程名 |
| 公共属性 | [ServoStatus](e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm) | 伺服状态，1：有使能，0：未使能 |
| 公共属性 | [SpeedStatus](f4f31109-f088-53ee-9496-0d577d350788.htm) | 读读状态，百分比（单位） |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PraseFrom](c6dacc82-4fc5-da4f-9584-9d8df17e20ac.htm) | 从新版本数据构造一个埃夫特机器人的数据类型 |
| 公共方法静态成员 | [PraseFromPrevious](7bfca75d-d243-f38a-02d8-da59b0a9dada.htm) | 从之前的版本数据构造一个埃夫特机器人的数据类型 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EfortData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/155d2b48-035c-59a4-7798-d7ad70205ed0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 构造函数](../html/155d2b48-035c-59a4-7798-d7ad70205ed0.htm "EfortData 构造函数 ")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[EfortData 方法](../html/6bf8cb08-fdd5-ce76-b5af-cb366990f8f4.htm "EfortData 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortData 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public EfortData()
```

```
Public Sub New
```

```
public:
EfortData()
```

```
new : unit -> EfortData
```

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EfortData 属性

[原文連結](http://api.hslcommunication.cn/html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortData 属性 |

[EfortData](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AuthorityStatus](97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm) | 权限状态，1：有权限，0：无权限 |
| 公共属性 | [AxisMoveStatus](e6082238-491c-708a-09b2-b3e1994980f3.htm) | 轴运动状态，1：有运动，0：未运动 |
| 公共属性 | [DbAxisAcc](4fc653f6-1979-57c5-480b-fd184236e397.htm) | 一到七轴的加速度 |
| 公共属性 | [DbAxisAccAcc](88bac00f-0011-58c5-973d-c582220bc9ce.htm) | 一到七轴的加加速度 |
| 公共属性 | [DbAxisDirCnt](b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm) | 轴反向计数 |
| 公共属性 | [DbAxisPos](2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm) | 一到七轴的角度 |
| 公共属性 | [DbAxisSpeed](047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm) | 一到七轴的速度 |
| 公共属性 | [DbAxisTime](ed4d34c0-894d-d267-cf3e-fac0676797c4.htm) | 轴工作总时长 |
| 公共属性 | [DbAxisTorque](7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm) | 一到七轴的力矩 |
| 公共属性 | [DbCartPos](f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm) | X,Y,Z,A,B,C方向，也叫笛卡尔坐标系 |
| 公共属性 | [DbDeviceTime](a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm) | 设备开机总时长 |
| 公共属性 | [ErrorStatus](d28c3c3e-b4ea-8064-3d62-10554dd80555.htm) | 报警状态，1:有报警，0：无报警 |
| 公共属性 | [ErrorText](bfd84369-977c-f9b1-da1d-e5d2722227d4.htm) | 错误信息 |
| 公共属性 | [HstopStatus](b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm) | 急停状态，1：无急停，0：有急停 |
| 公共属性 | [IoDIn](056aa978-0740-28a1-7591-ae38d28234c8.htm) | IoDIn状态 |
| 公共属性 | [IoDOut](5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm) | IoDOut状态 |
| 公共属性 | [IoIIn](8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm) | IoIIn状态 |
| 公共属性 | [IoIOut](14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm) | IoIOut状态 |
| 公共属性 | [ModeStatus](2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm) | 模式状态，1:手动，2:自动，3:远程 |
| 公共属性 | [PacketEnd](f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm) | 报文结束标记 |
| 公共属性 | [PacketHeartbeat](0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm) | 数据心跳 |
| 公共属性 | [PacketOrders](f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm) | 数据命令 |
| 公共属性 | [PacketStart](fa9c65dc-ee91-8616-6b04-bfa89285d082.htm) | 报文开始的字符串 |
| 公共属性 | [ProgHoldStatus](f0e77644-2606-422c-a2da-06aaadaf9d3e.htm) | 程序暂停状态，1：有暂停，0：无暂停 |
| 公共属性 | [ProgLoadStatus](795d513d-df00-ede3-898d-c8ee109bd3d1.htm) | 程序加载状态，1：有加载，0：无加载 |
| 公共属性 | [ProgMoveStatus](e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm) | 程序运行状态，1：有运行，0：未运行 |
| 公共属性 | [ProgramName](ae5a4e14-1c87-0710-60ef-3f8108efe575.htm) | 加载程序名 |
| 公共属性 | [ProjectName](3aadbce0-04f2-4ec8-e804-e628791978e8.htm) | 加载工程名 |
| 公共属性 | [ServoStatus](e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm) | 伺服状态，1：有使能，0：未使能 |
| 公共属性 | [SpeedStatus](f4f31109-f088-53ee-9496-0d577d350788.htm) | 读读状态，百分比（单位） |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AuthorityStatus 属性 

[原文連結](http://api.hslcommunication.cn/html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataAuthorityStatus 属性 |

权限状态，1：有权限，0：无权限

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte AuthorityStatus { get; set; }
```

```
Public Property AuthorityStatus As Byte
	Get
	Set
```

```
public:
property unsigned char AuthorityStatus {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member AuthorityStatus : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AxisMoveStatus 属性 

[原文連結](http://api.hslcommunication.cn/html/e6082238-491c-708a-09b2-b3e1994980f3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataAxisMoveStatus 属性 |

轴运动状态，1：有运动，0：未运动

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte AxisMoveStatus { get; set; }
```

```
Public Property AxisMoveStatus As Byte
	Get
	Set
```

```
public:
property unsigned char AxisMoveStatus {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member AxisMoveStatus : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DbAxisAcc 属性 

[原文連結](http://api.hslcommunication.cn/html/4fc653f6-1979-57c5-480b-fd184236e397.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataDbAxisAcc 属性 |

一到七轴的加速度

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public float[] DbAxisAcc { get; set; }
```

```
Public Property DbAxisAcc As Single()
	Get
	Set
```

```
public:
property array<float>^ DbAxisAcc {
	array<float>^ get ();
	void set (array<float>^ value);
}
```

```
member DbAxisAcc : float32[] with get, set
```

#### 属性值

类型：Single

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DbAxisAccAcc 属性 

[原文連結](http://api.hslcommunication.cn/html/88bac00f-0011-58c5-973d-c582220bc9ce.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataDbAxisAccAcc 属性 |

一到七轴的加加速度

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public float[] DbAxisAccAcc { get; set; }
```

```
Public Property DbAxisAccAcc As Single()
	Get
	Set
```

```
public:
property array<float>^ DbAxisAccAcc {
	array<float>^ get ();
	void set (array<float>^ value);
}
```

```
member DbAxisAccAcc : float32[] with get, set
```

#### 属性值

类型：Single

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DbAxisDirCnt 属性 

[原文連結](http://api.hslcommunication.cn/html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataDbAxisDirCnt 属性 |

轴反向计数

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int[] DbAxisDirCnt { get; set; }
```

```
Public Property DbAxisDirCnt As Integer()
	Get
	Set
```

```
public:
property array<int>^ DbAxisDirCnt {
	array<int>^ get ();
	void set (array<int>^ value);
}
```

```
member DbAxisDirCnt : int[] with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DbAxisPos 属性 

[原文連結](http://api.hslcommunication.cn/html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataDbAxisPos 属性 |

一到七轴的角度

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public float[] DbAxisPos { get; set; }
```

```
Public Property DbAxisPos As Single()
	Get
	Set
```

```
public:
property array<float>^ DbAxisPos {
	array<float>^ get ();
	void set (array<float>^ value);
}
```

```
member DbAxisPos : float32[] with get, set
```

#### 属性值

类型：Single

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DbAxisSpeed 属性 

[原文連結](http://api.hslcommunication.cn/html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataDbAxisSpeed 属性 |

一到七轴的速度

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public float[] DbAxisSpeed { get; set; }
```

```
Public Property DbAxisSpeed As Single()
	Get
	Set
```

```
public:
property array<float>^ DbAxisSpeed {
	array<float>^ get ();
	void set (array<float>^ value);
}
```

```
member DbAxisSpeed : float32[] with get, set
```

#### 属性值

类型：Single

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DbAxisTime 属性 

[原文連結](http://api.hslcommunication.cn/html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataDbAxisTime 属性 |

轴工作总时长

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int[] DbAxisTime { get; set; }
```

```
Public Property DbAxisTime As Integer()
	Get
	Set
```

```
public:
property array<int>^ DbAxisTime {
	array<int>^ get ();
	void set (array<int>^ value);
}
```

```
member DbAxisTime : int[] with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DbAxisTorque 属性 

[原文連結](http://api.hslcommunication.cn/html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataDbAxisTorque 属性 |

一到七轴的力矩

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public float[] DbAxisTorque { get; set; }
```

```
Public Property DbAxisTorque As Single()
	Get
	Set
```

```
public:
property array<float>^ DbAxisTorque {
	array<float>^ get ();
	void set (array<float>^ value);
}
```

```
member DbAxisTorque : float32[] with get, set
```

#### 属性值

类型：Single

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DbCartPos 属性 

[原文連結](http://api.hslcommunication.cn/html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataDbCartPos 属性 |

X,Y,Z,A,B,C方向，也叫笛卡尔坐标系

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public float[] DbCartPos { get; set; }
```

```
Public Property DbCartPos As Single()
	Get
	Set
```

```
public:
property array<float>^ DbCartPos {
	array<float>^ get ();
	void set (array<float>^ value);
}
```

```
member DbCartPos : float32[] with get, set
```

#### 属性值

类型：Single

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DbDeviceTime 属性 

[原文連結](http://api.hslcommunication.cn/html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataDbDeviceTime 属性 |

设备开机总时长

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int DbDeviceTime { get; set; }
```

```
Public Property DbDeviceTime As Integer
	Get
	Set
```

```
public:
property int DbDeviceTime {
	int get ();
	void set (int value);
}
```

```
member DbDeviceTime : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ErrorStatus 属性 

[原文連結](http://api.hslcommunication.cn/html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataErrorStatus 属性 |

报警状态，1:有报警，0：无报警

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte ErrorStatus { get; set; }
```

```
Public Property ErrorStatus As Byte
	Get
	Set
```

```
public:
property unsigned char ErrorStatus {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member ErrorStatus : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ErrorText 属性 

[原文連結](http://api.hslcommunication.cn/html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataErrorText 属性 |

错误信息

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string ErrorText { get; set; }
```

```
Public Property ErrorText As String
	Get
	Set
```

```
public:
property String^ ErrorText {
	String^ get ();
	void set (String^ value);
}
```

```
member ErrorText : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HstopStatus 属性 

[原文連結](http://api.hslcommunication.cn/html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataHstopStatus 属性 |

急停状态，1：无急停，0：有急停

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte HstopStatus { get; set; }
```

```
Public Property HstopStatus As Byte
	Get
	Set
```

```
public:
property unsigned char HstopStatus {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member HstopStatus : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IoDIn 属性 

[原文連結](http://api.hslcommunication.cn/html/056aa978-0740-28a1-7591-ae38d28234c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataIoDIn 属性 |

IoDIn状态

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] IoDIn { get; set; }
```

```
Public Property IoDIn As Byte()
	Get
	Set
```

```
public:
property array<unsigned char>^ IoDIn {
	array<unsigned char>^ get ();
	void set (array<unsigned char>^ value);
}
```

```
member IoDIn : byte[] with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IoDOut 属性 

[原文連結](http://api.hslcommunication.cn/html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataIoDOut 属性 |

IoDOut状态

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] IoDOut { get; set; }
```

```
Public Property IoDOut As Byte()
	Get
	Set
```

```
public:
property array<unsigned char>^ IoDOut {
	array<unsigned char>^ get ();
	void set (array<unsigned char>^ value);
}
```

```
member IoDOut : byte[] with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IoIIn 属性 

[原文連結](http://api.hslcommunication.cn/html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataIoIIn 属性 |

IoIIn状态

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int[] IoIIn { get; set; }
```

```
Public Property IoIIn As Integer()
	Get
	Set
```

```
public:
property array<int>^ IoIIn {
	array<int>^ get ();
	void set (array<int>^ value);
}
```

```
member IoIIn : int[] with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IoIOut 属性 

[原文連結](http://api.hslcommunication.cn/html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataIoIOut 属性 |

IoIOut状态

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int[] IoIOut { get; set; }
```

```
Public Property IoIOut As Integer()
	Get
	Set
```

```
public:
property array<int>^ IoIOut {
	array<int>^ get ();
	void set (array<int>^ value);
}
```

```
member IoIOut : int[] with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModeStatus 属性 

[原文連結](http://api.hslcommunication.cn/html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataModeStatus 属性 |

模式状态，1:手动，2:自动，3:远程

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ushort ModeStatus { get; set; }
```

```
Public Property ModeStatus As UShort
	Get
	Set
```

```
public:
property unsigned short ModeStatus {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member ModeStatus : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PacketEnd 属性 

[原文連結](http://api.hslcommunication.cn/html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataPacketEnd 属性 |

报文结束标记

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string PacketEnd { get; set; }
```

```
Public Property PacketEnd As String
	Get
	Set
```

```
public:
property String^ PacketEnd {
	String^ get ();
	void set (String^ value);
}
```

```
member PacketEnd : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PacketHeartbeat 属性 

[原文連結](http://api.hslcommunication.cn/html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataPacketHeartbeat 属性 |

数据心跳

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ushort PacketHeartbeat { get; set; }
```

```
Public Property PacketHeartbeat As UShort
	Get
	Set
```

```
public:
property unsigned short PacketHeartbeat {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member PacketHeartbeat : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PacketOrders 属性 

[原文連結](http://api.hslcommunication.cn/html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataPacketOrders 属性 |

数据命令

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ushort PacketOrders { get; set; }
```

```
Public Property PacketOrders As UShort
	Get
	Set
```

```
public:
property unsigned short PacketOrders {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member PacketOrders : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PacketStart 属性 

[原文連結](http://api.hslcommunication.cn/html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataPacketStart 属性 |

报文开始的字符串

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string PacketStart { get; set; }
```

```
Public Property PacketStart As String
	Get
	Set
```

```
public:
property String^ PacketStart {
	String^ get ();
	void set (String^ value);
}
```

```
member PacketStart : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ProgHoldStatus 属性 

[原文連結](http://api.hslcommunication.cn/html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataProgHoldStatus 属性 |

程序暂停状态，1：有暂停，0：无暂停

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte ProgHoldStatus { get; set; }
```

```
Public Property ProgHoldStatus As Byte
	Get
	Set
```

```
public:
property unsigned char ProgHoldStatus {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member ProgHoldStatus : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ProgLoadStatus 属性 

[原文連結](http://api.hslcommunication.cn/html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataProgLoadStatus 属性 |

程序加载状态，1：有加载，0：无加载

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte ProgLoadStatus { get; set; }
```

```
Public Property ProgLoadStatus As Byte
	Get
	Set
```

```
public:
property unsigned char ProgLoadStatus {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member ProgLoadStatus : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ProgMoveStatus 属性 

[原文連結](http://api.hslcommunication.cn/html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.EFORT](../html/a25bf846-86de-f385-03a1-570078f0c95b.htm "HslCommunication.Robot.EFORT")

[EfortData 类](../html/d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm "EfortData 类")

[EfortData 属性](../html/18e63891-4d82-a422-f7df-1c2bbee8f38e.htm "EfortData 属性")

[AuthorityStatus 属性](../html/97dd1460-9d6d-3953-fc93-ffa9dc116eee.htm "AuthorityStatus 属性 ")

[AxisMoveStatus 属性](../html/e6082238-491c-708a-09b2-b3e1994980f3.htm "AxisMoveStatus 属性 ")

[DbAxisAcc 属性](../html/4fc653f6-1979-57c5-480b-fd184236e397.htm "DbAxisAcc 属性 ")

[DbAxisAccAcc 属性](../html/88bac00f-0011-58c5-973d-c582220bc9ce.htm "DbAxisAccAcc 属性 ")

[DbAxisDirCnt 属性](../html/b78b2ec0-f38e-e43b-ba85-70b98fffcc91.htm "DbAxisDirCnt 属性 ")

[DbAxisPos 属性](../html/2e33f5be-c12e-dc08-b1f7-3f60ca23786a.htm "DbAxisPos 属性 ")

[DbAxisSpeed 属性](../html/047c4ef5-63c2-a34e-5bb4-c1e1530e4572.htm "DbAxisSpeed 属性 ")

[DbAxisTime 属性](../html/ed4d34c0-894d-d267-cf3e-fac0676797c4.htm "DbAxisTime 属性 ")

[DbAxisTorque 属性](../html/7c482eb6-b5bb-128c-fbb8-7b4d5ff7a593.htm "DbAxisTorque 属性 ")

[DbCartPos 属性](../html/f12e137a-a406-d6fd-7b1d-a3c4d712affc.htm "DbCartPos 属性 ")

[DbDeviceTime 属性](../html/a1b18710-88cc-6c8b-ece5-ff220c2f7fc5.htm "DbDeviceTime 属性 ")

[ErrorStatus 属性](../html/d28c3c3e-b4ea-8064-3d62-10554dd80555.htm "ErrorStatus 属性 ")

[ErrorText 属性](../html/bfd84369-977c-f9b1-da1d-e5d2722227d4.htm "ErrorText 属性 ")

[HstopStatus 属性](../html/b57dcee7-f769-da7f-2597-69d1f9ffe4c7.htm "HstopStatus 属性 ")

[IoDIn 属性](../html/056aa978-0740-28a1-7591-ae38d28234c8.htm "IoDIn 属性 ")

[IoDOut 属性](../html/5d41b966-1f96-5ab5-6d24-a3348f0808d7.htm "IoDOut 属性 ")

[IoIIn 属性](../html/8a7e4e6c-4017-1c87-4bcf-40dcc6dd753c.htm "IoIIn 属性 ")

[IoIOut 属性](../html/14fad9a8-4e91-4c17-70b0-b967e2e15ef2.htm "IoIOut 属性 ")

[ModeStatus 属性](../html/2112aaf8-50e4-ff3a-8517-410d3e42a4f1.htm "ModeStatus 属性 ")

[PacketEnd 属性](../html/f3782238-d4bb-3abb-cd4f-03db08ea9f14.htm "PacketEnd 属性 ")

[PacketHeartbeat 属性](../html/0584ab63-30e3-cd62-6c3d-db5172ce0c58.htm "PacketHeartbeat 属性 ")

[PacketOrders 属性](../html/f5866a57-0dbf-adac-5cfc-a6a4b57ed913.htm "PacketOrders 属性 ")

[PacketStart 属性](../html/fa9c65dc-ee91-8616-6b04-bfa89285d082.htm "PacketStart 属性 ")

[ProgHoldStatus 属性](../html/f0e77644-2606-422c-a2da-06aaadaf9d3e.htm "ProgHoldStatus 属性 ")

[ProgLoadStatus 属性](../html/795d513d-df00-ede3-898d-c8ee109bd3d1.htm "ProgLoadStatus 属性 ")

[ProgMoveStatus 属性](../html/e7c244f3-6c6b-bbba-39f9-0f70a8497069.htm "ProgMoveStatus 属性 ")

[ProgramName 属性](../html/ae5a4e14-1c87-0710-60ef-3f8108efe575.htm "ProgramName 属性 ")

[ProjectName 属性](../html/3aadbce0-04f2-4ec8-e804-e628791978e8.htm "ProjectName 属性 ")

[ServoStatus 属性](../html/e2c9706c-0e47-0df8-316e-5373ffc1ec87.htm "ServoStatus 属性 ")

[SpeedStatus 属性](../html/f4f31109-f088-53ee-9496-0d577d350788.htm "SpeedStatus 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| EfortDataProgMoveStatus 属性 |

程序运行状态，1：有运行，0：未运行

**命名空间：**
 [HslCommunication.Robot.EFORT](a25bf846-86de-f385-03a1-570078f0c95b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte ProgMoveStatus { get; set; }
```

```
Public Property ProgMoveStatus As Byte
	Get
	Set
```

```
public:
property unsigned char ProgMoveStatus {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member ProgMoveStatus : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[EfortData 类](d1d4a7ff-ad59-2858-0739-cd89b153a02b.htm)

[HslCommunication.Robot.EFORT 命名空间](a25bf846-86de-f385-03a1-570078f0c95b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)