# 🚀 pkg-go

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hchicken/pkg-go)](https://goreportcard.com/report/github.com/hchicken/pkg-go)
[![Release](https://img.shields.io/github/v/release/hchicken/pkg-go)](https://github.com/hchicken/pkg-go/releases)

> 🎯 **企业级 Go 语言基础工具库集合**  
> 提供高性能、易用的工具包和扩展功能，助力开发者快速构建生产级 Go 应用程序

## ✨ 特性

- 🏗️ **模块化设计** - 每个包独立发布，按需引入
- 🚀 **高性能** - 经过生产环境验证的高效实现
- 📚 **完整文档** - 详细的 API 文档和使用示例
- 🔧 **易于集成** - 简洁的 API 设计，快速上手
- 🛡️ **类型安全** - 充分利用 Go 的类型系统
- 🧪 **测试覆盖** - 完善的单元测试保障代码质量

## 📦 包概览

### 🔧 核心工具包

| 包名 | 最新版本 | 描述 | 主要功能 |
|------|----------|------|----------|
| [**arraryx**](#arraryx) | `latest` | 数组操作工具包 | 去重、差集、元素操作 |
| [**bytex**](#bytex) | `latest` | 字节操作工具包 | 大小转换、格式化 |
| [**stringx**](#stringx) | `v1.0.3` | 字符串处理工具包 | 加密、编码、压缩、随机生成 |
| [**structx**](#structx) | `v1.0.1` | 结构体操作工具包 | 结构体转换、序列化 |
| [**date**](#date) | `latest` | 时间日期处理工具包 | 时间戳转换、格式解析 |
| [**filex**](#filex) | `latest` | 文件操作工具包 | 路径获取、文件处理 |
| [**osx**](#osx) | `v1.0.1` | 操作系统相关工具包 | 环境变量处理 |

### 🌐 Web 框架扩展

| 包名 | 最新版本 | 描述 | 主要功能 |
|------|----------|------|----------|
| [**ginx**](#ginx) | `latest` | Gin 框架扩展包 | 参数绑定、响应处理、验证器 |
| [**httpx**](#httpx) | `latest` | HTTP 客户端工具包 | 基于 Resty 的 HTTP 请求 |
| [**jwtx**](#jwtx) | `latest` | JWT 认证工具包 | Token 生成与解析 |

### 🗄️ 数据库与缓存

| 包名 | 最新版本 | 描述 | 主要功能 |
|------|----------|------|----------|
| [**gormx**](#gormx) | `latest` | GORM 数据库扩展包 | 连接池、模型定义、查询构建 |
| [**cache**](#cache) | `latest` | 缓存操作工具包 | Redis 连接池管理 |

### 📊 中间件与消息队列

| 包名 | 最新版本 | 描述 | 主要功能 |
|------|----------|------|----------|
| [**logx**](#logx) | `v1.0.3` | 日志处理工具包 | 日志分割、格式化、多输出 |
| [**kafkax**](#kafkax) | `latest` | Kafka 消息队列工具包 | 生产者、消费者、SASL 认证 |

## 🚀 快速开始

### 📥 安装方式

#### 方式一：安装整个工具库
```bash
go get github.com/hchicken/pkg-go
```

#### 方式二：按需安装单个包（推荐）
```bash
# 安装字符串工具包
go get github.com/hchicken/pkg-go/stringx@v1.0.3

# 安装日志工具包
go get github.com/hchicken/pkg-go/logx@v1.0.3

# 安装其他包
go get github.com/hchicken/pkg-go/ginx
go get github.com/hchicken/pkg-go/gormx
```

### 🎯 基本使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/hchicken/pkg-go/stringx"
    "github.com/hchicken/pkg-go/date"
    "github.com/hchicken/pkg-go/bytex"
)

func main() {
    // 🔐 字符串加密和编码
    text := "Hello, World!"
    md5Hash := stringx.StrToMd5(text)
    base64Str := stringx.StrToBase64(text)
    uuid := stringx.UUID()
    
    fmt.Printf("原文: %s\n", text)
    fmt.Printf("MD5: %s\n", md5Hash)
    fmt.Printf("Base64: %s\n", base64Str)
    fmt.Printf("UUID: %s\n", uuid)
    
    // ⏰ 时间处理
    timestamp := int64(1609459200)
    timeStr := date.UnixToTime(timestamp)
    fmt.Printf("时间戳 %d 转换为: %s\n", timestamp, timeStr)
    
    // 📊 字节大小格式化
    fileSize := int64(1024 * 1024 * 500) // 500MB
    sizeStr := bytex.KbSize(fileSize)
    fmt.Printf("文件大小: %s\n", sizeStr)
    
    // 🎲 随机字符串生成
    randomStr := stringx.RandString(10)
    fmt.Printf("随机字符串: %s\n", randomStr)
}
```

### 🔧 版本管理

使用我们提供的 Makefile 命令来管理版本：

```bash
# 查看指定包的最新版本
make get-version
# 然后输入包名，如：stringx

# 发布新版本
make release
```

## 📚 API 文档

### 🔧 arraryx - 数组操作工具包

高效的数组操作工具，支持泛型操作。

```go
import "github.com/hchicken/pkg-go/arraryx"

// 🔍 计算两个数组的差集
slice1 := []int{1, 2, 3, 4, 5}
slice2 := []int{3, 4, 5, 6, 7}
diff := arraryx.Difference(slice1, slice2) // 结果: [1, 2]

// 📊 数组元素检查和操作
// 支持任意类型的数组操作
```

### 🔤 stringx - 字符串处理工具包

功能丰富的字符串处理工具，涵盖加密、编码、压缩等常用操作。

```go
import "github.com/hchicken/pkg-go/stringx"

// 🔐 加密和哈希
text := "hello world"
md5Hash := stringx.StrToMd5(text)           // MD5 加密
sha1Hash := stringx.Sha1("key", text)       // SHA1 HMAC
sha1Base64 := stringx.Sha1Base64("key", text) // SHA1 + Base64

// 🆔 唯一标识生成
uuid := stringx.UUID()                      // 生成 UUID v4
randomStr := stringx.RandString(16)         // 生成随机字符串

// 📝 编码和解码
base64Str := stringx.StrToBase64(text)      // Base64 编码
fileBase64, _ := stringx.FileToBase64("file.txt") // 文件转 Base64

// 🗜️ 压缩和解压
data := map[string]interface{}{"key": "value"}
compressed, _ := stringx.GzipEn(data)       // Gzip 压缩
decompressed, _ := stringx.GzipDe(compressed) // Gzip 解压

// 🔗 URL 参数构建
params := map[string]string{"name": "john", "age": "25"}
urlStr := stringx.MakeUrlStr(params)        // 生成: age=25&name=john
```

### ⏰ date - 时间日期处理工具包

灵活的时间处理工具，支持多种时间格式转换。

```go
import "github.com/hchicken/pkg-go/date"

// 🕐 时间戳转换
timestamp := int64(1609459200)
timeStr := date.UnixToTime(timestamp)       // "2021-01-01 00:00:00"

// 📅 字符串转时间戳
unixTime, _ := date.TimeToUnix("2021-01-01 00:00:00")

// 🎯 自定义格式转换
layout := "2006/01/02 15:04:05"
unixTime, _ := date.TimeToUnix1(layout, "2021/01/01 00:00:00")

// 🌍 智能时间解析（支持多种格式）
unixTime, _ := date.TimeToUnixV2("2021-01-01T00:00:00Z")
unixTime, _ := date.TimeToUnixV2("Jan 1, 2021")
```

### 📊 bytex - 字节操作工具包

字节大小格式化和转换工具。

```go
import "github.com/hchicken/pkg-go/bytex"

// 📏 文件大小格式化
size1 := int64(1024)                        // 1KB
size2 := int64(1024 * 1024)                 // 1MB
size3 := int64(1024 * 1024 * 1024)          // 1GB

fmt.Println(bytex.KbSize(size1))            // "1.00 KB"
fmt.Println(bytex.KbSize(size2))            // "1.00 MB"
fmt.Println(bytex.KbSize(size3))            // "1.00 GB"
```

### 🏗️ structx - 结构体操作工具包

结构体转换和序列化工具。

```go
import "github.com/hchicken/pkg-go/structx"

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// 🔄 结构体转换
user := User{Name: "John", Age: 25}
var userMap map[string]interface{}
err := structx.StructDecode(user, &userMap)

// ⚡ 高性能转换（使用 json-iterator）
err = structx.StructSpecialDecode(user, &userMap)
```

### 📁 filex - 文件操作工具包

文件路径获取和处理工具。

```go
import "github.com/hchicken/pkg-go/filex"

// 📍 获取当前执行文件的绝对路径
currentPath, err := filex.GetCurrentAbPath()
fmt.Println("当前路径:", currentPath)

// 兼容 go run 和编译后的可执行文件
```

### 🖥️ osx - 操作系统工具包

操作系统相关的实用工具。

```go
import "github.com/hchicken/pkg-go/osx"

// 🌍 环境变量获取（带默认值）
dbHost := osx.GetEnv("DB_HOST", "localhost")
dbPort := osx.GetEnv("DB_PORT", "5432")
```

### 🌐 ginx - Gin 框架扩展包

Gin 框架的增强工具，提供更便捷的开发体验。

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/hchicken/pkg-go/ginx/binding"
    "github.com/hchicken/pkg-go/ginx/response"
    "github.com/hchicken/pkg-go/ginx/protos"
)

// 📝 请求参数绑定
func CreateUser(c *gin.Context) {
    var req struct {
        Name string `json:"name" binding:"required"`
        Age  int    `json:"age" binding:"min=1,max=120"`
    }
    
    // JSON 参数绑定
    if err := binding.ShouldBindJSON(c, &req); err != nil {
        return // 自动返回错误响应
    }
    
    // Query 参数绑定
    var query protos.ReqQueryBase
    if err := binding.ShouldBindQuery(c, &query); err != nil {
        return
    }
    
    // 📤 统一响应格式
    response.Json(c, 
        response.Data(map[string]interface{}{
            "id": 123,
            "name": req.Name,
        }),
        response.Message("用户创建成功"),
    )
}

// 📄 分页查询基础结构
type UserQuery struct {
    protos.ReqQueryBase
    Name string `form:"name"`
}
```

### 🌍 httpx - HTTP 客户端工具包

基于 Resty 的高级 HTTP 客户端。

```go
import "github.com/hchicken/pkg-go/httpx"

// 🚀 创建 HTTP 客户端
client := httpx.NewHttpClient()

// 🔧 设置请求参数
client.SetHeaders(map[string]string{
    "Content-Type": "application/json",
    "Authorization": "Bearer token",
})

client.SetParam(map[string]string{
    "page": "1",
    "size": "10",
})

// 📨 发送请求
var result map[string]interface{}
client.SetResult(&result)

response, err := client.Get("https://api.example.com/users")
if err != nil {
    log.Fatal(err)
}

// 📊 处理响应
fmt.Printf("状态码: %d\n", response.StatusCode())
fmt.Printf("响应数据: %+v\n", result)
```

### 🔐 jwtx - JWT 认证工具包

JWT Token 生成和解析工具。

```go
import "github.com/hchicken/pkg-go/jwtx"

// 🎫 生成 JWT Token
token, err := jwtx.GenerateToken("john_doe", "user_password")
if err != nil {
    log.Fatal("Token 生成失败:", err)
}
fmt.Println("生成的 Token:", token)

// 🔍 解析 JWT Token
claims, err := jwtx.ParseToken(token)
if err != nil {
    log.Fatal("Token 解析失败:", err)
}

fmt.Printf("用户名: %s\n", claims.Username)
fmt.Printf("过期时间: %v\n", claims.ExpiresAt)
```

### 🗄️ gormx - GORM 数据库扩展包

GORM 的增强工具，提供更便捷的数据库操作。

```go
import (
    "github.com/hchicken/pkg-go/gormx"
    "gorm.io/gorm"
)

// 📋 使用基础模型
type User struct {
    ID   uint   `gorm:"primarykey"`
    Name string `gorm:"size:100;not null"`
    Email string `gorm:"uniqueIndex"`
    gormx.TabBaseModel // 自动包含创建时间、更新时间等字段
}

// 🔗 创建数据库实例
db := gormx.NewDatabase(
    gormx.WithConnPool(gormDB), // 传入 *gorm.DB 实例
)

// TabBaseModel 包含以下字段：
// - CreatedBy string    // 创建人
// - UpdatedBy string    // 更新人  
// - CreatedAt JsonTime  // 创建时间
// - UpdatedAt JsonTime  // 更新时间

// JsonTime 支持自定义时间格式序列化
```

### 📝 logx - 日志处理工具包

功能强大的日志处理工具，支持多种输出方式和格式。

```go
import "github.com/hchicken/pkg-go/logx"

// 🎛️ 创建日志记录器
logger := logx.NewLogger(
    logx.WithPath("./logs"),           // 日志目录
    logx.WithFile("app.log"),          // 日志文件名
    logx.WithLevel(logx.InfoLevel),    // 日志级别
    logx.WithWriterType(logx.LumberjackWriter), // 使用 Lumberjack 轮转
    logx.WithMaxSize(100),             // 最大文件大小 (MB)
    logx.WithMaxBackups(5),            // 保留文件数量
    logx.WithMaxAge(30),               // 保留天数
)

// 📊 使用不同级别的日志
loggerIns := logger.GetLogger("business") // 获取业务日志实例

loggerIns.Debug("调试信息")
loggerIns.Info("普通信息")
loggerIns.Warn("警告信息")
loggerIns.Error("错误信息")

// 🔄 支持日志轮转
// - RotateWriter: 按时间轮转
// - LumberjackWriter: 按大小轮转
```

### 💾 cache - 缓存操作工具包

Redis 连接池管理工具。

```go
import "github.com/hchicken/pkg-go/cache"

// 🏊 创建 Redis 连接池
pool, err := cache.NewCachePool(
    cache.WithHost("localhost:6379"),
    cache.WithPassword("your_password"),
    cache.WithDatabase(0),
    cache.WithMaxIdle(10),
    cache.WithMaxActive(100),
)
if err != nil {
    log.Fatal("Redis 连接池创建失败:", err)
}
defer pool.Close()

// 🔗 获取连接并使用
conn := pool.GetConnection()
defer conn.Close()

// 使用 Redis 命令
_, err = conn.Do("SET", "key", "value")
_, err = conn.Do("GET", "key")
```

### 📨 kafkax - Kafka 消息队列工具包

Kafka 生产者和消费者工具，支持 SASL 认证。

```go
import (
    "github.com/hchicken/pkg-go/kafkax"
    "github.com/hchicken/pkg-go/kafkax/writer"
    "github.com/hchicken/pkg-go/kafkax/reader"
)

// 🏭 创建 Kafka 客户端
client, err := kafkax.NewKafkaClient(
    kafkax.Address([]string{"localhost:9092"}),
    kafkax.User("username"),
    kafkax.Password("password"),
)
if err != nil {
    log.Fatal("Kafka 客户端创建失败:", err)
}

// 📤 创建生产者
producer := client.NewWriter(
    writer.Topic("my-topic"),
    writer.Balancer(&writer.LeastBytes{}),
)
defer producer.Close()

// 发送消息
err = producer.WriteMessages(context.Background(),
    kafka.Message{
        Key:   []byte("key1"),
        Value: []byte("message content"),
    },
)

// 📥 创建消费者
consumer := reader.NewReader(
    reader.Address([]string{"localhost:9092"}),
    reader.Topic("my-topic"),
    reader.GroupID("my-group"),
)
defer consumer.Close()

// 消费消息
for {
    msg, err := consumer.ReadMessage(context.Background())
    if err != nil {
        break
    }
    fmt.Printf("收到消息: %s\n", string(msg.Value))
}
```

## 🛠️ 开发工具

### 📦 包管理命令

项目提供了便捷的 Makefile 命令来管理包版本和发布：

```bash
# 🔍 查看指定包的最新版本
make get-version
# 交互式输入包名，如：stringx、logx 等

# 🚀 发布新版本（交互式）
make release
# 按提示选择要发布的包和版本号

# 🧹 清理构建文件
make clean

# ❓ 查看所有可用命令
make help
```

### 🎯 发布脚本详解

使用 `release.sh` 脚本进行版本发布：

```bash
# 📋 基本语法
./release.sh <包名> <版本号>

# 🔧 发布单个包
./release.sh stringx v1.0.4
./release.sh logx v1.0.4
./release.sh ginx v1.0.1

# 🌟 发布所有包（统一版本）
./release.sh all v1.0.0

# 📊 脚本特性
# ✅ 自动检查 Git 状态
# ✅ 验证版本号格式
# ✅ 创建 Git 标签
# ✅ 生成发布归档
# ✅ 彩色日志输出
```

### 🔧 开发环境设置

```bash
# 📥 克隆项目
git clone https://github.com/hchicken/pkg-go.git
cd pkg-go

# 🔍 查看项目结构
tree -L 2

# 📋 查看所有包的最新版本
git tag --sort=-version:refname | head -20

# 🧪 运行测试（如果有）
go test ./...

# 📝 代码格式化
go fmt ./...
```

### 📈 版本管理策略

- 🏷️ **语义化版本控制**: 遵循 `vX.Y.Z` 格式
- 📦 **独立版本管理**: 每个包独立发布版本
- 🔖 **Git 标签格式**: `包名/版本号`，如 `stringx/v1.0.3`
- 🚀 **发布流程**: 开发 → 测试 → 标签 → 发布

### 🎨 代码规范

- 📝 **注释**: 所有公开函数必须有注释
- 🧪 **测试**: 核心功能需要单元测试
- 📋 **文档**: API 变更需要更新文档
- 🔍 **代码检查**: 使用 `go vet` 和 `golint`

## 📋 系统要求

| 项目 | 要求 |
|------|------|
| **Go 版本** | >= 1.19 |
| **操作系统** | Linux、macOS、Windows |
| **架构支持** | amd64、arm64 |

## 🏗️ 项目架构

```
pkg-go/
├── 🔧 核心工具包/
│   ├── arraryx/     # 数组操作
│   ├── bytex/       # 字节处理
│   ├── stringx/     # 字符串工具
│   ├── structx/     # 结构体转换
│   ├── date/        # 时间处理
│   ├── filex/       # 文件操作
│   └── osx/         # 系统工具
├── 🌐 Web 扩展/
│   ├── ginx/        # Gin 框架扩展
│   ├── httpx/       # HTTP 客户端
│   └── jwtx/        # JWT 认证
├── 🗄️ 数据存储/
│   ├── gormx/       # GORM 扩展
│   └── cache/       # Redis 缓存
└── 📊 中间件/
    ├── logx/        # 日志处理
    └── kafkax/      # Kafka 消息队列
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！🎉

### 🚀 快速开始贡献

1. **🍴 Fork 项目**
   ```bash
   # 在 GitHub 上点击 Fork 按钮
   git clone https://github.com/your-username/pkg-go.git
   ```

2. **🌿 创建特性分支**
   ```bash
   git checkout -b feature/amazing-feature
   ```

3. **💻 进行开发**
   ```bash
   # 编写代码
   # 添加测试
   # 更新文档
   ```

4. **✅ 提交更改**
   ```bash
   git add .
   git commit -m "✨ Add amazing feature"
   ```

5. **📤 推送并创建 PR**
   ```bash
   git push origin feature/amazing-feature
   # 在 GitHub 上创建 Pull Request
   ```

### 📝 贡献类型

- 🐛 **Bug 修复** - 修复现有功能的问题
- ✨ **新功能** - 添加新的工具包或功能
- 📚 **文档改进** - 完善文档和示例
- 🎨 **代码优化** - 性能优化和代码重构
- 🧪 **测试增强** - 添加或改进测试用例

### 🔍 代码审查标准

- ✅ 代码风格符合 Go 规范
- ✅ 包含必要的单元测试
- ✅ 公开 API 有完整注释
- ✅ 向后兼容性考虑
- ✅ 性能影响评估

## 📊 项目统计

![GitHub stars](https://img.shields.io/github/stars/hchicken/pkg-go?style=social)
![GitHub forks](https://img.shields.io/github/forks/hchicken/pkg-go?style=social)
![GitHub issues](https://img.shields.io/github/issues/hchicken/pkg-go)
![GitHub pull requests](https://img.shields.io/github/issues-pr/hchicken/pkg-go)

## 📄 许可证

本项目采用 **MIT 许可证**。详情请参阅 [LICENSE](LICENSE) 文件。

```
MIT License - 自由使用、修改和分发
```

## 🔗 相关链接

| 链接 | 描述 |
|------|------|
| 🏠 [项目主页](https://github.com/hchicken/pkg-go) | GitHub 仓库主页 |
| 🐛 [问题反馈](https://github.com/hchicken/pkg-go/issues) | 提交 Bug 和功能请求 |
| 📦 [版本发布](https://github.com/hchicken/pkg-go/releases) | 查看发布历史 |
| 📖 [Wiki](https://github.com/hchicken/pkg-go/wiki) | 详细文档和教程 |
| 💬 [讨论区](https://github.com/hchicken/pkg-go/discussions) | 社区讨论 |

## 📞 联系方式

### 🤔 遇到问题？

1. **📖 查看文档** - 先查看本 README 和相关包的文档
2. **🔍 搜索 Issues** - 看看是否有人遇到过类似问题
3. **💬 参与讨论** - 在 Discussions 中提问
4. **🐛 提交 Issue** - 发现 Bug 请创建详细的 Issue

### 💡 有好想法？

- 💭 **功能建议** - 在 Issues 中提交功能请求
- 🚀 **直接贡献** - 提交 Pull Request
- 📧 **私下交流** - 发送邮件给项目维护者

---

<div align="center">

### 🌟 如果这个项目对你有帮助，请给它一个 Star！

[![Star History Chart](https://api.star-history.com/svg?repos=hchicken/pkg-go&type=Date)](https://star-history.com/#hchicken/pkg-go&Date)

**感谢所有贡献者的支持！** 🙏

</div>