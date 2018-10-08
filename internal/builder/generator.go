package builder

import (
	"fmt"
	templategenerator "github.com/payfazz/buildfazz/internal/builder/template"
	"github.com/payfazz/buildfazz/internal/circleci"
	"io/ioutil"
	"os"
	"strings"
)

// Generator payfazz builder generator
type Generator struct {
	Data        circleci.Data
	projectName string
}

func (g *Generator) generateFolders(folders []string) {
	for _, folder := range folders {
		path := fmt.Sprintf("./%s", folder)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, 0755); err != nil {
				panic(fmt.Sprintf("can't create ./%s folder", path))
			}
		}
	}
}

func (g *Generator) generateScript(filename string, filepath string, additionalTemplate *string) {
	var temp string
	if additionalTemplate == nil {
		temp = ""
	} else {
		temp = *additionalTemplate
	}
	var replacer = strings.NewReplacer("@project-name$", g.projectName)
	if _, err := os.Stat(fmt.Sprintf(`./docker/builder-bundle%s%s`, filepath, filename)); os.IsNotExist(err) {
		buildersh, err := ioutil.ReadFile(fmt.Sprintf("./internal/builder/template%s%s%s", filepath, filename, temp))
		if err != nil {
			panic(fmt.Sprintf("can't read %s%s%s file", filepath, filename, temp))
		}
		builderScript := replacer.Replace(string(buildersh))
		fo, _ := os.Create(fmt.Sprintf("./docker/builder-bundle%s%s", filepath, filename))
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
		if _, err := fo.Write([]byte(builderScript)); err != nil {
			panic(err)
		}
	}
}

func (g *Generator) generateScriptFromString(filename string, filepath string, additionalTemplate *string) {
	var replacer = strings.NewReplacer("@project-name$", g.projectName)
	var path = fmt.Sprintf(`./docker/builder-bundle%s%s`, filepath, filename)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		os.Remove(path)
	}
	builderScript := replacer.Replace(templategenerator.BuilderScript[filename])
	fo, _ := os.Create(fmt.Sprintf("./docker/builder-bundle%s%s", filepath, filename))
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := fo.Write([]byte(builderScript)); err != nil {
		panic(err)
	}
}

func (g *Generator) copyFile() {
	g.generateScriptFromString("build", "/", nil)
	g.generateScriptFromString("builder", "/", nil)
	g.generateScriptFromString("sync-output", "/", nil)
	g.generateScriptFromString("sync-src", "/", nil)
	g.generateScriptFromString("docker-alpine-builder", "/dockerfiles/", nil)

	for _, v := range g.Data.Jobs.Build.Dockers {
		// go languange
		if strings.Index(v.Image, "go") != -1 {
			lang := "-go"
			g.generateScriptFromString("docker-builder-go", "/dockerfiles/", &lang)
			g.generateScriptFromString("scripts-build-go", "/scripts/", &lang)
			g.generateScriptFromString("scripts-install-dep-go", "/scripts/", &lang)
			g.generateScriptFromString("scripts-start-app-go", "/scripts/", &lang)
		}
	}
}

// Generate builder from template
func (g *Generator) Generate() {
	folders := []string{"docker", "docker/builder-bundle", "docker/builder-bundle/dockerfiles",
		"docker/builder-bundle/output", "docker/builder-bundle/scripts", "docker"}
	g.generateFolders(folders)
	g.copyFile()
}

// NewGenerator new generator objects
func NewGenerator(data circleci.Data, projectName string) *Generator {
	return &Generator{Data: data, projectName: projectName}
}
