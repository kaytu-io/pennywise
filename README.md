<h1 align="center"> Pennywise </h1>

<p align="center">
    <img alt="GitHub repo size" src="https://img.shields.io/badge/License-Apache%202.0-blue?logo=github&style=for-the-badge&logo">
    <img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/kaytu-io/pennywise?logo=github&style=for-the-badge">
    <img alt="GitHub tag (with filter)" src="https://img.shields.io/github/v/tag/kaytu-io/pennywise?style=for-the-badge&logo=git">
    <img alt="GitHub go.mod Go version (subdirectory of monorepo)" src="https://img.shields.io/github/go-mod/go-version/kaytu-io/pennywise?style=for-the-badge&logo=go">
</p>

## Overview
Pennywise estimates the cost of cloud infrastructure before the actual deployment by analyzing Terraform and OpenTofu plan files. The current version supports commonly used resources in AWS and Azure.

![Cost Screenrecord](.github/assets/cost_screenrecord.gif)
## Getting Started

### Install


- Linux/macOS
    ```shell
    wget -qO - https://raw.githubusercontent.com/kaytu-io/pennywise/main/scripts/install.sh | sh
    ```
- Windows\
    Download and install manually from [releases](https://github.com/kaytu-io/pennywise/releases) 

### Sign-up for an account
Sign up for free on [kaytu](http://app.kaytu.io/)

### Login
Login to your kaytu account using:
```shell
pennywise login
``` 
this command will give you a link to open in your browser to help you sign-up and login into your kaytu account.

### Get costs
Run the following terraform commands to build the terraform plan json file:

```shell
terraform init
terraform plan -out tfplan.binary
terraform show -json tfplan.binary | jq > tfplan.json
```
And then estimate the project cost by passing the terraform plan json file to cost terraform command:
```shell
pennywise cost terraform --json-path path-to-json --usage path-to-usage-file
```

You can also specify the usage file path by usage tag.
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

To get a more detailed documents on CLI options and commands, please refer to [docs](./docs/pennywise.md)

## Contributing
Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pul requests to us.

## License
This project is licensed under the Apache License Version 2.0 - see the [LICENSE](LICENSE) file for details.
