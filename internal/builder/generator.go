package builder

import (
	"github.com/payfazz/buildfazz/internal/base"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Generator payfazz builder generator
type Generator struct {
	Data           base.Data
	projectName    string
	projectTag     string
	dockerfilePath string
	shPath         string
}

func (g *Generator) getWorkingPath() string {
	var replacer = strings.NewReplacer(os.Getenv("GOPATH"), "")
	result := replacer.Replace(g.Data.Pwd)
	return result
}

func (g *Generator) getAddOn() string{
	if strings.Index(g.Data.Base, "golang") != -1 {
		return `ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 $GOPATH/bin/dep
RUN chmod +x $GOPATH/bin/dep
RUN dep init
RUN dep ensure --vendor-only
`
	}
	return ""
}

func (g *Generator) generateDockerFile() {
	var replacer = strings.NewReplacer("$base", g.Data.Base,
		"$working_directory", g.Data.WorkingDirectory, "$path", g.getWorkingPath(), "$add-on", g.getAddOn())
	g.dockerfilePath = "./Dockerfile"
	if _, err := os.Stat(g.dockerfilePath); !os.IsNotExist(err) {
		os.Remove(g.dockerfilePath)
	}
	builderScript := replacer.Replace(tmpl)
	fo, _ := os.Create(g.dockerfilePath)
	defer func() {
		if err := fo.Close(); err != nil {
			log.Fatalf("can't create file %s, err : %s", g.dockerfilePath, err)
		}
	}()
	if _, err := fo.Write([]byte(builderScript)); err != nil {
		log.Fatalf("error while writting file %s, err : %s", g.dockerfilePath, err)
	}
}

func (g *Generator) generateSh() {
	var replacer = strings.NewReplacer("$projectName", g.projectName, "$projectTag", g.projectTag)

	g.shPath = "./docker.sh"
	if _, err := os.Stat(g.dockerfilePath); !os.IsNotExist(err) {
		os.Remove(g.shPath)
	}
	builderScript := replacer.Replace(shTmpl)
	fo, _ := os.Create(g.shPath)
	defer func() {
		if err := fo.Close(); err != nil {
			log.Fatalf("can't create file %s, err : %s", g.shPath, err)
		}
	}()
	if _, err := fo.Write([]byte(builderScript)); err != nil {
		log.Fatalf("error while writting file %s, err : %s", g.dockerfilePath, err)
	}
	os.Chmod(g.shPath, 0755)
}

func (g *Generator) execSh() {
	cmd := exec.Command("/bin/sh", g.shPath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error while running SH, err : %s", err)
	}
}

func (g *Generator) clearFiles() {
	os.Remove(g.dockerfilePath)
	os.Remove(g.shPath)
}

func (g *Generator) Start() {
	g.generateDockerFile()
	g.generateSh()
	//g.execSh()
	//defer func() {
	//	g.clearFiles()
	//}()
}

// NewGenerator new generator objects
func NewGenerator(data base.Data, projectName string, projectTag string) *Generator {
	return &Generator{Data: data, projectName: projectName, projectTag: projectTag}
}
