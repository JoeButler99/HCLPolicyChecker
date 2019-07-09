package main

// Code lifted from terraform source
// TODO - Find better way to access the private functions in Terraform source

import (
	"fmt"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

//

// Variable represents a "variable" block in a module or file.
type HCLObject struct {
	Name        string
	Description string
	Default     cty.Value
	Type        cty.Type
	ParsingMode configs.VariableParsingMode

	DescriptionSet bool

	DeclRange hcl.Range
}

func (h *HCLObject) FromVariable(v *configs.Variable) {
	h.Name = v.Name
	h.Description = v.Description
	h.Default = v.Default
	h.Type = v.Type
	h.ParsingMode = v.ParsingMode

}

func (h *HCLObject) FromOutput(o *configs.Output) {
	t, _ := o.Expr.Value(nil)
	h.Name = o.Name
	h.Description = o.Description
	h.Default = cty.Value{}
	h.Type = t.Type()
	h.ParsingMode = configs.VariableParseHCL

}



func (h *HCLObject) FromResource(r *configs.Resource) {
	fmt.Println(r)
	fmt.Println()
	fmt.Println(r.Config)

	a, _ := r.Config.JustAttributes()
	fmt.Println(a)

	for k, v := range a {
		fmt.Println(k, v)
		if k == "tags" {
			fmt.Println("TAGS")
			fmt.Printf("%+v\n", v)
			fmt.Println()
		}
	}
	//fmt.Println(r.Config.PartialContent(nil))

	h.Name = r.Name
	h.Description = ""
	h.Default = cty.Value{}
	h.Type = cty.Type{}
	h.ParsingMode = configs.VariableParseHCL

}
