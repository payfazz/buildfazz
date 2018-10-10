package main

import (
	"fmt"
	"github.com/payfazz/buildfazz/internal/base"
	"github.com/payfazz/buildfazz/internal/builder"
	"github.com/payfazz/buildfazz/internal/help"
	"os"
	"strings"
)

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
		fmt.Printf("\ninsufficient arguments to call buildfazz!\n %s", help.NewBasicHelp().GenerateHelp())
		os.Exit(1)
	}
	// show help
	if args[1] == "--help" {
		fmt.Println(help.NewBasicHelp().GenerateHelp())
		os.Exit(0)
	}
	// remove index 0
	removeStringFromArray(&args, 0)
	// search docker command
	for k, v := range args {
		switch v {
		case "build":
			if args[k+1] == "--help" {
				fmt.Println(help.NewBuildHelp().GenerateHelp())
				os.Exit(0)
			}
			mapper["type"] = v
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

func executeCommand(mapper map[string]string) builder.GeneratorInterface {
	switch mapper["type"] {
	case "build":
		return builder.NewBuilderGenerator(base.NewReaderConfig(mapper["pwd"]).Config, mapper["projectName"], mapper["projectPath"])
	}
	return nil
}

func main() {
	var pwd string
	// get args from user
	mapper := argsParser(os.Args)
	if len(mapper) < 1 {
		fmt.Printf("command not found!\n%s", help.NewBasicHelp().GenerateHelp())
		os.Exit(0)
	}
	// get current path
	pwd, _ = os.Getwd()
	mapper["pwd"] = fmt.Sprintf("%s/", pwd)
	if mapper["path"] != "" {
		mapper["pwd"] = mapper["path"]
	}
	// map command
	bld := executeCommand(mapper)
	// if error occur
	if bld == nil {
		fmt.Printf("command not found!\n %s", help.NewBasicHelp().GenerateHelp())
		os.Exit(1)
	}
	bld.Start()
}
