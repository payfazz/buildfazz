package base

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
)

// Reader ...
type Reader struct {
	Config Data
}

// NewReaderConfig ...
func NewReaderConfig(pwd string) *Reader {
	r := Reader{}
	fp := path.Join(pwd, "/buildfazz.yml")
	filename, err := filepath.Abs(fp)
	if err != nil {
		log.Fatalf("file not found")
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("can't read 'buildfazz.yml', please create and fill 'buildfazz.yml'")
	}
	err = yaml.Unmarshal(yamlFile, &r.Config)
	if err != nil {
		log.Fatalf("failed while parsing config, err : %s", err)
	}
	r.Config.Pwd = pwd
	return &r
}
