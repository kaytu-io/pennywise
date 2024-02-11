package diff

import (
	"github.com/kaytu-io/pennywise/pkg/cost"
	"github.com/kaytu-io/pennywise/pkg/schema"
)

type Action string

const (
	ActionCreate Action = "CREATE"
	ActionChange Action = "CHANGE"
	ActionRemove Action = "REMOVE"
)

// ResourceDiff type to show diff of a Resource
type ResourceDiff struct {
	Address     string
	Provider    schema.ProviderName
	Type        string
	Skipped     bool
	IsSupported bool

	ComponentDiffs []ComponentDiff
	Action         Action
}

// ComponentDiff type to show diff of a Component
type ComponentDiff struct {
	Component cost.Component
	Action    Action
}
