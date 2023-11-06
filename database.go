package main

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	var err error
	logger.Default.LogMode(logger.Info)
	db, err = gorm.Open(
		sqlite.Open(sqlite3DSN),
		&gorm.Config{
			TranslateError: true,
			Logger:         logger.Discard,
			// Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			// 	SlowThreshold:             00 * time.Millisecond,
			// 	LogLevel:                  logger.Info,
			// 	IgnoreRecordNotFoundError: false,
			// 	Colorful:                  true,
			// }),
		},
	)
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate((*MessageQueue)(nil), (*User)(nil), (*Form)(nil))
}
