package sql

import (
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/pkg/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	err      error
	DataBase *gorm.DB
)

func StartMysql() {
	DataBase, err = gorm.Open(mysql.Open(Dsn()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DataBase.InstanceSet("gorm:table_options", "ENGINE="+env.GetEnv("DB_ENGINE", "InnoDB"))
	sqlDB, _ := DataBase.DB()
	sqlDB.SetMaxIdleConns(10)
	err := DataBase.AutoMigrate()
	if err != nil {
		panic(err)
	}
}
func Dsn() string {
	return env.GetEnv("DB_USER", "root") + ":" + env.GetEnv("DB_PASSWORD", "password") + "@tcp(" + env.GetEnv("DB_HOST", "127.0.0.1") + ":" + env.GetEnv("DB_PORT", "3306") + ")/" + env.GetEnv("DB_NAME", "data_processor") + "?" + env.GetEnv("SQL_CONFIG", "parseTime=true")
}
func MysqlCheck() bool {
	if DataBase != nil {
		sqlDB, _ := DataBase.DB()
		err := sqlDB.Ping()
		if err != nil {
			log.Err(err).Msg("mysql health check failed")
		}
		return err == nil
	}
	return false
}

// close mysql connection
func CloseMysql() {
	if DataBase != nil {
		sqlDB, _ := DataBase.DB()
		err := sqlDB.Close()
		if err != nil {
			panic(err)
		}
	}
}
