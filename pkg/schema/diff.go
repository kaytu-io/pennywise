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
	PriorCost decimal.Decimal
	NewCost   decimal.Decimal
}

type ModularStateDiff struct {
	Resources    map[string]ResourceDiff
	ChildModules map[string]ModularStateDiff

	PriorCost decimal.Decimal
	NewCost   decimal.Decimal
	Action    Action
}

// ResourceDiff type to show diff of a Resource
type ResourceDiff struct {
	Address     string
	Provider    ProviderName
	Type        string
	Skipped     bool
	IsSupported bool

	ComponentDiffs map[string][]ComponentDiff
	PriorCost      decimal.Decimal
	NewCost        decimal.Decimal
	Action         Action
}

// ComponentDiff type to show diff of a Component
type ComponentDiff struct {
	Component cost.Component

	Current   *cost.Component
	CompareTo *cost.Component

	Action   Action
	CostDiff decimal.Decimal
}
