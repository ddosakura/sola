package sola

import (
	"github.com/jinzhu/gorm"
)

// CacheORM - *gorm.DB
func (s *Sola) CacheORM(key string, db *gorm.DB) {
	s.orm[key] = db
}

// ORM getter
func (s *Sola) ORM(key string) *gorm.DB {
	if s.devMode {
		return s.orm[key].Debug()
	}
	return s.orm[key]
}

// DefaultORM getter
func (s *Sola) DefaultORM() *gorm.DB {
	return s.ORM("default")
}
