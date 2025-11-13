package web

import "embed"

//go:embed dist/**/* dist/*
var assets embed.FS

func Assets() embed.FS {
	return assets
}
