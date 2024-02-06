## Usage


You can also specify the usage file path by --usage tag.
The usage file is responsible for getting usage details from user.
The usage file is supported in two types: `json` and `yaml`

The json file is as follows:
````json
{
  "azurerm_virtual_machine.windows": {
    "monthly_os_disk_operations": 1000000,
    "monthly_data_disk_operations": 2000000
  },
  "azurerm_virtual_machine.linux_withMonthlyHours": {
    "monthly_hrs": 100
  },
  "azurerm_virtual_machine.windows_withMonthlyHours": {
    "monthly_hrs": 100
  }
}
````
The yaml file is as follows:
````yaml
azurerm_virtual_machine.windows:
  monthly_os_disk_operations: 1000000
  monthly_data_disk_operations: 2000000
azurerm_virtual_machine.linux_withMonthlyHours:
  monthly_hrs: 100
azurerm_virtual_machine.windows_withMonthlyHours:
  monthly_hrs: 100
````
Also, here's the documents for supported usage parameters of each resource type:\
[aws-usage](./docs/aws-usage-parameters.md)\
[azure-usage](./docs/azure-usage-parameters.md)