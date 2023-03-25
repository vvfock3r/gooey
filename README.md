# Gooey

![Go Version](https://img.shields.io/badge/Go-1.20-blue)

## 介绍

Gooey是Go语言编写的一个简单的的用于快速开发命令行工具的模板

## 特性

* 提供非核心功能的解决方案，比如配置文件热加载
* 采用模块化设计，未加载的模块不会编译到二进制文件中

## 内置模块

* Config（配置文件）：
  * 支持多路径搜索配置文件
  * 支持命令行指定配置文件
  * 支持无配置文件模式运行
* Watch（配置文件监控）：
  * 监控配置文件，并对已注册的模块进行热更新
* Logger（日志）：
  * 支持console和json格式
  * 所有配置都支持热更新

## 版本要求
* Go: 1.20+

## 步骤

### 克隆代码

```bash
$ git clone https://github.com/vvfock3r/gooey.git
$ cd gooey
$ go mod tidy
```

### 设置 Git

```bash
# 设置Git Hooks
$ git config core.hooksPath .githooks

# 在每次提交前会执行.githooks目录下的钩子脚本，比如
$ git add * && git commit -m "test: git hooks" 
pre-commit
    RUN go mod tidy
    RUN gofmt -w -r "interface{} -> any" .
    RUN go vet .
[main 931a3e8] update
 1 file changed, 39 insertions(+), 232 deletions(-)
 rewrite README.md (94%)
```

### 修改配置

`load/modules.go`

```go
package load

import (
	"gooey/iface"
	"gooey/module/automaxprocs"
	"gooey/module/config"
	"gooey/module/help"
	"gooey/module/logger"
	"gooey/module/version"
	"gooey/module/watch"
)

// ModuleList 模块列表
var ModuleList = []iface.Module{
	// 独立的模块放在最上面
	&version.Version{},
	&help.Help{HiddenHelpCommand: true},

	// 配置文件相关的模块
	&config.Config{
		Name:      "etc/default",
		Exts:      []string{"yaml"},
		Path:      []string{".", "$HOME", "/etc"},
		MustExist: false,
	},
	&watch.Watch{[]iface.Module{
		&logger.Logger{AddCaller: true},
	}},

	// 依赖于配置文件的模块放到下面
	&logger.Logger{AddCaller: true},
	&automaxprocs.AutoMaxProcs{},
}
```

### 运行程序

```bash
$ go run .
```

