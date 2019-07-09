package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/hashicorp/terraform/command"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/configs/configload"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// CLI Flags
const DefaultDataDir = ".terraform"

var Usage = "HCLPolicyChecker -target <TerraformDirectory> -policy <PolicyFile.yml>"

type CliConfiguration struct {
	TargetDir  string
	PolicyName string
	Recurse    bool
}

func ParseCli() CliConfiguration {
	cc := CliConfiguration{}

	flag.StringVar(&cc.TargetDir, "target", "", "Target directory")
	flag.StringVar(&cc.PolicyName, "policy", "", "Policy File")

	flag.Parse()

	return cc
}

// TODO - Display config function

func ValidateConfig(c *CliConfiguration) {

	if c.TargetDir == "" {
		fmt.Println(Usage)
		QuitError(errors.New("Target Directory cannot be empty."), "", 1)
	}

	if c.PolicyName == "" {
		fmt.Println(Usage)
		QuitError(errors.New("Policy File cannot be empty."), "", 1)
	}

}

// YAML Policy Loading

type Check struct {
	CheckName string      `yaml:"CheckName"`
	Details   interface{} `yaml:"Details"`
}

type PolicyConfig struct {
	Variables []Check            `yaml:"Variables"`
	Outputs   []Check            `yaml:"Outputs"`
	Resources map[string][]Check `yaml:"Resources"`
}

func (p *PolicyConfig) LoadPolicy(filename string) error {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}

	return yaml.Unmarshal(yamlFile, p)
}

// HCL / Terraform Loading

type HclConfig struct {
	command.Meta
}

func (m *HclConfig) normalizePath(path string) string {
	var err error

	// First we will make it absolute so that we have a consistent place
	// to start.
	path, err = filepath.Abs(path)
	if err != nil {
		// We'll just accept what we were given, then.
		return path
	}

	cwd, err := os.Getwd()
	if err != nil || !filepath.IsAbs(cwd) {
		return path
	}

	ret, err := filepath.Rel(cwd, path)
	if err != nil {
		return path
	}

	return ret
}

func (c *HclConfig) LoadConfig(dir string) *configs.Config {
	dir = c.normalizePath(dir)

	loader, err := configload.NewLoader(&configload.Config{
		ModulesDir: filepath.Join(c.DataDir(), "modules"),
		Services:   c.Services,
	})
	QuitError(err, "Error Loading HCL config", 1)

	cfg, cfgDiags := loader.LoadConfig(dir)

	if cfgDiags.HasErrors() {
		QuitError(cfgDiags, "Error Loading config to analyse. ", 1)
	}

	return cfg
}
