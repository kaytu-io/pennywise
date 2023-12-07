module github.com/kaytu-io/pennywise/server

go 1.18

require (
	github.com/brpaz/echozap v1.1.3
	github.com/cycloidio/sqlr v1.0.0
	github.com/dmarkham/enumer v1.5.3
	github.com/go-sql-driver/mysql v1.6.0
	github.com/labstack/echo/v4 v4.11.3
	github.com/lopezator/migrator v0.3.0
	github.com/machinebox/progress v0.2.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/shopspring/decimal v1.3.1
	github.com/spf13/afero v1.10.0
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/zap v1.10.0
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/net v0.17.0
	golang.org/x/text v0.13.0
	golang.org/x/tools v0.6.0
	gopkg.in/go-playground/validator.v9 v9.31.0

)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lib/pq v1.10.5 // indirect
	github.com/matryer/is v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pascaldekloe/name v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.2.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hashicorp/terraform => github.com/cycloidio/terraform v1.1.9-cy

replace github.com/gruntwork-io/terragrunt => github.com/cycloidio/terragrunt v0.0.0-20230905115542-1fe1ff682fd9
