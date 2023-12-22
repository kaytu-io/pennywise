package azurerm

//go:generate enumer -type=Service -output=service_string.go -linecomment=true

// Service is the type defining the services
type Service uint8

// List of all the supported services
const (
	VirtualMachines Service = iota // Virtual Machines
)

var (
	// The list of all services is https://azure.microsoft.com/en-us/services/, the left side is
	// the Family and the main content is the Services
	services = map[string]struct{}{
		"Key Vault":                     {},
		"Virtual Machines":              {},
		"Storage":                       {},
		"Container Registry":            {},
		"Azure DNS":                     {},
		"Load Balancer":                 {},
		"Application Gateway":           {},
		"NAT Gateway":                   {},
		"VPN Gateway":                   {},
		"Content Delivery Network":      {},
		"Virtual Network":               {},
		"Azure Cosmos DB":               {},
		"Azure Database for MariaDB":    {},
		"Azure Database for MySQL":      {},
		"Azure Database for PostgreSQL": {},
		"SQL Database":                  {},
		"SQL Managed Instance":          {},
		"Azure Kubernetes Service":      {},
		"Automation":                    {},
		"Logic Apps":                    {},
	}
)
