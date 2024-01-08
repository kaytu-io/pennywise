package aws

// SupportedServices is a list of all AWS services that are supported by Terracost.
var supportedServices = map[string]struct{}{
	"AmazonEC2":         {},
	"AmazonEFS":         {},
	"AmazonEKS":         {},
	"AmazonFSx":         {},
	"AmazonRDS":         {},
	"AmazonElastiCache": {},
	"AmazonCloudWatch":  {},
	"AWSELB":            {},
}

// IsServiceSupported returns true if the AWS service is valid and supported by Terracost (e.g. for ingestion.)
func IsServiceSupported(service string) bool {
	_, ok := supportedServices[service]
	return ok
}
