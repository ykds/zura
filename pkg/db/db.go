package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)


const (
	Mysql = "mysql"
)

type Config struct {
	Driver   string `json:"driver" yaml:"driver"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	DBName   string `json:"db_name" yaml:"db_name"`
}

func DefaultConfig() *Config {
	return &Config{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "admin",
		Password: "123456",
		DBName:   "zira",
	}
}

type Option func(database *Database)

func WithConfig(c *Config) Option {
	return func(database *Database) {
		database.c = c
	}
}

type Database struct {
	*gorm.DB
	c *Config
}

func New(opt ...Option) *Database {
	c := DefaultConfig()
	database := &Database {
		c: c,
	}
	for _, o := range opt {
		o(database)
	}
	switch database.c.Driver {
	case Mysql:
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
		panic(err)
	}
	// TODO 如何决策？
	sqldb.SetConnMaxLifetime(10 * time.Minute)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetMaxOpenConns(50)
	return database
}