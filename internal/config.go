package internal

import (
	"flag"
	"fmt"
	"github.com/mannemsolutions/pgtester/pkg/pg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

/*
 * This module reads the config file and returns a config object with all entries from the config yaml file.
 */

const (
	defaultConfFile = "./tests.yaml"
)

type Tests []Test

type Test struct {
	Name    string     `yaml:"name"`
	Query   string     `yaml:"query"`
	Results pg.Results `yaml:"results"`
	Reverse bool       `yaml:"reverse"`
}

func (t *Test) Validate() (err error) {
	if t.Name == "" {
		t.Name = t.Query
	} else if t.Query == "" {
		// Let's hope it is a query
		t.Query = t.Name
	}
	if t.Name == "" {
		return fmt.Errorf("a defined test is missing the query and name arguments")
	}
	return nil
}

func (t *Test) IncreaseOnError() (increase int) {
	if t.Reverse {
		return 0
	}
	return 1
}

func (t *Test) IncreaseOnSuccess() (increase int) {
	return 1 - t.IncreaseOnError()
}

func (t *Test) MsgOnError() (msg string) {
	if t.Reverse {
		return "expected error"
	}
	return "unexpected error"
}

func (t *Test) MsgOnSuccess() (msg string) {
	if t.Reverse {
		return "unexpected success"
	}
	return "success as expected"
}

type Configs []Config

type Config struct {
	Path    string
	Debug   bool          `yaml:"debug"`
	Delay   time.Duration `yaml:"delay"`
	Retries uint          `yaml:"retries"`
	Tests   Tests         `yaml:"tests"`
	DSN     pg.Dsn        `yaml:"dsn"`
}

func NewConfig(path string) (c Config, err error) {
	// This only parsed as yaml, nothing else
	// #nosec
	yamlConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(yamlConfig, &c)
	if c.Delay.Nanoseconds() == 0 {
		c.Delay = time.Second
	}
	c.Path = path
	return c, nil
}

// ReadFromFileOrDir returns an array of Configs parsed from all yaml files, found while recursively walking
// through a directory, while following symlinks as needed.
func ReadFromFileOrDir(path string) (configs Configs, err error) {
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return Configs{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		return Configs{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return Configs{}, err
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if fileInfo.IsDir() {
		// file is a directory
		entries, err := file.ReadDir(0)
		if err != nil {
			_ = file.Close()
			return Configs{}, err
		}
		// I want the entries sorted, so adding them to a list of strings
		var entryNames []string
		for _, entry := range entries {
			entryNames = append(entryNames, entry.Name())
		}
		sort.Strings(entryNames)
		for _, entryName := range entryNames {
			newConfigs, err := ReadFromFileOrDir(filepath.Join(path, entryName))
			if err != nil {
				_ = file.Close()
				return Configs{}, err
			}
			configs = append(configs, newConfigs...)
		}
	} else {
		// file is not a directory
		config, err := NewConfig(path)
		if err != nil {
			_ = file.Close()
			return Configs{}, err
		}
		configs = Configs{config}
	}
	return configs, file.Close()
}

func GetConfigs() (configs Configs, err error) {
	var debug bool
	var version bool
	flag.BoolVar(&debug, "d", false, "Add debugging output")
	flag.BoolVar(&version, "v", false, "Show version information")

	flag.Parse()
	if version {
		fmt.Println(appVersion)
		os.Exit(0)
	}
	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{defaultConfFile}
	}
	for _, path := range paths {
		newConfigs, err := ReadFromFileOrDir(path)
		if err != nil {
			return Configs{}, nil
		}
		configs = append(configs, newConfigs...)
	}

	for _, config := range configs {
		config.Debug = config.Debug || debug
	}

	return configs, err
}
