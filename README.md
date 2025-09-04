# pkg-go

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hchicken/pkg-go)](https://goreportcard.com/report/github.com/hchicken/pkg-go)

一个功能丰富的 Go 语言基础工具库集合，提供了常用的工具包和扩展功能，帮助开发者快速构建 Go 应用程序。

## 📦 包列表

### 🔧 核心工具包

| 包名 | 版本 | 描述 |
|------|------|------|
| [arraryx](#arraryx) | - | 数组操作工具包 |
| [bytex](#bytex) | - | 字节操作工具包 |
| [stringx](#stringx) | v1.0.3 | 字符串处理工具包 |
| [structx](#structx) | v1.0.1 | 结构体操作工具包 |
| [date](#date) | - | 时间日期处理工具包 |
| [filex](#filex) | - | 文件操作工具包 |
| [osx](#osx) | v1.0.1 | 操作系统相关工具包 |

### 🌐 Web 框架扩展

| 包名 | 版本 | 描述 |
|------|------|------|
| [ginx](#ginx) | - | Gin 框架扩展包 |
| [httpx](#httpx) | - | HTTP 客户端工具包 |
| [jwtx](#jwtx) | - | JWT 认证工具包 |

### 🗄️ 数据库与缓存

| 包名 | 版本 | 描述 |
|------|------|------|
| [gormx](#gormx) | - | GORM 数据库扩展包 |
| [cache](#cache) | - | 缓存操作工具包 |

### 📊 中间件与消息队列

| 包名 | 版本 | 描述 |
|------|------|------|
| [logx](#logx) | v1.0.3 | 日志处理工具包 |
| [kafkax](#kafkax) | - | Kafka 消息队列工具包 |

## 🚀 快速开始

### 安装

```bash
go get github.com/hchicken/pkg-go
```

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/hchicken/pkg-go/stringx"
    "github.com/hchicken/pkg-go/date"
)

func main() {
    // 字符串工具
    uuid := stringx.UUID()
    fmt.Println("UUID:", uuid)
    
    // 时间工具
    timestamp := date.UnixToTime(1609459200)
    fmt.Println("Time:", timestamp)
}
```

## 📚 详细文档

### arraryx

数组操作工具包，提供数组去重、元素删除等功能。

```go
import "github.com/hchicken/pkg-go/arraryx"

// 数组去重
slice := []int{1, 2, 2, 3, 3, 4}
unique := arraryx.RemoveDuplicates(slice)

// 删除指定元素
result := arraryx.RemoveElement(slice, 2)
```

### stringx

字符串处理工具包，提供 MD5 加密、UUID 生成、Base64 编码、Gzip 压缩等功能。

```go
import "github.com/hchicken/pkg-go/stringx"

// MD5 加密
hash := stringx.StrToMd5("hello world")

// 生成 UUID
uuid := stringx.UUID()

// Base64 编码
encoded := stringx.StrToBase64("hello")

// 随机字符串
random := stringx.RandString(10)

// Gzip 压缩
compressed, err := stringx.GzipEn(data)
```

### date

时间日期处理工具包，提供时间戳转换、时间格式化等功能。

```go
import "github.com/hchicken/pkg-go/date"

// 时间戳转时间字符串
timeStr := date.UnixToTime(1609459200)

// 时间字符串转时间戳
timestamp, err := date.TimeToUnix("2021-01-01 00:00:00")

// 兼容多种时间格式
timestamp, err := date.TimeToUnixV2("2021-01-01T00:00:00Z")
```

### ginx

Gin 框架扩展包，提供参数绑定、响应处理、验证器等功能。

```go
import (
    "github.com/hchicken/pkg-go/ginx/binding"
    "github.com/hchicken/pkg-go/ginx/response"
    "github.com/hchicken/pkg-go/ginx/validator"
)

// 参数绑定
func handler(c *gin.Context) {
    var req RequestStruct
    if err := binding.ShouldBindJSON(c, &req); err != nil {
        return
    }
    
    // 响应处理
    response.Json(c, response.Success(data))
}

// 初始化验证器
validator.InitValidator()
```

### gormx

GORM 数据库扩展包，提供数据库连接池、模型定义、查询构建等功能。

```go
import "github.com/hchicken/pkg-go/gormx"

// 创建数据库实例
db := gormx.NewDatabase(
    gormx.WithConnPool(gormDB),
    gormx.WithConnDbModel(&User{}),
)

// 使用基础模型
type User struct {
    ID   uint   `gorm:"primarykey"`
    Name string `gorm:"size:100"`
    gormx.TabBaseModel // 包含创建时间、更新时间等字段
}
```

### logx

日志处理工具包，支持日志分割、格式化、多种输出方式。

```go
import "github.com/hchicken/pkg-go/logx"

// 创建日志记录器
logger := logx.NewLogger(
    logx.WithLevel("info"),
    logx.WithOutput("file"),
    logx.WithFilename("app.log"),
)

// 使用日志
logger.Info("这是一条信息日志")
logger.Error("这是一条错误日志")
```

### cache

缓存操作工具包，提供 Redis 连接池管理。

```go
import "github.com/hchicken/pkg-go/cache"

// 创建缓存池
pool, err := cache.NewCachePool(
    cache.WithHost("localhost:6379"),
    cache.WithPassword("password"),
)

// 获取连接
conn := pool.GetConnection()
defer conn.Close()
```

### jwtx

JWT 认证工具包，提供 Token 生成和解析功能。

```go
import "github.com/hchicken/pkg-go/jwtx"

// 生成 Token
token, err := jwtx.GenerateToken("username", "password")

// 解析 Token
claims, err := jwtx.ParseToken(token)
```

### httpx

HTTP 客户端工具包，基于 Resty 提供便捷的 HTTP 请求功能。

```go
import "github.com/hchicken/pkg-go/httpx"

// HTTP 请求
response, err := httpx.Get("https://api.example.com/data")
```

## 🛠️ 开发工具

### 构建和发布

项目提供了便捷的构建和发布工具：

```bash
# 构建项目
make build

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make vet

# 发布单个包
make release

# 清理构建文件
make clean
```

### 发布脚本

使用 `release.sh` 脚本可以方便地发布包版本：

```bash
# 发布单个包
./release.sh stringx v1.0.4

# 发布所有包
./release.sh all v1.0.0
```

## 📋 系统要求

- Go 1.19 或更高版本
- 支持的操作系统：Linux、macOS、Windows

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

### 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 🔗 相关链接

- [项目主页](https://github.com/hchicken/pkg-go)
- [问题反馈](https://github.com/hchicken/pkg-go/issues)
- [版本发布](https://github.com/hchicken/pkg-go/releases)

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 [Issue](https://github.com/hchicken/pkg-go/issues)
- 发送邮件至项目维护者

---

⭐ 如果这个项目对你有帮助，请给它一个 Star！