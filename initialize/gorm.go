package initialize

import (
	"database/sql"
	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"server/utils"
	"time"

	"server/global"
	modelAuthority "server/model/authority"
	modelFileM "server/model/fileM"
	modelMonitor "server/model/monitor"
	modelSysTool "server/model/sysTool"
)

type writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf 格式化打印日志
func (w *writer) Printf(message string, data ...interface{}) {
	if global.Config.Mysql.LogZap {
		global.Log.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}

func gormConfig() *gorm.Config {
	newLogger := logger.New(
		NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  logger.Warn,            // 日志级别
			IgnoreRecordNotFoundError: true,                   // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                  // 禁用彩色打印
		},
	)
	config := &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	return config
}

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() *gorm.DB {
	//先建立表 也可以手动先建立
	err := CreateDb()
	if err != nil {
		global.Log.Error("CreateDb  failed", zap.Error(err))
		os.Exit(1)
	}

	m := global.Config.Mysql
	if m.Dbname == "" {
		return nil
	}

	dsn := m.Username + ":" + m.Password + "@tcp(" + fmt.Sprintf("%s:%s", m.Host, m.Port) + ")/" + m.Dbname + "?" + m.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig()); err != nil {
		global.Log.Error("mysql连接失败", zap.Error(err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

// RegisterTables 初始化数据库表
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		// 权限
		modelAuthority.UserModel{},
		modelAuthority.RoleModel{},
		modelAuthority.MenuModel{},
		modelAuthority.ApiModel{},
		// 监控
		modelMonitor.OperationLogModel{},
		// fileM
		modelFileM.FileModel{},
		// 系统工具
		modelSysTool.CronModel{},
		gormadapter.CasbinRule{},
	)

	if err != nil {
		global.Log.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.Log.Info("register table success")
}

func mysqlEmptyDsn() string {
	host := global.Viper.Get("mysql.host")
	if host == "" {
		host = "127.0.0.1"
	}
	port := global.Viper.GetInt("mysql.port")
	if port == 0 {
		port = 3306
	}
	uName := global.Viper.Get("mysql.username")
	passWd := global.Viper.Get("mysql.password")
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/", uName, passWd, host, port)
}

func CreateDb() error {
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER "+
		"SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", global.Viper.Get("mysql.db-name"))
	dsn := mysqlEmptyDsn()
	fmt.Println(createSql, dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}
func MigrateData(db *gorm.DB) {
	exists, _ := utils.FileExists("lock.txt")
	fmt.Println(exists)
	if !exists {
		InitData(db)
		os.Create("lock.txt")
	}
}
