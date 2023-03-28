package load

import (
	"github.com/vvfock3r/gooey/iface"
	"github.com/vvfock3r/gooey/module/automaxprocs"
	"github.com/vvfock3r/gooey/module/config"
	"github.com/vvfock3r/gooey/module/help"
	"github.com/vvfock3r/gooey/module/logger"
	"github.com/vvfock3r/gooey/module/version"
	"github.com/vvfock3r/gooey/module/watch"
)

// ModuleList 包含所有内置模块的列表
var ModuleList = []iface.Module{
	// 默认静默的模块
	// &gops.Agent{},
	//&_gorm.GORM{
	//	CommandUseList:       []string{"gooey"},
	//	ReplaceDefaultLogger: true,
	//},

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
	&automaxprocs.AutoMaxProcs{},
}
