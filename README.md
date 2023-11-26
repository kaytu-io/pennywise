# Pennywise
Pennywise is an open-source program designed for calculating cloud resources' costs. It currently supports AWS and Azure. The project consists of a server, and a CLI program.

## Overview
### Server
The server component is intended for deployment on a server with MySQL database configuration. The server stores pricing data from AWS and Azure in a MySQL database. During startup, a migrator is run to set up the required tables and schema.

### CLI Program
The CLI program parses data from Terraform in three possible formats: a Terraform file, a Terraform plan file, or a Terraform plan JSON file. The cost request is then sent to the Pennywise server, and the result, comprising the cost (in the specified currency), is received.

## Getting Started
Follow these steps to get started with Pennywise:

### Server Deployment
Clone the Pennywise repository:

```shell
git clone https://github.com/kaytu-io/pennywise.git
```

```shell
cd pennywise/server
```

Set up the server configuration by editing config.yaml.

Run the migrator to set up the required database tables:

Start the Pennywise server:

```shell
go run server.go
```

### CLI Program
Clone the Pennywise repository (if not done already).

Navigate to the CLI program directory:

```shell
cd pennywise/cli
```

Create a Terraform plan file:

```shell
terraform plan -out planfile.tfplan
```

Generate a JSON file from the Terraform plan:

```shell
terraform show -json planfile.tfplan > tfplan.json
```
Run the CLI program with the generated Terraform plan file:

```shell
go run cli.go -f planfile.tfplan
```

and the run with the JSON file:

```shell
go run cli.go -json tfplan.json
```
