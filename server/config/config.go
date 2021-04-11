package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/jsonq"
)

type config struct {
	jq *jsonq.JsonQuery
}

var c *config = &config{}

func GetConfig() *jsonq.JsonQuery {
	getInstance()

	return c.jq
}

func HasKey(key ...string) bool {
	getInstance()

	_, err := c.jq.Interface(key...)
	if err != nil {
		return false
	}
	return true
}

func getInstance() {
	if c.jq == nil {
		load()
	}
}

func load() *jsonq.JsonQuery {
	// Get current directory path
	wkDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Get config file path
	filename := "development.json"
	if len(os.Getenv("ENV")) != 0 {
		filename = strings.ToLower(os.Getenv("ENV")) + ".json"
	}

	configPath := filepath.Join(wkDir, "../environment/", filename)

	// Read environment json
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	configJson := map[string]interface{}{}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.Decode(&configJson)
	jq := jsonq.NewQuery(configJson)
	c.jq = jq

	return jq
	//fmt.Println(data)
}
