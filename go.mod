module github.com/keitap/github-asana-request-review-action

go 1.16

replace bitbucket.org/mikehouston/asana-go => github.com/keitap/asana-go v0.0.0-20210425151158-ba9ef27944c8

require (
	bitbucket.org/mikehouston/asana-go v0.0.0-20201102222432-715318d0343a
	github.com/google/go-github/v35 v35.1.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/oauth2 v0.0.0-20210413134643-5e61552d6c78
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/yaml.v2 v2.3.0
)
