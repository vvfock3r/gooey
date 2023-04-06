package mysql

import (
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/vvfock3r/gooey/kernel/module/logger"
)

var DB *sqlx.DB

// MySQL implement the Module interface
type MySQL struct {
	AllowedCommands []string
}

func (m *MySQL) Register(*cobra.Command) {}

func (m *MySQL) MustCheck(*cobra.Command) {}

func (m *MySQL) Initialize(cmd *cobra.Command) error {
	// 判断命令是否需要加载数据库
	if !m.allow(cmd) {
		return nil
	}

	// 提取子树
	v := viper.Sub("settings.mysql")
	if v == nil {
		return fmt.Errorf("miss settings.mysql")
	}

	// 生成配置
	v.SetDefault("charset", "utf8mb4")
	v.SetDefault("collation", "utf8mb4_general_ci")
	mysqlConfig := mysql.Config{
		User:                 v.GetString("username"),
		Passwd:               v.GetString("password"),
		Net:                  "tcp",
		Addr:                 v.GetString("host") + ":" + v.GetString("port"),
		DBName:               v.GetString("dbname"),
		Params:               map[string]string{"charset": v.GetString("charset")},
		Collation:            v.GetString("collation"),
		Loc:                  time.Local,
		ParseTime:            true,
		Timeout:              v.GetDuration("conntimeout"),
		ReadTimeout:          v.GetDuration("readtimeout"),
		WriteTimeout:         v.GetDuration("writetimeout"),
		CheckConnLiveness:    true,
		AllowNativePasswords: true,
	}

	// 替换go-sql-driver/mysql内部的Logger
	err := mysql.SetLogger(&mysqlLogger{logger: zap.L()})
	if err != nil {
		panic(err)
	}

	// 连接数据库
	db, err := sqlx.Connect("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		logger.Error("connect database error", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("connect database success")

	// 设置连接池
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Second * 300)

	// 赋值给全局DB
	DB = db

	return nil
}

func (m *MySQL) allow(cmd *cobra.Command) bool {
	for _, use := range m.AllowedCommands {
		if use == cmd.Use {
			return true
		}
	}
	return false
}

type mysqlLogger struct {
	logger *zap.Logger
}

func (l *mysqlLogger) Print(v ...any) {
	l.logger.Error(fmt.Sprint(v...))
}
