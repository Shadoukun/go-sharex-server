package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var config = Config{}

// Config holds configuration values.
type Config struct {
	Host    string     `yaml:"host"`
	Port    string     `yaml:"port"`
	Admin   permission `yaml:"admin"`
	Default permission `yaml:"default"`
}

type permission struct {
	Password  string   `yaml:"password"`
	MaxUpload int64    `yaml:"maxupload"`
	FileTypes []string `yaml:"filetypes"`
}

// getConf reads config values from a file.
func (c *Config) getConf() *Config {

	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Panic("Failed to read config file:	#%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("failed to Unmarshal: %v", err)
	}

	c.getEnvVars()
	return c
}

// getEnvVars overrides config values with set environment variables.
func (c *Config) getEnvVars() *Config {

	var (
		host           = os.Getenv("SHAREX_HOST")
		port           = os.Getenv("SHAREX_PORT")
		pass           = os.Getenv("SHAREX_PASS")
		passAdmin      = os.Getenv("SHAREX_PASS_ADMIN")
		maxUpload      = os.Getenv("SHAREX_MAXUPLOAD")
		maxUploadAdmin = os.Getenv("SHAREX_MAXUPLOAD_ADMIN")
		filetypes      = os.Getenv("SHAREX_FILETYPES")
		filetypesAdmin = os.Getenv("SHAREX_FILETYPES_ADMIN")
	)

	mu, _ := strconv.ParseInt(maxUpload, 10, 64)
	mua, _ := strconv.ParseInt(maxUploadAdmin, 10, 64)

	if host != "" {
		c.Host = host
	}
	if port != "" {
		c.Port = port
	}

	// I hate these just as much as you do.
	// password
	switch {
	// pass + passAdmin
	case pass != "" && passAdmin != "":
		c.Default.Password = pass
		c.Admin.Password = passAdmin
		// pass only
	case pass != "" && passAdmin == "":
		c.Default.Password = pass
		c.Admin.Password = pass
		// passAdmin only
	case pass == "" && passAdmin != "":
		c.Default.Password = passAdmin
		c.Admin.Password = passAdmin
	}

	// maxUpload
	switch {
	// maxUpload + maxUploadAdmin
	case maxUpload != "" && maxUploadAdmin != "":
		c.Default.MaxUpload = mu
		c.Admin.MaxUpload = mua
	// maxUpload only
	case maxUpload != "" && maxUploadAdmin == "":
		c.Default.MaxUpload = mu
		c.Admin.MaxUpload = mu
	// maxUploadAdmin only
	case maxUpload == "" && maxUploadAdmin != "":
		c.Default.MaxUpload = mua
		c.Admin.MaxUpload = mua
	}

	// fileTypes
	switch {
	// filetypes * filetypesAdmin
	case filetypes != "" && filetypesAdmin != "":
		c.Default.FileTypes = strings.Split(filetypes, ",")
		c.Admin.FileTypes = strings.Split(filetypesAdmin, ",")
	// filetypes only
	case filetypes != "" && filetypesAdmin == "":
		c.Default.FileTypes = strings.Split(filetypes, ",")
		c.Admin.FileTypes = strings.Split(filetypes, ",")
	// filetypesAdmin only
	case filetypes == "" && filetypesAdmin != "":
		c.Default.FileTypes = strings.Split(filetypesAdmin, ",")
		c.Admin.FileTypes = strings.Split(filetypesAdmin, ",")
	}

	return c
}
