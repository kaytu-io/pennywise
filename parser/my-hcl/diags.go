package my_hcl

import "fmt"

type DiagType string

var (
	AttributeDiag DiagType = "attribute"
	BlockDiag     DiagType = "block"
	TfProjectDiag DiagType = "tf-project"
)

type Diags struct {
	Name       string
	Type       DiagType
	Errors     []error
	ChildDiags []*Diags
}

func (d Diags) Show() {
	fmt.Println("Diags for", d.Type, d.Name, ":")
	for _, err := range d.Errors {
		fmt.Println(err.Error())
	}
	for _, childDiag := range d.ChildDiags {
		childDiag.Show()
	}
}
