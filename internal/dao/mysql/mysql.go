package mysql

import (
	"fmt"
	"go_webapp/pkg/settings"
	"log"
	"os"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

// Init 初始化MySQL连接
func Init(cfg *settings.MySqlSettingS) (*gorm.DB, error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.UserName, cfg.Password,
		cfg.Host, cfg.DBName)
	//设置全局 logger 打印 SQL
	myLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 的阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: myLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: cfg.TablePrefix,
		},
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}
