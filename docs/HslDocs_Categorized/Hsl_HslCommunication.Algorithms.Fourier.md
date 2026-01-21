# HslCommunication - HslCommunication.Algorithms.Fourier

> 分類頁數: 20



---
## HslCommunication.Algorithms.Fourier

[原文連結](http://api.hslcommunication.cn/html/9b966345-70f7-c13e-968b-718c9cc61077.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTFilter 类](../html/08da094c-cd0b-c330-9c21-1f5ded563c50.htm "FFTFilter 类")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Algorithms.Fourier 命名空间 |

[缺少 "N:HslCommunication.Algorithms.Fourier" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [FFTFilter](08da094c-cd0b-c330-9c21-1f5ded563c50.htm) | 一个基于傅立叶变换的一个滤波算法 |
| 公共类 | [FFTHelper](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm) | 离散傅氏变换的快速算法，处理的信号，适合单周期信号数为2的N次方个，支持变换及逆变换 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFTFilter 类

[原文連結](http://api.hslcommunication.cn/html/08da094c-cd0b-c330-9c21-1f5ded563c50.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTFilter 类](../html/08da094c-cd0b-c330-9c21-1f5ded563c50.htm "FFTFilter 类")

[FFTFilter 构造函数](../html/c7ca3e7b-6586-425b-0a14-11ef03aaddfa.htm "FFTFilter 构造函数 ")

[FFTFilter 方法](../html/98aae6db-b7a5-0a69-d6b2-ed3988bb6192.htm "FFTFilter 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTFilter 类 |

一个基于傅立叶变换的一个滤波算法

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Algorithms.FourierFFTFilter

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FFTFilter
```

```
Public Class FFTFilter
```

```
public ref class FFTFilter
```

```
type FFTFilter =  class end
```

FFTFilter 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FFTFilter](c7ca3e7b-6586-425b-0a14-11ef03aaddfa.htm) | 初始化 FFTFilter 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [FillDataArrayT](7ba23974-a805-d21b-4133-540aab3fc2ff.htm) | 对指定的数据进行填充，方便的进行傅立叶计算 |
| 公共方法静态成员 | [FilterFFT(Double, Double)](5a41c250-c20e-a153-91de-d86bdefa5ad3.htm) | 对指定的原始数据进行滤波，并返回成功的数据值 |
| 公共方法静态成员 | [FilterFFT(Single, Double)](78c9c1c9-c4af-8896-04f6-0f69d4a2327b.htm) | 对指定的原始数据进行滤波，并返回成功的数据值 |
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

![](../icons/SectionExpanded.png)备注

非常感谢来自北京的monk网友，提供了完整的解决方法。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFTFilter 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c7ca3e7b-6586-425b-0a14-11ef03aaddfa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTFilter 类](../html/08da094c-cd0b-c330-9c21-1f5ded563c50.htm "FFTFilter 类")

[FFTFilter 构造函数](../html/c7ca3e7b-6586-425b-0a14-11ef03aaddfa.htm "FFTFilter 构造函数 ")

[FFTFilter 方法](../html/98aae6db-b7a5-0a69-d6b2-ed3988bb6192.htm "FFTFilter 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTFilter 构造函数 |

初始化 [FFTFilter](08da094c-cd0b-c330-9c21-1f5ded563c50.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FFTFilter()
```

```
Public Sub New
```

```
public:
FFTFilter()
```

```
new : unit -> FFTFilter
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTFilter 类](08da094c-cd0b-c330-9c21-1f5ded563c50.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFTFilter 方法

[原文連結](http://api.hslcommunication.cn/html/98aae6db-b7a5-0a69-d6b2-ed3988bb6192.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTFilter 类](../html/08da094c-cd0b-c330-9c21-1f5ded563c50.htm "FFTFilter 类")

[FFTFilter 方法](../html/98aae6db-b7a5-0a69-d6b2-ed3988bb6192.htm "FFTFilter 方法")

[FillDataArray(T) 方法](../html/7ba23974-a805-d21b-4133-540aab3fc2ff.htm "FillDataArray(T) 方法 ")

[FilterFFT 方法](../html/5f5ef1e6-038c-6acf-0fb4-089e77a2c039.htm "FilterFFT 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTFilter 方法 |

[FFTFilter](08da094c-cd0b-c330-9c21-1f5ded563c50.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [FillDataArrayT](7ba23974-a805-d21b-4133-540aab3fc2ff.htm) | 对指定的数据进行填充，方便的进行傅立叶计算 |
| 公共方法静态成员 | [FilterFFT(Double, Double)](5a41c250-c20e-a153-91de-d86bdefa5ad3.htm) | 对指定的原始数据进行滤波，并返回成功的数据值 |
| 公共方法静态成员 | [FilterFFT(Single, Double)](78c9c1c9-c4af-8896-04f6-0f69d4a2327b.htm) | 对指定的原始数据进行滤波，并返回成功的数据值 |
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

[FFTFilter 类](08da094c-cd0b-c330-9c21-1f5ded563c50.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FillDataArray(T) 方法 

[原文連結](http://api.hslcommunication.cn/html/7ba23974-a805-d21b-4133-540aab3fc2ff.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTFilter 类](../html/08da094c-cd0b-c330-9c21-1f5ded563c50.htm "FFTFilter 类")

[FFTFilter 方法](../html/98aae6db-b7a5-0a69-d6b2-ed3988bb6192.htm "FFTFilter 方法")

[FillDataArray(T) 方法](../html/7ba23974-a805-d21b-4133-540aab3fc2ff.htm "FillDataArray(T) 方法 ")

[FilterFFT 方法](../html/5f5ef1e6-038c-6acf-0fb4-089e77a2c039.htm "FilterFFT 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTFilterFillDataArrayT 方法 |

对指定的数据进行填充，方便的进行傅立叶计算

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static List<T> FillDataArray<T>(
	List<T> source,
	out int putLength
)
```

```
Public Shared Function FillDataArray(Of T) ( 
	source As List(Of T),
	<OutAttribute> ByRef putLength As Integer
) As List(Of T)
```

```
public:
generic<typename T>
static List<T>^ FillDataArray(
	List<T>^ source, 
	[OutAttribute] int% putLength
)
```

```
static member FillDataArray : 
        source : List<'T> * 
        putLength : int byref -> List<'T> 
```

#### 参数

source
:   类型：System.Collections.GenericListT  
    数据源

putLength
:   类型：SystemInt32  
    输出的长度

#### 类型参数

T
:   数据的数据类型

#### 返回值

类型：ListT  
填充结果

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTFilter 类](08da094c-cd0b-c330-9c21-1f5ded563c50.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FilterFFT 方法 

[原文連結](http://api.hslcommunication.cn/html/5f5ef1e6-038c-6acf-0fb4-089e77a2c039.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTFilter 类](../html/08da094c-cd0b-c330-9c21-1f5ded563c50.htm "FFTFilter 类")

[FFTFilter 方法](../html/98aae6db-b7a5-0a69-d6b2-ed3988bb6192.htm "FFTFilter 方法")

[FilterFFT 方法](../html/5f5ef1e6-038c-6acf-0fb4-089e77a2c039.htm "FilterFFT 方法 ")

[FilterFFT 方法 (Double[], Double)](../html/5a41c250-c20e-a153-91de-d86bdefa5ad3.htm "FilterFFT 方法 (Double[], Double)")

[FilterFFT 方法 (Single[], Double)](../html/78c9c1c9-c4af-8896-04f6-0f69d4a2327b.htm "FilterFFT 方法 (Single[], Double)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTFilterFilterFFT 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [FilterFFT(Double, Double)](5a41c250-c20e-a153-91de-d86bdefa5ad3.htm) | 对指定的原始数据进行滤波，并返回成功的数据值 |
| 公共方法静态成员 | [FilterFFT(Single, Double)](78c9c1c9-c4af-8896-04f6-0f69d4a2327b.htm) | 对指定的原始数据进行滤波，并返回成功的数据值 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTFilter 类](08da094c-cd0b-c330-9c21-1f5ded563c50.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FilterFFT 方法 (Double[], Double)

[原文連結](http://api.hslcommunication.cn/html/5a41c250-c20e-a153-91de-d86bdefa5ad3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTFilter 类](../html/08da094c-cd0b-c330-9c21-1f5ded563c50.htm "FFTFilter 类")

[FFTFilter 方法](../html/98aae6db-b7a5-0a69-d6b2-ed3988bb6192.htm "FFTFilter 方法")

[FilterFFT 方法](../html/5f5ef1e6-038c-6acf-0fb4-089e77a2c039.htm "FilterFFT 方法 ")

[FilterFFT 方法 (Double[], Double)](../html/5a41c250-c20e-a153-91de-d86bdefa5ad3.htm "FilterFFT 方法 (Double[], Double)")

[FilterFFT 方法 (Single[], Double)](../html/78c9c1c9-c4af-8896-04f6-0f69d4a2327b.htm "FilterFFT 方法 (Single[], Double)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTFilterFilterFFT 方法 (Double, Double) |

对指定的原始数据进行滤波，并返回成功的数据值

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static double[] FilterFFT(
	double[] source,
	double filter
)
```

```
Public Shared Function FilterFFT ( 
	source As Double(),
	filter As Double
) As Double()
```

```
public:
static array<double>^ FilterFFT(
	array<double>^ source, 
	double filter
)
```

```
static member FilterFFT : 
        source : float[] * 
        filter : float -> float[] 
```

#### 参数

source
:   类型：SystemDouble  
    数据源，数组的长度需要为2的n次方。

filter
:   类型：SystemDouble  
    滤波值：最大值为1，不能低于0，越接近1，滤波强度越强，也可能会导致失去真实信号，为0时没有滤波效果。

#### 返回值

类型：Double  
滤波后的数据值

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTFilter 类](08da094c-cd0b-c330-9c21-1f5ded563c50.htm)

[FilterFFT 重载](5f5ef1e6-038c-6acf-0fb4-089e77a2c039.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FilterFFT 方法 (Single[], Double)

[原文連結](http://api.hslcommunication.cn/html/78c9c1c9-c4af-8896-04f6-0f69d4a2327b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTFilter 类](../html/08da094c-cd0b-c330-9c21-1f5ded563c50.htm "FFTFilter 类")

[FFTFilter 方法](../html/98aae6db-b7a5-0a69-d6b2-ed3988bb6192.htm "FFTFilter 方法")

[FilterFFT 方法](../html/5f5ef1e6-038c-6acf-0fb4-089e77a2c039.htm "FilterFFT 方法 ")

[FilterFFT 方法 (Double[], Double)](../html/5a41c250-c20e-a153-91de-d86bdefa5ad3.htm "FilterFFT 方法 (Double[], Double)")

[FilterFFT 方法 (Single[], Double)](../html/78c9c1c9-c4af-8896-04f6-0f69d4a2327b.htm "FilterFFT 方法 (Single[], Double)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTFilterFilterFFT 方法 (Single, Double) |

对指定的原始数据进行滤波，并返回成功的数据值

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static float[] FilterFFT(
	float[] source,
	double filter
)
```

```
Public Shared Function FilterFFT ( 
	source As Single(),
	filter As Double
) As Single()
```

```
public:
static array<float>^ FilterFFT(
	array<float>^ source, 
	double filter
)
```

```
static member FilterFFT : 
        source : float32[] * 
        filter : float -> float32[] 
```

#### 参数

source
:   类型：SystemSingle  
    数据源，数组的长度需要为2的n次方。

filter
:   类型：SystemDouble  
    滤波值：最大值为1，不能低于0，越接近1，滤波强度越强，也可能会导致失去真实信号，为0时没有滤波效果。

#### 返回值

类型：Single  
滤波后的数据值

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTFilter 类](08da094c-cd0b-c330-9c21-1f5ded563c50.htm)

[FilterFFT 重载](5f5ef1e6-038c-6acf-0fb4-089e77a2c039.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFTHelper 类

[原文連結](http://api.hslcommunication.cn/html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 构造函数](../html/44585a66-169a-adac-20c0-03c62a99542b.htm "FFTHelper 构造函数 ")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelper 类 |

离散傅氏变换的快速算法，处理的信号，适合单周期信号数为2的N次方个，支持变换及逆变换

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Algorithms.FourierFFTHelper

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FFTHelper
```

```
Public Class FFTHelper
```

```
public ref class FFTHelper
```

```
type FFTHelper =  class end
```

FFTHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FFTHelper](44585a66-169a-adac-20c0-03c62a99542b.htm) | 初始化 FFTHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [FFT(Double)](43b16b32-060b-f66b-97cb-619854d1ea24.htm) | 快速傅立叶变换 |
| 公共方法静态成员 | [FFT(Double, Double)](e415bbf5-ab43-9851-a4c3-c976b2f8535d.htm) | 快速傅立叶变换 |
| 公共方法静态成员 | [FFT(Single, Single)](419273bc-ee92-08b1-fa78-019634ada3f7.htm) | 快速傅立叶变换 |
| 公共方法静态成员 | [FFTValue](9bc6019e-5c5c-1400-609e-0348ddea0b55.htm) | 快速傅立叶变换 |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetFFTImage](3ad3be89-d9e1-bfab-8b41-3519c253834b.htm) | 获取FFT变换后的显示图形，需要指定图形的相关参数 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法静态成员 | [IFFT(Double, Double)](bdda41a2-3734-96fa-67e7-03629268ae91.htm) | 快速傅立叶变换的逆变换 |
| 公共方法静态成员 | [IFFT(Single, Single)](498e8e16-b4ca-cedf-a9ee-c2020d0b26f7.htm) | 快速傅立叶变换的逆变换 |
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

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFTHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/44585a66-169a-adac-20c0-03c62a99542b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 构造函数](../html/44585a66-169a-adac-20c0-03c62a99542b.htm "FFTHelper 构造函数 ")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelper 构造函数 |

初始化 [FFTHelper](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FFTHelper()
```

```
Public Sub New
```

```
public:
FFTHelper()
```

```
new : unit -> FFTHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFTHelper 方法

[原文連結](http://api.hslcommunication.cn/html/8753e47d-9e22-449b-1430-938520027d20.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[FFT 方法](../html/8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm "FFT 方法 ")

[FFTValue 方法](../html/9bc6019e-5c5c-1400-609e-0348ddea0b55.htm "FFTValue 方法 ")

[GetFFTImage 方法](../html/3ad3be89-d9e1-bfab-8b41-3519c253834b.htm "GetFFTImage 方法 ")

[IFFT 方法](../html/b5adeeff-683f-ef66-2395-2386a7e56041.htm "IFFT 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelper 方法 |

[FFTHelper](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [FFT(Double)](43b16b32-060b-f66b-97cb-619854d1ea24.htm) | 快速傅立叶变换 |
| 公共方法静态成员 | [FFT(Double, Double)](e415bbf5-ab43-9851-a4c3-c976b2f8535d.htm) | 快速傅立叶变换 |
| 公共方法静态成员 | [FFT(Single, Single)](419273bc-ee92-08b1-fa78-019634ada3f7.htm) | 快速傅立叶变换 |
| 公共方法静态成员 | [FFTValue](9bc6019e-5c5c-1400-609e-0348ddea0b55.htm) | 快速傅立叶变换 |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetFFTImage](3ad3be89-d9e1-bfab-8b41-3519c253834b.htm) | 获取FFT变换后的显示图形，需要指定图形的相关参数 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法静态成员 | [IFFT(Double, Double)](bdda41a2-3734-96fa-67e7-03629268ae91.htm) | 快速傅立叶变换的逆变换 |
| 公共方法静态成员 | [IFFT(Single, Single)](498e8e16-b4ca-cedf-a9ee-c2020d0b26f7.htm) | 快速傅立叶变换的逆变换 |
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

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFT 方法 

[原文連結](http://api.hslcommunication.cn/html/8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[FFT 方法](../html/8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm "FFT 方法 ")

[FFT 方法 (Double[])](../html/43b16b32-060b-f66b-97cb-619854d1ea24.htm "FFT 方法 (Double[])")

[FFT 方法 (Double[], Double[])](../html/e415bbf5-ab43-9851-a4c3-c976b2f8535d.htm "FFT 方法 (Double[], Double[])")

[FFT 方法 (Single[], Single[])](../html/419273bc-ee92-08b1-fa78-019634ada3f7.htm "FFT 方法 (Single[], Single[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelperFFT 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [FFT(Double)](43b16b32-060b-f66b-97cb-619854d1ea24.htm) | 快速傅立叶变换 |
| 公共方法静态成员 | [FFT(Double, Double)](e415bbf5-ab43-9851-a4c3-c976b2f8535d.htm) | 快速傅立叶变换 |
| 公共方法静态成员 | [FFT(Single, Single)](419273bc-ee92-08b1-fa78-019634ada3f7.htm) | 快速傅立叶变换 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFT 方法 (Double[])

[原文連結](http://api.hslcommunication.cn/html/43b16b32-060b-f66b-97cb-619854d1ea24.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[FFT 方法](../html/8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm "FFT 方法 ")

[FFT 方法 (Double[])](../html/43b16b32-060b-f66b-97cb-619854d1ea24.htm "FFT 方法 (Double[])")

[FFT 方法 (Double[], Double[])](../html/e415bbf5-ab43-9851-a4c3-c976b2f8535d.htm "FFT 方法 (Double[], Double[])")

[FFT 方法 (Single[], Single[])](../html/419273bc-ee92-08b1-fa78-019634ada3f7.htm "FFT 方法 (Single[], Single[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelperFFT 方法 (Double) |

快速傅立叶变换

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static double[] FFT(
	double[] xreal
)
```

```
Public Shared Function FFT ( 
	xreal As Double()
) As Double()
```

```
public:
static array<double>^ FFT(
	array<double>^ xreal
)
```

```
static member FFT : 
        xreal : float[] -> float[] 
```

#### 参数

xreal
:   类型：SystemDouble  
    实数部分

#### 返回值

类型：Double  
变换后的数组值

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[FFT 重载](8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFT 方法 (Double[], Double[])

[原文連結](http://api.hslcommunication.cn/html/e415bbf5-ab43-9851-a4c3-c976b2f8535d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[FFT 方法](../html/8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm "FFT 方法 ")

[FFT 方法 (Double[])](../html/43b16b32-060b-f66b-97cb-619854d1ea24.htm "FFT 方法 (Double[])")

[FFT 方法 (Double[], Double[])](../html/e415bbf5-ab43-9851-a4c3-c976b2f8535d.htm "FFT 方法 (Double[], Double[])")

[FFT 方法 (Single[], Single[])](../html/419273bc-ee92-08b1-fa78-019634ada3f7.htm "FFT 方法 (Single[], Single[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelperFFT 方法 (Double, Double) |

快速傅立叶变换

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static int FFT(
	double[] xreal,
	double[] ximag
)
```

```
Public Shared Function FFT ( 
	xreal As Double(),
	ximag As Double()
) As Integer
```

```
public:
static int FFT(
	array<double>^ xreal, 
	array<double>^ ximag
)
```

```
static member FFT : 
        xreal : float[] * 
        ximag : float[] -> int 
```

#### 参数

xreal
:   类型：SystemDouble  
    实数部分，数组长度最好为2的n次方

ximag
:   类型：SystemDouble  
    虚数部分，数组长度最好为2的n次方

#### 返回值

类型：Int32  
变换后的数组值

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[FFT 重载](8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFT 方法 (Single[], Single[])

[原文連結](http://api.hslcommunication.cn/html/419273bc-ee92-08b1-fa78-019634ada3f7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[FFT 方法](../html/8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm "FFT 方法 ")

[FFT 方法 (Double[])](../html/43b16b32-060b-f66b-97cb-619854d1ea24.htm "FFT 方法 (Double[])")

[FFT 方法 (Double[], Double[])](../html/e415bbf5-ab43-9851-a4c3-c976b2f8535d.htm "FFT 方法 (Double[], Double[])")

[FFT 方法 (Single[], Single[])](../html/419273bc-ee92-08b1-fa78-019634ada3f7.htm "FFT 方法 (Single[], Single[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelperFFT 方法 (Single, Single) |

快速傅立叶变换

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static int FFT(
	float[] xreal,
	float[] ximag
)
```

```
Public Shared Function FFT ( 
	xreal As Single(),
	ximag As Single()
) As Integer
```

```
public:
static int FFT(
	array<float>^ xreal, 
	array<float>^ ximag
)
```

```
static member FFT : 
        xreal : float32[] * 
        ximag : float32[] -> int 
```

#### 参数

xreal
:   类型：SystemSingle  
    实数部分，数组长度最好为2的n次方

ximag
:   类型：SystemSingle  
    虚数部分，数组长度最好为2的n次方

#### 返回值

类型：Int32  
变换后的数组值

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[FFT 重载](8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FFTValue 方法 

[原文連結](http://api.hslcommunication.cn/html/9bc6019e-5c5c-1400-609e-0348ddea0b55.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[FFT 方法](../html/8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm "FFT 方法 ")

[FFTValue 方法](../html/9bc6019e-5c5c-1400-609e-0348ddea0b55.htm "FFTValue 方法 ")

[GetFFTImage 方法](../html/3ad3be89-d9e1-bfab-8b41-3519c253834b.htm "GetFFTImage 方法 ")

[IFFT 方法](../html/b5adeeff-683f-ef66-2395-2386a7e56041.htm "IFFT 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelperFFTValue 方法 |

快速傅立叶变换

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static double[] FFTValue(
	double[] xreal,
	double[] ximag,
	bool isSqrtDouble = false
)
```

```
Public Shared Function FFTValue ( 
	xreal As Double(),
	ximag As Double(),
	Optional isSqrtDouble As Boolean = false
) As Double()
```

```
public:
static array<double>^ FFTValue(
	array<double>^ xreal, 
	array<double>^ ximag, 
	bool isSqrtDouble = false
)
```

```
static member FFTValue : 
        xreal : float[] * 
        ximag : float[] * 
        ?isSqrtDouble : bool 
(* Defaults:
        let _isSqrtDouble = defaultArg isSqrtDouble false
*)
-> float[] 
```

#### 参数

xreal
:   类型：SystemDouble  
    实数部分，数组长度最好为2的n次方

ximag
:   类型：SystemDouble  
    虚数部分，数组长度最好为2的n次方

isSqrtDouble (Optional)
:   类型：SystemBoolean  
    是否开两次根，显示的噪点信息会更新明显

#### 返回值

类型：Double  
变换后的数组值

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetFFTImage 方法 

[原文連結](http://api.hslcommunication.cn/html/3ad3be89-d9e1-bfab-8b41-3519c253834b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[FFT 方法](../html/8acf543f-e7ce-8836-4cec-898ca15cf2dd.htm "FFT 方法 ")

[FFTValue 方法](../html/9bc6019e-5c5c-1400-609e-0348ddea0b55.htm "FFTValue 方法 ")

[GetFFTImage 方法](../html/3ad3be89-d9e1-bfab-8b41-3519c253834b.htm "GetFFTImage 方法 ")

[IFFT 方法](../html/b5adeeff-683f-ef66-2395-2386a7e56041.htm "IFFT 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelperGetFFTImage 方法 |

获取FFT变换后的显示图形，需要指定图形的相关参数

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Bitmap GetFFTImage(
	double[] xreal,
	int width,
	int heigh,
	Color lineColor,
	bool isSqrtDouble = false
)
```

```
Public Shared Function GetFFTImage ( 
	xreal As Double(),
	width As Integer,
	heigh As Integer,
	lineColor As Color,
	Optional isSqrtDouble As Boolean = false
) As Bitmap
```

```
public:
static Bitmap^ GetFFTImage(
	array<double>^ xreal, 
	int width, 
	int heigh, 
	Color lineColor, 
	bool isSqrtDouble = false
)
```

```
static member GetFFTImage : 
        xreal : float[] * 
        width : int * 
        heigh : int * 
        lineColor : Color * 
        ?isSqrtDouble : bool 
(* Defaults:
        let _isSqrtDouble = defaultArg isSqrtDouble false
*)
-> Bitmap 
```

#### 参数

xreal
:   类型：SystemDouble  
    实数部分的值

width
:   类型：SystemInt32  
    图形的宽度

heigh
:   类型：SystemInt32  
    图形的高度

lineColor
:   类型：System.DrawingColor  
    线条颜色

isSqrtDouble (Optional)
:   类型：SystemBoolean  
    是否开两次根，显示的噪点信息会更新明显

#### 返回值

类型：Bitmap  
等待呈现的图形

![](../icons/SectionExpanded.png)备注

| 警告 **警告：** |
| --- |
| .net standrard2.0 下不支持。 |

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IFFT 方法 

[原文連結](http://api.hslcommunication.cn/html/b5adeeff-683f-ef66-2395-2386a7e56041.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[IFFT 方法](../html/b5adeeff-683f-ef66-2395-2386a7e56041.htm "IFFT 方法 ")

[IFFT 方法 (Double[], Double[])](../html/bdda41a2-3734-96fa-67e7-03629268ae91.htm "IFFT 方法 (Double[], Double[])")

[IFFT 方法 (Single[], Single[])](../html/498e8e16-b4ca-cedf-a9ee-c2020d0b26f7.htm "IFFT 方法 (Single[], Single[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelperIFFT 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [IFFT(Double, Double)](bdda41a2-3734-96fa-67e7-03629268ae91.htm) | 快速傅立叶变换的逆变换 |
| 公共方法静态成员 | [IFFT(Single, Single)](498e8e16-b4ca-cedf-a9ee-c2020d0b26f7.htm) | 快速傅立叶变换的逆变换 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IFFT 方法 (Double[], Double[])

[原文連結](http://api.hslcommunication.cn/html/bdda41a2-3734-96fa-67e7-03629268ae91.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[IFFT 方法](../html/b5adeeff-683f-ef66-2395-2386a7e56041.htm "IFFT 方法 ")

[IFFT 方法 (Double[], Double[])](../html/bdda41a2-3734-96fa-67e7-03629268ae91.htm "IFFT 方法 (Double[], Double[])")

[IFFT 方法 (Single[], Single[])](../html/498e8e16-b4ca-cedf-a9ee-c2020d0b26f7.htm "IFFT 方法 (Single[], Single[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelperIFFT 方法 (Double, Double) |

快速傅立叶变换的逆变换

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static int IFFT(
	double[] xreal,
	double[] ximag
)
```

```
Public Shared Function IFFT ( 
	xreal As Double(),
	ximag As Double()
) As Integer
```

```
public:
static int IFFT(
	array<double>^ xreal, 
	array<double>^ ximag
)
```

```
static member IFFT : 
        xreal : float[] * 
        ximag : float[] -> int 
```

#### 参数

xreal
:   类型：SystemDouble  
    实数部分，数组长度最好为2的n次方

ximag
:   类型：SystemDouble  
    虚数部分，数组长度最好为2的n次方

#### 返回值

类型：Int32  
2的多少次方

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[IFFT 重载](b5adeeff-683f-ef66-2395-2386a7e56041.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IFFT 方法 (Single[], Single[])

[原文連結](http://api.hslcommunication.cn/html/498e8e16-b4ca-cedf-a9ee-c2020d0b26f7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.Fourier](../html/9b966345-70f7-c13e-968b-718c9cc61077.htm "HslCommunication.Algorithms.Fourier")

[FFTHelper 类](../html/47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm "FFTHelper 类")

[FFTHelper 方法](../html/8753e47d-9e22-449b-1430-938520027d20.htm "FFTHelper 方法")

[IFFT 方法](../html/b5adeeff-683f-ef66-2395-2386a7e56041.htm "IFFT 方法 ")

[IFFT 方法 (Double[], Double[])](../html/bdda41a2-3734-96fa-67e7-03629268ae91.htm "IFFT 方法 (Double[], Double[])")

[IFFT 方法 (Single[], Single[])](../html/498e8e16-b4ca-cedf-a9ee-c2020d0b26f7.htm "IFFT 方法 (Single[], Single[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FFTHelperIFFT 方法 (Single, Single) |

快速傅立叶变换的逆变换

**命名空间：**
 [HslCommunication.Algorithms.Fourier](9b966345-70f7-c13e-968b-718c9cc61077.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static int IFFT(
	float[] xreal,
	float[] ximag
)
```

```
Public Shared Function IFFT ( 
	xreal As Single(),
	ximag As Single()
) As Integer
```

```
public:
static int IFFT(
	array<float>^ xreal, 
	array<float>^ ximag
)
```

```
static member IFFT : 
        xreal : float32[] * 
        ximag : float32[] -> int 
```

#### 参数

xreal
:   类型：SystemSingle  
    实数部分，数组长度最好为2的n次方

ximag
:   类型：SystemSingle  
    虚数部分，数组长度最好为2的n次方

#### 返回值

类型：Int32  
2的多少次方

![](../icons/SectionExpanded.png)参见

#### 引用

[FFTHelper 类](47af23b4-9eb6-d1f4-f93b-0c4145013a4e.htm)

[IFFT 重载](b5adeeff-683f-ef66-2395-2386a7e56041.htm)

[HslCommunication.Algorithms.Fourier 命名空间](9b966345-70f7-c13e-968b-718c9cc61077.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)