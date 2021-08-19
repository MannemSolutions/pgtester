package internal

import (
	"flag"
	"fmt"
	"github.com/mannemsolutions/pgtester/pkg/pg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

/*
 * This module reads the config file and returns a config object with all entries from the config yaml file.
 */

const (
	envConfName     = "PGTESTS"
	defaultConfFile = "./tests.yaml"
)

type Tests []Test

type Test struct {
	Name   string            `yaml:"name"`
	Query  string            `yaml:"query"`
	Results pg.OneFieldResults `yaml:"results"`
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

type Config struct {
	Debug   bool	      `yaml:"debug"`
	Delay   time.Duration `yaml:"delay"`
	Retries uint          `yaml:"retries"`
	Tests   Tests         `yaml:"tests"`
	DSN     pg.Dsn        `yaml:"dsn"`
}

func (c *Config) Defaults() {
	if c.Delay.Nanoseconds() == 0 {
		c.Delay = time.Second
	}
}

func NewConfig() (config Config, err error) {
	var configFile string
	var debug bool
	flag.StringVar(&configFile, "f", os.Getenv(envConfName), "Specify file with tests")
	flag.BoolVar(&debug, "d", false, "Add debugging output")
	flag.Parse()
	 if configFile == "" {
		 configFile = defaultConfFile
	 }
	configFile, err = filepath.EvalSymlinks(configFile)
	if err != nil {
		config.Debug = config.Debug || debug
		return config, err
	}

	// This only parsed as yaml, nothing else
	// #nosec
	yamlConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(yamlConfig, &config)
	config.Defaults()
	return config, err
}
