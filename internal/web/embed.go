package web

import (
	"embed"
	"io/fs"
)

//go:embed all:static
var staticFS embed.FS

func StaticFS() fs.FS {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
