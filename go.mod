module github.com/romberli/go-template

go 1.16

replace github.com/spf13/pflag v1.0.5 => github.com/romberli/pflag v1.0.6-alpha

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/asaskevich/govalidator v0.0.0-20200819183940-29e1ff8eb0bb
	github.com/gin-gonic/gin v1.5.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/romberli/go-util v0.3.12-0.20211008113222-c672e84666e3
	github.com/romberli/log v1.0.20
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.6-0.20200529100950-7c765ddd0476
	go.uber.org/zap v1.16.0
)
