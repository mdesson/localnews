package main

import (
	"github.com/mdesson/localnews/app"
)

const (
	SourcesFile = "sources.json"
)

func main() {
	a, err := app.NewApp(SourcesFile)
	if err != nil {
		panic(err)
	}
	a.Start()
}
