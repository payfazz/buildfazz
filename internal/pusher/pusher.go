package pusher

import (
	"bytes"
	"fmt"
	"github.com/payfazz/buildfazz/internal/builder"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Generator struct {
	projectName string
	projectTag  string
	shPath      string
	deployer    string
}

func (g *Generator) generateSh() {
	var dep = ""
	if g.deployer == "mac" {
		dep = "docker.for.mac."
	}
	var replacer = strings.NewReplacer("${deployer}", dep)
	g.shPath = "pusher.sh"
	if _, err := os.Stat(g.shPath); !os.IsNotExist(err) {
		os.Remove(g.shPath)
	}
	fo, _ := os.Create(g.shPath)
	defer func() {
		if err := fo.Close(); err != nil {
			log.Fatalf("can't create file %s, err : %s", g.shPath, err)
		}
	}()
	builderScript := replacer.Replace(template)
	if _, err := fo.Write([]byte(builderScript)); err != nil {
		log.Fatalf("error while writting file %s, err : %s", g.shPath, err)
	}
	os.Chmod(g.shPath, 0755)
}

func (g *Generator) clearFiles() {
	os.Remove(g.shPath)
}

func (g *Generator) execSh() {
	var stdoutBuf bytes.Buffer
	var errStdout error
	proj := fmt.Sprintf("%s:%s", g.projectName, g.projectTag)
	cmd := exec.Command("/bin/sh", g.shPath, proj)
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

func (g *Generator) Start() {
	g.generateSh()
	g.execSh()

	defer func() {
		g.clearFiles()
	}()
}

func NewPusherGenerator(projectName string, projectTag string, deployer string) builder.GeneratorInterface {
	return &Generator{
		projectName: projectName,
		projectTag:  projectTag,
		deployer:    deployer,
	}
}
