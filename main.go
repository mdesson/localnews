package main

import (
	"embed"

	"github.com/mdesson/localnews/app"
)

//go:embed sources.json
var sourcesFile []byte

//go:embed static/*
var staticFolder embed.FS

//go:embed templates/*
var templatesFolder embed.FS

func main() {
	a, err := app.NewApp(sourcesFile, templatesFolder)
	if err != nil {
		panic(err)
	}
	a.Start(staticFolder)
}
