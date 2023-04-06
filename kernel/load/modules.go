package load

import (
	"github.com/vvfock3r/gooey/kernel/iface"
	"github.com/vvfock3r/gooey/kernel/module/config"
	"github.com/vvfock3r/gooey/kernel/module/help"
	"github.com/vvfock3r/gooey/kernel/module/logger"
	"github.com/vvfock3r/gooey/kernel/module/maxprocs"
	"github.com/vvfock3r/gooey/kernel/module/mysql"
	"github.com/vvfock3r/gooey/kernel/module/version"
	"github.com/vvfock3r/gooey/kernel/module/watch"
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