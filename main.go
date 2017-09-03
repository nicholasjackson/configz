package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/hcl"
	input "github.com/tcnksm/go-input"
)

func main() {

	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	data, err := ioutil.ReadFile("./config.hcl")
	if err != nil {
		log.Fatalln(err)
	}

	file, err := hcl.Parse(string(data))
	if err != nil {
		log.Fatalf("unable to read config %v", err)
	}

	var config map[string]interface{}
	hcl.DecodeObject(&config, file.Node)

	sections := config["section"].([]map[string]interface{})
	for _, section := range sections {
		for key, value := range section {
			switch key {
			case "title":
				doTitle(value.([]map[string]interface{})[0])
			case "input":
				doInput(value.([]map[string]interface{})[0], ui)
			case "command":
				doCommand(value.([]map[string]interface{})[0])
			}
		}
	}
}

func doInput(section map[string]interface{}, ui *input.UI) {
	query := section["question"].(string)

	mask := false
	if section["mask"] != nil {
		mask = section["mask"].(bool)
	}

	value, err := ui.Ask(query, &input.Options{
		Required:    true,
		Mask:        mask,
		MaskDefault: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv(section["env_var"].(string), value)
}

func doTitle(section map[string]interface{}) {
	fmt.Println(section["value"])
}

func doCommand(section map[string]interface{}) {
	log.Println(section["title"].(string))

	stringArgs := make([]string, 0)

	if section["args"] != nil {
		args := section["args"].([]interface{})
		for _, arg := range args {
			stringArgs = append(stringArgs, arg.(string))
		}
	}

	cmd := exec.Command(section["command"].(string), stringArgs...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
