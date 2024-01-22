package hcl

import "fmt"

type DiagType string

var (
	AttributeDiag DiagType = "attribute"
	BlockDiag     DiagType = "block"
	TfProjectDiag DiagType = "tf-project"
)

// Diags to store diagnostics
// only remain saved after the last run
type Diags struct {
	Name       string
	Type       DiagType
	Errors     []error
	ChildDiags []Diags
}

// Show shows the diags if any error exists
func (d Diags) Show() (string, bool) {
	hasError := false
	str := fmt.Sprintf("Diags for %s %s :\n", d.Type, d.Name)
	for _, err := range d.Errors {
		str = str + err.Error() + "\n"
		hasError = true
	}
	for _, childDiag := range d.ChildDiags {
		if childStr, ok := childDiag.Show(); ok {
			hasError = true
			str = str + childStr + "\n"
		}
	}
	return str, hasError
}
