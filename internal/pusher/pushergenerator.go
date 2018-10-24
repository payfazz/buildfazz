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

// Generator ...
type Generator struct {
	projectName string
	projectTag  string
	shPath      string
	deployer    string
	server      string
	ssh         string
}

func (g *Generator) generateSh() {
	var replacer = strings.NewReplacer("${deployer}", g.deployer, "${server}", g.server, "${ssh}", g.ssh)
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
	var stderrBuf bytes.Buffer
	proj := fmt.Sprintf("%s:%s", g.projectName, g.projectTag)
	cmd := exec.Command("/bin/sh", g.shPath, proj)
	stderrIn, err := cmd.StderrPipe()
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	go func() {
		io.Copy(stderr, stderrIn)
	}()
	err = cmd.Wait()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			log.Fatalf("%#v\n", string(err.Error()))
		}
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

// Start ...
func (g *Generator) Start() {
	g.generateSh()
	fmt.Printf("\n\nWARNING, DO NOT CLOSE YOUR APPLICATION!\nYOUR APPS WILL STUCK IF YOU DO THAT!\nDOCKER PUSH ON PROGRESS\n\n")
	g.execSh()

	defer func() {
		g.clearFiles()
		fmt.Println("PUSH SUCCESS\nImages ", g.projectName, ":", g.projectTag, " pushed to : ", g.ssh)
		os.Exit(0)
	}()
}

// NewPusherGenerator ...
func NewPusherGenerator(mapper map[string]string) builder.GeneratorInterface {
	if mapper["port"] == "" {
		mapper["port"] = "5000"
	}
	if mapper["target"] == "" {
		mapper["target"] = fmt.Sprintf("localhost:%s", mapper["port"])
	}
	if mapper["env"] == "mac" {
		mapper["env"] = "docker.for.mac."
	}
	if mapper["projectTag"] == "" {
		mapper["projectTag"] = "latest"
	}
	return &Generator{
		projectName: mapper["projectName"],
		projectTag:  mapper["projectTag"],
		deployer:    mapper["env"],
		server:      mapper["target"],
		ssh:         mapper["ssh"],
	}
}
