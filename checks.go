package main

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
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

func CheckHasKey(e interface{}, FieldName string) bool {
	ValueIface := reflect.ValueOf(e)

	// Check if the passed interface is a pointer
	if ValueIface.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueIface = reflect.New(reflect.TypeOf(e))
	}

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
	}
}
