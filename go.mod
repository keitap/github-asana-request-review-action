module github.com/keitap/github-asana-request-review-action

go 1.25.5

replace bitbucket.org/mikehouston/asana-go => github.com/keitap/asana-go v0.0.0-20210425173123-936f596fb971

require (
	bitbucket.org/mikehouston/asana-go v0.0.0-20250814164459-d578140040b5
	github.com/google/go-github/v74 v74.0.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/oauth2 v0.35.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/xid v1.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
