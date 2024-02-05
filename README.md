<h1 align="center"> Pennywise </h1>

[![Build Status](https://travis-ci.org/kaytu-io/pennywise.svg?branch=master)](https://travis-ci.org/kaytu-io/pennywise)
[![GitHub release](https://img.shields.io/github/release/kaytu-io/pennywise.svg)](https://github.com/kaytu-io/pennywise/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Overview
Pennywise is an open-source program designed for calculating cloud resources' costs. It currently supports AWS and Azure. The project consists of a server, and a CLI program.

The CLI program parses data from Terraform in three possible formats: a Terraform file, a Terraform plan file, or a Terraform plan JSON file. The cost request is then sent to the Pennywise server, and the result, comprising the cost, is received.

## Getting Started

### Install pennywise client

Download the binary from [releases](https://github.com/kaytu-io/pennywise/releases).\
If you are using Linux you can run this command to install the CLI: 
```shell
wget -qO - https://raw.githubusercontent.com/kaytu-io/pennywise/main/scripts/install.sh | bash
```

### Usage

Login to your kaytu account using:
```shell
pennywise login
``` 
this command will give you a link to open in your browser to help you sign-up and login into your kaytu account.

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
