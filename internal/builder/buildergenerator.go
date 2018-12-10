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
type Generator struct {
	Data           base.Data
	projectName    string
	projectTag     string
	dockerfilePath string
	shPath         string
	os             string
}

func (g *Generator) getWorkingPath() string {
	var result string
	log.Println(g.Data.Pwd)
	if strings.Index(g.Data.Base, "golang") != -1 {
		gopath := os.Getenv("GOPATH")
		if gopath != "" {
			var replacer = strings.NewReplacer(os.Getenv("GOPATH"), "")
			result = replacer.Replace(g.Data.Pwd)
			result = fmt.Sprintf("%s/%s", "WORKDIR $GOPATH", result)
		} else {
			split := strings.Split(g.Data.Pwd, "/")
			ln := len(split)
			projectName := split[ln-2]
			projectGroup := split[ln-3]
			result = fmt.Sprintf("WORKDIR github.com/%s/%s", projectGroup, projectName)
		}
	}
	return result
}

func (g *Generator) getAddOn() string {
	if strings.Index(g.Data.Base, "golang") != -1 {
		return `ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 $GOPATH/bin/dep
RUN chmod +x $GOPATH/bin/dep
RUN dep init
RUN dep ensure --vendor-only
`
	}
	return ""
}

func (g *Generator) getRunningScript(main string) string {
	if strings.Index(g.Data.Base, "golang") != -1 {
		return fmt.Sprintf(`
COPY . ./
RUN go build -o /app %s/*.go
RUN rm -rf $GOPATH/bin/dep`, main)
	}
	return ""
}

func (g *Generator) generateDockerFile() {
	var replacer = strings.NewReplacer("$base", g.Data.Base,
		"$path", g.getWorkingPath(), "$add-on", g.getAddOn(), "$run", g.getRunningScript(g.Data.Main), "$os", g.os)
	g.dockerfilePath = fmt.Sprintf("%s%s", g.Data.Pwd, "Dockerfile")
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

func (g *Generator) execSh() {
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

func (g *Generator) clearFiles() {
	os.Remove(g.dockerfilePath)
	os.Remove(g.shPath)
}

// Start start generator
func (g *Generator) Start() {
	g.generateSh()
	g.generateDockerFile()
	g.execSh()
	defer func() {
		g.clearFiles()
		fmt.Println("build success")
		os.Exit(0)
	}()
}

// NewBuilderGenerator new builder generator objects
func NewBuilderGenerator(data base.Data, mapper map[string]string) GeneratorInterface {
	if mapper["os"] == "" {
		mapper["os"] = "debian"
	}
	if data.Version != "" && mapper["projectTag"] == "" {
		var ref, _ = base.GetRef(data.Pwd)
		mapper["projectTag"] = data.Version + "-" + ref
	}
	if mapper["projectTag"] == "" {
		mapper["projectTag"] = "latest"
	}
	return &Generator{Data: data, projectName: mapper["projectName"], projectTag: mapper["projectTag"], os: mapper["os"]}
}
