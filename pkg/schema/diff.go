package schema

import (
	"github.com/kaytu-io/pennywise/pkg/cost"
	"github.com/shopspring/decimal"
)

type Action string

const (
	ActionCreate Action = "CREATE"
	ActionModify Action = "MODIFY"
	ActionRemove Action = "REMOVE"
)

type StateDiff struct {
	Resources map[string]ResourceDiff
	CostDiff  decimal.Decimal
}

// ResourceDiff type to show diff of a Resource
type ResourceDiff struct {
	Address     string
	Provider    ProviderName
	Type        string
	Skipped     bool
	IsSupported bool

	ComponentDiffs map[string][]ComponentDiff
	Action         Action
	CostDiff       decimal.Decimal
}

// ComponentDiff type to show diff of a Component
type ComponentDiff struct {
	Component cost.Component
	Action    Action
	CostDiff  decimal.Decimal
}
