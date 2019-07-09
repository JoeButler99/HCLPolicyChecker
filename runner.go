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

func RunChecks(p *PolicyConfig, loadedHcl *configs.Config) {

	// TODO - This is a static mapping of string to function
	//
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
			fmt.Printf("Checking Resource: %s\n", hclVar.Name)
			for resourceType, checkList := range p.Resources {
				if hclVar.Type == resourceType {
					for _, check := range checkList {
						var hclObj HCLObject
						hclObj.FromResource(hclVar)
						RunCheck(&check, &results, &hclObj, "Resource")
					}
				}
			}
		}
	}

	results.DisplayResults()

}

func (r *Results) AddPass(msg string) {
	r.PassCheckNumbers = append(r.PassCheckNumbers, r.TotalChecks)
	r.TotalChecks += 1
	r.PassedChecks += 1
	r.CheckText = append(r.CheckText, msg)
	if r.LogDuringRun {
		color.Green(msg)
	}
}

func (r *Results) AddFail(msg string) {
	r.FailedCheckNumbers = append(r.FailedCheckNumbers, r.TotalChecks)
	r.TotalChecks += 1
	r.FailedChecks += 1
	r.CheckText = append(r.CheckText, msg)
	if r.LogDuringRun {
		color.Red(msg)
	}
}

func (r *Results) AddError(msg string) {
	r.ErroredCheckNumbers = append(r.ErroredCheckNumbers, r.TotalChecks)
	r.TotalChecks += 1
	r.ErroredChecks += 1
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
