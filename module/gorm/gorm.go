package gorm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	_gorm "gorm.io/gorm"
)

// TODO
// 1、gorm使用zap
// 2、默认参数调整
// 3、区分MySQL和sqlite
// 4、贼拉多

var DB *_gorm.DB

// GORM implement the Module interface
type GORM struct {
	CommandUseList []string
}

type MySQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBname   string
	Charset  string
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
		return fmt.Errorf("miss settings.mysql2")
	}

	// 创建实例
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&charset=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBname,
		c.Charset)
	DB, err = _gorm.Open(mysql.Open(dsn), &_gorm.Config{})
	return err
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
