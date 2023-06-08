package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	MysqlDriver = "mysql"
)

type Database struct {
	*gorm.DB
	c     *Config
	debug bool
}

type Option func(*Database)

func WithDebug(debug bool) Option {
	return func(l *Database) {
		l.debug = debug
	}
}

func New(c *Config, opts ...Option) *Database {
	database := &Database{
		c: c,
	}
	for _, opt := range opts {
		opt(database)
	}
	switch database.c.Driver {
	case MysqlDriver:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			database.c.Username, database.c.Password, database.c.Host, database.c.Port, database.c.DBName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("连接mysql失败：%v", err))
		}
		database.DB = db
	default:
		panic("不支持该数据库类型")
	}

	sqldb, err := database.DB.DB()
	if err != nil {
		panic(fmt.Errorf("get sql db error: %v", err))
	}
	// TODO 如何决策？
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetMaxOpenConns(10)

	level := logger.Warn
	if database.debug {
		level = logger.Info
	}
	database.DB = database.DB.Session(&gorm.Session{
		Logger: logger.New(log.New(os.Stdout, "", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  level,
		}),
	})
	return database
}
