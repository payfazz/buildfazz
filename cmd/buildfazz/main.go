package main

import (
	"fmt"
	"github.com/payfazz/buildfazz/internal/base"
	"github.com/payfazz/buildfazz/internal/builder"
	"github.com/payfazz/buildfazz/internal/help"
	"github.com/payfazz/buildfazz/internal/pusher"
	"os"
	"strings"
)

// check is array isset or not
func isset(arr []string, index int) bool {
	return (len(arr) > index)
}

// option mapper helper
func mapOptions(args *[]string, mapper *map[string]string, val string) bool {
	idx := 0
	if isset(*args, idx+1) && (*args)[idx+1] != "" {
		(*mapper)[val] = (*args)[idx+1]
		removeStringFromArray(args, idx, 2)
		return true
	}
	return false
}

// get command options
func getBuildOption(args *[]string, mapper *map[string]string) {
	for stat := true; stat; stat = len(*args) > 1 {
		switch (*args)[0] {
		case "-p":
			if !mapOptions(args, mapper, "path") {
				fmt.Println("your path format is wrong! please use: -p [path]")
				os.Exit(1)
			}
			break
		case "-os":
			if !mapOptions(args, mapper, "os") {
				fmt.Println("your path format is wrong! please use: -os [debian/ubuntu/scratch]")
				os.Exit(1)
			}
			break
		}
	}
}

// get push options
func getPushOption(args *[]string, mapper *map[string]string) {
	for stat := true; stat; stat = len(*args) > 1 {
		switch (*args)[0] {
		case "-e":
			if !mapOptions(args, mapper, "env") {
				fmt.Println("your path format is wrong! please use: -e [mac]")
				os.Exit(1)
			}
			break
		case "-t":
			if !mapOptions(args, mapper, "target") {
				fmt.Println("your path format is wrong! please use: -t [server target]")
				os.Exit(1)
			}
			break
		case "-ssh":
			if !mapOptions(args, mapper, "ssh") {
				fmt.Println("your path format is wrong! please use: -ssh [ssh target]")
				os.Exit(1)
			}
		case "-p":
			if !mapOptions(args, mapper, "port") {
				fmt.Println("your path format is wrong! please use: -p [port]")
				os.Exit(1)
			}
			break
		}
	}
}

// splice array
func removeStringFromArray(args *[]string, i int, length int) {
	(*args)[i] = ""
	*args = append((*args)[:i], (*args)[i+length:]...)
}

// args mapper helper
func mapArgs(args *[]string, mapper *map[string]string, idx int, key string, val string) {
	(*mapper)[key] = val
	removeStringFromArray(args, idx, 1)
}

func getProjectProp(args *[]string, mapper *map[string]string) {
	project := (*args)[0]
	fmt.Println(project)
	if project == "" {
		fmt.Println("what is your docker image? put: {docker-name}:[docker-tag]")
		os.Exit(1)
	}
	temp := strings.Split(project, ":")
	(*mapper)["projectName"] = strings.ToLower(temp[0])
	(*mapper)["projectTag"] = "latest"
	if isset(temp, 1) && temp[1] != "" {
		(*mapper)["projectTag"] = strings.ToLower(temp[1])
	}
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
	removeStringFromArray(&args, 0, 1)
	// search docker command
	for k, v := range args {
		switch v {
		case "build":
			if isset(args, k+1) && args[k+1] == "--help" {
				fmt.Println(help.NewBuildHelp().GenerateHelp())
				os.Exit(0)
			}
			mapArgs(&args, &mapper, k, "type", v)
			getBuildOption(&args, &mapper)
			break
		case "push":
			if isset(args, k+1) && args[k+1] == "--help" {
				fmt.Println(help.NewPushHelp().GenerateHelp())
				os.Exit(0)
			}
			mapArgs(&args, &mapper, k, "type", v)
			getPushOption(&args, &mapper)
			break
		}
	}
	getProjectProp(&args, &mapper)
	return mapper
}

// execute command
func executeCommand(mapper map[string]string) builder.GeneratorInterface {
	switch mapper["type"] {
	case "build":
		return builder.NewBuilderGenerator(base.NewReaderConfig(mapper["pwd"]).Config, mapper)
	case "push":
		return pusher.NewPusherGenerator(mapper)
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
		fmt.Printf("error occur!\n %s", help.NewBasicHelp().GenerateHelp())
		os.Exit(1)
	}
	fmt.Println(mapper)
	bld.Start()
}
