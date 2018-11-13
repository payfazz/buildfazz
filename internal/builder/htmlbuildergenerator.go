package builder

import (
	"bytes"
	"fmt"
	"github.com/payfazz/buildfazz/internal/base"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Generator payfazz builder generator
type HtmlGenerator struct {
	Data           base.Data
	projectName    string
	projectTag     string
	dockerfilePath string
	shPath         string
}

func (g *HtmlGenerator) generateDockerFile() {
	g.dockerfilePath = fmt.Sprintf("%s%s", g.Data.Pwd, "Dockerfile")
	if _, err := os.Stat(g.dockerfilePath); !os.IsNotExist(err) {
		os.Remove(g.dockerfilePath)
	}
	fo, _ := os.Create(g.dockerfilePath)
	defer func() {
		if err := fo.Close(); err != nil {
			log.Fatalf("can't create file %s, err : %s", g.dockerfilePath, err)
		}
	}()
	if _, err := fo.Write([]byte(htmlTmpl)); err != nil {
		log.Fatalf("error while writting file %s, err : %s", g.dockerfilePath, err)
	}
}

func (g *HtmlGenerator) generateSh() {
	var replacer = strings.NewReplacer("$projectName", g.projectName, "$projectTag", g.projectTag)
	g.shPath = fmt.Sprintf("%s%s", g.Data.Pwd, "docker.sh")
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

func (g *HtmlGenerator) clearFiles() {
	os.Remove(g.dockerfilePath)
	os.Remove(g.shPath)
}

func (g *HtmlGenerator) execSh() {
	var stdoutBuf bytes.Buffer
	var errStdout error
	cmd := exec.Command("/bin/sh", g.shPath)
	stdoutIn, _ := cmd.StdoutPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	outStr := string(stdoutBuf.Bytes())
	fmt.Printf("\nout:\n%s\n", outStr)
}

// Start start generator
func (g *HtmlGenerator) Start() {
	g.generateSh()
	g.generateDockerFile()
	g.execSh()
	defer func() {
		g.clearFiles()
		fmt.Println("build success")
		os.Exit(0)
	}()
}

// NewHtmlBuilderGenerator new builder generator objects
func NewHtmlBuilderGenerator(data base.Data, mapper map[string]string) GeneratorInterface {
	if data.Version != "" && mapper["projectTag"] == "" {
		mapper["projectTag"] = data.Version
	}
	if mapper["projectTag"] == "" {
		mapper["projectTag"] = "latest"
	}
	return &HtmlGenerator{Data: data, projectName: mapper["projectName"], projectTag: mapper["projectTag"]}
}

