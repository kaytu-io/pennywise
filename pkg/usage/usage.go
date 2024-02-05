package usage

import "strings"

const (
	// Key is the key used to set the usage
	// on the values passed to the resources
	Key string = "pennywise_usage"
)

// Default is the default Usage that will be used if none is configured
var Default = Usage{
	"aws_eks_node_group": map[string]interface{}{
		"instances":                        15,
		"operating_system":                 "linux",
		"reserved_instance_type":           "standard",
		"reserved_instance_term":           "1_year",
		"reserved_instance_payment_option": "partial_upfront",
		"monthly_cpu_credit_hrs":           350,
		"vcpu_count":                       2,
	},
	"aws_efs_file_system": map[string]interface{}{
		"storage_gb":                         180,
		"infrequent_access_storage_gb":       10,
		"monthly_infrequent_access_read_gb":  20,
		"monthly_infrequent_access_write_gb": 30,
	},
	"aws_fsx_openzfs_file_system": map[string]interface{}{
		"backup_storage_gb": 1024,
	},
	"aws_fsx_windows_file_system": map[string]interface{}{
		"backup_storage_gb": 1024,
	},
	"aws_fsx_ontap_file_system": map[string]interface{}{
		"backup_storage_gb": 1024,
	},
	"aws_fsx_lustre_file_system": map[string]interface{}{
		"backup_storage_gb": 1024,
	},
	"aws_nat_gateway": map[string]interface{}{
		"monthly_data_processed_gb": 10,
	},
	"azurerm_virtual_machine": map[string]interface{}{
		"monthly_os_disk_operations":   1000000,
		"monthly_data_disk_operations": 2000000,
		"monthly_hours":                730,
	},
	"azurerm_managed_disk": map[string]interface{}{
		"monthly_disk_operations": 20000,
	},
	"azurerm_linux_virtual_machine": map[string]interface{}{
		"monthly_hours": 730,
	},
	"azurerm_windows_virtual_machine": map[string]interface{}{
		"monthly_hours": 730,
	},
	"azurerm_lb": map[string]interface{}{
		"monthly_data_proceed": 1000,
	},
}

// Usage is the struct defining all the configure usages
type Usage map[string]map[string]interface{}

// GetUsage will return the usage from the resource rt (ex: aws_instance)
func (u Usage) GetUsage(rt string, addr string) map[string]interface{} {
	us, ok := u[addr]
	if ok {
		return us
	}
	addr = strings.Split(addr, "[")[0]
	us, ok = u[addr]
	if ok {
		return us
	}
	us, ok = u[rt]
	if ok {
		return us
	}

	return nil
}
