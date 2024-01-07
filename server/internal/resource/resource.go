package resource

import (
	"github.com/kaytu-io/pennywise/server/resource"
)

// Resource represents a single cloud resource. It has a unique Address and a collection of multiple
// Component queries.
type Resource struct {
	// Address uniquely identifies this cloud Resource.
	Address string

	// Provider is the cloud provider that this Resource belongs to.
	Provider resource.ProviderName

	// Type describes the type of the Resource.
	Type string

	// Components is a list of price components that make up this Resource. If it is empty, the resource
	// is considered to be skipped.
	Components []resource.Component
}
