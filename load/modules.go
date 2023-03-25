package load

import (
	"gooey/iface"
	"gooey/module"
	"gooey/module/logger"
)

// ModuleList 模块列表
var ModuleList = []iface.Module{
	// 独立的模块放在最上面
	&module.Version{},
	&module.Help{HiddenHelpCommand: true},

	// 配置文件相关的模块
	&module.Config{
		Name:      "etc/default",
		Exts:      []string{"yaml"},
		Path:      []string{".", "$HOME", "/etc"},
		MustExist: false,
	},
	&module.Watch{List: []iface.Module{&logger.Logger{}}},

	// 依赖于配置文件的模块放到下面
	&logger.Logger{AddCaller: true},
}
