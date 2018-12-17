package pusher

import (
	"fmt"
	"github.com/payfazz/buildfazz/internal/builder"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
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

// Create new docker tag
func (g *Generator) createTag() string {
	oldtag := fmt.Sprintf("%s:%s", g.projectName, g.projectTag)
	newtag := fmt.Sprintf("%s%s/%s", g.deployer, g.server, oldtag)

	cmd := exec.Command("docker", "tag", oldtag, newtag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to create docker tag: %s\n", err)
	}

	log.Printf("created docker tag %s\n", newtag)
	return newtag
}

// Mash port until it's open, otherwise return error
func waitPort(proto, address string, duration time.Duration) error {
	var (
		err  error
		conn net.Conn
	)
	timeLimit := time.Now().Add(duration)
	for {
		conn, err = net.Dial(proto, address)
		if err == nil || time.Now().After(timeLimit) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if conn != nil {
		return nil
	} else {
		return err
	}
}

// Starts local SSH tunnel in port 5000 to the registry
func (g *Generator) startTunnel() *exec.Cmd {
	tun := fmt.Sprintf("5000:%s", g.server)
	cmd := exec.Command("ssh", "-NTL", tun, g.ssh)

	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to create tunnel: %s\n", err)
	} else {
		log.Printf("starting tunnel to %s\n", g.ssh)
	}

	timeout := 10 * time.Second
	err := waitPort("tcp", "localhost:5000", timeout)
	if err != nil {
		log.Fatalf("creating tunnel timed out after %v\n", timeout)
	}

	log.Println("started tunnel")

	return cmd
}

// Stops SSH tunnel
func (g *Generator) stopTunnel(cmd *exec.Cmd) {
	if err := cmd.Process.Kill(); err != nil {
		log.Fatalf("failed to stop SSH tunnel: %s\n", err)
	}

	log.Printf("stopped tunnel to %s\n", g.ssh)
}

// Push tag to registry
func (g *Generator) pushTag(tag string) {
	cmd := exec.Command("docker", "push", tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("failed to push tag %s: %s", tag, err)
	} else {
		log.Printf("pushed tag %s", tag)
	}
}

func (g *Generator) execCommands() {
	tag := g.createTag()
	tun := g.startTunnel()
	g.pushTag(tag)
	defer g.stopTunnel(tun)

}

// Start ...
func (g *Generator) Start() {
	fmt.Printf("\n\nWARNING, DO NOT CLOSE YOUR APPLICATION!\nYOUR APPS WILL STUCK IF YOU DO THAT!\nDOCKER PUSH ON PROGRESS\n\n")
	g.execCommands()

	defer func() {
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
