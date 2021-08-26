package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/encoding/yaml"
	"github.com/cuigh/auxo/ext/files"
	"github.com/joho/godotenv"
)

var loader = NewLoader()

type ParseFunc func(filename string) (interface{}, error)

type Loader struct {
	parsers map[string]ParseFunc
}

func NewLoader() *Loader {
	return &Loader{
		parsers: map[string]ParseFunc{
			"env":  ParseEnv,
			"json": ParseJSON,
			"yaml": ParseYaml,
		},
	}
}

func (l *Loader) Load(dir, profile string) (string, error) {
	configs := make(map[string]string)
	configs[filepath.Join(dir, "config", "app."+profile+".json")] = "json"
	configs[filepath.Join(dir, "config", "app."+profile+".yaml")] = "yaml"
	configs[filepath.Join(dir, "config", "app."+profile+".yml")] = "yaml"
	configs[filepath.Join(dir, ".env."+profile)] = "env"

	for f, t := range configs {
		if files.NotExist(f) {
			continue
		}

		parser, ok := l.parsers[t]
		if !ok {
			return "", errors.New("not supported config type: " + t)
		}

		m, err := parser(f)
		d, err := json.Marshal(m)
		if err != nil {
			return "", err
		}
		return string(d), nil
	}
	return "", nil
}

func ParseEnv(filename string) (interface{}, error) {
	return godotenv.Read(filename)
}

func ParseJSON(filename string) (interface{}, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var m data.Map
	err = json.Unmarshal(d, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func ParseYaml(filename string) (interface{}, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var m data.Map
	err = yaml.Unmarshal(d, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
