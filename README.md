# Gooey ![GoVersionRequire](https://img.shields.io/badge/go%20require-1.20+-blue) [![PkgGoDev](https://pkg.go.dev/github.com/vvfock3r/gooey)](https://pkg.go.dev/github.com/vvfock3r/gooey) [![Go Report Card](https://goreportcard.com/report/github.com/vvfock3r/gooey)](https://goreportcard.com/report/github.com/vvfock3r/gooey)

## 介绍

Gooey是Go语言编写的一个简单的的用于快速开发命令行工具的模板

## 特性

* 提供非核心功能的解决方案，比如配置文件热加载
* 采用模块化设计，未加载的模块不会编译到二进制文件中

## 要求

* Go: 1.20+

## 内置模块

* version：添加选项 -v / --version
* help：添加选项 -h / --help
* config（配置文件）：
  * 支持多路径搜索配置文件
  * 支持命令行指定配置文件
  * 支持无配置文件模式运行
* watch（配置文件监控）：
  * 监控配置文件，并对已注册的模块进行热更新
* logger（日志）：
  * 支持console和json格式
  * 所有配置都支持热更新
* automaxprocs（uber开源的自动调整P的数量以更好的适用于容器运行）
* gops（google开源的一个用于列出和诊断当前在您的系统上运行的Go进程的命令）

## 目录结构

```bash
[root@ap-hongkong gooey]# tree -L 2
.
├── cmd
│   └── root.go
├── Dockerfile
├── etc
│   └── default.yaml
├── go.mod
├── go.sum
├── iface             # 定义模块接口
│   └── iface.go
├── LICENSE
├── load              # 加载的模块列表
│   └── modules.go
├── main.go
├── module            # 所有的内置模块
│   ├── automaxprocs
│   ├── config
│   ├── help
│   ├── logger
│   ├── version
│   └── watch
└── README.md
```

## 功能演示

<details>
    <summary>点击查看详情</summary>
    <p>

```bash
$ go run .# 1、克隆代码
$ git clone https://github.com/vvfock3r/gooey.git
$ cd gooey
$ go mod tidy

# 2、设置Git Hooks(可选)
# 在每次提交前会执行.githooks目录下的钩子脚本，比如
$ git config core.hooksPath .githooks
$ git add * && git commit -m "git hooks test"
pre-commit
    RUN go mod tidy
    RUN gofmt -w -r "interface{} -> any" .
    RUN go vet .
[main 931a3e8] update
 1 file changed, 39 insertions(+), 232 deletions(-)
 rewrite README.md (94%)
 
# 3、根据实际情况修改要加载的模块
# load/modules.go
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

# 4、测试：动态修改日志配置
$ go run .              # 运行
$ vim etc/default.yaml  # 修改log.level为error

{"level":"info","time":"2023-03-26 19:37:46","caller":"cmd/root.go:37","message":"2023-03-26 19:37:46"}
{"level":"warn","time":"2023-03-26 19:37:46","caller":"cmd/root.go:38","message":"2023-03-26 19:37:46"}
{"level":"error","time":"2023-03-26 19:37:46","caller":"cmd/root.go:39","message":"2023-03-26 19:37:46"}

{"level":"warn","time":"2023-03-26 19:37:46","caller":"watch/watch.go:48","message":"config update trigger","operation":"write","filename":"/root/gooey/etc/default.yaml"}
{"level":"error","time":"2023-03-26 19:37:47","caller":"cmd/root.go:39","message":"2023-03-26 19:37:47"}

{"level":"error","time":"2023-03-26 19:37:48","caller":"cmd/root.go:39","message":"2023-03-26 19:37:48"}

{"level":"error","time":"2023-03-26 19:37:49","caller":"cmd/root.go:39","message":"2023-03-26 19:37:49"}
```

</p>
</details>

