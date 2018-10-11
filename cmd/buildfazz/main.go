package main

import (
	"fmt"
	"github.com/payfazz/buildfazz/internal/base"
	"github.com/payfazz/buildfazz/internal/builder"
	"github.com/payfazz/buildfazz/internal/help"
	"github.com/payfazz/buildfazz/internal/pusher"
	"log"
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
				log.Fatalf("your path format is wrong! please use: -p [path]")
			}
			break
		case "-os":
			if !mapOptions(args, mapper, "os") {
				log.Fatalf("your path format is wrong! please use: -os [debian/ubuntu/scratch]")
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
				log.Fatalf("your path format is wrong! please use: -e [mac]")
			}
			break
		case "-t":
			if !mapOptions(args, mapper, "target") {
				log.Fatalf("your path format is wrong! please use: -t [server target]")
			}
			break
		case "--ssh":
			if !mapOptions(args, mapper, "ssh") {
				log.Fatalf("your path format is wrong! please use: --ssh [ssh target]")
			}
		case "-p":
			if !mapOptions(args, mapper, "port") {
				log.Fatalf("your path format is wrong! please use: -p [port]")
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
	if project == "" {
		log.Fatalf("what is your docker image? put: {docker-name}:[docker-tag]")
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
		log.Fatalf("\ninsufficient arguments to call buildfazz!\n %s", help.NewBasicHelp().GenerateHelp())
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
		log.Fatalf("command not found!\n%s", help.NewBasicHelp().GenerateHelp())
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
		log.Fatalf("error occur!\n %s", help.NewBasicHelp().GenerateHelp())
	}
	bld.Start()
}
