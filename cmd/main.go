package main

import (
	"fmt"
	"github.com/payfazz/buildfazz/internal/base"
	"github.com/payfazz/buildfazz/internal/builder"
	"os"
	"strings"
)

func getHelp() string {
	return `
Usage: buildfazz COMMAND [OPTIONS] {docker-name}:[docker-tag]

Commands:
	build		Build docker image
	Options:
		-p		Set buildfazz working directory

`
}

// check is array isset or not
func isset(arr []string, index int) bool {
	return (len(arr) > index)
}

// get command options
func getOption(args []string, mapper *map[string]string) {
	for k, v := range args {
		switch v {
		case "-p":
			if isset(args, k+1) && args[k+1] != "" {
				(*mapper)["path"] = args[k+1]
				removeStringFromArray(&args, k+1)
				removeStringFromArray(&args, k)
			} else {
				fmt.Println("your path format is wrong! please use: --path [path]")
				os.Exit(1)
			}
			break
		}
	}
}

// splice array
func removeStringFromArray(args *[]string, i int) {
	(*args)[i] = ""
	*args = append((*args)[:i], (*args)[i+1:]...)
}

// parse arguments from user
func argsParser(args []string) map[string]string {
	mapper := make(map[string]string)
	// insufficient command params
	if len(args) < 2 {
		fmt.Println("insufficient arguments to call buildfazz")
		os.Exit(1)
	}
	// show help
	if args[1] == "--help" {
		fmt.Println(getHelp())
		os.Exit(0)
	}
	// remove index 0
	removeStringFromArray(&args, 0)
	// search docker command
	for k, v := range args {
		switch v {
		case "build":
			removeStringFromArray(&args, k)
			getOption(args, &mapper)
			project := args[0]
			temp := strings.Split(project, ":")
			mapper["projectName"] = strings.ToLower(temp[0])
			mapper["projectPath"] = "latest"
			if isset(temp, 1) && temp[1] != "" {
				mapper["projectPath"] = strings.ToLower(temp[1])
			}

			break
		}
	}
	return mapper
}

func main() {
	var pwd string
	mapper := argsParser(os.Args)
	pwd, _ = os.Getwd()
	if mapper["path"] != "" {
		pwd = mapper["path"]
	}
	bld := builder.NewGenerator(base.NewReaderConfig(pwd).Config, mapper["projectName"], mapper["projectPath"])
	bld.Start()
}
