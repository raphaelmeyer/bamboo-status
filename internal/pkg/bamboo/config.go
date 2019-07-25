package bamboo

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"syscall"
)

type configuration struct {
	Server     string   `json:"server"`
	Username   string   `json:"username"`
	BuildPlans []string `json:"buildPlans"`
}

func readConfig(filename string) (configuration, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return configuration{}, err
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return configuration{}, err
	}

	var config configuration
	err = json.Unmarshal([]byte(bytes), &config)
	if err != nil {
		return configuration{}, err
	}

	return config, nil
}

func getPassword(username string) (string, error) {
	fmt.Printf("Password for user %v: ", username)
	password, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println()
	if err != nil {
		return "", err
	}

	return string(password), nil
}
