package main

import (
	"fmt"
	"os"
)

func QuitError(err error, msg string, exitCode int) {
	if err != nil {
		fmt.Println(err.Error())
		if msg != "" {
			fmt.Println(msg)
		}
		os.Exit(exitCode)
	}
}

func main() {

	// Load CLI Config
	conf := ParseCli()
	ValidateConfig(&conf)

	// Load Policy Config
	pol := PolicyConfig{}
	QuitError(pol.LoadPolicy(conf.PolicyName), fmt.Sprintf("Error loading policy file: %s", conf.PolicyName), 1)

	// Load HCL Files
	hclCfg := HclConfig{}
	loadedConfig := hclCfg.LoadConfig(conf.TargetDir)

	// Run the checks
	RunChecks(&pol, loadedConfig)

}
