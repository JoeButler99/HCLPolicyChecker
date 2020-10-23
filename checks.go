package main

import (
	"fmt"
	"github.com/JoeButler99/HCLPolicyChecker/lookup"
	"github.com/mitchellh/mapstructure"
	"os"
	"reflect"
	"strings"
)

type HasKeyDetails struct {
	Items []string
}

type KeyValueLengthDetails struct {
	KeyName              string
	MinLength, MaxLength int
}

var escapeNounWords = []string{"this", "public", "private"}

func CheckHasKey(e interface{}, FieldName string) bool {
	ValueIface := reflect.ValueOf(e)

	// Check if the passed interface is a pointer
	if ValueIface.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueIface = reflect.New(reflect.TypeOf(e))
	}
	// log.Println(ValueIface.Elem())
	Field := ValueIface.Elem().FieldByName(FieldName)

	if !Field.IsValid() {
		return false
	}

	if Field.String() == "" {
		return false
	}

	return true
}

func GetFieldString(v interface{}, field string) string {
	// TODO error checking below
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return fmt.Sprint(f)
}

func CheckStringLength(min, max int, inputString string) (bool, string) {
	l := len(inputString)
	if l < min {
		return false, fmt.Sprintf(" length is %d, minimum length: %d", l, min)
	}
	if l > max {
		return false, fmt.Sprintf(" length is %d, maximun length: %d", l, min)
	}

	return true, "sd"
}

// StringInSlice takes a string and a slice of strings and asserts if the string is in slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RunCheck(check *Check, results *Results, hclVar *HCLObject, typeName string) {
	// Run the checks
	switch check.CheckName {
	case "HasKey":

		details := HasKeyDetails{}
		err := mapstructure.Decode(check.Details, &details)
		if err != nil {
			results.addError(err.Error())
			break
		}

		for _, keyName := range details.Items {
			if CheckHasKey(hclVar, keyName) || CheckHasKey(hclVar, strings.Title(keyName)) {
				results.addPass(fmt.Sprintf("Found Key %s in %s %s", strings.Title(keyName), typeName, hclVar.Name))
			} else {
				results.addFail(fmt.Sprintf("Missing or Empty Key %s in %s %s", strings.Title(keyName), typeName, hclVar.Name))
			}
		}
	case "KeyValueLength":
		details := []KeyValueLengthDetails{}
		err := mapstructure.Decode(check.Details, &details)
		if err != nil {
			results.addError(err.Error())
			break
		}

		for _, keyValueLengthCheck := range details {
			kn := strings.Title(keyValueLengthCheck.KeyName) // TODO - strings.Title could be something better...
			if CheckHasKey(hclVar, kn) {
				s := GetFieldString(hclVar, kn)
				passed, msg := CheckStringLength(keyValueLengthCheck.MinLength, keyValueLengthCheck.MaxLength, s)
				if passed {
					results.addPass(fmt.Sprintf("%s %s.%s length between %d and %d chars", typeName, hclVar.Name, kn, keyValueLengthCheck.MinLength, keyValueLengthCheck.MaxLength))
				} else {
					results.addFail(fmt.Sprintf("Fail %s %s.%s %s. - %s", typeName, hclVar.Name, kn, msg, hclVar.DeclRange))
				}
			} else {
				// TODO - THis can be a duplicate of above - we might want a nicer way to store the results to de-dup
				results.addFail(fmt.Sprintf("Missing or Empty Key %s in %s %s", kn, typeName, hclVar.Name))
			}
		}
	case "LowerCaseName":
		if hclVar.Name == strings.ToLower(hclVar.Name) {
			results.addPass("Name found to be lowercase")
		} else {
			results.addFail(fmt.Sprintf("Name: %s contains uppercase characters", hclVar.Name))
		}
	case "NoHyphens":
		if hclVar.Name == strings.ReplaceAll(hclVar.Name, "-", "") {
			results.addPass("No Hyphens found in name")
		} else {
			results.addFail(fmt.Sprintf("Name: %s contains hyphens", hclVar.Name))
		}
	case "ResourceNameNoun":
	// TODO does it check every word or just check if one noun
		resourceNames := strings.Split(hclVar.Name, "_")
		failed := false
		for _, resourceName := range resourceNames {
			if StringInSlice(resourceName, escapeNounWords) {
				continue
			}
			
			content, err := lookup.GetWord(resourceName)
			if err != nil || content["tags"] == nil {
				results.addWarning("Resource name is not a noun")
				failed = true 
				break 
			}
			tags := content["tags"].([]interface{})
			noun := false
			for _, tag := range tags {
				if tag == "n" {
					noun = true
				}
			}
			if !noun {
				results.addFail("Resource name is not a noun")
				failed = true 
				break
			}
		}
		if !failed {
			results.addPass("Resource name is a noun")
		}	
	case "NoResourceTypeName":
		// https://www.terraform-best-practices.com/naming#resource-and-data-source-arguments
		resourceNames := strings.Split(hclVar.Name, "_")
		fail := false
		for _, word := range strings.Split(hclVar.Type.GoString(), "_") {
			if StringInSlice(word, resourceNames) {
				fail = true
			}
		}
		if fail {
			results.addFail(fmt.Sprintf("Resource Name: %s contains similar words to the Resource type %s", hclVar.Name, hclVar.Type.GoString()))
		} else {
			results.addPass("Resource names not similar to resource type names")
		}
	case "ResourceCountFirst":
		if hclVar.ValidCountPosition {
			results.addPass("Resource has count defined at top of statement with blank line padding")
		} else {
			results.addFail(fmt.Sprintf("Resource name %s does not have a count defined with blank line padded", hclVar.Name))
		}

	default:
		fmt.Printf("Did not recognise check: %s\n", check.CheckName)
		os.Exit(5)
	}
}
