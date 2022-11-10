module github.com/romberli/sql-parser-go

go 1.16

replace github.com/spf13/pflag v1.0.5 => github.com/romberli/pflag v1.0.6-alpha

require (
	github.com/asaskevich/govalidator v0.0.0-20200819183940-29e1ff8eb0bb
	github.com/hashicorp/go-multierror v1.1.0
	github.com/pingcap/errors v0.11.5-0.20210425183316-da1aaba5fb63
	github.com/romberli/go-util v0.3.13-0.20211223065033-fdb92c74739a
	github.com/romberli/log v1.0.20
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
)
