package main

// Code lifted from terraform source
// TODO - Find better way to access the private functions in Terraform source

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/lang/funcs"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"log"
)

type HCLObject struct {
	Name        string
	Description string
	Default     cty.Value
	Type        cty.Type
	ParsingMode configs.VariableParsingMode

	DescriptionSet bool

	DeclRange hcl.Range

	Count              hcl.Expression
	ValidCountPosition bool
}

func (h *HCLObject) FromVariable(v *configs.Variable) {
	h.Name = v.Name
	h.Description = v.Description
	h.Default = v.Default
	h.Type = v.Type
	h.ParsingMode = v.ParsingMode
}

func (h *HCLObject) FromLocal(v *configs.Local) {
	t, _ := v.Expr.Value(nil)
	h.Name = v.Name
	h.Description = ""
	h.Default = cty.Value{}
	h.Type = t.Type()
	h.ParsingMode = configs.VariableParseHCL
}

func (h *HCLObject) FromOutput(o *configs.Output) {
	t, _ := o.Expr.Value(nil)
	h.Name = o.Name
	h.Description = o.Description
	h.Default = cty.Value{}
	h.Type = t.Type()
	h.ParsingMode = configs.VariableParseHCL
}

func (h *HCLObject) FromResource(r *configs.Resource, module *configs.Module) {

	a, _ := r.Config.JustAttributes()
	countLine := r.Count.Range().Start.Line
	validCount := true
	for k, v := range a {
		// TODO check name range is at least -- 2 lines after count
		if v.NameRange.Start.Line < countLine+2 {
			validCount = false
		}
		// make sure count is passed
		if k == "tags" {
			fmt.Println("NEW")
			//vars := make(map[string]cty.Value)

			for a, b := range v.Expr.Variables() {
				log.Println(a, b.RootName())
			}
			//fmt.Println(vars)

			evalContext := &hcl.EvalContext{
				Variables: map[string]cty.Value{},
				//Variables: variables
				Functions: map[string]function.Function{
					"merge": funcs.MergeFunc,
				},
			}

			v.Expr.Value(evalContext)
			//fmt.Println(val.AsValueSet())
			//fmt.Println(val.AsValueMap())
		}
	}
	//fmt.Println(r.Config.PartialContent(nil))

	h.Name = r.Name
	h.Description = ""
	h.Default = cty.Value{}
	h.Type = cty.Type{}
	h.ParsingMode = configs.VariableParseHCL
	h.Count = r.Count
	h.ValidCountPosition = validCount
}
