package main

import (
	"flag"
	"scanner/app"
	"scanner/env"
)

var envFlag = flag.String("env", "./env.toml", "env not found")

func main() {
	flag.Parse()
	e := env.NewEnv(*envFlag)
	app.NewApp(e)
}
