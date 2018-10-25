package base

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"path/filepath"
)

// Reader ...
type Reader struct {
	Config Data
}

// NewReaderConfig ...
func NewReaderConfig(pwd string) (*Reader, error) {
	r := Reader{}
	fp := path.Join(pwd, "/buildfazz.yml")
	filename, err := filepath.Abs(fp)
	if err != nil {
		return nil, errors.New("file not found")
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("can't read 'buildfazz.yml', please create and fill 'buildfazz.yml'")
	}
	err = yaml.Unmarshal(yamlFile, &r.Config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed while parsing config, err : %s", err))
	}
	r.Config.Pwd = pwd
	return &r, nil
}
