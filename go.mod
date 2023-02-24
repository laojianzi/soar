module github.com/laojianzi/soar

go 1.15

require (
	github.com/CorgiMan/json2 v0.0.0-20150213135156-e72957aba209
	github.com/astaxie/beego v1.12.3
	github.com/dchest/uniuri v0.0.0-20200228104902-7aecb25e1fe5
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813
	github.com/go-sql-driver/mysql v1.6.0
	github.com/kr/pretty v0.3.1
	github.com/percona/go-mysql v0.0.0-20210427141028-73d29c6da78c
	github.com/pingcap/parser v0.0.0-20210525032559-c37778aff307
	github.com/russross/blackfriday v1.6.0
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca
	github.com/tidwall/gjson v1.14.3
	gopkg.in/yaml.v2 v2.4.0
	vitess.io/vitess v0.0.0-20200325000816-eda961851d63
)

require (
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/pingcap/check v0.0.0-20200212061837-5e12011dc712 // indirect
	github.com/pingcap/errors v0.11.5-0.20201126102027-b0a155152ca3 // indirect
	github.com/pingcap/log v0.0.0-20210317133921-96f4fcab92a4 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/lint v0.0.0-20200130185559-910be7a94367 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	golang.org/x/sys v0.0.0-20220818161305-2296e01440c6 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20200224152610-e50cd9704f63 // indirect
	google.golang.org/grpc v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	honnef.co/go/tools v0.2.0 // indirect
)

replace (
	// fix potential security issue(CVE-2020-26160) introduced by indirect dependency.
	github.com/dgrijalva/jwt-go => github.com/form3tech-oss/jwt-go v3.2.6-0.20210809144907-32ab6a8243d7+incompatible
	// fix potential security issue(CVE-2020-27813) introduced by indirect dependency.
	github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.4.1
	// fix potential security issue(CVE-2018-1099) introduced by indirect dependency.
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738 => go.etcd.io/etcd v0.0.0-20200824191128-ae9734ed278b
)
