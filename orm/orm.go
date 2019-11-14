package orm

import (
	"github.com/jinzhu/gorm"
)

// Init ORM
func Init(dialect string, args ...interface{}) *gorm.DB {
	db, err := gorm.Open(dialect, args...)
	if err != nil {
		panic("Failed to connect to database!")
	}
	return db
}
