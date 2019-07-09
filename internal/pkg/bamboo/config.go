package bamboo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", fmt.Errorf("Missing input")
	}
	return scanner.Text(), nil
}
