package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/terraform/configs"
)

type Results struct {
	PassedChecks, FailedChecks, ErroredChecks, TotalChecks    int
	CheckText                                                 []string
	PassCheckNumbers, FailedCheckNumbers, ErroredCheckNumbers []int // Which indexes of CheckText have a failed item
	LogDuringRun                                              bool
}

func RunChecks(p *PolicyConfig, loadedHcl *configs.Config) *Results {

	// TODO - This is a static mapping of string to function
	results := Results{
		LogDuringRun: true,
	}

	// Do checks against variables
	if len(p.Variables) > 0 && len(loadedHcl.Module.Variables) > 0 {
		fmt.Println("Performing checks against type: 'variable'")

		for _, hclVar := range loadedHcl.Module.Variables {
			fmt.Printf("Checking Variable: %s\n", hclVar.Name)

			for _, check := range p.Variables {
				var hclObj HCLObject
				hclObj.FromVariable(hclVar)
				RunCheck(&check, &results, &hclObj, "Variable")
			}
		}
	}

	if len(p.Outputs) > 0 && len(loadedHcl.Module.Outputs) > 0 {
		fmt.Println("Performing checks against type: 'output'")
		for _, hclVar := range loadedHcl.Module.Outputs {
			fmt.Printf("Checking Output: %s\n", hclVar.Name)
			for _, check := range p.Outputs {
				var hclObj HCLObject
				hclObj.FromOutput(hclVar)
				RunCheck(&check, &results, &hclObj, "Output")
			}
		}
	}

	if len(p.Resources) > 0 && len(loadedHcl.Module.ManagedResources) > 0 {
		fmt.Println("Performing checks against type: 'resource'")
		for _, hclVar := range loadedHcl.Module.ManagedResources {
			fmt.Printf("Checking Resource: %s  %s\n", hclVar.Type, hclVar.Name)
			for resourceType, checkList := range p.Resources {

				if hclVar.Type == resourceType || resourceType == "*" {
					for _, check := range checkList {
						var hclObj HCLObject
						hclObj.FromResource(hclVar, loadedHcl.Module)
						RunCheck(&check, &results, &hclObj, "Resource")
					}
				}
			}
		}
	}

	if len(p.Data) > 0 && len(loadedHcl.Module.DataResources) > 0 {
		fmt.Println("Performing checks against type: 'data'")
		for _, hclData := range loadedHcl.Module.DataResources {
			fmt.Printf("Checking Data Resource: %s  %s\n", hclData.Type, hclData.Name)
			for _, check := range p.Data {
				var hclObj HCLObject
				hclObj.FromResource(hclData, loadedHcl.Module)
				RunCheck(&check, &results, &hclObj, "Data")
			}
		}
	}

	if len(p.Locals) > 0 && len(loadedHcl.Module.Locals) > 0 {
		fmt.Println("Performing checks against type: 'local'")
		for _, hclLocal := range loadedHcl.Module.Locals {
			fmt.Printf("Checking Local: %s\n", hclLocal.Name)
			for _, check := range p.Data {
				var hclObj HCLObject
				hclObj.FromLocal(hclLocal)
				RunCheck(&check, &results, &hclObj, "Local")
			}
		}
	}

	results.DisplayResults()
	return &results
}

func (r *Results) addPass(msg string) {
	r.PassCheckNumbers = append(r.PassCheckNumbers, r.TotalChecks)
	r.TotalChecks++
	r.PassedChecks++
	r.CheckText = append(r.CheckText, msg)
	if r.LogDuringRun {
		color.Green(msg)
	}
}

func (r *Results) addFail(msg string) {
	r.FailedCheckNumbers = append(r.FailedCheckNumbers, r.TotalChecks)
	r.TotalChecks++
	r.FailedChecks++
	r.CheckText = append(r.CheckText, msg)
	if r.LogDuringRun {
		color.Red(msg)
	}
}

func (r *Results) addError(msg string) {
	r.ErroredCheckNumbers = append(r.ErroredCheckNumbers, r.TotalChecks)
	r.TotalChecks++
	r.ErroredChecks++
	r.CheckText = append(r.CheckText, msg)
	if r.LogDuringRun {
		color.Magenta(msg)
	}
}

func (r *Results) DisplayResults() {
	fmt.Println("\nResults: ")
	fmt.Println("Total Passed Checks: ", r.PassedChecks)
	fmt.Println("Total Failed Checks: ", r.FailedChecks)
	fmt.Println("Total Errored Checks: ", r.ErroredChecks)
}
