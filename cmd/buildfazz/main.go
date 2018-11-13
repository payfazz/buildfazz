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

// splice array
func removeStringFromArray(args *[]string, i int, length int) {
	(*args)[i] = ""
	*args = append((*args)[:i], (*args)[i+length:]...)
}

// swap array position to the end of elements
func swapArrayToEnd (args *[]string, i int, length int) {
	temp := (*args)[i]
	*args = append((*args)[(i+1):length], temp)
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
		//case "-os":
		//	if !mapOptions(args, mapper, "os") {
		//		log.Fatalf("your path format is wrong! please use: -os [debian/ubuntu/scratch]")
		//	}
		//	break
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
		default:
			swapArrayToEnd(args, 0, len(*args))
			break
		}
	}
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
	if isset(temp, 1) && temp[1] != "" {
		(*mapper)["projectTag"] = strings.ToLower(temp[1])
	}
}

func isInArray(haystack []string, needle string) bool {
	for _,v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
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
	for stat := true; stat; stat = len(args) > 1 {
		switch args[0] {
		case "build":
			if isset(args, 1) && args[1] == "--help" {
				fmt.Println(help.NewBuildHelp().GenerateHelp())
				os.Exit(0)
			}
			mapArgs(&args, &mapper, 0, "type", args[0])
			if len(args) > 0 {
				getBuildOption(&args, &mapper)
			}
			if len(args) > 0 {
				getProjectProp(&args, &mapper)
			}
			break
		case "push":
			if isset(args, 1) && args[1] == "--help" {
				fmt.Println(help.NewPushHelp().GenerateHelp())
				os.Exit(0)
			}
			mapArgs(&args, &mapper, 0, "type", args[0])
			if !isInArray(args, "--ssh") {
				fmt.Println("--ssh not found")
				fmt.Println(help.NewPushHelp().GenerateHelp())
				os.Exit(0)
			}
			if len(args) > 0 {
				getPushOption(&args, &mapper)
			}
			if len(args) > 0 {
				getProjectProp(&args, &mapper)
			}
			break
		default:
			fmt.Println(help.NewBasicHelp().GenerateHelp())
			os.Exit(0)
		}
	}
	return mapper
}

// execute command
func executeCommand(mapper map[string]string) builder.GeneratorInterface {
	tempCfg, err := base.NewReaderConfig(mapper["pwd"])
	switch mapper["type"] {
	case "build":
		if err != nil {
			log.Fatalf(err.Error())
		}
		cfg := tempCfg.Config
		if mapper["projectName"] != "" || cfg.ProjectName != "" {
			if mapper["projectName"] == "" {
				mapper["projectName"] = cfg.ProjectName
			}
			if cfg.Base == "html" {
				return builder.NewHtmlBuilderGenerator(cfg, mapper)
			}
			return builder.NewBuilderGenerator(cfg, mapper)
		}
		fmt.Println(help.NewBasicHelp().GenerateHelp())
		os.Exit(0)
	case "push":
		cfg := base.Data{}
		if tempCfg != nil {
			cfg = tempCfg.Config
		}
		if mapper["projectName"] == "" && cfg.ProjectName != ""{
			mapper["projectName"] = cfg.ProjectName
		}
		if mapper["projectTag"] == "" && cfg.Version != ""{
			mapper["projectTag"] = cfg.Version
		}
		if mapper["projectName"] == "" {
			fmt.Println("you need to define project name")
			fmt.Println(help.NewBasicHelp().GenerateHelp())
			os.Exit(0)
		}
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
	log.Println(mapper)
	bld.Start()
}
