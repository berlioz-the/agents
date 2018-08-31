package main

import (
	"encoding/json"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"os"
)

func writeYamlConfig(config interface{}, path string) {
	y, err := yaml.Marshal(config)
	if err != nil {
		return
	}

	fmt.Printf("***** WRITING YAML to %s, CONTENTS: \n", path)
	fmt.Println(string(y))

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	f.WriteString(string(y))
	f.Sync()
}

func writeJsonConfig(config interface{}, path string) {
	y, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return
	}

	fmt.Printf("***** WRITING JSON to %s, CONTENTS: \n", path)
	fmt.Println(string(y))

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	f.WriteString(string(y))
	f.Sync()
}
