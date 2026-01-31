package main

import (
	"fmt"

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
	fmt.Println(a.Sources)
}
