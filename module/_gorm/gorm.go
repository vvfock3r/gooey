package _gorm

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"

	"github.com/vvfock3r/gooey/module/logger"
)

// TODO
// 1、gorm使用zap
// 2、默认参数调整
// 3、区分MySQL和sqlite
// 4、贼拉多
// 该模块并没有完善，暂时不建议使用

var DB *gorm.DB

// GORM implement the Module interface
type GORM struct {
	CommandUseList       []string
	ReplaceDefaultLogger bool
}

type MySQLConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DBname       string
	Charset      string
	Conntimeout  string
	Readtimeout  string
	WriteTimeout string
}

func (o *GORM) Register(*cobra.Command) {}

func (o *GORM) MustCheck(*cobra.Command) {}

func (o *GORM) Initialize(cmd *cobra.Command) error {
	// 判断命令是否需要加载数据库
	if !o.inWhiteList(cmd) {
		return nil
	}

	// 判断配置项是否存在
	v := viper.Sub("settings.mysql")
	if v == nil {
		return fmt.Errorf("miss settings.mysql")
	}

	// 解析到结构体
	c := MySQLConfig{}
	err := v.Unmarshal(&c)
	if err != nil {
		return fmt.Errorf("unmarshal error: " + "settings.mysql")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&charset=%s&timeout=%s&readTimeout=%s&writeTimeout=%s",
		c.Username, c.Password, c.Host, c.Port, c.DBname, c.Charset, c.Conntimeout, c.Readtimeout, c.WriteTimeout)

	// 替换默认的Logger
	gormConfig := gorm.Config{}
	if o.ReplaceDefaultLogger {
		newLogger := zapgorm2.New(zap.L())
		newLogger.SetAsDefault()
		gormConfig.Logger = newLogger
	}

	// 创建实例
	DB, err = gorm.Open(mysql.Open(dsn), &gormConfig)

	// 错误会在gorm内部输出
	if err != nil {
		os.Exit(1)
	}
	logger.Info("connect database success")

	return nil
}

// inWhiteList 并非所有的命令/子命令都需要连接数据库,这里判断命令是否需要连接数据库
func (o *GORM) inWhiteList(cmd *cobra.Command) bool {
	for _, use := range o.CommandUseList {
		if use == cmd.Use {
			return true
		}
	}
	return false
}
