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

// ModuleList 包含所有内置模块的列表
var ModuleList = []iface.Module{
	// 默认静默的模块
	// &gops.Agent{},

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
