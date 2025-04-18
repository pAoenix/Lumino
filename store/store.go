package store

import (
	"Lumino/common"
	"Lumino/common/logger"
	"Lumino/model"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	pgLogger "gorm.io/gorm/logger"
	"strings"
	"time"
)

// DB -
type DB struct {
	*gorm.DB
	config string
	dirver string
}

// NewPgDB -
func NewPgDB() *DB {
	return NewPgDBWithConfig(model.PgDBName)
}

// NewPgDBWithConfig -
func NewPgDBWithConfig(pgConfigName string) *DB {
	user := viper.GetString(pgConfigName + ".user")
	passwd := viper.GetString(pgConfigName + ".passwd")
	host := viper.GetString(pgConfigName + ".host")
	port := viper.GetInt(pgConfigName + ".port")
	dbname := viper.GetString(pgConfigName + ".dbname")
	pgConnectSer := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, user, passwd, dbname, port)
	// 日志打印
	inDebug := common.Mode() == common.DebugMode
	var logLevel pgLogger.LogLevel
	if inDebug {
		logLevel = pgLogger.Info
	} else {
		logLevel = pgLogger.Warn
	}
	pgDB := newDB(model.PgDBName, pgConnectSer, logLevel)
	return pgDB
}

// newDB 创建通用的数据库连接
func newDB(driver, connectStr string, logLevel pgLogger.LogLevel) *DB {
	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{
		Logger: pgLogger.Default.LogMode(logLevel),
	})
	if err != nil {
		logger.Fatalf("gorm open error, %s, connect: %s", err, connectStr)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(viper.GetInt("postgresql.maxOpenConns")) //最大连接数
	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	if err != nil {
		logger.Fatalf("set maxOpenConns and connMaxLifetime error, %s, connect: %s", err, connectStr)
	}
	return &DB{
		DB:     db,
		dirver: driver,
		config: connectStr,
	}
}

// IsDuplicateError 唯一索引判断
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	// 检查GORM v2的通用错误
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// 检查PostgreSQL错误
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return true
	}

	// 检查MySQL/SQLite等错误特征
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "duplicate entry") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "duplicate key value") ||
		strings.Contains(errStr, "violates unique constraint")
}
