package circleci

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"runtime"
)

// Reader ...
type Reader struct {
	Config Data
}

// NewCircleCIReader ...
func NewCircleCIReader(pwd string) *Reader {
	r := Reader{}
	_, filename, _, _ := runtime.Caller(1)
	fp := path.Join(path.Dir(filename), "../.circleci/config.yml")
	filename, err := filepath.Abs(fp)
	log.Println(filename)
	if err != nil {
		panic("file not found")
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &r.Config)
	if err != nil {
		panic(err)
	}
	return &r
}
