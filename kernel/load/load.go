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

	// 独立的模块放
	&version.Version{},
	&help.Help{HiddenHelpCommand: true},
	&config.Config{
		Name:      "etc/default",
		Exts:      []string{"yaml"},
		Path:      []string{".", "$HOME", "/etc"},
		MustExist: false,
	},

	// 具有依赖关系的模块,详情可以查看模块的import部分
	&logger.Logger{AddCaller: true},
	&maxprocs.AutoMaxProcs{},
	&watch.Watch{List: []iface.Module{&logger.Logger{AddCaller: true}}},
	&mysql.MySQL{AllowedCommands: []string{"gooey"}},
}
