package proxy

import (
	"fmt"

	"github.com/ddosakura/sola/v2"
	lua "github.com/yuin/gopher-lua"
)

// script(s)
const (
	ScriptBackup = `function handle()
	set_header("X-Sola-Script", "` + lua.LuaVersion + `")
	return 301, "%s"..URL
end`
	ScriptFavicon = `function handle()
	if (URL == "/favicon.ico")
	then
		set_header("X-Sola-Script", "` + lua.LuaVersion + `")
		return 301, "%s"
	end
end`
)

// Favicon Middleware
func Favicon(url string) sola.Middleware {
	return New(fmt.Sprintf(ScriptFavicon, url))
}

// Backup Middleware
func Backup(addr string) sola.Middleware {
	return New(fmt.Sprintf(ScriptBackup, addr))
}

// BackupSola App
func BackupSola(addr string) *sola.Sola {
	app := sola.New()
	app.Use(Backup(addr))
	return app
}
