package main

import (
	"github.com/payfazz/buildfazz/internal/builder"
	"github.com/payfazz/buildfazz/internal/circleci"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		panic("send argument when you build please")
	}
	param := args[1]
	pwd, _ := os.Getwd()
	builder := builder.NewGenerator(circleci.NewCircleCIReader(pwd).Config, param)
	builder.Generate()
}
