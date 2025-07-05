# Go 工具包使用文档

## 目录
- [哈希与加密](#哈希与加密)
- [随机生成](#随机生成)
- [字符串处理](#字符串处理)
- [环境变量](#环境变量)
- [路径处理](#路径处理)
- [正则表达式](#正则表达式)
- [进度显示](#进度显示)

## 哈希与加密
### 🔐 密码加密相关
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `HashAndSalt` | `value string` | `(string, error)` | 使用 bcrypt 算法对密码进行加密，推荐用于密码存储 |
| `ComparePasswords` | `hashedPassword, password string` | `bool` | 验证密码是否匹配，用于登录验证 |

### 🔑 基础哈希算法
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `MD5Hash` | `text string` | `string` | 计算文本的 MD5 值（不推荐用于密码） |
| `SHA1Hash` | `text string` | `string` | 计算文本的 SHA1 值 |
| `SHA256Hash` | `text string` | `string` | 计算文本的 SHA256 值 |
| `SHA512Hash` | `text string` | `string` | 计算文本的 SHA512 值（最安全） |

### 🛡️ 高级哈希功能
| 函数名 | 参数 | 返回��� | 描述 |
|-------|------|--------|------|
| `GenerateSalt` | `length int` | `string` | 生成指定长度的随机盐值 |
| `HashWithSalt` | `password, salt string` | `string` | 将密码与盐值组合后进行 SHA256 哈希 |
| `HashWithPepper` | `password, pepper string` | `string` | 将密码与固定盐值(pepper)组合后哈希 |
| `DoubleHash` | `text string` | `string` | 对文本进行双重 SHA256 哈希，提升安全性 |

### 🔍 哈希验证
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `VerifyHash` | `text, hash, algorithm string` | `bool` | 验证文本与哈希值是否匹配 |

支持的算法：
- `"md5"`
- `"sha1"`
- `"sha256"`
- `"sha512"`

### 📝 Base64 编解码
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `Base64Encode` | `text string` | `string` | 将文本转换为 Base64 编码 |
| `Base64Decode` | `encodedText string` | `(string, error)` | 解码 Base64 字符串 |

## 随机生成
### 🎲 基础随机函数
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `RandomCode` | `length int` | `string` | 生成指定长度的随机数字码 |
| `RandomString` | `length int` | `string` | 生成指定长度的随机字符串(含大小写字母和数字) |
| `RandomInt` | `min, max int` | `int` | 生成指定范围内的随机整数 |
| `RandomFloat` | `min, max float64` | `float64` | 生成指定范围内的随机浮点数 |
| `RandomBool` | - | `bool` | 生成随机布尔值 |

### 🔤 字符串随机生成
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `RandomLowerString` | `length int` | `string` | 生成指定长度的随机小写字母字符串 |
| `RandomUpperString` | `length int` | `string` | 生成指定长度的随机大写字母字符串 |
| `RandomChoice[T any]` | `slice []T` | `T` | 从切片中随机选择一个元素 |
| `RandomTime` | `start, end time.Time` | `time.Time` | 生成指定范围内的随机时间 |

## 字符串处理
### 📝 大小写转换
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `LcFirst` | `str string` | `string` | 首字母小写 |
| `Lc` | `str string` | `string` | 转小写 |
| `UcFirst` | `str string` | `string` | 首字母大写 |

### 🔄 命名转换
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `SnakeCase` | `str string` | `string` | 转换为下划线命名 (例如: HelloWorld -> hello_world) |
| `CamelCase` | `str string` | `string` | 转换为驼峰命名，支持特殊单词 (例如: user_id -> userID, ip_location -> ipLocation) |
| `PascalCase` | `str string` | `string` | 转换为帕斯卡命名，首字母大写的驼峰 (例如: user_id -> UserID) |
| `CamelColon` | `str string` | `string` | 驼峰转路径格式 (例如: ProcessTask -> process::task) |

#### 特殊单词处理
以下单词在转换时会保持大写：
- ID, IP, URL, URI, API, UUID, HTML, XML, JSON
- YAML, CSS, SQL, HTTP(S), FTP, SSH, SSL
- TCP, UDP, GUI, UI, CDN, DNS, CPU, GPU
- RAM, SDK, JWT, OAuth

示例：

### 🔠 其他字符串工具
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `StrPtr` | `str string` | `*string` | 返回字符串的指针 |
| `Pluralize` | `word string` | `string` | 将单词转换为复数形式 |
| `Able` | `name string` | `string` | 添加 "able" 后缀 |
| `IsInternalType` | `name string` | `bool` | 检查是否为内部类型(以__开头) |

## 环境变量
### 🌍 环境变量获取
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `GetEnv` | `key string, defaultValue ...string` | `string` | 获取环境变量值 |
| `GetEnvInt` | `key string, defaultValue ...int` | `int` | 获取整数类型环境变量 |
| `GetEnvInt64` | `key string, defaultValue ...int64` | `int64` | 获取 int64 类型环境变量 |
| `GetEnvBool` | `key string, defaultValue ...bool` | `bool` | 获取布尔类型环境变量 |
| `GetEnvFloat64` | `key string, defaultValue ...float64` | `float64` | 获取浮点类型环境变量 |

### 📚 数组类型环境变量
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `GetEnvArray` | `key, sep string, defaultValue ...[]string` | `[]string` | 获取字符串数组环境变量 |
| `GetEnvIntArray` | `key, sep string, defaultValue ...[]int` | `[]int` | 获取整数数组环境变量 |
| `GetEnvInt64Array` | `key, sep string, defaultValue ...[]int64` | `[]int64` | 获取 int64 数组环境变量 |
| `GetEnvFloat64Array` | `key, sep string, defaultValue ...[]float64` | `[]float64` | 获取浮点数组环境变量 |

## 路径处理
### 📂 路径工具
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `GetProjectPath` | - | `(string, error)` | 获取项目根路径 |
| `GetModPath` | `projectPath *string` | `(string, error)` | 获取模块路径 |
| `GetPkgPath` | `projectPath, filePath string` | `(string, error)` | 获取包路径 |
| `GetGoPath` | - | `string` | 获取 GOPATH |
| `GetFilePath` | `path string` | `(string, error)` | 获取文件路径 |
| `GetFileDir` | `path string` | `(string, error)` | 获取文件所在目录 |
| `MkdirAll` | `path string` | `error` | 创建多级目录 |

## 正则表达式
### 📋 格式验证
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `IsValidPhoneNumber` | `phone string` | `bool` | 验证中国手机号 |
| `IsValidIP` | `ip string` | `bool` | 验证 IPv4 地址 |
| `IsValidIPv6` | `ip string` | `bool` | 验证 IPv6 地址 |
| `IsValidIDCard` | `id string` | `bool` | 验证身份证号(支持15位、18位和外国人永久居留证) |
| `IsValidBankCard` | `card string` | `bool` | 验证银行卡号(16-19位) |
| `IsValidPostCode` | `code string` | `bool` | 验证中国大陆邮编 |

### 🕒 时间日期验证
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `IsValidDate` | `date string` | `bool` | 验证日期格式(YYYY-MM-DD) |
| `IsValidTime` | `timeStr string` | `bool` | 验证时间格式(HH:MM:SS) |
| `IsValidDateTime` | `datetime string` | `bool` | 验证日期时间格式(YYYY-MM-DD HH:MM:SS) |
| `IsValidAmount` | `amount string` | `bool` | 验证金额格式(123.45) |

## 进度显示
### 🔄 进度条
| 函数名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| `SmoothProgress` | `start, end int, status string, duration time.Duration, keepVisible bool` | - | 显示平滑进度条动画 |

