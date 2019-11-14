package main

import (
	"github.com/ddosakura/sola/middleware/auth"
	"github.com/ddosakura/sola/orm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db = orm.Init("sqlite3", "test.db")

	_sign = auth.Sign(auth.AuthJWT, []byte("sola_key"))
	_auth = auth.Auth(auth.AuthJWT, []byte("sola_key"))
)
