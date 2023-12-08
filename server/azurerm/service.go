package azurerm

// Service is the type defining the services
type Service uint8

var (
	// The list of all services is https://azure.microsoft.com/en-us/services/, the left side is
	// the Family and the main content is the Services
	services = map[string]struct{}{
		"Virtual Machines":    {},
		"Storage":             {},
		"Container Registry":  {},
		"Azure DNS":           {},
		"Load Balancer":       {},
		"Application Gateway": {},
		"NAT Gateway":         {},
		"Virtual Network":     {},
	}
)
