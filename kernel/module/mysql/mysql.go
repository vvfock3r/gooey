package mysql

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/vvfock3r/gooey/kernel/module/logger"
)

var DB *sqlx.DB

var (
	defaultHostKey   = "settings.mysql.host"
	defaultHostValue = "127.0.0.1"

	defaultPortKey   = "settings.mysql.port"
	defaultPortValue = "3306"

	defaultUserKey   = "settings.mysql.user"
	defaultUserValue = "root"

	defaultPasswordKey   = "settings.mysql.password"
	defaultPasswordValue = ""

	defaultDatabaseKey   = "settings.mysql.database"
	defaultDatabaseValue = ""

	defaultCharsetKey   = "settings.mysql.charset"
	defaultCharsetValue = "utf8mb4"

	defaultCollationKey   = "settings.mysql.collation"
	defaultCollationValue = "utf8mb4_general_ci"

	defaultConntimeoutKey   = "settings.mysql.connect_timeout"
	defaultConntimeoutValue = "5s"

	defaultReadtimeoutKey   = "settings.mysql.read_timeout"
	defaultReadtimeoutValue = "30s"

	defaultWritetimeoutKey   = "settings.mysql.write_timeout"
	defaultWritetimeoutValue = "30s"

	defaultMaxAllowedPacketKey   = "settings.mysql.max_allowed_packet"
	defaultMaxAllowedPacketValue = "16MB"
)

// MySQL implement the Module interface
type MySQL struct {
	AddFlag         bool
	AllowedCommands []string
}

func (m *MySQL) Register(cmd *cobra.Command) {
	if !m.AddFlag {
		// default
		viper.SetDefault(defaultHostKey, defaultHostValue)
		viper.SetDefault(defaultPortKey, defaultPortValue)
		viper.SetDefault(defaultUserKey, defaultUserValue)
		viper.SetDefault(defaultPasswordKey, defaultPasswordValue)
		viper.SetDefault(defaultDatabaseKey, defaultDatabaseValue)
		viper.SetDefault(defaultCharsetKey, defaultCharsetValue)
		viper.SetDefault(defaultCollationKey, defaultCollationValue)
		viper.SetDefault(defaultConntimeoutKey, defaultConntimeoutValue)
		viper.SetDefault(defaultReadtimeoutKey, defaultReadtimeoutValue)
		viper.SetDefault(defaultWritetimeoutKey, defaultWritetimeoutValue)
		viper.SetDefault(defaultMaxAllowedPacketKey, defaultMaxAllowedPacketValue)
		return
	}

	// flags
	cmd.PersistentFlags().StringP("host", "h", defaultHostValue, "mysql host")
	cmd.PersistentFlags().StringP("port", "P", defaultPortValue, "mysql port")
	cmd.PersistentFlags().StringP("user", "u", defaultUserValue, "mysql user")
	cmd.PersistentFlags().StringP("password", "p", defaultPasswordValue, "mysql password")
	cmd.PersistentFlags().StringP("database", "d", defaultDatabaseValue, "mysql database")
	cmd.PersistentFlags().String("charset", defaultCharsetValue, "mysql charset")
	cmd.PersistentFlags().String("collation", defaultCollationValue, "mysql collation")
	cmd.PersistentFlags().String("connect_timeout", defaultConntimeoutValue, "mysql connect timeout")
	cmd.PersistentFlags().String("read_timeout", defaultReadtimeoutValue, "mysql read timeout")
	cmd.PersistentFlags().String("write_timeout", defaultWritetimeoutValue, "mysql write timeout")
	cmd.PersistentFlags().String("max_allowed_packet", defaultMaxAllowedPacketValue, "mysql max allowed packet")

	// bind
	err := viper.BindPFlag(defaultHostKey, cmd.PersistentFlags().Lookup("host"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultPortKey, cmd.PersistentFlags().Lookup("port"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultUserKey, cmd.PersistentFlags().Lookup("user"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultPasswordKey, cmd.PersistentFlags().Lookup("password"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultDatabaseKey, cmd.PersistentFlags().Lookup("database"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultCharsetKey, cmd.PersistentFlags().Lookup("charset"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultCollationKey, cmd.PersistentFlags().Lookup("collation"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultConntimeoutKey, cmd.PersistentFlags().Lookup("connect_timeout"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultReadtimeoutKey, cmd.PersistentFlags().Lookup("read_timeout"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultWritetimeoutKey, cmd.PersistentFlags().Lookup("write_timeout"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(defaultMaxAllowedPacketKey, cmd.PersistentFlags().Lookup("max_allowed_packet"))
	if err != nil {
		panic(err)
	}
}

func (m *MySQL) MustCheck(*cobra.Command) {}

func (m *MySQL) Initialize(cmd *cobra.Command) error {
	// allow connection to database
	if !m.allow(cmd) {
		return nil
	}

	// enable interactive password
	if strings.TrimSpace(viper.GetString(defaultPasswordKey)) == "" {
		fmt.Printf("Password: ")
		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		fmt.Println()
		viper.Set(defaultPasswordKey, password)
	}

	// long parameter
	addr := viper.GetString(defaultHostKey) + ":" + viper.GetString(defaultPortKey)
	params := map[string]string{"charset": viper.GetString(defaultCharsetKey)}

	// max_allowed_packet
	packet := strings.ToLower(viper.GetString(defaultMaxAllowedPacketKey))
	if !strings.HasSuffix(packet, "mb") {
		logger.Error("the max_allowed_packet parameter must specify the mb unit")
		os.Exit(1)
	}
	maxAllowedPacketMB, err := strconv.Atoi(strings.TrimSuffix(packet, "mb"))
	if err != nil {
		logger.Error("the max_allowed_packet parameter type conversion failed")
		os.Exit(1)
	}

	// build configuration
	mysqlConfig := mysql.Config{
		User:                 viper.GetString(defaultUserKey),
		Passwd:               viper.GetString(defaultPasswordKey),
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               viper.GetString(defaultDatabaseKey),
		Params:               params,
		Collation:            viper.GetString(defaultCollationKey),
		Loc:                  time.Local,
		ParseTime:            true,
		Timeout:              viper.GetDuration(defaultConntimeoutKey),
		ReadTimeout:          viper.GetDuration(defaultReadtimeoutKey),
		WriteTimeout:         viper.GetDuration(defaultWritetimeoutKey),
		CheckConnLiveness:    true,
		AllowNativePasswords: true,
		MaxAllowedPacket:     maxAllowedPacketMB << 20, // N MiB
	}

	// replace the Logger inside go-sql-driver/mysql
	err = mysql.SetLogger(&mysqlLogger{logger: logger.DefaultLogger})
	if err != nil {
		panic(err)
	}

	// connect to the database
	db, err := sqlx.Connect("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		logger.Error("connect database error", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("connect database success")

	// set up connection pool
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Second * 300)

	// global DB
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
