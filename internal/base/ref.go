package base

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
)

// Get current git commit shortref
func GetRef(pwd string) (string, error) {
	headPath, _ := filepath.Abs(path.Join(pwd, "/.git/HEAD"))
	headBytes, err := ioutil.ReadFile(headPath)
	if err != nil {
		return "", errors.New("the project directory needs to be a git repository")
	}
	headStr := string(headBytes)

	// Check if HEAD is in a detached state
	match, _ := regexp.MatchString("^[0-9a-f]{40}", headStr)

	if match {
		headStr = headStr[32:40]
		return headStr, nil
	} else {
		headStr = headStr[5 : len(headStr)-1]
		refPath, _ := filepath.Abs(path.Join(pwd, "/.git/", headStr))
		refBytes, err := ioutil.ReadFile(refPath)
		if err != nil {
			return "", errors.New("failed getting current head")
		}
		headStr = string(refBytes)

		ref := headStr[:8]
		return ref, nil
	}
}
