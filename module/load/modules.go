package load

import (
	"github.com/vvfock3r/gooey/module/iface"
	"github.com/vvfock3r/gooey/module/item/config"
	"github.com/vvfock3r/gooey/module/item/help"
	"github.com/vvfock3r/gooey/module/item/logger"
	"github.com/vvfock3r/gooey/module/item/maxprocs"
	"github.com/vvfock3r/gooey/module/item/mysql"
	"github.com/vvfock3r/gooey/module/item/version"
	"github.com/vvfock3r/gooey/module/item/watch"
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
	&watch.Watch{List: []iface.Module{
		&logger.Logger{AddCaller: true},
	}},

	// 依赖于配置文件的模块放到下面
	&logger.Logger{AddCaller: true},
	&maxprocs.AutoMaxProcs{},
	&mysql.MySQL{
		AllowedCommands: []string{"gooey"},
	},
}