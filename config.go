package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/encoding/yaml"
	"github.com/cuigh/auxo/ext/files"
)

var loader = NewLoader()

type ParseFunc func(filename string) (data.Map, error)

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

func ParseEnv(filename string) (data.Map, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	m := data.Map{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") {
			pair := strings.SplitN(line, "=", 2)
			m.Set(pair[0], pair[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func ParseJSON(filename string) (data.Map, error) {
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

func ParseYaml(filename string) (data.Map, error) {
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
