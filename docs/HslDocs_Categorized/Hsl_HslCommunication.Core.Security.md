# HslCommunication - HslCommunication.Core.Security

> 分類頁數: 30



---
## HslCommunication.Core.Security

[原文連結](http://api.hslcommunication.cn/html/9951c682-3526-ca1b-21eb-2be212e912ea.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[HslCertificate 类](../html/53cc213c-50c4-52f8-0a00-76009a75efd6.htm "HslCertificate 类")

[ICryptography 接口](../html/492b9453-132b-0879-1f74-f51c559ff6b7.htm "ICryptography 接口")

[OpenSslNative 类](../html/1dca84e8-db21-669e-2937-46bb11dfdcd8.htm "OpenSslNative 类")

[OpenSslNative.SSL\_CTX\_keylog\_cb\_func 委托](../html/886ab729-b1d9-04f0-82b8-e649a490dea1.htm "OpenSslNative.SSL_CTX_keylog_cb_func 委托")

[RSAHelper 类](../html/e99257fe-3067-1445-30c6-f0bb2e400ed1.htm "RSAHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Core.Security 命名空间 |

[缺少 "N:HslCommunication.Core.Security" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [AesCryptography](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm) | 实例化一个AES加密解密的对象，默认 ECB 模式的对象 |
| 公共类 | [DesCryptography](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm) | DES加密解密的对象 |
| 公共类代码示例 | [HslCertificate](53cc213c-50c4-52f8-0a00-76009a75efd6.htm) | 基于RSA加密模型的证书，支持自定义的颁发证书，以及校验证书合法性  Certificates based on RSA encryption model support custom issuance certificates and verify certificate legitimacy |
| 公共类 | [OpenSslNative](1dca84e8-db21-669e-2937-46bb11dfdcd8.htm) | SSL相关的类方法名称，本代码来源 https://github.com/thomas-v2/S7CommPlusDriver |
| 公共类 | [RSAHelper](e99257fe-3067-1445-30c6-f0bb2e400ed1.htm) | RSA加密解密算法的辅助方法，可以用PEM格式的密钥创建公钥，或是私钥对象，然后用来加解密操作。 |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [ICryptography](492b9453-132b-0879-1f74-f51c559ff6b7.htm) | 加密解密的数据接口  Encrypted and decrypted data interface |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [OpenSslNativeSSL\_CTX\_keylog\_cb\_func](886ab729-b1d9-04f0-82b8-e649a490dea1.htm) |  |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AesCryptography 类

[原文連結](http://api.hslcommunication.cn/html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 构造函数](../html/bde6b3c3-b2ff-4936-19be-72f0b2625f50.htm "AesCryptography 构造函数 ")

[AesCryptography 属性](../html/b91ded20-0f8e-5b27-b355-8170b278e3bd.htm "AesCryptography 属性")

[AesCryptography 方法](../html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm "AesCryptography 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptography 类 |

实例化一个AES加密解密的对象，默认 ECB 模式的对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Core.SecurityAesCryptography

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class AesCryptography : ICryptography
```

```
Public Class AesCryptography
	Implements ICryptography
```

```
public ref class AesCryptography : ICryptography
```

```
type AesCryptography =  
    class
        interface ICryptography
    end
```

AesCryptography 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [AesCryptography](bde6b3c3-b2ff-4936-19be-72f0b2625f50.htm) | 使用指定的密钥实例化一个AES加密解密的对象，密钥由32位数字或字母组成，例如 12345678123456781234567812345678 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Key](2d6bb5ee-71c0-a330-0817-1b4e930cd9d0.htm) | 当前加密的密钥信息  currently encrypted key information |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Decrypt(Byte)](d289fecc-1f34-420f-47e0-d566610e24a3.htm) | 对原始的数据进行解密的操作，返回解密之后的二进制原始数据  Decrypt the original data and return the decrypted binary original data |
| 公共方法 | [Decrypt(String)](3d8cf489-ed1e-1ca6-e78e-fca1994fcb34.htm) | 针对Base64字符串进行解密操作，转为二进制数据后进行解密，解密之后使用UTF8编码获取最终的字符串数据  Decrypt the Base64 string, convert it to binary data, and decrypt it. After decryption, use UTF8 encoding to obtain the final string data. |
| 公共方法 | [Encrypt(Byte)](245c095c-94f9-65f4-1716-0131d405a870.htm) | 对原始的数据进行加密的操作，返回加密之后的二进制原始数据  Encrypt the original data and return the encrypted binary original data |
| 公共方法 | [Encrypt(String)](dfdc3add-1b6a-2c1e-4845-1f73875941bd.htm) | 针对字符串进行加密，并返回加密后的字符串数据，字符串的编码默认为UTF8，加密后返回Base64编码  Encrypt the string and return the encrypted string data. The encoding of the string is UTF8 by default. After encryption, the Base64 encoding is returned. |
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

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AesCryptography 构造函数 

[原文連結](http://api.hslcommunication.cn/html/bde6b3c3-b2ff-4936-19be-72f0b2625f50.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 构造函数](../html/bde6b3c3-b2ff-4936-19be-72f0b2625f50.htm "AesCryptography 构造函数 ")

[AesCryptography 属性](../html/b91ded20-0f8e-5b27-b355-8170b278e3bd.htm "AesCryptography 属性")

[AesCryptography 方法](../html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm "AesCryptography 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptography 构造函数 |

使用指定的密钥实例化一个AES加密解密的对象，密钥由32位数字或字母组成，例如 12345678123456781234567812345678

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public AesCryptography(
	string key,
	CipherMode mode = CipherMode.ECB
)
```

```
Public Sub New ( 
	key As String,
	Optional mode As CipherMode = CipherMode.ECB
)
```

```
public:
AesCryptography(
	String^ key, 
	CipherMode mode = CipherMode::ECB
)
```

```
new : 
        key : string * 
        ?mode : CipherMode 
(* Defaults:
        let _mode = defaultArg mode CipherMode.ECB
*)
-> AesCryptography
```

#### 参数

key
:   类型：SystemString  
    密钥

mode (Optional)
:   类型：System.Security.CryptographyCipherMode  
    加密的模式，默认为 ECB

![](../icons/SectionExpanded.png)参见

#### 引用

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AesCryptography 属性

[原文連結](http://api.hslcommunication.cn/html/b91ded20-0f8e-5b27-b355-8170b278e3bd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 属性](../html/b91ded20-0f8e-5b27-b355-8170b278e3bd.htm "AesCryptography 属性")

[Key 属性](../html/2d6bb5ee-71c0-a330-0817-1b4e930cd9d0.htm "Key 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptography 属性 |

[AesCryptography](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Key](2d6bb5ee-71c0-a330-0817-1b4e930cd9d0.htm) | 当前加密的密钥信息  currently encrypted key information |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Key 属性 

[原文連結](http://api.hslcommunication.cn/html/2d6bb5ee-71c0-a330-0817-1b4e930cd9d0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 属性](../html/b91ded20-0f8e-5b27-b355-8170b278e3bd.htm "AesCryptography 属性")

[Key 属性](../html/2d6bb5ee-71c0-a330-0817-1b4e930cd9d0.htm "Key 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptographyKey 属性 |

当前加密的密钥信息  
currently encrypted key information

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Key { get; }
```

```
Public ReadOnly Property Key As String
	Get
```

```
public:
virtual property String^ Key {
	String^ get () sealed;
}
```

```
abstract Key : string with get
override Key : string with get
```

#### 属性值

类型：String

#### 实现

[ICryptographyKey](084cec8f-91d4-d896-dd9a-1fc04efeafaa.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AesCryptography 方法

[原文連結](http://api.hslcommunication.cn/html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 方法](../html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm "AesCryptography 方法")

[Decrypt 方法](../html/561b83f3-bccd-8cdc-2f0d-9953e54afb69.htm "Decrypt 方法 ")

[Encrypt 方法](../html/97dae3a6-6db4-603e-f101-ffd24685a33d.htm "Encrypt 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptography 方法 |

[AesCryptography](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Decrypt(Byte)](d289fecc-1f34-420f-47e0-d566610e24a3.htm) | 对原始的数据进行解密的操作，返回解密之后的二进制原始数据  Decrypt the original data and return the decrypted binary original data |
| 公共方法 | [Decrypt(String)](3d8cf489-ed1e-1ca6-e78e-fca1994fcb34.htm) | 针对Base64字符串进行解密操作，转为二进制数据后进行解密，解密之后使用UTF8编码获取最终的字符串数据  Decrypt the Base64 string, convert it to binary data, and decrypt it. After decryption, use UTF8 encoding to obtain the final string data. |
| 公共方法 | [Encrypt(Byte)](245c095c-94f9-65f4-1716-0131d405a870.htm) | 对原始的数据进行加密的操作，返回加密之后的二进制原始数据  Encrypt the original data and return the encrypted binary original data |
| 公共方法 | [Encrypt(String)](dfdc3add-1b6a-2c1e-4845-1f73875941bd.htm) | 针对字符串进行加密，并返回加密后的字符串数据，字符串的编码默认为UTF8，加密后返回Base64编码  Encrypt the string and return the encrypted string data. The encoding of the string is UTF8 by default. After encryption, the Base64 encoding is returned. |
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

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Decrypt 方法 

[原文連結](http://api.hslcommunication.cn/html/561b83f3-bccd-8cdc-2f0d-9953e54afb69.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 方法](../html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm "AesCryptography 方法")

[Decrypt 方法](../html/561b83f3-bccd-8cdc-2f0d-9953e54afb69.htm "Decrypt 方法 ")

[Decrypt 方法 (Byte[])](../html/d289fecc-1f34-420f-47e0-d566610e24a3.htm "Decrypt 方法 (Byte[])")

[Decrypt 方法 (String)](../html/3d8cf489-ed1e-1ca6-e78e-fca1994fcb34.htm "Decrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptographyDecrypt 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Decrypt(Byte)](d289fecc-1f34-420f-47e0-d566610e24a3.htm) | 对原始的数据进行解密的操作，返回解密之后的二进制原始数据  Decrypt the original data and return the decrypted binary original data |
| 公共方法 | [Decrypt(String)](3d8cf489-ed1e-1ca6-e78e-fca1994fcb34.htm) | 针对Base64字符串进行解密操作，转为二进制数据后进行解密，解密之后使用UTF8编码获取最终的字符串数据  Decrypt the Base64 string, convert it to binary data, and decrypt it. After decryption, use UTF8 encoding to obtain the final string data. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Decrypt 方法 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/d289fecc-1f34-420f-47e0-d566610e24a3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 方法](../html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm "AesCryptography 方法")

[Decrypt 方法](../html/561b83f3-bccd-8cdc-2f0d-9953e54afb69.htm "Decrypt 方法 ")

[Decrypt 方法 (Byte[])](../html/d289fecc-1f34-420f-47e0-d566610e24a3.htm "Decrypt 方法 (Byte[])")

[Decrypt 方法 (String)](../html/3d8cf489-ed1e-1ca6-e78e-fca1994fcb34.htm "Decrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptographyDecrypt 方法 (Byte) |

对原始的数据进行解密的操作，返回解密之后的二进制原始数据  
Decrypt the original data and return the decrypted binary original data

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Decrypt(
	byte[] data
)
```

```
Public Function Decrypt ( 
	data As Byte()
) As Byte()
```

```
public:
virtual array<unsigned char>^ Decrypt(
	array<unsigned char>^ data
) sealed
```

```
abstract Decrypt : 
        data : byte[] -> byte[] 
override Decrypt : 
        data : byte[] -> byte[]
```

#### 参数

data
:   类型：SystemByte  
    等待解密的数据

#### 返回值

类型：Byte  
解密之后的原始二进制数据

#### 实现

[ICryptographyDecrypt(Byte)](f7f478fb-b074-39f1-5620-99ab1137ff01.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[Decrypt 重载](561b83f3-bccd-8cdc-2f0d-9953e54afb69.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Decrypt 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/3d8cf489-ed1e-1ca6-e78e-fca1994fcb34.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 方法](../html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm "AesCryptography 方法")

[Decrypt 方法](../html/561b83f3-bccd-8cdc-2f0d-9953e54afb69.htm "Decrypt 方法 ")

[Decrypt 方法 (Byte[])](../html/d289fecc-1f34-420f-47e0-d566610e24a3.htm "Decrypt 方法 (Byte[])")

[Decrypt 方法 (String)](../html/3d8cf489-ed1e-1ca6-e78e-fca1994fcb34.htm "Decrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptographyDecrypt 方法 (String) |

针对Base64字符串进行解密操作，转为二进制数据后进行解密，解密之后使用UTF8编码获取最终的字符串数据  
Decrypt the Base64 string, convert it to binary data, and decrypt it. After decryption, use UTF8 encoding to obtain the final string data.

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Decrypt(
	string data
)
```

```
Public Function Decrypt ( 
	data As String
) As String
```

```
public:
virtual String^ Decrypt(
	String^ data
) sealed
```

```
abstract Decrypt : 
        data : string -> string 
override Decrypt : 
        data : string -> string
```

#### 参数

data
:   类型：SystemString  
    base64编码的字符串数据

#### 返回值

类型：String  
最终的解析完成的字符串

#### 实现

[ICryptographyDecrypt(String)](35508463-b243-692d-5d1a-3e001ca0dbbd.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[Decrypt 重载](561b83f3-bccd-8cdc-2f0d-9953e54afb69.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Encrypt 方法 

[原文連結](http://api.hslcommunication.cn/html/97dae3a6-6db4-603e-f101-ffd24685a33d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 方法](../html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm "AesCryptography 方法")

[Encrypt 方法](../html/97dae3a6-6db4-603e-f101-ffd24685a33d.htm "Encrypt 方法 ")

[Encrypt 方法 (Byte[])](../html/245c095c-94f9-65f4-1716-0131d405a870.htm "Encrypt 方法 (Byte[])")

[Encrypt 方法 (String)](../html/dfdc3add-1b6a-2c1e-4845-1f73875941bd.htm "Encrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptographyEncrypt 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Encrypt(Byte)](245c095c-94f9-65f4-1716-0131d405a870.htm) | 对原始的数据进行加密的操作，返回加密之后的二进制原始数据  Encrypt the original data and return the encrypted binary original data |
| 公共方法 | [Encrypt(String)](dfdc3add-1b6a-2c1e-4845-1f73875941bd.htm) | 针对字符串进行加密，并返回加密后的字符串数据，字符串的编码默认为UTF8，加密后返回Base64编码  Encrypt the string and return the encrypted string data. The encoding of the string is UTF8 by default. After encryption, the Base64 encoding is returned. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Encrypt 方法 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/245c095c-94f9-65f4-1716-0131d405a870.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 方法](../html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm "AesCryptography 方法")

[Encrypt 方法](../html/97dae3a6-6db4-603e-f101-ffd24685a33d.htm "Encrypt 方法 ")

[Encrypt 方法 (Byte[])](../html/245c095c-94f9-65f4-1716-0131d405a870.htm "Encrypt 方法 (Byte[])")

[Encrypt 方法 (String)](../html/dfdc3add-1b6a-2c1e-4845-1f73875941bd.htm "Encrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptographyEncrypt 方法 (Byte) |

对原始的数据进行加密的操作，返回加密之后的二进制原始数据  
Encrypt the original data and return the encrypted binary original data

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Encrypt(
	byte[] data
)
```

```
Public Function Encrypt ( 
	data As Byte()
) As Byte()
```

```
public:
virtual array<unsigned char>^ Encrypt(
	array<unsigned char>^ data
) sealed
```

```
abstract Encrypt : 
        data : byte[] -> byte[] 
override Encrypt : 
        data : byte[] -> byte[]
```

#### 参数

data
:   类型：SystemByte  
    等待加密的数据

#### 返回值

类型：Byte  
加密之后的二进制数据

#### 实现

[ICryptographyEncrypt(Byte)](2870f77e-0e42-a69e-277d-9c1c7a9c699d.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[Encrypt 重载](97dae3a6-6db4-603e-f101-ffd24685a33d.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Encrypt 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/dfdc3add-1b6a-2c1e-4845-1f73875941bd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[AesCryptography 类](../html/6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm "AesCryptography 类")

[AesCryptography 方法](../html/302b4e83-714d-b85f-1a01-d5325ce6b76d.htm "AesCryptography 方法")

[Encrypt 方法](../html/97dae3a6-6db4-603e-f101-ffd24685a33d.htm "Encrypt 方法 ")

[Encrypt 方法 (Byte[])](../html/245c095c-94f9-65f4-1716-0131d405a870.htm "Encrypt 方法 (Byte[])")

[Encrypt 方法 (String)](../html/dfdc3add-1b6a-2c1e-4845-1f73875941bd.htm "Encrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AesCryptographyEncrypt 方法 (String) |

针对字符串进行加密，并返回加密后的字符串数据，字符串的编码默认为UTF8，加密后返回Base64编码  
Encrypt the string and return the encrypted string data. The encoding of the string is UTF8 by default. After encryption, the Base64 encoding is returned.

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Encrypt(
	string data
)
```

```
Public Function Encrypt ( 
	data As String
) As String
```

```
public:
virtual String^ Encrypt(
	String^ data
) sealed
```

```
abstract Encrypt : 
        data : string -> string 
override Encrypt : 
        data : string -> string
```

#### 参数

data
:   类型：SystemString  
    等待加密的字符串

#### 返回值

类型：String  
加密后的Base64编码

#### 实现

[ICryptographyEncrypt(String)](1c9b3294-3a48-b8f6-7d7c-42113ea43c58.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[AesCryptography 类](6f683ec7-bda8-3ac1-0081-f949bb8ed363.htm)

[Encrypt 重载](97dae3a6-6db4-603e-f101-ffd24685a33d.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DesCryptography 类

[原文連結](http://api.hslcommunication.cn/html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 构造函数](../html/cb410a1d-4e49-5d38-3c0b-beabe8014a3b.htm "DesCryptography 构造函数 ")

[DesCryptography 属性](../html/59728e89-8a18-8680-df5d-a4ac8835a6eb.htm "DesCryptography 属性")

[DesCryptography 方法](../html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm "DesCryptography 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptography 类 |

DES加密解密的对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Core.SecurityDesCryptography

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DesCryptography : ICryptography
```

```
Public Class DesCryptography
	Implements ICryptography
```

```
public ref class DesCryptography : ICryptography
```

```
type DesCryptography =  
    class
        interface ICryptography
    end
```

DesCryptography 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DesCryptography](cb410a1d-4e49-5d38-3c0b-beabe8014a3b.htm) | 使用指定的密钥来实例化一个加密对象，该密钥右8位的字符和数字组成，例如 12345678 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Key](4d805539-5eb0-7d5f-f650-b6a4da87bc2a.htm) | 当前加密的密钥信息  currently encrypted key information |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Decrypt(Byte)](8242f829-258b-6ac5-de66-eefd967379bc.htm) | 对原始的数据进行解密的操作，返回解密之后的二进制原始数据  Decrypt the original data and return the decrypted binary original data |
| 公共方法 | [Decrypt(String)](04da0e28-1e94-6078-3c76-3a91b59c72db.htm) | 针对Base64字符串进行解密操作，转为二进制数据后进行解密，解密之后使用UTF8编码获取最终的字符串数据  Decrypt the Base64 string, convert it to binary data, and decrypt it. After decryption, use UTF8 encoding to obtain the final string data. |
| 公共方法 | [Encrypt(Byte)](7fa9e921-f6fb-699a-5340-4b1efe4c07b1.htm) | 对原始的数据进行加密的操作，返回加密之后的二进制原始数据  Encrypt the original data and return the encrypted binary original data |
| 公共方法 | [Encrypt(String)](43d7b059-2fc0-265a-85c6-7d26a5146f5e.htm) | 针对字符串进行加密，并返回加密后的字符串数据，字符串的编码默认为UTF8，加密后返回Base64编码  Encrypt the string and return the encrypted string data. The encoding of the string is UTF8 by default. After encryption, the Base64 encoding is returned. |
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

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DesCryptography 构造函数 

[原文連結](http://api.hslcommunication.cn/html/cb410a1d-4e49-5d38-3c0b-beabe8014a3b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 构造函数](../html/cb410a1d-4e49-5d38-3c0b-beabe8014a3b.htm "DesCryptography 构造函数 ")

[DesCryptography 属性](../html/59728e89-8a18-8680-df5d-a4ac8835a6eb.htm "DesCryptography 属性")

[DesCryptography 方法](../html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm "DesCryptography 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptography 构造函数 |

使用指定的密钥来实例化一个加密对象，该密钥右8位的字符和数字组成，例如 12345678

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DesCryptography(
	string key
)
```

```
Public Sub New ( 
	key As String
)
```

```
public:
DesCryptography(
	String^ key
)
```

```
new : 
        key : string -> DesCryptography
```

#### 参数

key
:   类型：SystemString  
    密钥

![](../icons/SectionExpanded.png)参见

#### 引用

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DesCryptography 属性

[原文連結](http://api.hslcommunication.cn/html/59728e89-8a18-8680-df5d-a4ac8835a6eb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 属性](../html/59728e89-8a18-8680-df5d-a4ac8835a6eb.htm "DesCryptography 属性")

[Key 属性](../html/4d805539-5eb0-7d5f-f650-b6a4da87bc2a.htm "Key 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptography 属性 |

[DesCryptography](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Key](4d805539-5eb0-7d5f-f650-b6a4da87bc2a.htm) | 当前加密的密钥信息  currently encrypted key information |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Key 属性 

[原文連結](http://api.hslcommunication.cn/html/4d805539-5eb0-7d5f-f650-b6a4da87bc2a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 属性](../html/59728e89-8a18-8680-df5d-a4ac8835a6eb.htm "DesCryptography 属性")

[Key 属性](../html/4d805539-5eb0-7d5f-f650-b6a4da87bc2a.htm "Key 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptographyKey 属性 |

当前加密的密钥信息  
currently encrypted key information

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Key { get; }
```

```
Public ReadOnly Property Key As String
	Get
```

```
public:
virtual property String^ Key {
	String^ get () sealed;
}
```

```
abstract Key : string with get
override Key : string with get
```

#### 属性值

类型：String

#### 实现

[ICryptographyKey](084cec8f-91d4-d896-dd9a-1fc04efeafaa.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DesCryptography 方法

[原文連結](http://api.hslcommunication.cn/html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 方法](../html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm "DesCryptography 方法")

[Decrypt 方法](../html/a2aadcde-a76e-18dc-d5e2-38f9d8adc95f.htm "Decrypt 方法 ")

[Encrypt 方法](../html/9d91dd13-f73a-d7dc-3b20-df0c13b8435e.htm "Encrypt 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptography 方法 |

[DesCryptography](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Decrypt(Byte)](8242f829-258b-6ac5-de66-eefd967379bc.htm) | 对原始的数据进行解密的操作，返回解密之后的二进制原始数据  Decrypt the original data and return the decrypted binary original data |
| 公共方法 | [Decrypt(String)](04da0e28-1e94-6078-3c76-3a91b59c72db.htm) | 针对Base64字符串进行解密操作，转为二进制数据后进行解密，解密之后使用UTF8编码获取最终的字符串数据  Decrypt the Base64 string, convert it to binary data, and decrypt it. After decryption, use UTF8 encoding to obtain the final string data. |
| 公共方法 | [Encrypt(Byte)](7fa9e921-f6fb-699a-5340-4b1efe4c07b1.htm) | 对原始的数据进行加密的操作，返回加密之后的二进制原始数据  Encrypt the original data and return the encrypted binary original data |
| 公共方法 | [Encrypt(String)](43d7b059-2fc0-265a-85c6-7d26a5146f5e.htm) | 针对字符串进行加密，并返回加密后的字符串数据，字符串的编码默认为UTF8，加密后返回Base64编码  Encrypt the string and return the encrypted string data. The encoding of the string is UTF8 by default. After encryption, the Base64 encoding is returned. |
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

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Decrypt 方法 

[原文連結](http://api.hslcommunication.cn/html/a2aadcde-a76e-18dc-d5e2-38f9d8adc95f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 方法](../html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm "DesCryptography 方法")

[Decrypt 方法](../html/a2aadcde-a76e-18dc-d5e2-38f9d8adc95f.htm "Decrypt 方法 ")

[Decrypt 方法 (Byte[])](../html/8242f829-258b-6ac5-de66-eefd967379bc.htm "Decrypt 方法 (Byte[])")

[Decrypt 方法 (String)](../html/04da0e28-1e94-6078-3c76-3a91b59c72db.htm "Decrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptographyDecrypt 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Decrypt(Byte)](8242f829-258b-6ac5-de66-eefd967379bc.htm) | 对原始的数据进行解密的操作，返回解密之后的二进制原始数据  Decrypt the original data and return the decrypted binary original data |
| 公共方法 | [Decrypt(String)](04da0e28-1e94-6078-3c76-3a91b59c72db.htm) | 针对Base64字符串进行解密操作，转为二进制数据后进行解密，解密之后使用UTF8编码获取最终的字符串数据  Decrypt the Base64 string, convert it to binary data, and decrypt it. After decryption, use UTF8 encoding to obtain the final string data. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Decrypt 方法 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/8242f829-258b-6ac5-de66-eefd967379bc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 方法](../html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm "DesCryptography 方法")

[Decrypt 方法](../html/a2aadcde-a76e-18dc-d5e2-38f9d8adc95f.htm "Decrypt 方法 ")

[Decrypt 方法 (Byte[])](../html/8242f829-258b-6ac5-de66-eefd967379bc.htm "Decrypt 方法 (Byte[])")

[Decrypt 方法 (String)](../html/04da0e28-1e94-6078-3c76-3a91b59c72db.htm "Decrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptographyDecrypt 方法 (Byte) |

对原始的数据进行解密的操作，返回解密之后的二进制原始数据  
Decrypt the original data and return the decrypted binary original data

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Decrypt(
	byte[] data
)
```

```
Public Function Decrypt ( 
	data As Byte()
) As Byte()
```

```
public:
virtual array<unsigned char>^ Decrypt(
	array<unsigned char>^ data
) sealed
```

```
abstract Decrypt : 
        data : byte[] -> byte[] 
override Decrypt : 
        data : byte[] -> byte[]
```

#### 参数

data
:   类型：SystemByte  
    等待解密的数据

#### 返回值

类型：Byte  
解密之后的原始二进制数据

#### 实现

[ICryptographyDecrypt(Byte)](f7f478fb-b074-39f1-5620-99ab1137ff01.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[Decrypt 重载](a2aadcde-a76e-18dc-d5e2-38f9d8adc95f.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Decrypt 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/04da0e28-1e94-6078-3c76-3a91b59c72db.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 方法](../html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm "DesCryptography 方法")

[Decrypt 方法](../html/a2aadcde-a76e-18dc-d5e2-38f9d8adc95f.htm "Decrypt 方法 ")

[Decrypt 方法 (Byte[])](../html/8242f829-258b-6ac5-de66-eefd967379bc.htm "Decrypt 方法 (Byte[])")

[Decrypt 方法 (String)](../html/04da0e28-1e94-6078-3c76-3a91b59c72db.htm "Decrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptographyDecrypt 方法 (String) |

针对Base64字符串进行解密操作，转为二进制数据后进行解密，解密之后使用UTF8编码获取最终的字符串数据  
Decrypt the Base64 string, convert it to binary data, and decrypt it. After decryption, use UTF8 encoding to obtain the final string data.

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Decrypt(
	string data
)
```

```
Public Function Decrypt ( 
	data As String
) As String
```

```
public:
virtual String^ Decrypt(
	String^ data
) sealed
```

```
abstract Decrypt : 
        data : string -> string 
override Decrypt : 
        data : string -> string
```

#### 参数

data
:   类型：SystemString  
    base64编码的字符串数据

#### 返回值

类型：String  
最终的解析完成的字符串

#### 实现

[ICryptographyDecrypt(String)](35508463-b243-692d-5d1a-3e001ca0dbbd.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[Decrypt 重载](a2aadcde-a76e-18dc-d5e2-38f9d8adc95f.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Encrypt 方法 

[原文連結](http://api.hslcommunication.cn/html/9d91dd13-f73a-d7dc-3b20-df0c13b8435e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 方法](../html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm "DesCryptography 方法")

[Encrypt 方法](../html/9d91dd13-f73a-d7dc-3b20-df0c13b8435e.htm "Encrypt 方法 ")

[Encrypt 方法 (Byte[])](../html/7fa9e921-f6fb-699a-5340-4b1efe4c07b1.htm "Encrypt 方法 (Byte[])")

[Encrypt 方法 (String)](../html/43d7b059-2fc0-265a-85c6-7d26a5146f5e.htm "Encrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptographyEncrypt 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Encrypt(Byte)](7fa9e921-f6fb-699a-5340-4b1efe4c07b1.htm) | 对原始的数据进行加密的操作，返回加密之后的二进制原始数据  Encrypt the original data and return the encrypted binary original data |
| 公共方法 | [Encrypt(String)](43d7b059-2fc0-265a-85c6-7d26a5146f5e.htm) | 针对字符串进行加密，并返回加密后的字符串数据，字符串的编码默认为UTF8，加密后返回Base64编码  Encrypt the string and return the encrypted string data. The encoding of the string is UTF8 by default. After encryption, the Base64 encoding is returned. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Encrypt 方法 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/7fa9e921-f6fb-699a-5340-4b1efe4c07b1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 方法](../html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm "DesCryptography 方法")

[Encrypt 方法](../html/9d91dd13-f73a-d7dc-3b20-df0c13b8435e.htm "Encrypt 方法 ")

[Encrypt 方法 (Byte[])](../html/7fa9e921-f6fb-699a-5340-4b1efe4c07b1.htm "Encrypt 方法 (Byte[])")

[Encrypt 方法 (String)](../html/43d7b059-2fc0-265a-85c6-7d26a5146f5e.htm "Encrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptographyEncrypt 方法 (Byte) |

对原始的数据进行加密的操作，返回加密之后的二进制原始数据  
Encrypt the original data and return the encrypted binary original data

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Encrypt(
	byte[] data
)
```

```
Public Function Encrypt ( 
	data As Byte()
) As Byte()
```

```
public:
virtual array<unsigned char>^ Encrypt(
	array<unsigned char>^ data
) sealed
```

```
abstract Encrypt : 
        data : byte[] -> byte[] 
override Encrypt : 
        data : byte[] -> byte[]
```

#### 参数

data
:   类型：SystemByte  
    等待加密的数据

#### 返回值

类型：Byte  
加密之后的二进制数据

#### 实现

[ICryptographyEncrypt(Byte)](2870f77e-0e42-a69e-277d-9c1c7a9c699d.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[Encrypt 重载](9d91dd13-f73a-d7dc-3b20-df0c13b8435e.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Encrypt 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/43d7b059-2fc0-265a-85c6-7d26a5146f5e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[DesCryptography 类](../html/6b69c85c-60c2-1133-5d65-e683a2fdb994.htm "DesCryptography 类")

[DesCryptography 方法](../html/c82980e5-bf27-f27d-d97b-d11a8c3a229a.htm "DesCryptography 方法")

[Encrypt 方法](../html/9d91dd13-f73a-d7dc-3b20-df0c13b8435e.htm "Encrypt 方法 ")

[Encrypt 方法 (Byte[])](../html/7fa9e921-f6fb-699a-5340-4b1efe4c07b1.htm "Encrypt 方法 (Byte[])")

[Encrypt 方法 (String)](../html/43d7b059-2fc0-265a-85c6-7d26a5146f5e.htm "Encrypt 方法 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DesCryptographyEncrypt 方法 (String) |

针对字符串进行加密，并返回加密后的字符串数据，字符串的编码默认为UTF8，加密后返回Base64编码  
Encrypt the string and return the encrypted string data. The encoding of the string is UTF8 by default. After encryption, the Base64 encoding is returned.

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Encrypt(
	string data
)
```

```
Public Function Encrypt ( 
	data As String
) As String
```

```
public:
virtual String^ Encrypt(
	String^ data
) sealed
```

```
abstract Encrypt : 
        data : string -> string 
override Encrypt : 
        data : string -> string
```

#### 参数

data
:   类型：SystemString  
    等待加密的字符串

#### 返回值

类型：String  
加密后的Base64编码

#### 实现

[ICryptographyEncrypt(String)](1c9b3294-3a48-b8f6-7d7c-42113ea43c58.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DesCryptography 类](6b69c85c-60c2-1133-5d65-e683a2fdb994.htm)

[Encrypt 重载](9d91dd13-f73a-d7dc-3b20-df0c13b8435e.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslCertificate 类

[原文連結](http://api.hslcommunication.cn/html/53cc213c-50c4-52f8-0a00-76009a75efd6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[HslCertificate 类](../html/53cc213c-50c4-52f8-0a00-76009a75efd6.htm "HslCertificate 类")

[HslCertificate 构造函数](../html/227d4a4b-3e7d-a398-b6f8-65b1ea756b7e.htm "HslCertificate 构造函数 ")

[HslCertificate 属性](../html/61fc01fe-39b3-c58c-0508-19c0f4c693f4.htm "HslCertificate 属性")

[HslCertificate 方法](../html/f4059711-23ca-ce26-12ce-403c216f4208.htm "HslCertificate 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCertificate 类 |

基于RSA加密模型的证书，支持自定义的颁发证书，以及校验证书合法性  
Certificates based on RSA encryption model support custom issuance certificates and verify certificate legitimacy

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Core.SecurityHslCertificate

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class HslCertificate
```

```
Public Class HslCertificate
```

```
public ref class HslCertificate
```

```
type HslCertificate =  class end
```

HslCertificate 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HslCertificate](2100e805-d87d-efa8-bbc9-aa292bd7eb47.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [HslCertificate(Byte, Byte)](daf2e719-932a-3b87-24b9-938142b88b39.htm) | 使用指定的公钥，私钥来实例化一个的对象  An object is instantiated using the specified public key and private key |
| 公共方法 | [HslCertificate(RSACryptoServiceProvider, RSACryptoServiceProvider)](e0a147f7-89bb-ced4-6687-83d890ed9d5b.htm) | 使用指定的公钥，私钥来实例化一个的对象  An object is instantiated using the specified public key and private key |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CreateTime](9580f339-e61d-3811-2b6f-2f055fa21a9c.htm) | 发证日期 |
| 公共属性 | [Descriptions](d766685d-0618-a18b-d99b-31e5b665ffac.htm) | 证书的其他描述信息 |
| 公共属性 | [EffectiveHours](8c9459a3-417b-3964-fa64-22e2ed97a5f3.htm) | 有效小时数，小于等于0 表示无期限 |
| 公共属性 | [From](694e4b92-a8c2-1256-fb2c-7bc02de2deec.htm) | 证书的颁发者 |
| 公共属性 | [KeyWord](fd7cc849-a59b-3213-d235-db6f8d330381.htm) | 获取或设置当前证书的关键字，可以用来给证书做分类 |
| 公共属性 | [NotAfter](759bf0fd-3c4b-c11c-abc3-7520c70aa477.htm) | 证书有效的截止时间 |
| 公共属性 | [NotBefore](2ef9832a-f2e8-b69a-4e75-50974bb4b498.htm) | 证书有效的起始时间 |
| 公共属性 | [PublicKey](b6e41bb0-dbdb-a719-48fe-fbd5d1e4a766.htm) | 证书的公钥信息 |
| 公共属性 | [To](9f1e3aff-0343-e43c-189d-190c69784a8c.htm) | 证书的持有者 |
| 公共属性 | [UniqueID](f6f1da82-cddb-53da-eb21-f1714c8465df.htm) | 获取或设置当前证书的唯一编号信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [CreateFrom](c4617903-7a8b-259d-0099-9b7870fd4c75.htm) | 从证书的原始字节创建一个HslCertificate对象，方便浏览证书的基本信息。  Create a HslCertificate object from the original bytes of the certificate to facilitate browsing the basic information of the certificate. |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetSaveBytes](621577f1-d985-5ab1-b4e5-a17ebc19c2b3.htm) | 获取当前证书的原始字节信息，可以存储到文件中，必须提供私钥信息，否则无法进行签名的操作  Gets the raw byte information of the current certificate, which can be stored in a file, and the private key information must be provided, otherwise the signing operation cannot be performed |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [LoadFrom](902551ae-ae6f-5095-85d0-4024b9e5f5da.htm) | 从文件的二进制数据中加载相关的参数 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [VerifyCer](47e8484f-45e8-d3e2-dce5-b4a4bde97856.htm) | 使用给定的公钥，校验当前的证书是否合法的，如果公钥为 null，则直接校验证书本身是否合法。  Use the given public key to verify whether the current certificate is valid, and if the public key is null, directly verify whether the certificate itself is valid. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)备注

证书可以用于接口的权限认证，不需要修改接口源代码或是配置文件，颁发证书就可以修改用户的权限，而且只要保密好私钥，那么证书本身就无法伪造，具有极高的安全性，具体用法参考示例代码。  
The certificate can be used for the permission authentication of the interface, no need to modify the interface source code or configuration file,
the issuance of the certificate can modify the user's permissions, and as long as the private key is kept secret, then the certificate itself cannot be forged,
with extremely high security, the specific usage refer to the sample code.

![](../icons/SectionExpanded.png)示例

证书这部分的功能主要分为，制作证书，以及校验证书，至于为什么不使用X509Certificate2证书的形式，
因为这种证书都是要授信机构颁发的，自己颁发的证书验签不了，所以在本库里提供一个用于自己颁发，自己验签的证书。  
假设我们有一些API接口需要使用证书来控制权限，有调用时间检验的，或是按接口名称校验的，或是按调用次数来校验的，接口见下面的代码。

接口的权限控制示例

[复制](# "复制")

```
// 先使用hslcommunicationdemo，生成公私钥的密钥对，例如下面的
private string pubKey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCwAs1K7u" +
    "qRb89J0DJZ/rFz5dHHWC0vTiGDayy//cc8jvd3MHY3aUiH2l0Th69rbBY81C" +
    "I5dzaUKMml393VOf1qTzt+CLV0XrVOPnZ3I9SuTmBh8dyB2pX/9wm/X3JL/Jg9tE" +
    "OofqmD6ChLPn01Nf5xiaRpyG1fntxJO8jS3+ckJQIDAQAB";
private string priKey = "MIICXgIBAAKBgQCwAs1K7uqRb89J0DJZ/rFz5dHHWC0vTiGDayy//cc8jvd3" +
    "MHY3aUiH2l0Th69rbBY81CI5dzaUKMml393VOf1qTzt+CLV0XrVOPnZ3I9SuTmBh8dyB2pX/9wm/X3JL/Jg9t" +
    "EOofqmD6ChLPn01Nf5xiaRpyG1fntxJO8jS3+ckJQIDAQABAoGBAKJt84ZzdYVB6cn0mIL0yE5siGuutAQz1jmx" +
    "blq/mF4KkdUso+v+/TBiFAMx9HDuDpeULQFbQsL34R9yuQVNvs4U+Oo0vqAD8aXrS9bDZWuPvpKM4fCLof5gxVnLK" +
    "5c2MV/Nq060ztcNBnzEpLhGLbcPZbTOShlAVnml4UJxsesBAkEA19JXl2nEEjvSIEzzXjsl0cp18OXO7C4UNIE5sEpQ" +
    "eNbJbgfcHRavd82OvrLzYKy4aHRQZvCEQBo+1ja7EdrOtwJBANDHJ9/rs88+hb9oWs0zeedpnr+5weijefhtLzKFBr/b6MP" +
    "iSAvYIUGie7UD+MedC37MHpQmiyPFUR3wd1onCAMCQEIB5yA8DOe2vBF894H+fRao4WGIJ708apmUXFx6nhoXNwwcA0oyQgBWMi" +
    "5I9P7AVLY9vMl/DjO80zMZEiXn0gkCQQCNeyjuedA0VW5Vs+Y0wPmPdOf7jWp2uLHjRDqDxdo5ElG+zQ192U1QKUgj9zkULLyph" +
    "LIEO8p2xzSdUWPQWDr5AkEAp9FfsofKAYIa6+rb4EWn9zavWz+gfIX9J/Zz3d7Mi0JnVX8dgC5o6yOxhM1h6oDgeHxeqZPTvtRUtI4U5dAmHQ==";

// 假设下面的接口就是使用证书校验的，不是谁都可以访问的
// Assume that the following interface is authenticated using a certificate and cannot be accessed by anyone
[HslMqttApi]
public OperateResult<int> ApiAdd( string certificate, int a, int b )
{
    // 下面演示如何针对 certificate 校验对象
    if (!HslCommunication.Core.Security.HslCertificate.VerifyCer( Convert.FromBase64String( pubKey ), Convert.FromBase64String( certificate )))
    {
        // 调用失败
        return new OperateResult<int>( "Called api failed, check certificate failed" );
    }

    // 如果证书校验通过，则开始解析证书内容
    HslCommunication.Core.Security.HslCertificate cert = HslCommunication.Core.Security.HslCertificate.CreateFrom( Convert.FromBase64String( certificate ) );
    if (cert.NotBefore <= DateTime.Now && DateTime.Now <= cert.NotAfter)
    {
        // 一般证书都是有时间有效期的，所以这里可以根据有效期决定是否开放接口
        return new OperateResult<int>( "Called api failed, the current certificate validity period has expired." );
    }

    return OperateResult.CreateSuccessResult( a + b );
}

// 假设下面的接口就是使用证书校验的，不仅要校验时间是否符合，还要校验系统名称是否匹配，证书的权限是否包含当前的接口名称
// Assume that the following interface is verified using a certificate, not only to verify whether the time is in line,
// but also to verify whether the system name matches, and whether the certificate's authority contains the current interface name
[HslMqttApi]
public OperateResult<int> ApiMulti( string certificate, int a, int b )
{
    // 下面演示如何针对 certificate 校验对象
    if (!HslCommunication.Core.Security.HslCertificate.VerifyCer( Convert.FromBase64String( pubKey ), Convert.FromBase64String( certificate ) ))
    {
        // 调用失败
        return new OperateResult<int>( "Called api failed, check certificate failed" );
    }

    // 如果证书校验通过，则开始解析证书内容
    HslCommunication.Core.Security.HslCertificate cert = HslCommunication.Core.Security.HslCertificate.CreateFrom( Convert.FromBase64String( certificate ) );
    if (cert.NotBefore <= DateTime.Now && DateTime.Now <= cert.NotAfter)
    {
        // 一般证书都是有时间有效期的，所以这里可以根据有效期决定是否开放接口
        return new OperateResult<int>( "Called api failed, the current certificate validity period has expired." );
    }

    // 如果希望多个系统，接口，软件的证书都用同一套的公私钥，那么不同系统的证书校验，可以使用KeyWord来区分，例如下面的 HslCommunication 软件校验
    // 否则只能每一套系统，分配一个公私钥密钥对了
    if (cert.KeyWord != "HslCommunication") return new OperateResult<int>( "Certificate key word not correct!" );

    // 如果希望你得证书支持接口列表信息，可以使用证书里的集合数据功能
    if (cert.Descriptions != null && cert.Descriptions.ContainsKey( "ApiMulti" ))
    {
        // 证书当前的接口列表里，包含了当前的接口名称，有权限调用
        return OperateResult.CreateSuccessResult( a * b );
    }
    else
    {
        return new OperateResult<int>( "Called ApiMulti failed, the current certificate does not have access." );
    }
}

// 假设下面的接口就是使用证书校验的，不仅要校验时间是否符合，还希望对证书持有者限制调用接口的次数，总次数写入到证书里去，这样比较灵活，不需要修改代码，甚至配置文件，只要颁发不同的证书即可
// Assuming that the following interface is using certificate verification, not only to verify whether the time is met, but also to limit the number of times the interface is called to the certificate holder,
// the total number of times written to the certificate, which is more flexible, no need to modify the code, or even the configuration file, as long as different certificates are issued
[HslMqttApi]
public OperateResult<int> ApiDivision( string certificate, int a, int b )
{
    // 下面演示如何针对 certificate 校验对象
    if (!HslCommunication.Core.Security.HslCertificate.VerifyCer( Convert.FromBase64String( pubKey ), Convert.FromBase64String( certificate ) ))
    {
        // 调用失败
        return new OperateResult<int>( "Called api failed, check certificate failed" );
    }

    // 如果证书校验通过，则开始解析证书内容
    HslCommunication.Core.Security.HslCertificate cert = HslCommunication.Core.Security.HslCertificate.CreateFrom( Convert.FromBase64String( certificate ) );
    if (cert.NotBefore <= DateTime.Now && DateTime.Now <= cert.NotAfter)
    {
        // 一般证书都是有时间有效期的，所以这里可以根据有效期决定是否开放接口
        return new OperateResult<int>( "Called api failed, the current certificate validity period has expired." );
    }

    // 如果希望多个系统，接口，软件的证书都用同一套的公私钥，那么不同系统的证书校验，可以使用KeyWord来区分，例如下面的 HslCommunication 软件校验
    // 否则只能每一套系统，分配一个公私钥密钥对了
    if (cert.KeyWord != "HslCommunication") return new OperateResult<int>( "Certificate key word not correct!" );


    // cert.EffectiveHours 是个 int 数据，用来存放接口的调用次数，当然使用另一个字符串属性 cert.UniqueID 也可以。
    bool enable = false; // 检查是否通过
    lock (dictLock)
    {
        // 没有
        if (!apiCalledCount.ContainsKey( cert.To )) apiCalledCount.Add( cert.To, 0 );
        enable = apiCalledCount[cert.To] < cert.EffectiveHours;

        if (enable) apiCalledCount[cert.To]++;
    }

    if (!enable) return new OperateResult<int>( "Called ApiDivision failed, The number of calls has reached the maximum." );
    return OperateResult.CreateSuccessResult( a / b ); // 返回正常的调用信息
}

// 当前系统的所有API调用次数，初始化时，应该从本地或是数据库加载之前的次数，这次简略处理
private Dictionary<string, long> apiCalledCount = new Dictionary<string, long>( );
private object dictLock = new object( );
```

当然我们还可以自己颁发证书，注意，这时候的私钥就非常有用了，私钥丢了，就发不了证书了。如果要重新生成公私钥，那么之前发出去的证书都失效了。

颁发证书的例子

[复制](# "复制")

```
// 如果需要实现发证的话，就是下面的代码，当然也可以使用 hslcommunicationdemo 来颁发证书的
public void CreateCert( )
{
    // 还是用了上个方法创建的公私钥对
    HslCommunication.Core.Security.HslCertificate cert = new HslCommunication.Core.Security.HslCertificate( Convert.FromBase64String( pubKey ), Convert.FromBase64String( priKey ) );
    cert.From = "Company A";   // 发证方
    cert.To = "Company B";     // 使用方，持有者
    cert.CreateTime = DateTime.Today;  // 发证日期
    cert.NotBefore = new DateTime( 2023, 2, 16 ); // 起始日期
    cert.NotAfter = new DateTime( 2023, 3, 15 );  // 假设一个月有效期
    cert.KeyWord = "HslCommunication";    // 可以根据实际情况指定
    cert.Descriptions = new Dictionary<string, string>( )
    {
        { "ApiMulti", "1" }
    };

    // 获取证书的字节
    byte[] source = cert.GetSaveBytes( );

    // 证书写入到文件
    System.IO.File.WriteAllBytes( "C:\\cert.cert", source );

    // 变成字符串发给别人使用
    string str = Convert.ToBase64String( source );
    System.IO.File.WriteAllText( "C:\\cert.txt", str, Encoding.UTF8 );
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslCertificate 构造函数 

[原文連結](http://api.hslcommunication.cn/html/227d4a4b-3e7d-a398-b6f8-65b1ea756b7e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[HslCertificate 类](../html/53cc213c-50c4-52f8-0a00-76009a75efd6.htm "HslCertificate 类")

[HslCertificate 构造函数](../html/227d4a4b-3e7d-a398-b6f8-65b1ea756b7e.htm "HslCertificate 构造函数 ")

[HslCertificate 构造函数](../html/2100e805-d87d-efa8-bbc9-aa292bd7eb47.htm "HslCertificate 构造函数 ")

[HslCertificate 构造函数 (Byte[], Byte[])](../html/daf2e719-932a-3b87-24b9-938142b88b39.htm "HslCertificate 构造函数 (Byte[], Byte[])")

[HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider)](../html/e0a147f7-89bb-ced4-6687-83d890ed9d5b.htm "HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCertificate 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HslCertificate](2100e805-d87d-efa8-bbc9-aa292bd7eb47.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [HslCertificate(Byte, Byte)](daf2e719-932a-3b87-24b9-938142b88b39.htm) | 使用指定的公钥，私钥来实例化一个的对象  An object is instantiated using the specified public key and private key |
| 公共方法 | [HslCertificate(RSACryptoServiceProvider, RSACryptoServiceProvider)](e0a147f7-89bb-ced4-6687-83d890ed9d5b.htm) | 使用指定的公钥，私钥来实例化一个的对象  An object is instantiated using the specified public key and private key |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCertificate 类](53cc213c-50c4-52f8-0a00-76009a75efd6.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslCertificate 构造函数 

[原文連結](http://api.hslcommunication.cn/html/2100e805-d87d-efa8-bbc9-aa292bd7eb47.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[HslCertificate 类](../html/53cc213c-50c4-52f8-0a00-76009a75efd6.htm "HslCertificate 类")

[HslCertificate 构造函数](../html/227d4a4b-3e7d-a398-b6f8-65b1ea756b7e.htm "HslCertificate 构造函数 ")

[HslCertificate 构造函数](../html/2100e805-d87d-efa8-bbc9-aa292bd7eb47.htm "HslCertificate 构造函数 ")

[HslCertificate 构造函数 (Byte[], Byte[])](../html/daf2e719-932a-3b87-24b9-938142b88b39.htm "HslCertificate 构造函数 (Byte[], Byte[])")

[HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider)](../html/e0a147f7-89bb-ced4-6687-83d890ed9d5b.htm "HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCertificate 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslCertificate()
```

```
Public Sub New
```

```
public:
HslCertificate()
```

```
new : unit -> HslCertificate
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCertificate 类](53cc213c-50c4-52f8-0a00-76009a75efd6.htm)

[HslCertificate 重载](227d4a4b-3e7d-a398-b6f8-65b1ea756b7e.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslCertificate 构造函数 (Byte[], Byte[])

[原文連結](http://api.hslcommunication.cn/html/daf2e719-932a-3b87-24b9-938142b88b39.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[HslCertificate 类](../html/53cc213c-50c4-52f8-0a00-76009a75efd6.htm "HslCertificate 类")

[HslCertificate 构造函数](../html/227d4a4b-3e7d-a398-b6f8-65b1ea756b7e.htm "HslCertificate 构造函数 ")

[HslCertificate 构造函数](../html/2100e805-d87d-efa8-bbc9-aa292bd7eb47.htm "HslCertificate 构造函数 ")

[HslCertificate 构造函数 (Byte[], Byte[])](../html/daf2e719-932a-3b87-24b9-938142b88b39.htm "HslCertificate 构造函数 (Byte[], Byte[])")

[HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider)](../html/e0a147f7-89bb-ced4-6687-83d890ed9d5b.htm "HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCertificate 构造函数 (Byte, Byte) |

使用指定的公钥，私钥来实例化一个的对象  
An object is instantiated using the specified public key and private key

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslCertificate(
	byte[] pubKey,
	byte[] priKey
)
```

```
Public Sub New ( 
	pubKey As Byte(),
	priKey As Byte()
)
```

```
public:
HslCertificate(
	array<unsigned char>^ pubKey, 
	array<unsigned char>^ priKey
)
```

```
new : 
        pubKey : byte[] * 
        priKey : byte[] -> HslCertificate
```

#### 参数

pubKey
:   类型：SystemByte  
    公钥的二进制数据

priKey
:   类型：SystemByte  
    私钥的二进制数据

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCertificate 类](53cc213c-50c4-52f8-0a00-76009a75efd6.htm)

[HslCertificate 重载](227d4a4b-3e7d-a398-b6f8-65b1ea756b7e.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider)

[原文連結](http://api.hslcommunication.cn/html/e0a147f7-89bb-ced4-6687-83d890ed9d5b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[HslCertificate 类](../html/53cc213c-50c4-52f8-0a00-76009a75efd6.htm "HslCertificate 类")

[HslCertificate 构造函数](../html/227d4a4b-3e7d-a398-b6f8-65b1ea756b7e.htm "HslCertificate 构造函数 ")

[HslCertificate 构造函数](../html/2100e805-d87d-efa8-bbc9-aa292bd7eb47.htm "HslCertificate 构造函数 ")

[HslCertificate 构造函数 (Byte[], Byte[])](../html/daf2e719-932a-3b87-24b9-938142b88b39.htm "HslCertificate 构造函数 (Byte[], Byte[])")

[HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider)](../html/e0a147f7-89bb-ced4-6687-83d890ed9d5b.htm "HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCertificate 构造函数 (RSACryptoServiceProvider, RSACryptoServiceProvider) |

使用指定的公钥，私钥来实例化一个的对象  
An object is instantiated using the specified public key and private key

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HslCertificate(
	RSACryptoServiceProvider pubKey,
	RSACryptoServiceProvider priKey
)
```

```
Public Sub New ( 
	pubKey As RSACryptoServiceProvider,
	priKey As RSACryptoServiceProvider
)
```

```
public:
HslCertificate(
	RSACryptoServiceProvider^ pubKey, 
	RSACryptoServiceProvider^ priKey
)
```

```
new : 
        pubKey : RSACryptoServiceProvider * 
        priKey : RSACryptoServiceProvider -> HslCertificate
```

#### 参数

pubKey
:   类型：System.Security.CryptographyRSACryptoServiceProvider  
    公钥的对象

priKey
:   类型：System.Security.CryptographyRSACryptoServiceProvider  
    私钥的对象

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCertificate 类](53cc213c-50c4-52f8-0a00-76009a75efd6.htm)

[HslCertificate 重载](227d4a4b-3e7d-a398-b6f8-65b1ea756b7e.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HslCertificate 属性

[原文連結](http://api.hslcommunication.cn/html/61fc01fe-39b3-c58c-0508-19c0f4c693f4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[HslCertificate 类](../html/53cc213c-50c4-52f8-0a00-76009a75efd6.htm "HslCertificate 类")

[HslCertificate 属性](../html/61fc01fe-39b3-c58c-0508-19c0f4c693f4.htm "HslCertificate 属性")

[CreateTime 属性](../html/9580f339-e61d-3811-2b6f-2f055fa21a9c.htm "CreateTime 属性 ")

[Descriptions 属性](../html/d766685d-0618-a18b-d99b-31e5b665ffac.htm "Descriptions 属性 ")

[EffectiveHours 属性](../html/8c9459a3-417b-3964-fa64-22e2ed97a5f3.htm "EffectiveHours 属性 ")

[From 属性](../html/694e4b92-a8c2-1256-fb2c-7bc02de2deec.htm "From 属性 ")

[KeyWord 属性](../html/fd7cc849-a59b-3213-d235-db6f8d330381.htm "KeyWord 属性 ")

[NotAfter 属性](../html/759bf0fd-3c4b-c11c-abc3-7520c70aa477.htm "NotAfter 属性 ")

[NotBefore 属性](../html/2ef9832a-f2e8-b69a-4e75-50974bb4b498.htm "NotBefore 属性 ")

[PublicKey 属性](../html/b6e41bb0-dbdb-a719-48fe-fbd5d1e4a766.htm "PublicKey 属性 ")

[To 属性](../html/9f1e3aff-0343-e43c-189d-190c69784a8c.htm "To 属性 ")

[UniqueID 属性](../html/f6f1da82-cddb-53da-eb21-f1714c8465df.htm "UniqueID 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCertificate 属性 |

[HslCertificate](53cc213c-50c4-52f8-0a00-76009a75efd6.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CreateTime](9580f339-e61d-3811-2b6f-2f055fa21a9c.htm) | 发证日期 |
| 公共属性 | [Descriptions](d766685d-0618-a18b-d99b-31e5b665ffac.htm) | 证书的其他描述信息 |
| 公共属性 | [EffectiveHours](8c9459a3-417b-3964-fa64-22e2ed97a5f3.htm) | 有效小时数，小于等于0 表示无期限 |
| 公共属性 | [From](694e4b92-a8c2-1256-fb2c-7bc02de2deec.htm) | 证书的颁发者 |
| 公共属性 | [KeyWord](fd7cc849-a59b-3213-d235-db6f8d330381.htm) | 获取或设置当前证书的关键字，可以用来给证书做分类 |
| 公共属性 | [NotAfter](759bf0fd-3c4b-c11c-abc3-7520c70aa477.htm) | 证书有效的截止时间 |
| 公共属性 | [NotBefore](2ef9832a-f2e8-b69a-4e75-50974bb4b498.htm) | 证书有效的起始时间 |
| 公共属性 | [PublicKey](b6e41bb0-dbdb-a719-48fe-fbd5d1e4a766.htm) | 证书的公钥信息 |
| 公共属性 | [To](9f1e3aff-0343-e43c-189d-190c69784a8c.htm) | 证书的持有者 |
| 公共属性 | [UniqueID](f6f1da82-cddb-53da-eb21-f1714c8465df.htm) | 获取或设置当前证书的唯一编号信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCertificate 类](53cc213c-50c4-52f8-0a00-76009a75efd6.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CreateTime 属性 

[原文連結](http://api.hslcommunication.cn/html/9580f339-e61d-3811-2b6f-2f055fa21a9c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Security](../html/9951c682-3526-ca1b-21eb-2be212e912ea.htm "HslCommunication.Core.Security")

[HslCertificate 类](../html/53cc213c-50c4-52f8-0a00-76009a75efd6.htm "HslCertificate 类")

[HslCertificate 属性](../html/61fc01fe-39b3-c58c-0508-19c0f4c693f4.htm "HslCertificate 属性")

[CreateTime 属性](../html/9580f339-e61d-3811-2b6f-2f055fa21a9c.htm "CreateTime 属性 ")

[Descriptions 属性](../html/d766685d-0618-a18b-d99b-31e5b665ffac.htm "Descriptions 属性 ")

[EffectiveHours 属性](../html/8c9459a3-417b-3964-fa64-22e2ed97a5f3.htm "EffectiveHours 属性 ")

[From 属性](../html/694e4b92-a8c2-1256-fb2c-7bc02de2deec.htm "From 属性 ")

[KeyWord 属性](../html/fd7cc849-a59b-3213-d235-db6f8d330381.htm "KeyWord 属性 ")

[NotAfter 属性](../html/759bf0fd-3c4b-c11c-abc3-7520c70aa477.htm "NotAfter 属性 ")

[NotBefore 属性](../html/2ef9832a-f2e8-b69a-4e75-50974bb4b498.htm "NotBefore 属性 ")

[PublicKey 属性](../html/b6e41bb0-dbdb-a719-48fe-fbd5d1e4a766.htm "PublicKey 属性 ")

[To 属性](../html/9f1e3aff-0343-e43c-189d-190c69784a8c.htm "To 属性 ")

[UniqueID 属性](../html/f6f1da82-cddb-53da-eb21-f1714c8465df.htm "UniqueID 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCertificateCreateTime 属性 |

发证日期

**命名空间：**
 [HslCommunication.Core.Security](9951c682-3526-ca1b-21eb-2be212e912ea.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DateTime CreateTime { get; set; }
```

```
Public Property CreateTime As DateTime
	Get
	Set
```

```
public:
property DateTime CreateTime {
	DateTime get ();
	void set (DateTime value);
}
```

```
member CreateTime : DateTime with get, set
```

#### 属性值

类型：DateTime

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCertificate 类](53cc213c-50c4-52f8-0a00-76009a75efd6.htm)

[HslCommunication.Core.Security 命名空间](9951c682-3526-ca1b-21eb-2be212e912ea.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)